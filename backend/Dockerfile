FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/zoneskill-api ./cmd/api

FROM alpine:3.21
RUN apk --no-cache add ca-certificates && \
    adduser -D -u 10001 app
WORKDIR /app
RUN mkdir -p /app/data /data && \
    chown -R app:app /app /data
COPY --from=build --chown=app:app /out/zoneskill-api /app/zoneskill-api
USER app
EXPOSE 8080
ENTRYPOINT ["/app/zoneskill-api"]

