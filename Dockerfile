FROM alpine:latest

ADD build/marnixkade /marnixkade

ENTRYPOINT ["/marnixkade"]
