package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

func (d *daemon) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", d.httpAdd)
	mux.HandleFunc("/list", d.httpList)
	mux.HandleFunc("/filter/", d.httpFilter)
	mux.ServeHTTP(w, r)
}

func renderErr(err interface{}) string {
	return fmt.Sprintf(`{"error":"%s"}`, err)
}

func (d *daemon) httpAdd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var f struct {
		Name      string
		Size      uint64
		Hashcount int
	}
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, renderErr(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}
	if err := d.add(f.Name, f.Size, f.Hashcount); err != nil {
		http.Error(w, renderErr(err), http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(f); err != nil {
		http.Error(w, renderErr(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}
}

func (d *daemon) httpList(w http.ResponseWriter, r *http.Request) {
	filters := d.list()
	filtersSaturation := make(map[string]float64)
	for _, name := range filters {
		f := d.filter(name)
		filtersSaturation[name] = f.Saturation()
	}
	if err := json.NewEncoder(w).Encode(filtersSaturation); err != nil {
		http.Error(w, renderErr(err), http.StatusInternalServerError)
		return
	}
}

func (d *daemon) httpFilter(w http.ResponseWriter, r *http.Request) {
	_, name := path.Split(r.URL.Path)
	f := d.filter(name)
	if f == nil {
		http.Error(w, renderErr(http.StatusText(http.StatusNotFound)), http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, renderErr(err), http.StatusInternalServerError)
			return
		}
		f.Add(string(body))
		w.WriteHeader(http.StatusNoContent)
	case "GET":
		resp := struct {
			Name string
			OK   bool
		}{name, f.Has(r.URL.Query().Get("body"))}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, renderErr(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
			return
		}

	case "DELETE":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, renderErr(err), http.StatusInternalServerError)
			return
		}
		f.Del(string(body))
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Add("Allow", "GET")
		w.Header().Add("Allow", "POST")
		w.Header().Add("Allow", "DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
