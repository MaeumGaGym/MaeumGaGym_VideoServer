FROM golang:1.21.5-alpine AS build

RUN apk add -U --no-cache ca-certificates gcc musl-dev ffmpeg ffmpeg-libs && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o ffmpeg_app

FROM alpine

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/bin/ffmpeg /usr/local/bin/ffmpeg
COPY --from=build /usr/lib/* /usr/local/lib/

RUN chmod +x /usr/local/bin/ffmpeg

WORKDIR /app

COPY --from=build /app/ffmpeg_app .
COPY --from=build /app/db/jobs.db /app/db/

ENTRYPOINT ["./ffmpeg_app"]