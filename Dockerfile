FROM golang:alpine

WORKDIR /go/src/github.com/textileio/textile-go
COPY . .

RUN go build -i -o textile textile.go

CMD ["textile-go -d -n -g=127.0.0.1:9000 --cafe-bind-addr=0.0.0.0:8000 --cafe-token-secret=swarmmmmm --cafe-referral-key=woohoo --cafe-db-hosts=0.0.0.0:27017 --cafe-db-name=textile_db"]
