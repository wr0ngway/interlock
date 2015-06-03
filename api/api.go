package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	router.HandleFunc("/plugins", a.apiPlugins).Methods("GET")
	router.HandleFunc("/signal/{action:.*}", a.apiSignal).Methods("POST")
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

func (a *Api) apiSignal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]

	// get post body for params
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error decoding parameters: %s", err), http.StatusBadRequest)
		return
	}

	var params map[string]string
	// if not an empty post, unmarshal
	if len(data) > 0 {
		if err := json.Unmarshal(data, &params); err != nil {
			http.Error(w, "error decoding parameters", http.StatusBadRequest)
			return
		}
	}

	a.config.Manager.Signal(action, params)
}
