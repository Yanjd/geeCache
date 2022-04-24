package geeCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/geecache/"

type HttpPool struct {
	self     string
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HttpPool) ServeHttp(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%v %v", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path, "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	g := GetGroup(groupName)
	if g == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	v, err := g.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(v.ByteSlice())
}
