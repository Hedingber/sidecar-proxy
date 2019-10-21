package sidecar_proxy

import (
	"errors"
	"flag"
	"os"

	sidecarproxy "github.com/v3io/sidecar-proxy/pkg/sidecar-proxy"
	"github.com/v3io/sidecar-proxy/pkg/sidecar-proxy/common"
	"github.com/v3io/sidecar-proxy/pkg/sidecar-proxy/metricshandler"
	"github.com/v3io/sidecar-proxy/pkg/sidecar-proxy/metricshandler/jupyterkernelbusyness"
	"github.com/v3io/sidecar-proxy/pkg/sidecar-proxy/metricshandler/numofrequests"

	"github.com/sirupsen/logrus"
)

func main() {

	var metricNames common.StringArrayFlag

	// args
	listenAddress := flag.String("listen-addr", os.Getenv("PROXY_LISTEN_ADDRESS"), "Port to listen on")
	forwardAddress := flag.String("forward-addr", os.Getenv("PROXY_FORWARD_ADDRESS"), "IP /w port to forward to (without protocol)")
	namespace := flag.String("namespace", os.Getenv("PROXY_NAMESPACE"), "Kubernetes namespace")
	serviceName := flag.String("service-name", os.Getenv("PROXY_SERVICE_NAME"), "Service which the proxy serves")
	instanceName := flag.String("instance-name", os.Getenv("PROXY_INSTANCE_NAME"), "Deployment instance name")
	logLevel := flag.String("log-level", os.Getenv("LOG_LEVEL"), "Set proxy's log level")
	flag.Var(&metricNames, "metric-names", "Set which metrics to collect")
	flag.Parse()

	// logger conf
	var logger = logrus.New()
	parsedLogLevel, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(parsedLogLevel)

	if len(metricNames) == 0 {
		panic(errors.New("at least one metric name should be given"))
	}

	// proxy server start
	proxyServer, err := sidecarproxy.NewProxyServer(logger, *listenAddress, *forwardAddress, *namespace, *serviceName, *instanceName, metricNames)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create a proxy server")
	}
	if err = proxyServer.Start(); err != nil {
		panic(err)
	}
}
