FROM alpine:3.13.1

RUN apk -U upgrade
 
RUN mkdir -p /tmp

WORKDIR /tmp

EXPOSE 9379

COPY check /tmp/
COPY privkey.pem /tmp/
COPY cert.pem /tmp/

CMD ["/tmp/check"]
