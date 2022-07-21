FROM golang:1.18.4-alpine3.15 as builder

LABEL maintainer="12306killer@gmail.com"

WORKDIR /usr/local/baidupan-upload

ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ${WORKDIR}

RUN cd cmd/baidupanupload \
    && export GIN_MODE=release \
    && CGO_ENABLED=0 go build -o /tmp/baidupanupload -mod=readonly

FROM scratch as runtime

WORKDIR /opt/app

COPY --from=builder /tmp/baidupanupload /usr/local/bin/baidupanupload
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/usr/local/bin/baidupanupload"]