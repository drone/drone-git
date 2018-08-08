FROM alpine:3.6
RUN apk add --no-cache ca-certificates git openssh curl perl

ADD posix/* /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/clone"]
