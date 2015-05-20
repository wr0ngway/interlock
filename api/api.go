package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/version"
	"github.com/gorilla/mux"
)

type Api struct {
	config *ApiConfig
}

func NewApi(cfg *ApiConfig) *Api {
	return &Api{
		config: cfg,
	}
}

func (a *Api) Run() error {
	globalMux := http.NewServeMux()

	router := mux.NewRouter()
	router.HandleFunc("/", a.apiIndex).Methods("GET")
	router.HandleFunc("/plugins", a.apiPlugins)
	globalMux.Handle("/", router)

	return http.ListenAndServe(a.config.ListenAddr, globalMux)
}

func (a *Api) apiIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("interlock %s\n", version.FULL_VERSION)))
}

func (a *Api) apiPlugins(w http.ResponseWriter, r *http.Request) {
	plugins := a.config.Manager.Plugins()

	info := []*interlock.PluginInfo{}
	for _, p := range plugins {
		info = append(info, p.Info())
	}

	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
