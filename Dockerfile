FROM debian:bullseye-slim

LABEL maintainer="jp <pntxall100@gmail.com@gmail.com>"

RUN apt update \
  && apt install -y \
  ca-certificates \
  libtesseract-dev=4.1.1-2.1 \
  tesseract-ocr=4.1.1-2.1

RUN wget https://go.dev/dl/go1.17.linux-amd64.tar.gz
RUN tar -xzf go1.17.linux-amd64.tar.gz

ENV GOROOT=/app/go
ENV GO111MODULE=on
ENV GOPATH=${HOME}/go
ENV PATH=${PATH}:${GOROOT}/bin:${GOPATH}/bin

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /ocrserver

CMD ["/ocrserver"]
