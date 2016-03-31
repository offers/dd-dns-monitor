FROM phusion/baseimage:0.9.18

ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
ENV GO15VENDOREXPERIMENT=1
ENV MONITOR_ROOT=/go/src/github.com/offers/dd-dns-monitor

RUN apt-get update \
    && apt-get install -y wget git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# install and setup go
RUN wget https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.6.linux-amd64.tar.gz \
    && rm go1.6.linux-amd64.tar.gz \
    && mkdir /go

# install glide
RUN mkdir /glide \
    && cd /glide \
    && wget https://github.com/Masterminds/glide/releases/download/0.10.1/glide-0.10.1-linux-amd64.tar.gz \
    && tar -xzf glide-0.10.1-linux-amd64.tar.gz \
    && mv linux-amd64/glide /usr/local/bin/ \
    && rm -rf /glide

# build dd-dns-monitor
ADD . /go/src/github.com/offers/dd-dns-monitor
RUN cd $MONITOR_ROOT \
    && glide install \
    && go install

RUN mkdir /etc/service/dd-dns-monitor \
    && cp $MONITOR_ROOT/run /etc/service/dd-dns-monitor/run
