package locate

import (
	"encoding/json"
	"github.com/congliqiang/distributed_object_demo/v6/data_server/types"
	"github.com/congliqiang/distributed_object_demo/v6/rabbitmq"
	"github.com/congliqiang/distributed_object_demo/v6/rs"
	"net/http"
	"os"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	_, _ = w.Write(b)
}

func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(os.Getenv("RABBIT_MQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		_ = json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}
