FROM alpine:latest

LABEL Maintainer="Myriad-Dreamin <camiyoru@gmail.com>"

# ustc source
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# access control list
# RUN apk add --no-cache acl

# c/c++
RUN apk add --no-cache libstdc++ && apk add --no-cache libgcc
# python2/python3
RUN apk add --no-cache python  && apk add --no-cache python3
# java openjdk8
RUN apk add openjdk8
ENV PATH $PATH:/usr/lib/jvm/java-1.8-openjdk/bin

# ENV PATH $PATH:$HOME/bin:$JAVA_HOME/bin


VOLUME /codes /checker_tools /judger_tools /problems

RUN mkdir -p /judger_tools/socks


ENTRYPOINT [ "/run.sh" ]

COPY run.sh /run.sh
COPY judger /judger
# COPY sos-checker /sos-checker


