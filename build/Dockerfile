FROM golang:1.21.0-alpine

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/user-segmentation-service

ENTRYPOINT [ "app" ]