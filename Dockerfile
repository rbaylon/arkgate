# Build Stage
FROM golang:1.21.1-alpine3.17 AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /sangapi

# Run Stage
FROM alpine:3.17

WORKDIR /app

COPY --from=build-env /sangapi /sangapi
COPY .env .

EXPOSE 3333

ENTRYPOINT [ "/sangapi" ]
