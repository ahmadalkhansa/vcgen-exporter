FROM docker.io/golang:1.23

WORKDIR /usr/src/app

RUN --mount=type=bind,source=./,target=/usr/src/app go build -v -o /usr/local/bin/vcgen-exporter ./

EXPOSE 8080

CMD ["vcgen-exporter"]
