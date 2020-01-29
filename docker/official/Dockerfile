FROM alpine:3.11 AS downloader

RUN apk add -U --no-cache curl

WORKDIR /etc/v2ray
# Configurable arguments defined at build time. User can override them by docker build --build-arg ARG_NAME=ARG_VALUE
# Default OS type = linux
ARG OS=linux
# Default architecture = 64
ARG ARCH=64
# Download and unpack latest release from official github release
RUN V2RAY_FILE=v2ray-${OS}-${ARCH}.zip \
    VERSION=$(curl -s https://github.com/v2ray/v2ray-core/releases/latest |grep -oE '\d\.\d+\.\d+') \
    && echo "Got latest version = ${VERSION}, file = ${V2RAY_FILE}" \
    && curl -Lo $V2RAY_FILE https://github.com/v2ray/v2ray-core/releases/download/v${VERSION}/$V2RAY_FILE \
    && unzip $V2RAY_FILE \
    && rm $V2RAY_FILE

# Use distroless as minimal base image to package the executable binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# Note: Use gcr.azk8s.cn as repo-mirror if gcr.io is blocked by China GFW.
FROM gcr.io/distroless/static:latest
WORKDIR /etc/v2ray
COPY --from=downloader /etc/v2ray .
COPY config.json /etc/v2ray/config.json
ENTRYPOINT ["/etc/v2ray/v2ray"]
CMD ["-config=/etc/v2ray/config.json"]
