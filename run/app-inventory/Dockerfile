FROM golang:1.8-alpine

WORKDIR /go/src/app-inventory
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app-inventory"]
