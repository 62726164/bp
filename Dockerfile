FROM debian:latest

RUN apt-get update && apt-get -y upgrade
RUN apt-get install ca-certificates -y

EXPOSE 9379

COPY check /usr/local/bin/
COPY privkey.pem /tmp/
COPY cert.pem /tmp/

CMD ["/usr/local/bin/check"]
