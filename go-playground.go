package main

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"net/http"
	"log"
	mux "github.com/gorilla/mux"
	"encoding/json"
)

func printStructAsJSON(v interface{}, w http.ResponseWriter) {
	out, err := json.Marshal(v)
	if err != nil {
		panic (err)
	}
	fmt.Fprintln(w, string(out))
}

func CPUStatsHandler(w http.ResponseWriter, r *http.Request) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	for _, s := range stat.CPUStats {
		printStructAsJSON(s,w)
	}
}

func DiskStatsHandler(w http.ResponseWriter, r *http.Request) {
	d, err := linuxproc.ReadDisk("/")
	if err != nil {
		log.Fatal("stat read fail")
	}

	// Convert to MB
	d.All = d.All / (1024 ^ 2)
	d.Free = d.Free / (1024 ^ 2)
	d.Used = d.Used / (1024 ^ 2)

	printStructAsJSON(d,w)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola GO!"))
}

func routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/cpu", CPUStatsHandler)
	r.HandleFunc("/disk", DiskStatsHandler)
	r.HandleFunc("/", HomeHandler)
	return r
}

func main() {
	r := routes()
	http.ListenAndServe(":8080", r)
}
