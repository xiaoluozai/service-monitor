package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/51idc/service-monitor/agent/g"
	"github.com/toolkits/file"
)

func configPageRoutes() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			if !file.IsExist(filepath.Join(g.Root, "/public", r.URL.Path, "index.html")) {
				http.NotFound(w, r)
				return
			}
		}
		http.FileServer(http.Dir(filepath.Join(g.Root, "/public"))).ServeHTTP(w, r)
	})

}
