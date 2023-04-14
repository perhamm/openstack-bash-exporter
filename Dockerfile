FROM golang:1.13
ADD cmd/bash-exporter ./src/cmd/bash-exporter
WORKDIR /go/src/cmd/bash-exporter
RUN go get -d 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bash-exporter .

FROM alpine:3.7
WORKDIR /root/
COPY --from=0 /go/src/cmd/bash-exporter/bash-exporter .
COPY ./examples/* /scripts/
CMD ["./bash-exporter"]
