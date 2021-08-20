FROM alpine:3.14

USER nobody
COPY riley /riley

ENTRYPOINT ["/riley"]
