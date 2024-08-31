FROM docker.io/golang:1.23-alpine3.20 AS base

WORKDIR /usr/src/app

RUN --mount=type=bind,source=./,target=/usr/src/app go build -buildvcs=false -v -o /usr/local/bin/vcgen-exporter ./

FROM docker.io/alpine:3.20

LABEL vcgen-exporter.version="v0.5.1"
LABEL vcgen-exporter.image.author="ahmadkhansa95@gmail.com"

COPY --from=base /usr/local/bin/vcgen-exporter /usr/local/bin/vcgen-exporter

EXPOSE 8080

CMD ["vcgen-exporter"]
