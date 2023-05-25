package servers

import (
	"encoding/json"
	"goCache/caches"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HTTPServer struct{
	cache *caches.Cache
}

func NewHTTPServer(cache *caches.Cache)*HTTPServer{
	return &HTTPServer{
		cache: cache,
	}
}

func (hs* HTTPServer) Run(address string) error {
	return http.ListenAndServe(address, hs.routerHandler())
}

func (hs* HTTPServer) routerHandler() http.Handler {
	router := httprouter.New()
	router.GET("/cache/:key", hs.getHandler)
	router.PUT("/cache/:key", hs.setHandler)
	router.DELETE("/cache/:key", hs.deleteHandler)
	router.GET("/status", hs.statusHandler)
	return router

}

func (hs* HTTPServer) getHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params){
	key := params.ByName("key")
	value, ok := hs.cache.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write(value)
}

func (hs* HTTPServer) setHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params){
	key := params.ByName("key")
	value , err := ioutil.ReadAll(request.Body)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hs.cache.Set(key, value)
}

func (hs* HTTPServer) deleteHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params){
	key := params.ByName("key")
	hs.cache.Delete(key)
}

func (hs* HTTPServer) statusHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params){
	status, err := json.Marshal(map[string]interface{}{
		"count" : hs.cache.Count(),
	})
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(status)
}