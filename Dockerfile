FROM debian:bullseye-slim

LABEL maintainer="jp <pntxall100@gmail.com@gmail.com>"

RUN apt update \
  && apt install -y \
  ca-certificates \
  libtesseract-dev=4.1.1-2.1 \
  tesseract-ocr=4.1.1-2.1 \
  golang

ENV GO111MODULE=on
ENV GOPATH=${HOME}/go
ENV PATH=${PATH}:${GOPATH}/bin


WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /ocrserver

CMD ["/ocrserver"]


