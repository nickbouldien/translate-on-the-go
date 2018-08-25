package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"encoding/json"
)

type App struct {
	Client *translate.Client
	//Ctx *context.Context
	Router *mux.Router
}

func (a *App) Init() {
	ctx := context.Background()

	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	a.Client = client
	defer client.Close()

	a.Router = mux.NewRouter()
	a.initRoutes()
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/list-languages", a.listLangs).Methods("GET")
}

func (a *App) listLangs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	lang := language.English

	langs, err := a.Client.SupportedLanguages(ctx, lang)
	if err != nil {
		log.Fatalf("Failed to get supported languages: %v", err)

	}

	respondWithJSON(w, http.StatusOK, langs)
	return
}


func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
