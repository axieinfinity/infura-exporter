package main

import (
	"flag"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/namcxn/infura-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
)

type EthBlockNumber struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

type queryConfig map[string]string

var (
//url = "https://mainnet.infura.io/v3/xxxxxxfdb3834ebe8bacea9xxxxxxx"
)

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {

	flag.Parse()
	url := config.Config.Infura.Url + config.Config.Infura.Token
	var rpc, err = rpc.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(
		NewEthBlockNumber(rpc),
	)

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.New(os.Stderr, log.Prefix(), log.Flags()),
		ErrorHandling: promhttp.ContinueOnError,
	})

	http.Handle("/metrics", handler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)

		io.WriteString(w, "Not found\n")
	})

	http.ListenAndServe(":9877", LogRequest(http.DefaultServeMux))
}

func NewEthBlockNumber(rpc *rpc.Client) *EthBlockNumber {
	return &EthBlockNumber{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"infura_block_number",
			"the number of most recent block",
			nil,
			nil,
		),
	}
}

func (collector *EthBlockNumber) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthBlockNumber) Collect(ch chan<- prometheus.Metric) {
	var result hexutil.Uint64
	if err := collector.rpc.Call(&result, "eth_blockNumber"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	value := float64(result)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
