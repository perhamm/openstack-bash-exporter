package main

import (
	// "encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/perhamm/openstack-bash-exporter/pkg/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/handlers"
)

var (
	verbMetricsVolMax  *prometheus.GaugeVec
	verbMetricsVolUsed *prometheus.GaugeVec
	verbMetricsMemMax  *prometheus.GaugeVec
	verbMetricsMemUsed *prometheus.GaugeVec
	verbMetricsCpuMax  *prometheus.GaugeVec
	verbMetricsCpuUsed *prometheus.GaugeVec
)

type customCheck struct {
	script string
}

func main() {
	addr := flag.String("web.listen-address", ":9300", "Address on which to expose metrics")
	interval := flag.Int("interval", 300, "Interval for metrics collection in seconds")
	path := flag.String("path", "./scripts", "path to directory with bash scripts")
	labels := flag.String("labels", "project_name,tenant_id", "additioanal labels")
	// prefix := flag.String("prefix", "openstack_limits", "Prefix for metrics")
	debug := flag.Bool("debug", false, "Debug log level")
	flag.Parse()

	var labelsArr []string

	labelsArr = strings.Split(*labels, ",")
	labelsArr = append(labelsArr, "verb", "job")

	verbMetricsVolMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s", "openstack_cinder_limits_volume_max_gb"),
			Help: "openstack_cinder_limits_volume_max_gb",
		},
		// []string{"verb", "job"},
		labelsArr,
	)
	prometheus.MustRegister(verbMetricsVolMax)

	verbMetricsVolUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s", "openstack_cinder_limits_volume_used_gb"),
			Help: "openstack_cinder_limits_volume_used_gb",
		},
		// []string{"verb", "job"},
		labelsArr,
	)
	prometheus.MustRegister(verbMetricsVolUsed)

	verbMetricsMemMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s", "openstack_nova_limits_memory_maxs"),
			Help: "openstack_nova_limits_memory_maxs",
		},
		// []string{"verb", "job"},
		labelsArr,
	)
	prometheus.MustRegister(verbMetricsMemMax)

	verbMetricsMemUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s", "openstack_nova_limits_memory_used"),
			Help: "openstack_nova_limits_memory_used",
		},
		// []string{"verb", "job"},
		labelsArr,
	)
	prometheus.MustRegister(verbMetricsMemUsed)

	verbMetricsCpuMax = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s", "openstack_nova_limits_vcpus_max"),
			Help: "openstack_nova_limits_vcpus_max",
		},
		// []string{"verb", "job"},
		labelsArr,
	)
	prometheus.MustRegister(verbMetricsCpuMax)

	verbMetricsCpuUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s", "openstack_nova_limits_vcpus_used"),
			Help: "openstack_nova_limits_vcpus_used",
		},
		// []string{"verb", "job"},
		labelsArr,
	)
	prometheus.MustRegister(verbMetricsCpuUsed)

	// files, err := ioutil.ReadDir(*path)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var namesVolMax []string
	namesVolMax = append(namesVolMax, "openstack_cinder_limits_volume_max_gb.sh")
	var namesVolUsed []string
	namesVolUsed = append(namesVolUsed, "openstack_cinder_limits_volume_used_gb.sh")
	var namesMemMax []string
	namesMemMax = append(namesMemMax, "openstack_nova_limits_memory_max.sh")
	var namesMemUsed []string
	namesMemUsed = append(namesMemUsed, "openstack_nova_limits_memory_used.sh")
	var namesCpuMax []string
	namesCpuMax = append(namesCpuMax, "openstack_nova_limits_vcpus_max.sh")
	var namesCpuUsed []string
	namesCpuUsed = append(namesCpuUsed, "openstack_nova_limits_vcpus_used.sh")
	// for _, f := range files {
	// 	if f.Name()[0:1] != "." {
	// 		names = append(names, f.Name())
	// 	}
	// }

	h := health.New()
	cc := &customCheck{script: *path}
	// Add the checks to the health instance
	h.AddChecks([]*health.Config{
		{
			Name:     "scripts-check",
			Checker:  cc,
			Interval: time.Duration(2) * time.Second,
			Fatal:    true,
		},
	})
	if err := h.Start(); err != nil {
		log.Fatalf("Unable to start healthcheck: %v", err)
	}

	http.HandleFunc("/health", handlers.NewJSONHandlerFunc(h, nil))
	http.Handle("/metrics", promhttp.Handler())
	go RunVolMax(int(*interval), *path, namesVolMax, labelsArr, *debug)
	go RunVolUsed(int(*interval), *path, namesVolUsed, labelsArr, *debug)
	go RunMemMax(int(*interval), *path, namesMemMax, labelsArr, *debug)
	go RunMemUsed(int(*interval), *path, namesMemUsed, labelsArr, *debug)
	go RunCpuMax(int(*interval), *path, namesCpuMax, labelsArr, *debug)
	go RunCpuUsed(int(*interval), *path, namesCpuUsed, labelsArr, *debug)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func RunVolMax(interval int, path string, names []string, labelsArr []string, debug bool) {
	for {
		var wg sync.WaitGroup
		oArr := []*run.Output{}
		wg.Add(len(names))
		for _, name := range names {
			o := run.Output{}
			o.Job = strings.Split(name, ".")[0]
			oArr = append(oArr, &o)
			thisPath := path + "/" + name
			p := run.Params{UseWg: true, Wg: &wg, Path: &thisPath}
			go o.RunJob(&p)
		}
		wg.Wait()
		// if debug == true {
		// 	ser, err := json.Marshal(o)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Println(string(ser))
		// }
		verbMetricsVolMax.Reset()
		for _, o := range oArr {

			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job
				fmt.Println(o.Schema.Labels)
				verbMetricsVolMax.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func RunVolUsed(interval int, path string, names []string, labelsArr []string, debug bool) {
	for {
		var wg sync.WaitGroup
		oArr := []*run.Output{}
		wg.Add(len(names))
		for _, name := range names {
			o := run.Output{}
			o.Job = strings.Split(name, ".")[0]
			oArr = append(oArr, &o)
			thisPath := path + "/" + name
			p := run.Params{UseWg: true, Wg: &wg, Path: &thisPath}
			go o.RunJob(&p)
		}
		wg.Wait()
		// if debug == true {
		// 	ser, err := json.Marshal(o)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Println(string(ser))
		// }
		verbMetricsVolUsed.Reset()
		for _, o := range oArr {

			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job
				fmt.Println(o.Schema.Labels)
				verbMetricsVolUsed.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func RunMemMax(interval int, path string, names []string, labelsArr []string, debug bool) {
	for {
		var wg sync.WaitGroup
		oArr := []*run.Output{}
		wg.Add(len(names))
		for _, name := range names {
			o := run.Output{}
			o.Job = strings.Split(name, ".")[0]
			oArr = append(oArr, &o)
			thisPath := path + "/" + name
			p := run.Params{UseWg: true, Wg: &wg, Path: &thisPath}
			go o.RunJob(&p)
		}
		wg.Wait()
		// if debug == true {
		// 	ser, err := json.Marshal(o)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Println(string(ser))
		// }
		verbMetricsMemMax.Reset()
		for _, o := range oArr {

			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job
				fmt.Println(o.Schema.Labels)
				verbMetricsMemMax.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func RunMemUsed(interval int, path string, names []string, labelsArr []string, debug bool) {
	for {
		var wg sync.WaitGroup
		oArr := []*run.Output{}
		wg.Add(len(names))
		for _, name := range names {
			o := run.Output{}
			o.Job = strings.Split(name, ".")[0]
			oArr = append(oArr, &o)
			thisPath := path + "/" + name
			p := run.Params{UseWg: true, Wg: &wg, Path: &thisPath}
			go o.RunJob(&p)
		}
		wg.Wait()
		// if debug == true {
		// 	ser, err := json.Marshal(o)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Println(string(ser))
		// }
		verbMetricsMemUsed.Reset()
		for _, o := range oArr {

			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job
				fmt.Println(o.Schema.Labels)
				verbMetricsMemUsed.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func RunCpuMax(interval int, path string, names []string, labelsArr []string, debug bool) {
	for {
		var wg sync.WaitGroup
		oArr := []*run.Output{}
		wg.Add(len(names))
		for _, name := range names {
			o := run.Output{}
			o.Job = strings.Split(name, ".")[0]
			oArr = append(oArr, &o)
			thisPath := path + "/" + name
			p := run.Params{UseWg: true, Wg: &wg, Path: &thisPath}
			go o.RunJob(&p)
		}
		wg.Wait()
		// if debug == true {
		// 	ser, err := json.Marshal(o)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Println(string(ser))
		// }
		verbMetricsCpuMax.Reset()
		for _, o := range oArr {

			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job
				fmt.Println(o.Schema.Labels)
				verbMetricsCpuMax.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func RunCpuUsed(interval int, path string, names []string, labelsArr []string, debug bool) {
	for {
		var wg sync.WaitGroup
		oArr := []*run.Output{}
		wg.Add(len(names))
		for _, name := range names {
			o := run.Output{}
			o.Job = strings.Split(name, ".")[0]
			oArr = append(oArr, &o)
			thisPath := path + "/" + name
			p := run.Params{UseWg: true, Wg: &wg, Path: &thisPath}
			go o.RunJob(&p)
		}
		wg.Wait()
		// if debug == true {
		// 	ser, err := json.Marshal(o)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Println(string(ser))
		// }
		verbMetricsCpuUsed.Reset()
		for _, o := range oArr {

			for metric, value := range o.Schema.Results {
				for _, label := range labelsArr {
					if _, ok := o.Schema.Labels[label]; !ok {
						o.Schema.Labels[label] = ""
					}
				}
				o.Schema.Labels["verb"] = metric
				o.Schema.Labels["job"] = o.Job
				fmt.Println(o.Schema.Labels)
				verbMetricsCpuUsed.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func (c *customCheck) Status() (interface{}, error) {
	files, err := ioutil.ReadDir(c.script)
	if err != nil {
		return nil, err
	} else if len(files) <= 0 {
		return nil, fmt.Errorf("Script dir is empty")
	} else {
		return nil, nil
	}
}
