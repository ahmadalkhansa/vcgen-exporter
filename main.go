package main

import (
	"io"
	"log"
	"net/http"
	"flag"
	"strconv"
)

var (
	cl = temp{pmic: true}
	vl = volts{sdramc: true, sdrami: true, sdramp: true}
	ad = adc{}
	ck = clock{arm: true, gpu: true, uart: true, emmc: true}
	th = throttle{}
)

var port int

func init() {
	const (
		portDefault = 8080
		portUsage = "Exporter's Listening port"
	)
	flag.IntVar(&port, "p", portDefault, portUsage)
}

func main() {
	flag.Parse()
	sm := func(w http.ResponseWriter, _ *http.Request) {
		var col string
		col += PromOut(cl)
		col += PromOut(vl)
		col += PromOut(ad)
		col += PromOut(ck)
		col += PromOut(th)
		io.WriteString(w, col)
	}

	http.HandleFunc("/metrics", sm)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
