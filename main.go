package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"fybrik.io/fybrik/pkg/helm"
	"github.com/montanaflynn/stats"
)

const (
	RELEASE_NAME = "my-notebook-fybrik-notebook-sample-arrow-flight-aef23"
	CHART_NAME   = "ghcr.io/fybrik/arrow-flight-module-chart:0.7.0"
	NAMESPACE    = "default"
)

type Perf struct {
	getConfig    []int64
	pull         []int64
	load         []int64
	install      []int64
	upgrade      []int64
	uninstall    []int64
	status1      []int64
	status2      []int64
	getResources []int64
	isInstalled1 []int64
	isInstalled2 []int64
}

func newPerf(iterations int) *Perf {
	return &Perf{getConfig: make([]int64, iterations),
		pull:         make([]int64, iterations),
		load:         make([]int64, iterations),
		install:      make([]int64, iterations),
		upgrade:      make([]int64, iterations),
		uninstall:    make([]int64, iterations),
		status1:      make([]int64, iterations),
		status2:      make([]int64, iterations),
		getResources: make([]int64, iterations),
		isInstalled1: make([]int64, iterations),
		isInstalled2: make([]int64, iterations)}
}

func toMilliseconds(t time.Time) int64 {
	return int64(time.Since(t).Milliseconds())
}

func run(perf *Perf, iter int, namespace, chartName string) {
	helmer := helm.NewHelmerImpl(true)
	cfgG, err := helmer.GetConfig(namespace, log.Printf)
	if err != nil {
		log.Fatalln(err)
	}
	helmer.Uninstall(cfgG, RELEASE_NAME)

	for i := 0; i < iter; i++ {
		fmt.Printf("iteration %d\n", i)
		start := time.Now()
		cfg, err := helmer.GetConfig(namespace, log.Printf)
		if err != nil {
			log.Fatalln(err)
		}
		perf.getConfig[i] = toMilliseconds(start)
		tmpDir, err := ioutil.TempDir("", "fybrik-helm-")
		log.Printf("before Pull %s\n", chartName)

		start = time.Now()
		err = helmer.Pull(cfg, chartName, tmpDir)
		if err != nil {
			log.Fatalln(err)
		}
		perf.pull[i] = toMilliseconds(start)

		log.Println("before Load %s", chartName)
		start = time.Now()
		chart, err := helmer.Load(chartName, tmpDir)
		if err != nil {
			log.Fatal(err)
		}
		perf.load[i] = toMilliseconds(start)

		log.Println("before isInsatalled1")
		start = time.Now()
		inst, err := helmer.IsInstalled(cfg, RELEASE_NAME)
		if err != nil {
			log.Fatal(err)
		}
		perf.isInstalled1[i] = toMilliseconds(start)
		log.Printf("is Installed 1 %t\n", inst)

		log.Println("before Status1")
		start = time.Now()
		_, err = helmer.Status(cfg, RELEASE_NAME)

		perf.status1[i] = toMilliseconds(start)

		log.Println("before Install")
		start = time.Now()
		rel, err := helmer.Install(context.Background(), cfg, chart, namespace, RELEASE_NAME, map[string]interface{}{})
		if err != nil {
			log.Fatal(err)
		}
		perf.install[i] = toMilliseconds(start)

		log.Printf("Release status %s\n", rel.Info.Status)
		log.Println("before Status")
		start = time.Now()
		status, err := helmer.Status(cfg, RELEASE_NAME)
		if err != nil {
			log.Fatal(err)
		}
		perf.status2[i] = toMilliseconds(start)

		log.Println("before GetResources")
		start = time.Now()
		_, err = helmer.GetResources(cfg, status.Manifest)
		if err != nil {
			log.Fatal(err)
		}
		perf.getResources[i] = toMilliseconds(start)

		start = time.Now()
		inst, err = helmer.IsInstalled(cfg, RELEASE_NAME)
		if err != nil {
			log.Fatal(err)
		}
		perf.isInstalled2[i] = toMilliseconds(start)
		log.Printf("is Installed 2 %t\n", inst)

		log.Println("before Upgrade")
		start = time.Now()
		_, err = helmer.Upgrade(context.Background(), cfg, chart, namespace, RELEASE_NAME, map[string]interface{}{})
		if err != nil {
			log.Fatal(err)
		}
		perf.upgrade[i] = toMilliseconds(start)

		log.Println("before uninstall")
		start = time.Now()
		_, err = helmer.Uninstall(cfg, RELEASE_NAME)
		if err != nil {
			log.Fatal(err)
		}
		perf.uninstall[i] = toMilliseconds(start)
	}

}

func printStatistics(perf *Perf) {
	log.Printf("perf = %#v\n", *perf)
	fmt.Println("Helmer performance measurements, all times in ms ")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)
	fmt.Fprintf(w, "\nOperation\tMin\tMax\tMean\tMedian (P50)\tP90\n")
	printFunction(w, perf.getConfig, "GetConfig")
	printFunction(w, perf.pull, "Pull")
	printFunction(w, perf.load, "Load")
	printFunction(w, perf.install, "Install")
	printFunction(w, perf.isInstalled1, "IsInstalled (false)")
	printFunction(w, perf.isInstalled2, "IsInstalled (true)")
	printFunction(w, perf.status1, "Status (false)")
	printFunction(w, perf.status2, "Status (true)")
	printFunction(w, perf.getResources, "GetResources")
	printFunction(w, perf.upgrade, "Upgrade")
	printFunction(w, perf.uninstall, "Uninstall")
	w.Flush()
}

func printFunction(w io.Writer, data []int64, name string) {
	d := stats.LoadRawData(data)
	min, _ := d.Min()
	max, _ := d.Max()
	mean, _ := d.Mean()
	median, _ := d.Median()
	p90, _ := d.Percentile(90)
	fmt.Fprintf(w, "%s\t%v\t%v\t%v\t%v\t%v\n", name, min, max, mean, median, p90)
}

func main() {

	chart := flag.String("chart", CHART_NAME, "name of the chart")
	iterations := flag.Int("n", 100, "# of iterations")
	namespace := flag.String("ns", NAMESPACE, "namespace")
	debug := flag.Bool("d", false, "print debug messages")
	flag.Parse()
	if !*debug {
		log.SetOutput(ioutil.Discard)
	}

	fmt.Printf("Start helmer performance test, # of iterations = %d, namespace %s, chart %s\n", *iterations, *namespace, *chart)
	perf := newPerf(*iterations)
	run(perf, *iterations, *namespace, *chart)
	printStatistics(perf)

}
