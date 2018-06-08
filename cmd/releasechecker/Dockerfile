FROM arm64v8/golang:1.10.2

COPY . /go/src/gitlab.com/nerzhul/bot
RUN cd /go/src/gitlab.com/nerzhul/bot \
        && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 make releasechecker

FROM arm64v8/alpine:3.7
RUN addgroup -g 778 -S releasechecker && adduser -u 778 -D -S -G releasechecker releasechecker
USER releasechecker
WORKDIR /var/lib/releasechecker
COPY --from=0 /artifacts/linux_arm64/releasechecker /usr/bin/releasechecker
COPY res/migrations /var/lib/releasechecker/migrations

CMD ["/usr/bin/releasechecker", "--config=/etc/bot/releasechecker.yml"]
