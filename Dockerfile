FROM golang:alpine AS build
ENV GOPATH=/home/go
RUN apk add upx
WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o ./app
RUN upx ./app

# Install hal
RUN rm -f go.mod
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go install -ldflags "-s -w -extldflags '-static'" github.com/dansimau/hal/cmd/hal@latest
RUN ls -la $GOPATH/bin/linux_amd64
RUN cp $GOPATH/bin/linux_amd64/hal ./hal
RUN upx ./hal

FROM scratch
COPY --from=build /build/app /app
COPY --from=build /build/hal /hal

ENTRYPOINT ["/app"]
