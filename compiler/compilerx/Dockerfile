FROM alpine:latest

LABEL Maintainer="Myriad-Dreamin <camiyoru@gmail.com>"

# ustc source
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# c/c++
RUN apk add --no-cache g++
# java openjdk8
RUN apk add openjdk8
ENV PATH $PATH:/usr/lib/jvm/java-1.8-openjdk/bin



VOLUME [ /codes /compiler_tools ]

EXPOSE 23366

ENTRYPOINT [ "/run.sh" ]

COPY run.sh /run.sh
COPY compiler /compiler



