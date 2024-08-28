package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
)

var hostname, _ = os.Hostname()

type command interface {
	measure() []string
	metric() (m string, l []string)
}

type (
	temp struct {
		pmic bool
	}
	volts struct {
		sdramc, sdrami, sdramp bool
	}
	adc   struct{}
	clock struct {
		arm, gpu, uart, emmc bool
	}
	throttle struct {}
)

func (t temp) measure() []string {
	var lcomm []string
	var lres []string
	command := "measure_temp"
	lcomm = append(lcomm, command)
	if t.pmic {
		command += " pmic"
		lcomm = append(lcomm, command)

	}
	for _, v := range lcomm {
		temp, err := VcComm(v)
		if err != nil {
			log.Fatal("IOCTL call has returned Errno ", err)
		}
		_, temp, _ = strings.Cut(temp, "=")
		temp, _, _ = strings.Cut(temp, "'")
		lres = append(lres, temp)
	}
	return lres
}

func (v volts) measure() []string {
	var lcomm []string
	var lres []string
	command := "measure_volts"
	lcomm = append(lcomm, command)
	switch {
	case v.sdramc:
		ccommand := command + " sdram_c"
		lcomm = append(lcomm, ccommand)
		fallthrough
	case v.sdrami:
		icommand := command + " sdram_i"
		lcomm = append(lcomm, icommand)
		fallthrough
	case v.sdramp:
		pcommand := command + " sdram_p"
		lcomm = append(lcomm, pcommand)
	}
	for _, v := range lcomm {
		volts, err := VcComm(v)
		if err != nil {
			log.Fatal("IOCTL call has returned Errno ", err)
		}
		_, volts, _ = strings.Cut(volts, "=")
		volts, _, _ = strings.Cut(volts, "V")
		lres = append(lres, volts)
	}
	return lres
}

func (p adc) measure() []string {
	command := "pmic_read_adc"
	var lres []string
	power, err := VcComm(command)
	if err != nil {
		log.Fatal("IOCTL call has returned Errno ", err)
	}
	pslice := strings.Split(power, "\n")
	for _, i := range pslice {
		var power string
		var found bool
		_, power, _ = strings.Cut(i, "=")
		power, _, found = strings.Cut(power, "V")
		if found != true {
			power, _, _ = strings.Cut(power, "A")
		}
		lres = append(lres, power)
	}
	return lres
}

func (c clock) measure() []string {
	var lres []string
	var lcomm []string
	command := "measure_clock"
	switch {
	case c.arm:
		ccommand := command + " arm"
		lcomm = append(lcomm, ccommand)
		fallthrough
	case c.gpu:
		icommand := command + " gpu"
		lcomm = append(lcomm, icommand)
		fallthrough
	case c.uart:
		pcommand := command + " uart"
		lcomm = append(lcomm, pcommand)
		fallthrough
	case c.emmc:
		pcommand := command + " emmc"
		lcomm = append(lcomm, pcommand)
	}
	for _, v := range lcomm {
		var hrtz string
		clo, err := VcComm(v)
		if err != nil {
			log.Fatal("IOCTL call has returned Errno ", err)
		}
		hrtz = strings.Split(clo, "=")[1]
		//clean null ascii
		hrtz = strings.ReplaceAll(hrtz, "\x00", "")
		lres = append(lres, hrtz)
	}
	return lres
}

func (o throttle) measure() []string {
	var lres []string
	command := "get_throttled"
	bs := make([]string, 20)
	for i := range bs {
		bs[i] = "0"
	}
	clo, err := VcComm(command)
	if err != nil {
		log.Fatal("IOCTL call has returned Errno ", err)
	}
	//clean null ascii
	clo = strings.ReplaceAll(clo, "\x00", "")
	clo = strings.Split(clo, "=")[1]
	v, errs := strconv.ParseUint(clo, 0, 32)
	if errs != nil {
		log.Fatal(errs)
	}
	s := strconv.FormatUint(v, 2)
	for i, j := range s {
		libs := len(bs) - 1 - i
		bs[libs] = string(j)
	}
	for i, j := 0, len(bs) - 1; i < 4; i++ {
		lres = append(lres, bs[i], bs[j])
		j--
	}
	return lres
}

func (t temp) metric() (m string, l[]string) {
	var lmetric []string
	lmetric = append(lmetric, param("component", "soc"))
	if t.pmic {
		lmetric = append(lmetric, param("component", "pmic"))
	}
	return "rpi_temp", lmetric
}

func (v volts) metric() (m string, l []string) {
	var lmetric []string
	lmetric = append(lmetric, param("component", "video_core"))
	switch {
	case v.sdramc:
		lmetric = append(lmetric, param("component", "sdram_core"))
		fallthrough
	case v.sdrami:
		lmetric = append(lmetric, param("component", "sdram_io"))
		fallthrough
	case v.sdramp:
		lmetric = append(lmetric, param("component", "sdram_phy"))
	}
	return "rpi_volt", lmetric
}

func (p adc) metric() (m string, l []string) {
	var lmetric []string
	command := "pmic_read_adc"
	power, err := VcComm(command)
	if err != nil {
		log.Fatal("IOCTL call has returned an Errno", err)
	}
	pslice := strings.Split(power, "\n")
	for _, i := range pslice {
		var metric string
		metric = strings.TrimSpace(i)
		metric = strings.Split(metric, " ")[0]
		metric = strings.ToLower(metric)
		component := param("component", metric[:len(metric)-2])
		t := param("type", metric[len(metric)-1:])
		metric = component+","+t
		lmetric = append(lmetric, metric)
	}
	return "rpi_pmic_adc", lmetric
}

func (c clock) metric() (m string, l []string) {
	var lmetric []string
	switch {
	case c.arm:
		lmetric = append(lmetric, param("component", "arm"))
		fallthrough
	case c.gpu:
		lmetric = append(lmetric, param("component", "gpu"))
		fallthrough
	case c.uart:
		lmetric = append(lmetric, param("component", "uart"))
		fallthrough
	case c.emmc:
		lmetric = append(lmetric, param("component", "emmc"))
	}
	return "rpi_clock", lmetric
}

func (t throttle) metric() (m string, l []string) {
	var lmetric = []string{
		"under-voltage detected",
		"soft temperature limit occured",
		"arm frequency capped",
		"throttling occured",
		"throttling",
		"arm frequency capped occured",
		"soft temperature limit active",
		"under-voltage occured",
	}
	for i, j := range lmetric {
		l := param("state", j)
		lmetric[i] = l
	}
	return "rpi_throttled", lmetric
}

func param(l, p string) string {
	return l + "=\"" + p + "\""
}

func PromOut(c command) string {
	var format string
	hlabel := fmt.Sprintf("host=\"%s\"", hostname)
	lres := c.measure()
	metric, lmetric := c.metric()
	for i := range lmetric {
		format += metric + "{" + lmetric[i] + "," + hlabel + "}" + lres[i] + "\n"
	}
	return format
}
