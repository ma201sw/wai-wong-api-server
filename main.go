package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"go-wai-wong/internal/config"
	"go-wai-wong/internal/golib"
	"go-wai-wong/internal/provider/jsonprovider"
	"go-wai-wong/internal/route"
	"go-wai-wong/internal/tokenhelper"
)

func defaultRouter() *chi.Mux {
	r := chi.NewRouter()

	return r
}

func main() {
	config.LoadConfig()

	r := defaultRouter()

	jsonProviderSrv := jsonprovider.New()
	tokenHelperSrv := tokenhelper.New()

	myGoLibsSrv := golib.New()

	r.Use(golib.Inject(myGoLibsSrv))
	r.Use(tokenhelper.Inject(tokenHelperSrv))
	r.Use(jsonprovider.Inject(jsonProviderSrv))
	route.Install(r)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server because: %v", err)
	}
}
