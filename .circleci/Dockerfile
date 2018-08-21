FROM golang:1.10.3

# replace shell with bash so we can source files
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# install dependencies
RUN apt-get update \
    && apt-get install -y gcc \
    && apt-get install -y curl \
    && apt-get install -y mingw-w64 \
    && apt-get -y autoclean

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# install gx
RUN go get -u github.com/whyrusleeping/gx \
    && go get -u github.com/whyrusleeping/gx-go
