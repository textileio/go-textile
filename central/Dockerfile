FROM golang:alpine

WORKDIR /go/src/github.com/textileio/textile-go/central
COPY . .

RUN go install -v ./...

CMD ["central"]
