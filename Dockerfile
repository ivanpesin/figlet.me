FROM golang AS builder

RUN apt-get update
RUN cd /go/src && git clone https://github.com/ivanpesin/figlet.me.git
RUN cd /go/src/figlet.me && go-wrapper download && go-wrapper install

FROM alpine
COPY --from=builder /go/bin/figlet.me /bin
RUN apk install figlet

ENV FIGLET_BIN /usr/bin/figlet
ENV FIGLET_DIR=/usr/share/figlet
EXPOSE 8080

CMD /bin/figlet.me