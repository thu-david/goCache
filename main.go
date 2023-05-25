package main

import (
	"goCache/caches"
	"goCache/servers"
)

func main(){
	cache := caches.NewCache()
	err := servers.NewHTTPServer(cache).Run(":5837")
	if err!=nil {
		panic(err)
	}
}