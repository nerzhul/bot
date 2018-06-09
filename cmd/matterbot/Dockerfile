FROM arm64v8/golang:1.10.3

COPY . /go/src/gitlab.com/nerzhul/bot
RUN cd /go/src/gitlab.com/nerzhul/bot \
        && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make matterbot

FROM arm64v8/alpine:3.7
RUN addgroup -g 751 -S matterbot && adduser -u 750 -D -S -G matterbot matterbot
USER matterbot
WORKDIR /var/lib/matterbot
COPY --from=0 /artifacts/linux_arm64/matterbot /usr/bin/matterbot

CMD ["/usr/bin/matterbot", "--config=/etc/bot/matterbot.yml"]
