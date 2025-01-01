package main

import (
	"io"
	"log"
	"net/http"
	"flag"
	"strconv"
	"strings"
)

type constructor struct {
	metric strings.Builder
	err error
}

var (
	cl = temp{enabled: true, pmic: true}
	vl = volts{enabled: true, sdramc: true, sdrami: true, sdramp: true}
	ad = adc{enabled: true}
	ck = clock{enabled: true, arm: true, gpu: true, uart: true, emmc: true}
	th = throttle{enabled: true}
)

var port int

func (cs *constructor) write(s string, e error) {
	if e != nil {
		log.Println(e)
	}
	_, cs.err = cs.metric.WriteString(s)
	if cs.err != nil {
		log.Println(cs.err)
	}
}

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
		var resp constructor
		log.Printf("%s %s request to %s", r.RemoteAddr, r.Method, r.URL.RequestURI())
		resp.write(PromOut(cl))
		resp.write(PromOut(vl))
		resp.write(PromOut(ad))
		resp.write(PromOut(ck))
		resp.write(PromOut(th))
		io.WriteString(w, resp.metric.String())
	}

	http.HandleFunc("/metrics", sm)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
