package es

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s", "version":%d,"size":%d,"hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create", os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}
