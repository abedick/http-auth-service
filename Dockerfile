FROM golang:1.11.0-alpine3.8

# Configure NGINX
ADD conf/nginx.conf /etc/nginx/nginx.conf
RUN set -x ; \
  addgroup -g 82 -S www-data ; \
  adduser -u 82 -D -S -G www-data www-data && exit 0 ; exit 1

# Configure supervisord
ADD conf/supervisord.conf /etc/supervisord.conf
RUN mkdir -p /var/log/supervisor/

# Gather Dependencies
RUN apk add nginx supervisor curl git mercurial \
	&& go get github.com/dgrijalva/jwt-go \
	&& go get github.com/gorilla/mux \
	&& go get github.com/tomasen/realip \
	&& go get github.com/BurntSushi/toml \
	&& go get github.com/sethvargo/go-password/password 

# Configure and build cabal-service
ENV DIR=/opt/service
ENV BUILDDIR=/tmp/service

RUN mkdir -p $DIR \
    && mkdir -p $BUILDDIR

ADD ./src $BUILDDIR/src
ADD ./conf/auth.toml $DIR/conf/auth.toml
ADD ./scripts/startService.sh $DIR/startService.sh

RUN go build -o cabal-service $BUILDDIR/src/*.go \
    && mv cabal-service $DIR \
    && chmod +x $DIR/startService.sh \
    && rm -rf $BUILDDIR

EXPOSE 3000
WORKDIR /
CMD		["supervisord", "--nodaemon", "--configuration", "/etc/supervisord.conf"]