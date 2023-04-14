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
	verbMetrics *prometheus.GaugeVec
)

type customCheck struct {
	script string
}

func main() {
	addr := flag.String("web.listen-address", ":9300", "Address on which to expose metrics")
	interval := flag.Int("interval", 300, "Interval for metrics collection in seconds")
	path := flag.String("path", "./scripts", "path to directory with bash scripts")
	labels := flag.String("labels", "projectname,tenantid", "additioanal labels")
	prefix := flag.String("prefix", "openstack_limits", "Prefix for metrics")
	debug := flag.Bool("debug", false, "Debug log level")
	flag.Parse()

	var labelsArr []string

	labelsArr = strings.Split(*labels, ",")
	labelsArr = append(labelsArr, "verb", "job")

	verbMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s", *prefix),
			Help: "bash exporter metrics",
		},
		// []string{"verb", "job"},
		labelsArr,
	)
	prometheus.MustRegister(verbMetrics)

	files, err := ioutil.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	for _, f := range files {
		if f.Name()[0:1] != "." {
			names = append(names, f.Name())
		}
	}

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
	go Run(int(*interval), *path, names, labelsArr, *debug)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func Run(interval int, path string, names []string, labelsArr []string, debug bool) {
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
		verbMetrics.Reset()
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
				verbMetrics.With(prometheus.Labels(o.Schema.Labels)).Set(float64(value))
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
