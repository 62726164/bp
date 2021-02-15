FROM alpine:3.13
 
RUN mkdir -p /tmp

WORKDIR /tmp

EXPOSE 9379

COPY check /tmp/
COPY privkey.pem /tmp/
COPY cert.pem /tmp/

CMD ["/tmp/check"]
