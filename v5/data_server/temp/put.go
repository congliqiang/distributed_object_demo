package temp

import (
	"github.com/congliqiang/distributed_object_demo/v5/data_server/locate"
	"github.com/congliqiang/distributed_object_demo/v5/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempinfo.Size {
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempinfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(datFile, tempinfo)
}

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempInfo *tempInfo) {
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempInfo.Name+"."+d)
	locate.Add(tempInfo.hash(), tempInfo.id())
}
