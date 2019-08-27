FROM alpine:latest

LABEL Maintainer="Myriad-Dreamin <camiyoru@gmail.com>"

# RUN su -
# RUN groupadd core-qwqs
# RUN useradd -g core-qwqs core-qwqer
# USER core-qwqer:coreqwqs



RUN apk add --no-cache libstdc++ && apk add --no-cache libgcc


VOLUME [ /codes /checker_tools ]

ENTRYPOINT [ "/run.sh" ]

COPY run.sh /run.sh
COPY judger /judger
COPY sos-checker /sos-checker


