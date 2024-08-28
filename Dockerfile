FROM docker.io/golang:1.23

LABEL vcgen-exporter.version="v0.5.0"
LABEL vcgen-exporter.image.author="ahmadkhansa95@gmail.com"

WORKDIR /usr/src/app

RUN --mount=type=bind,source=./,target=/usr/src/app go build -buildvcs=false -v -o /usr/local/bin/vcgen-exporter ./

EXPOSE 8080

CMD ["vcgen-exporter"]
