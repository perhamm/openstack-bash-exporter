FROM golang:1.20.3
ADD . /go/src/
WORKDIR /go/src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o openstack-bash-exporter ./cmd/openstack-bash-exporter/

FROM python:alpine3.17
WORKDIR /root/

RUN apk update \
    && apk add --no-cache \
        curl \
        jq \
        vim \
        bash \
        bash-doc \
        bash-completion \
    && apk add --no-cache --virtual .build-deps \
        gcc \
        git \
        libffi-dev \
        linux-headers \
        musl-dev \
        openssl-dev \
    && pip install --upgrade \
        gnocchiclient \
        pip \
        python-barbicanclient \
        python-ceilometerclient \
        python-cinderclient \
        python-cloudkittyclient \
        python-designateclient \
        python-fuelclient \
        python-glanceclient \
        python-heatclient \
        python-magnumclient \
        python-manilaclient \
        python-mistralclient \
        python-monascaclient \
        python-muranoclient \
        python-neutronclient \
        python-novaclient \
        python-openstackclient \
        python-saharaclient \
        python-senlinclient \
        python-swiftclient \
        python-troveclient \
    && apk del .build-deps \
    && rm -fr /var/cache/apk/*

COPY --from=0 /go/src/openstack-bash-exporter .
COPY ./scripts/* /root/scripts/
CMD ["./openstack-bash-exporter"]
