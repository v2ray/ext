FROM golang:latest as builder

MAINTAINER admin@v2ray.com

RUN go get -u v2ray.com/core/...
RUN go get -u v2ray.com/ext/...
RUN go install v2ray.com/ext/tools/build/vbuild
RUN vbuild -dir /usr/bin/v2ray


FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
# Change TimeZone
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
# Clean APK cache
RUN rm -rf /var/cache/apk/*

RUN mkdir /usr/bin/v2ray/
RUN mkdir /etc/v2ray/
RUN mkdir /var/log/v2ray/

COPY --from=builder /usr/bin/v2ray  /usr/bin/v2ray

ENV PATH /usr/bin/v2ray/:$PATH

EXPOSE 8000
COPY config.json /etc/v2ray/config.json

CMD ["/usr/bin/v2ray/v2ray", "-config=/etc/v2ray/config.json"]
