FROM golang AS builder

RUN apt-get update
RUN cd /go/src && git clone https://github.com/ivanpesin/figlet.me.git
ENV CGO_ENABLED 0
RUN cd /go/src/figlet.me && go-wrapper download && go-wrapper install

FROM alpine
RUN apk update
COPY --from=builder /go/bin/figlet.me /bin
RUN apk add figlet
COPY --from=builder /go/src/figlet.me/fonts /usr/share/figlet/fonts

ENV FIGLET_BIN /usr/bin/figlet
ENV FIGLET_DIR=/usr/share/figlet
EXPOSE 8080

CMD /bin/figlet.me