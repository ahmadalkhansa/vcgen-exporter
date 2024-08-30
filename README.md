# Vcgen-Exporter

Vcgen-exporter is a prometheus exporter that implements [vcgencmd](https://www.raspberrypi.com/documentation/computers/os.html#vcgencmd).

The software has been tested on Raspberry Pi 5 Model B Rev 1.0

| Model                          | Notes                                    |
|--------------------------------|------------------------------------------|
| Raspberry Pi 5 Model B Rev 1.0 |                                          |
| Raspberry Pi 4 Model B Rev 1.2 | does not support `pmic_read_adc command` |

![grafana dashboard](https://github.com/ahmadalkhansa/vcgen-exporter/blob/main/docs/images/RaspberryPi-Dashboard-Grafana.png?raw=true)

## Docker

The exporter listens on port 8080 and export metrics at path /metrics. A docker command to launch the container:

`docker run -p 8080:8080 --privileged docker.io/2281995/vcgen-exporter`

## Supported Commands

The exporter outputs metrics of the following vcgencmd commands:

```
vcgencmd get_throttled
vcgencmd measure_temp
vcgencmd measure_temp pmic
vcgencmd measure_volts [sdram_c,sdram_i,sdram_p]
vcgencmd pmic_read_adc
vcgencmd clock [arm,gpu,uart,emmc]
```

## Prometheus Metrics

The exported metrics can be seen below and they can be visualized using [grafana dashboard](https://github.com/ahmadalkhansa/vcgen-exporter/blob/main/grafana.json)

```
rpi_temp{component="soc",host="<hostname>"}49.9
rpi_temp{component="pmic",host="<hostname>"}51.2
rpi_volt{component="video_core",host="<hostname>"}0.9135
rpi_volt{component="sdram_core",host="<hostname>"}0.6000
rpi_volt{component="sdram_io",host="<hostname>"}0.6000
rpi_volt{component="sdram_phy",host="<hostname>"}1.1000
rpi_pmic_adc{component="3v7_wl_sw",type="a",host="<hostname>"}0.09564114
rpi_pmic_adc{component="3v3_sys",type="a",host="<hostname>"}0.07709848
rpi_pmic_adc{component="1v8_sys",type="a",host="<hostname>"}0.18933040
rpi_pmic_adc{component="ddr_vdd2",type="a",host="<hostname>"}0.00390372
rpi_pmic_adc{component="ddr_vddq",type="a",host="<hostname>"}0.00000000
rpi_pmic_adc{component="1v1_sys",type="a",host="<hostname>"}0.27716410
rpi_pmic_adc{component="0v8_sw",type="a",host="<hostname>"}0.33571990
rpi_pmic_adc{component="vdd_core",type="a",host="<hostname>"}1.75028000
rpi_pmic_adc{component="3v3_dac",type="a",host="<hostname>"}0.00030525
rpi_pmic_adc{component="3v3_adc",type="a",host="<hostname>"}0.00054945
rpi_pmic_adc{component="0v8_aon",type="a",host="<hostname>"}0.00549450
rpi_pmic_adc{component="hdmi",type="a",host="<hostname>"}0.01159950
rpi_pmic_adc{component="3v7_wl_sw",type="v",host="<hostname>"}3.70070400
rpi_pmic_adc{component="3v3_sys",type="v",host="<hostname>"}3.31374500
rpi_pmic_adc{component="1v8_sys",type="v",host="<hostname>"}1.80170800
rpi_pmic_adc{component="ddr_vdd2",type="v",host="<hostname>"}1.10439400
rpi_pmic_adc{component="ddr_vddq",type="v",host="<hostname>"}0.61355250
rpi_pmic_adc{component="1v1_sys",type="v",host="<hostname>"}1.10586000
rpi_pmic_adc{component="0v8_sw",type="v",host="<hostname>"}0.80219700
rpi_pmic_adc{component="vdd_core",type="v",host="<hostname>"}0.91518840
rpi_pmic_adc{component="3v3_dac",type="v",host="<hostname>"}3.30860500
rpi_pmic_adc{component="3v3_adc",type="v",host="<hostname>"}3.30494200
rpi_pmic_adc{component="0v8_aon",type="v",host="<hostname>"}0.79765500
rpi_pmic_adc{component="hdmi",type="v",host="<hostname>"}5.10272000
rpi_pmic_adc{component="ext5v",type="v",host="<hostname>"}5.10138000
rpi_pmic_adc{component="batt",type="v",host="<hostname>"}0.00683760
rpi_clock{component="arm",host="<hostname>"}2400020480
rpi_clock{component="gpu",host="<hostname>"}0
rpi_clock{component="uart",host="<hostname>"}43996948
rpi_clock{component="emmc",host="<hostname>"}200005008
rpi_throttled{state="under-voltage detected",host="<hostname>"}0
rpi_throttled{state="soft temperature limit occured",host="<hostname>"}1
rpi_throttled{state="arm frequency capped",host="<hostname>"}0
rpi_throttled{state="throttling occured",host="<hostname>"}1
rpi_throttled{state="throttling",host="<hostname>"}0
rpi_throttled{state="arm frequency capped occured",host="<hostname>"}1
rpi_throttled{state="soft temperature limit active",host="<hostname>"}0
rpi_throttled{state="under-voltage occured",host="<hostname>"}0
```

