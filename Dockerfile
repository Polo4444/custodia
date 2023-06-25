FROM alpine:latest
LABEL developers="Polo"
LABEL code="v1"
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

LABEL developer1="polo@keystonesolutions.io"
COPY ./app ./

EXPOSE 8080
RUN chmod +x /app
CMD [ "/app" ]