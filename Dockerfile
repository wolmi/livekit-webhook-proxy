FROM golang:alpine as build

RUN apk add --update --no-cache upx

WORKDIR /go/src/livekit-webhook-proxy

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY cmd ./cmd
COPY utils ./utils
COPY types ./type
RUN go build -ldflags="-w -s" -o /go/bin/livekit-webhook-proxy .
RUN upx /go/bin/livekit-webhook-proxy

FROM alpine:latest

WORKDIR /app
RUN addgroup -S livekit-webhook-proxy && adduser -S livekit-webhook-proxy -G livekit-webhook-proxy && chown livekit-webhook-proxy:livekit-webhook-proxy /app
USER livekit-webhook-proxy

COPY --from=build /go/bin/livekit-webhook-proxy .

ENTRYPOINT [ "/app/livekit-webhook-proxy" ]
