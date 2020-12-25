package main

import (
	"github.com/congliqiang/distributed_object_demo/v2/interface_server/heartbeat"
	"github.com/congliqiang/distributed_object_demo/v2/interface_server/locate"
	"github.com/congliqiang/distributed_object_demo/v2/interface_server/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
