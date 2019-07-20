package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofoody/api-gateway/pkg/config"
	"github.com/gofoody/api-gateway/pkg/ctrl"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New()

	initLogger(cfg.GetLogLevel())

	router := mountEndpoints(cfg)
	startService(cfg.GetHttpPort(), router)
}

func initLogger(logLevel string) {
	level, _ := log.ParseLevel(logLevel)
	log.SetLevel(level)
	log.SetOutput(os.Stdout)
}

func mountEndpoints(cfg *config.Config) *mux.Router {
	r := mux.NewRouter()

	statusCtrl := ctrl.NewStatusCtrl()
	r.HandleFunc("/api/status", statusCtrl.Show)

	gatewayCtrl := ctrl.NewGatewayCtrl(cfg)
	r.PathPrefix("/").HandlerFunc(gatewayCtrl.Router)

	return r
}

func startService(port int, router *mux.Router) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Infof("api gateway running at:%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("failed to start api gateway, error:%v", err)
	}
}
