FROM golang:alpine as build-env
RUN apk --no-cache add tzdata
RUN apk add build-base
RUN apk add -U --no-cache ca-certificates
WORKDIR /AppleProductMonitor
ADD . /AppleProductMonitor
RUN cd /AppleProductMonitor/bot && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot.app

FROM scratch
WORKDIR /app
COPY --from=build-env /AppleProductMonitor/bot/bot.app /app
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/bot.app"]

