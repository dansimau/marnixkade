FROM golang:alpine AS build
RUN apk add upx
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o ./app
RUN upx ./app

FROM scratch
COPY --from=build /build/app /app
COPY hal.yaml /

ENTRYPOINT ["/app"]
