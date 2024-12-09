FROM golang:alpine AS build
RUN apk add upx
WORKDIR /build
COPY . .
RUN sed -i -e '/^replace/s/^/\/\//' go.mod
RUN go get github.com/dansimau/hal
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o ./app
RUN upx ./app

FROM scratch
COPY --from=build /build/app /app

ENTRYPOINT ["/app"]
