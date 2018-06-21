FROM reginleiff/hb-go:latest as builder
WORKDIR /go/src/github.com/reginleiff/go-tic-tac-toe
ADD . /go/src/github.com/reginleiff/go-tic-tac-toe
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/reginleiff/go-tic-tac-toe/app .
EXPOSE 3000
ENTRYPOINT ["./app"]

