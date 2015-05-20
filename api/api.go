package api

import (
	"fmt"
	"net/http"

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
	apiRouter := mux.NewRouter()

	rootRouter := mux.NewRouter()
	rootRouter.HandleFunc("/", a.apiRoot).Methods("GET")
	globalMux.Handle("/", rootRouter)

	apiRouter.HandleFunc("/api", a.apiIndex).Methods("GET")
	globalMux.Handle("/api", apiRouter)

	return http.ListenAndServe(a.config.ListenAddr, globalMux)
}

func (a *Api) apiRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api", http.StatusFound)
}

func (a *Api) apiIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("interlock %s\n", version.FULL_VERSION)))
}
