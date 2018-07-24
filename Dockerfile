FROM golang:alpine

WORKDIR /go/src/github.com/textileio/textile-go
COPY . .

RUN go install -v ./...

CMD ["textile-go -d -n -g=0.0.0.0:9000 --swarm-ports=4001,4002 --cafe-bind-addr=0.0.0.0:8000 --cafe-token-secret=swarmmmmm --cafe-referral-key=woohoo --cafe-db-hosts=mongo:27017 --cafe-db-name=textile_db"]
