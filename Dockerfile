FROM golang:latest as builder
WORKDIR /go/src/github.com/reginleiff/go-tic-tac-toe
ADD . /go/src/github.com/reginleiff/go-tic-tac-toe
RUN go get -d -v ./... 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/github.com/reginleiff/go-tic-tac-toe/app .
COPY ./config.toml .
EXPOSE 3000
ENTRYPOINT ["./app"]

