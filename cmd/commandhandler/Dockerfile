FROM arm64v8/golang:1.10.2

COPY . /go/src/gitlab.com/nerzhul/bot
RUN cd /go/src/gitlab.com/nerzhul/bot \
        && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make commandhandler

FROM arm64v8/alpine:3.7
RUN addgroup -g 751 -S commandhandler && adduser -u 754 -D -S -G commandhandler commandhandler
USER commandhandler
WORKDIR /var/lib/commandhandler
COPY --from=0 /artifacts/linux_arm64/commandhandler /usr/bin/commandhandler

CMD ["/usr/bin/commandhandler", "--config=/etc/bot/commandhandler.yml"]
