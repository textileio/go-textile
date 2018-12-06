FROM golang:1.11-stretch
MAINTAINER Sander Pick <sander@textile.io>

# replace shell with bash so we can source files
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# install dependencies
RUN apt-get update \
    && apt-get install -y curl \
    && apt-get -y autoclean

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# install gx
RUN go get -u github.com/whyrusleeping/gx \
    && go get -u github.com/whyrusleeping/gx-go

# get source
ENV SRC_DIR /go/src/github.com/textileio/textile-go
COPY . $SRC_DIR

# build source
RUN cd $SRC_DIR \
  && make setup \
  && make build

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

# Get the textile binary, entrypoint script, and TLS CAs from the build container.
ENV SRC_DIR /go/src/github.com/textileio/textile-go
COPY --from=0 $SRC_DIR/textile /usr/local/bin/textile
#COPY --from=0 $SRC_DIR/bin/container_daemon /usr/local/bin/start_ipfs
COPY --from=0 /tmp/su-exec/su-exec /sbin/su-exec
COPY --from=0 /tmp/tini /sbin/tini
COPY --from=0 /etc/ssl/certs /etc/ssl/certs

# This shared lib (part of glibc) doesn't seem to be included with busybox.
COPY --from=0 /lib/x86_64-linux-gnu/libdl-2.24.so /lib/libdl.so.2

# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 40600
# Web Gateway; can be exposed publicly with a proxy, e.g. as https://ipfs.example.org
EXPOSE 5050
# Swarm Websockets; must be exposed publicly when the node is listening using the websocket transport (/ipX/.../tcp/8081/ws).
EXPOSE 8081

# Create the fs-repo directory and switch to a non-privileged user.
ENV TEXTILE_PATH /data/textile
RUN mkdir -p $TEXTILE_PATH \
  && adduser -D -h $TEXTILE_PATH -u 1000 -G users textile \
  && chown textile:users $TEXTILE_PATH

# Expose the fs-repo as a volume.
# start_ipfs initializes an fs-repo if none is mounted.
# Important this happens after the USER directive so permission are correct.
VOLUME $TEXTILE_PATH

# The default logging level
#ENV IPFS_LOGGING ""

# This just makes sure that:
# 1. There's an fs-repo, and initializes one if there isn't.
# 2. The API and Gateway are accessible from outside the container.
#ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/start_ipfs"]

# Execute the daemon subcommand by default
CMD ["daemon"]
