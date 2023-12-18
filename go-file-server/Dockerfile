FROM golang:1.21.5-alpine AS build

RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app

FROM alpine

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

COPY --from=build /app/app .

ENTRYPOINT ["./app"]