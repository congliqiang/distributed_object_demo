package versions

import (
	"encoding/json"
	"github.com/congliqiang/distributed_object_demo/v5/es"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, e := es.SearchAllVersion(name, from, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			_, _ = w.Write(b)
			_, _ = w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
