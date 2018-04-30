FROM golang:1.10

# install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN adduser --disabled-password --gecos '' api
USER api

WORKDIR /go/src/app
COPY . .

#RUN dep ensure
RUN go install -v

CMD ["app"]

#RUN go get github.com/pilu/fresh
#RUN go-wrapper download
#RUN go-wrapper install

#CMD [ "fresh" ]
