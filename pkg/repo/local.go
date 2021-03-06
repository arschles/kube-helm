package repo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kubernetes/helm/pkg/chart"
)

var localRepoPath string

// StartLocalRepo starts a web server and serves files from the given path
func StartLocalRepo(path string) {
	fmt.Println("Now serving you on localhost:8879...")
	localRepoPath = path
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/charts/", indexHandler)
	http.ListenAndServe(":8879", nil)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Kubernetes Package manager!\nBrowse charts on localhost:8879/charts!")
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[len("/charts/"):]
	if len(strings.Split(file, ".")) > 1 {
		serveFile(w, r, file)
	} else if file == "" {
		fmt.Fprintf(w, "list of charts should be here at some point")
	} else if file == "index" {
		fmt.Fprintf(w, "index file data should be here at some point")
	} else {
		fmt.Fprintf(w, "Ummm... Nothing to see here folks")
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, file string) {
	http.ServeFile(w, r, filepath.Join(localRepoPath, file))
}

// AddChartToLocalRepo saves a chart in the given path and then reindexes the index file
func AddChartToLocalRepo(ch *chart.Chart, path string) error {
	_, err := chart.Save(ch, path)
	if err != nil {
		return err
	}
	return Reindex(ch, path+"/index.yaml")
}

// Reindex adds an entry to the index file at the given path
func Reindex(ch *chart.Chart, path string) error {
	name := ch.Chartfile().Name + "-" + ch.Chartfile().Version
	y, err := LoadIndexFile(path)
	if err != nil {
		return err
	}
	found := false
	for k := range y.Entries {
		if k == name {
			found = true
			break
		}
	}
	if !found {
		url := "localhost:8879/charts/" + name + ".tgz"

		out, err := y.addEntry(name, url)
		if err != nil {
			return err
		}

		ioutil.WriteFile(path, out, 0644)
	}
	return nil
}
