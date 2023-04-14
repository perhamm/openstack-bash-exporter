FROM golang:1.20.3
ADD cmd/openstack-bash-exporter ./src/cmd/openstack-bash-exporter
WORKDIR /go/src/cmd/openstack-bash-exporter
RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o openstack-bash-exporter .

FROM alpine:3.17
WORKDIR /root/

RUN apk add --update \
  bash \
  bash-doc \
  bash-completion \
  bind-tools \
  python-dev \
  py-pip \
  py-setuptools \
  ca-certificates \
  gcc \
  libffi-dev \
  openssl-dev \
  musl-dev \
  linux-headers \
  && pip install --upgrade --no-cache-dir pip setuptools python-openstackclient \
  && apk del gcc musl-dev linux-headers libffi-dev \
  && rm -rf /var/cache/apk/*

COPY --from=0 /go/src/cmd/openstack-bash-exporter/openstack-bash-exporter .
COPY ./scripts/* /root/scripts/
CMD ["./openstack-bash-exporter"]
