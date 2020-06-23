package objects

import (
	"distributed_object_demo/v2/interface_server/locate"
	"distributed_object_demo/v2/interface_server/objectstream"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = io.Copy(w, stream)
}

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}
