package ctrl

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofoody/api-gateway/pkg/config"

	log "github.com/sirupsen/logrus"
)

type GatewayCtrl interface {
	Router(w http.ResponseWriter, r *http.Request)
}

type gatewayCtrl struct {
	cfg *config.Config
}

func NewGatewayCtrl(cfg *config.Config) GatewayCtrl {
	c := &gatewayCtrl{cfg: cfg}
	return c
}

func (c *gatewayCtrl) Router(w http.ResponseWriter, r *http.Request) {
	reqPath := r.URL.Path
	log.Infof("reqeust path:'%s'", reqPath)

	switch reqPath {
	case "/api/login":
		forward(w, r, c.cfg.GetAuthURL())
	default:
		w.Write([]byte(fmt.Sprintf("routing not found, path:'%s'", reqPath)))
	}
}

func forward(w http.ResponseWriter, r *http.Request, host string) {
	// create request
	fwURL := fmt.Sprintf("%s%s", host, r.RequestURI)
	log.Debugf("fwURL=%s", fwURL)
	fwReq, err := http.NewRequest(r.Method, fwURL, r.Body)
	if err != nil {
		log.Errorf("failed to create request, url:'%s'", fwURL)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set headers
	fwReq.Header = make(http.Header)
	for h, val := range r.Header {
		fwReq.Header[h] = val
	}

	// execute request
	httpClient := http.Client{}
	fwRes, err := httpClient.Do(fwReq)
	if err != nil {
		log.Errorf("failed to make request, url:'%s'", fwURL)
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// process response
	defer fwRes.Body.Close()
	body, err := ioutil.ReadAll(fwRes.Body)
	if err != nil {
		log.Errorf("failed to parse response, url:'%s'", fwURL)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// headerResp := strings.Join(fwRes.Header["Content-Type"], "")
	// w.Header().Set("Content-Type", headerResp)
	w.WriteHeader(fwRes.StatusCode)
	w.Write(body)
}
