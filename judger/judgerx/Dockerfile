FROM alpine:latest

LABEL Maintainer="Myriad-Dreamin <camiyoru@gmail.com>"


RUN apk add --no-cache acl
RUN apk add --no-cache libstdc++ && apk add --no-cache libgcc

VOLUME /codes /checker_tools /judger_tools /problems

RUN mkdir -p /judger_tools/socks


ENTRYPOINT [ "/run.sh" ]

COPY run.sh /run.sh
COPY judger /judger
# COPY sos-checker /sos-checker


