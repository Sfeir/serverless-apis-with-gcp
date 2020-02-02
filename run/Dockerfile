FROM golang:1.13-alpine

# Set the Current Working Directory inside the container
WORKDIR /go/src/app-inventory
COPY . .

RUN apk update && apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app-inventory"]
