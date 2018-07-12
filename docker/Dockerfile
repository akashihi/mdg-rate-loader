FROM octoblu/alpine-ca-certificates:latest
LABEL Description="MDG currency rates loader image"
COPY mdg-rate-loader /
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENTRYPOINT ["/mdg-rate-loader","-host","postgres"]