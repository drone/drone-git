FROM ppc64le/alpine:3.14
RUN apk add --no-cache ca-certificates git git-lfs openssh curl perl aws-cli sudo

ADD posix/* /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/clone"]
