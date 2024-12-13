# Stage 1: Dependency Download
FROM golang:1.22-alpine as dep
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Build
FROM dep as build

ARG BUILD_TIME
ARG COMMIT
ARG VERSION

WORKDIR /app
COPY --from=dep /app /app
COPY . .
RUN apk add --no-cache make
#RUN go build -o dnsMasqAPI main.go
RUN BUILD_TIME=$BUILD_TIME COMMIT=$COMMIT VERSION=$VERSION make build

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/dnsMasqAPI /app/

ENTRYPOINT ["./dnsMasqAPI", "server"]
