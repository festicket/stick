FROM alpine:3.5

RUN apk add --no-cache go=1.7.3-r0 bash make
ADD . /src/

RUN export GOPATH=$(pwd)
WORKDIR src
