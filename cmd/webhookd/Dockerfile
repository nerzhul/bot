FROM arm64v8/golang:1.10.3

COPY . /go/src/gitlab.com/nerzhul/bot
RUN cd /go/src/gitlab.com/nerzhul/bot \
        && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make webhook

FROM arm64v8/alpine:3.7
RUN addgroup -g 750 -S bot-webhook && adduser -u 750 -D -S -G bot-webhook bot-webhook
USER bot-webhook
WORKDIR /var/lib/bot-webhook
COPY --from=0 /artifacts/linux_arm64/webhook /usr/bin/bot-webhook
EXPOSE 8080

CMD ["/usr/bin/bot-webhook", "--config=/etc/bot/webhook.yml"]

