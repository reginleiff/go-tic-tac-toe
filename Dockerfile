FROM alpine:latest

ADD . /src
WORKDIR /src

EXPOSE 3000

ENTRYPOINT ["/src/app"]


