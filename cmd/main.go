package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	db "github.com/sound-of-destiny/infiCombo_exporter/mongo"
)

const (
	database = "infiCombo"
)

var (
	addr = flag.String("listen-address", ":10000", "The address to listen on for HTTP requests.")

	temp_valid_a = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "temp_valid_a",
		Help: "infiCombo.",
	})

	sm = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "sm",
		Help: "infiCombo.",
	})

	pressure = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pressure",
		Help: "infiCombo.",
	})
)

type Beiqin_efire_data struct {
	Devtype  string         `json:"devtype"`
	Datarows []Beiqin_efire `json:"datarows"`
}

type Beiqin_efire struct {
	Cur_valid_a float64 `json:"cur_valid_a"`
	Cur_valid_b float64 `json:"cur_valid_b"`
	Cur_valid_c float64 `json:"cur_valid_c"`
	Cur_valid_n float64 `json:"cur_valid_n"`
	Lcur_valid  float64 `json:"lcur_valid"`

	Temp_valid_a float64 `json:"temp_valid_a"`
	Temp_valid_b float64 `json:"temp_valid_b"`
	Temp_valid_c float64 `json:"temp_valid_c"`
	Temp_valid_n float64 `json:"temp_valid_n"`

	Vol_valid_a float64 `json:"vol_valid_a"`
	Vol_valid_b float64 `json:"vol_valid_b"`
	Vol_valid_c float64 `json:"vol_valid_c"`

	Deveui      string `json:"deveui"`
	Collecttime string `json:"collecttime"`
}

type Xinshanying_smoke_v2_data struct {
	Devtype  string                 `json:"devtype"`
	Datarows []Xinshanying_smoke_v2 `json:"datarows"`
}

type Xinshanying_smoke_v2 struct {
	Dpower float64 `json:"dpower"`
	Bpower float64 `json:"bpower"`
	Sm     float64 `json:"sm"`

	Deveui      string `json:"deveui"`
	Collecttime string `json:"collecttime"`
}

type Wanzhong_pressure_data struct {
	Devtype  string              `json:"devtype"`
	Datarows []Wanzhong_pressure `json:"datarows"`
}

type Wanzhong_pressure struct {
	Pressure float64 `json:"pressure"`
	Battery  float64 `json:"battery"`

	Deveui      string `json:"deveui"`
	Collecttime string `json:"collecttime"`
}

func init() {
	prometheus.MustRegister(temp_valid_a)
	prometheus.MustRegister(sm)
	prometheus.MustRegister(pressure)
}

func beiqin_efire_func(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()

	var beiqin_efire_data Beiqin_efire_data
	err := json.Unmarshal(b, &beiqin_efire_data)
	if err != nil {
		log.Fatal(err)
	}

	v := beiqin_efire_data.Datarows[0].Temp_valid_a
	temp_valid_a.Set(v)

	db.Insert(database, beiqin_efire_data.Devtype+"_"+beiqin_efire_data.Datarows[0].Deveui, beiqin_efire_data.Datarows[0])

	log.Println(beiqin_efire_data)
}

func xinshanying_smoke_v2_func(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()

	var xinshanying_smoke_v2_data Xinshanying_smoke_v2_data
	err := json.Unmarshal(b, &xinshanying_smoke_v2_data)
	if err != nil {
		log.Fatal(err)
	}

	v := xinshanying_smoke_v2_data.Datarows[0].Sm
	sm.Set(v)

	db.Insert(database, xinshanying_smoke_v2_data.Devtype+"_"+xinshanying_smoke_v2_data.Datarows[0].Deveui, xinshanying_smoke_v2_data.Datarows[0])

	log.Println(xinshanying_smoke_v2_data)
}

func wanzhong_pressure_func(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := buf.Bytes()

	var wanzhong_pressure_data Wanzhong_pressure_data
	err := json.Unmarshal(b, &wanzhong_pressure_data)
	if err != nil {
		log.Fatal(err)
	}

	v := wanzhong_pressure_data.Datarows[0].Pressure
	pressure.Set(v)

	db.Insert(database, wanzhong_pressure_data.Devtype+"_"+wanzhong_pressure_data.Datarows[0].Deveui, wanzhong_pressure_data.Datarows[0])

	log.Println(wanzhong_pressure_data)
}

func main() {

	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/beiqin_efire", beiqin_efire_func)

	http.HandleFunc("/xinshanying_smoke_v2", xinshanying_smoke_v2_func)

	http.HandleFunc("/wanzhong_pressure", wanzhong_pressure_func)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Node Exporter</title></head>
			<body>
			<h1>infiCombo Exporter</h1>
			<p><a href="` + *addr + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
