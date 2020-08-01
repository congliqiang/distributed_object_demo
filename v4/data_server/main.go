package main

import (
	"distributed_object_demo/v4/data_server/heartbeat"
	"distributed_object_demo/v4/data_server/locate"
	"distributed_object_demo/v4/data_server/objects"
	"distributed_object_demo/v4/data_server/temp"
	"log"
	"net/http"
	"os"
)

func main() {
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
