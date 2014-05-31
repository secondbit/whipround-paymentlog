FROM ubuntu:12.04
RUN apt-get update -q
RUN DEBIAN_FRONTEND=noninteractive apt-get install -qy build-essential curl git
RUN curl -s https://go.googlecode.com/files/go1.2.1.src.tar.gz | tar -v -C /usr/local -xz
RUN cd /usr/local/go/src && ./make.bash --no-clean 2>&1
ENV PATH /usr/local/go/bin:$PATHi
ENV GOPATH /opt/go/
ADD . /opt/go/src/code.whipround.net/paymentlogs
RUN cd /opt/go/src/code.whipround.net/paymentlogs && go get -d -v ./... && go build -v ./... && go test .
