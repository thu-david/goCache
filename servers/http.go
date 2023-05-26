package servers

import (
	"encoding/json"
	"goCache/caches"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

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

func wrapUriWithVersion(uri string) string {
	return path.Join("/", APIVersion, uri)
}

func (hs* HTTPServer) routerHandler() http.Handler {

	router := httprouter.New()
	router.GET(wrapUriWithVersion("/cache/:key"), hs.getHandler)
	router.PUT(wrapUriWithVersion("/cache/:key"), hs.setHandler)
	router.DELETE(wrapUriWithVersion("/cache/:key"), hs.deleteHandler)
	router.GET(wrapUriWithVersion("/status"), hs.statusHandler)
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
	ttl, err := ttlOf(request)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = hs.cache.SetWithTTL(key, value, ttl)
	if err!=nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte("ERROR:"+err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ttlOf(request *http.Request) (int64, error) {
	ttls, ok := request.Header["Ttl"]
	if !ok || len(ttls) < 1 {
		return caches.NeverDie, nil
	}
	return strconv.ParseInt(ttls[0], 10, 64)
}

func (hs* HTTPServer) deleteHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params){
	key := params.ByName("key")
	hs.cache.Delete(key)
}

func (hs* HTTPServer) statusHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params){
	status, err := json.Marshal(hs.cache.Status())
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(status)
}