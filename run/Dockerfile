FROM golang:1.8-alpine

WORKDIR /go/src/hello-devfest
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["hello-devfest"]
