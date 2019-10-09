FROM golang:1.13.1-stretch
MAINTAINER Sander Pick <sander@textile.io>

# This is (in large part) copied (with love) from
# https://hub.docker.com/r/ipfs/go-ipfs/dockerfile

# Get source
ENV SRC_DIR /go-textile

# Download packages first so they can be cached.
COPY go.mod go.sum $SRC_DIR/
RUN cd $SRC_DIR \
  && go mod download

COPY . $SRC_DIR

# build source
RUN cd $SRC_DIR \
  && go install github.com/ahmetb/govvv \
  && make textile

# Get su-exec, a very minimal tool for dropping privileges,
# and tini, a very minimal init daemon for containers
ENV SUEXEC_VERSION v0.2
ENV TINI_VERSION v0.16.1
RUN set -x \
  && cd /tmp \
  && git clone https://github.com/ncopa/su-exec.git \
  && cd su-exec \
  && git checkout -q $SUEXEC_VERSION \
  && make \
  && cd /tmp \
  && wget -q -O tini https://github.com/krallin/tini/releases/download/$TINI_VERSION/tini \
  && chmod +x tini

# Get the TLS CA certificates, they're not provided by busybox.
RUN apt-get update && apt-get install -y ca-certificates

# Now comes the actual target image, which aims to be as small as possible.
FROM busybox:1-glibc
MAINTAINER Sander Pick <sander@textile.io>

# Get the ipfs binary, entrypoint script, and TLS CAs from the build container.
ENV SRC_DIR /go-textile
COPY --from=0 /go/bin/textile /usr/local/bin/textile
COPY --from=0 $SRC_DIR/bin/container_daemon /usr/local/bin/start_textile
COPY --from=0 /tmp/su-exec/su-exec /sbin/su-exec
COPY --from=0 /tmp/tini /sbin/tini
COPY --from=0 /etc/ssl/certs /etc/ssl/certs

# This shared lib (part of glibc) doesn't seem to be included with busybox.
COPY --from=0 /lib/x86_64-linux-gnu/libdl-2.24.so /lib/libdl.so.2

# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Swarm Websockets; must be exposed publicly when the node is listening using the websocket transport (/ipX/.../tcp/8081/ws).
EXPOSE 8081
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 40600
# Web Gateway;
EXPOSE 5050
# Profiling API;
EXPOSE 6060

# Create the fs-repo directory
ENV TEXTILE_PATH /data/textile
RUN mkdir -p $TEXTILE_PATH \
  && adduser -D -h $TEXTILE_PATH -u 1000 -G users textile \
  && chown textile:users $TEXTILE_PATH

# Switch to a non-privileged user
USER textile

# Expose the fs-repo as a volume.
# start_textile initializes an fs-repo if none is mounted.
# Important this happens after the USER directive so permission are correct.
VOLUME $TEXTILE_PATH

# Init opts
ENV INIT_ARGS \
  --repo=$TEXTILE_PATH \
  --swarm-ports=4001,8081 \
  --api-bind-addr=0.0.0.0:40600 \
  --gateway-bind-addr=0.0.0.0:5050 \
  --profile-bind-addr=0.0.0.0:6060 \
  --debug

# This just makes sure that:
# 1. There's an fs-repo, and initializes one if there isn't.
# 2. The API and Gateway are accessible from outside the container.
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/start_textile"]

# Execute the daemon subcommand by default
CMD ["daemon", "--repo=/data/textile", "--debug"]
