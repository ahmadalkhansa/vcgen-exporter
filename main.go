package main

import (
	"io"
	"log"
	"net/http"
	"flag"
	"strconv"
	"strings"
)

var (
	cl = temp{enabled: true, pmic: true}
	vl = volts{enabled: true, sdramc: true, sdrami: true, sdramp: true}
	ad = adc{enabled: true}
	ck = clock{enabled: true, arm: true, gpu: true, uart: true, emmc: true}
	th = throttle{enabled: true}
)

var port int

func init() {
	log.Println("vcgen-exporter initializing...")
	const (
		portDefault = 8080
		portUsage = "Exporter's Listening port"
	)
	flag.IntVar(&port, "p", portDefault, portUsage)
}

func main() {
	flag.Parse()
	sm := func(w http.ResponseWriter, r *http.Request) {
		var resp strings.Builder
		var col string
		var err error
		var erri error
		log.Printf("%s %s request to %s", r.RemoteAddr, r.Method, r.URL.RequestURI())
		col, err = PromOut(cl)
		_, erri = resp.WriteString(col)
		col, err = PromOut(vl)
		_, erri = resp.WriteString(col)
		col, err = PromOut(ad)
		_, erri = resp.WriteString(col)
		col, err = PromOut(ck)
		_, erri = resp.WriteString(col)
		col, err = PromOut(th)
		_, erri = resp.WriteString(col)
		if err != nil {
			log.Println(err)
		}
		if erri != nil {
			log.Println(erri)
		}
		io.WriteString(w, resp.String())
	}

	http.HandleFunc("/metrics", sm)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
