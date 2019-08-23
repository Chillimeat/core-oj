FROM alpine:latest

LABEL Maintainer="Myriad-Dreamin <camiyoru@gmail.com>"

RUN apk add --no-cache g++

VOLUME [ /codes /compiler_tools ]

EXPOSE 23366

ENTRYPOINT [ "/run.sh" ]

COPY run.sh /run.sh
COPY compiler /compiler



