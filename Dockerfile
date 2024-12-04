FROM scratch
ADD build/marnixkade /marnixkade
ENTRYPOINT ["/marnixkade"]
