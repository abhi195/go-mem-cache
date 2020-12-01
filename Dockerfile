FROM golang:1.15.5-alpine as GO_BUILD

# set working dir
WORKDIR /tmp/go-mem-cache

COPY . .

# build go executable
RUN go test \
    && go build -o ./out/cache-server server/*

# start fresh with only alpine
FROM alpine:3.12.1

# copy go executable from previous build
COPY --from=GO_BUILD /tmp/go-mem-cache/out/cache-server /app/cache-server

# expose port
EXPOSE 8080

# start server
ENTRYPOINT [ "/app/cache-server" ]