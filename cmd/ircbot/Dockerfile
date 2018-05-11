FROM arm64v8/golang:1.10.2

COPY . /go/src/gitlab.com/nerzhul/bot
RUN cd /go/src/gitlab.com/nerzhul/bot \
        && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make ircbot

FROM arm64v8/alpine:3.7
RUN addgroup -g 751 -S ircbot && adduser -u 750 -D -S -G ircbot ircbot
USER ircbot
WORKDIR /var/lib/bot-ircbot
COPY --from=0 /artifacts/linux_arm64/ircbot /usr/bin/ircbot
COPY --from=0 /etc/ssl/certs /etc/ssl/certs
COPY --from=0 /usr/share/ca-certificates /usr/share/ca-certificates

CMD ["/usr/bin/ircbot", "--config=/etc/bot/ircbot.yml"]
