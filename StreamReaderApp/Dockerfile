FROM golang:alpine as build

RUN apk update &&\
    apk add --no-progress --no-cache \
    gcc \
    librdkafka-dev \
    openssl-libs-static \
    zlib-static \
    zstd-libs \
    libsasl \
    librdkafka-static \
    lz4-dev \
    lz4-static \
    zstd-static \
    libc-dev \
    musl-dev

ENV PROJECT_PATH="/go/src/github.com/selcux/data-pipeline-sample/StreamReaderApp"

WORKDIR $PROJECT_PATH

ENV GOOS=linux GOARCH=amd64

COPY go.mod .

RUN go mod download

COPY . $PROJECT_PATH

RUN go build -tags musl -ldflags '-w -extldflags "-static"' -o /app ./cmd

FROM scratch
COPY --from=build /app /app
ENTRYPOINT ["/app"]