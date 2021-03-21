package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"github.com/roylee0704/gron"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	APIEndPoint    = "https://api.binance.com"
	APITimeout     = 5
	APITickerPrice = "/api/v3/ticker/price"
)

func main() {

	var (
		webConfig     = webflag.AddFlags(kingpin.CommandLine)
		listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9101").String()
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		intervalQuery = kingpin.Flag("binance.api-interval", "Interval query Binance API.").Default("10").Int()
		testUpTrend   = kingpin.Flag("binance.testUpTrend", "Test uptrend trigger.").Default("false").Bool()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("binance_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting binance_exporter")

	binance := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "binance_usdt",
		Help: "binance symbol pricing in usdt",
	}, []string{"symbol"})

	c := gron.New()
	level.Info(logger).Log("cron", "Interval query Binance API: "+strconv.Itoa(*intervalQuery))

	c.AddFunc(gron.Every(time.Duration(*intervalQuery)*time.Second), func() {
		result := []struct {
			Symbol string `json:"symbol"`
			Price  string `json:"price"`
		}{}

		response, err := binanceAPIGet(APITickerPrice, logger)
		if err != nil {
			level.Error(logger).Log("cron", "Query Binance API error", err.Error())
			return
		}

		json.Unmarshal([]byte(response), &result)
		for _, item := range result {
			if strings.Contains(item.Symbol, "USDT") {
				price, _ := strconv.ParseFloat(item.Price, 64)
				binance.WithLabelValues(item.Symbol).Set(price)
			}
		}

		if *testUpTrend {
			binance.WithLabelValues("TEST-USDT").Set(float64(time.Now().UTC().Unix()))
		}
	})
	c.Start()
	defer c.Stop()

	// Registration Indicator Information
	prometheus.MustRegister(binance)

	// Expose
	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	srv := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(srv, *webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}
}
