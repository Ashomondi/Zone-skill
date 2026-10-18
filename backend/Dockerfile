FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o /out/zoneskill-api ./cmd/api

FROM alpine:3.21
RUN adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=build /out/zoneskill-api /app/zoneskill-api
USER app
EXPOSE 8080
ENTRYPOINT ["/app/zoneskill-api"]
