package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-resty/resty/v2"
)

func binanceAPIGet(api string, logger log.Logger) (string, error) {
	client := resty.New()
	client.SetTimeout(time.Duration(APITimeout) * time.Second)

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		ForceContentType("application/json").
		Get(APIEndPoint + api)

	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		level.Error(logger).Log("binance", "API Status Code: "+strconv.Itoa(resp.StatusCode()))
		level.Error(logger).Log("binance", "Status:"+resp.Status())
		level.Error(logger).Log("binance", "Proto: "+resp.Proto())
		level.Error(logger).Log("binance", "Time: "+resp.Time().String())
		level.Error(logger).Log("binance", "Received At: "+resp.ReceivedAt().String())
		level.Error(logger).Log("binance", "x-mbx-uuid: "+resp.Header().Get("x-mbx-uuid"))
		level.Error(logger).Log("binance", "x-mbx-used-weight: "+resp.Header().Get("x-mbx-used-weight"))
		level.Error(logger).Log("binance", "x-mbx-used-weight-1m: "+resp.Header().Get("x-mbx-used-weight-1m"))
		level.Error(logger).Log("binance", "Body", resp.String())
		return "", errors.New("Status code " + resp.Status())
	}
	//client.GetClient().CloseIdleConnections()

	level.Info(logger).Log("binance", "API Status Code: "+strconv.Itoa(resp.StatusCode()))
	level.Info(logger).Log("binance", "x-mbx-uuid: "+resp.Header().Get("x-mbx-uuid"))
	level.Info(logger).Log("binance", "x-mbx-used-weight: "+resp.Header().Get("x-mbx-used-weight"))
	level.Info(logger).Log("binance", "x-mbx-used-weight-1m: "+resp.Header().Get("x-mbx-used-weight-1m"))
	return resp.String(), nil
}
