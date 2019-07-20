package ctrl

import (
	"net/http"
)

type StatusCtrl interface {
	Show(w http.ResponseWriter, r *http.Request)
}

type statusCtrl struct {
}

func NewStatusCtrl() StatusCtrl {
	c := &statusCtrl{}
	return c
}

func (c *statusCtrl) Show(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("api gateway status OK."))
}
