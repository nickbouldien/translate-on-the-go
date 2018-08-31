package main

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"

	"github.com/gorilla/mux"
	"translate_go/utils"
)

type App struct {
	Client *translate.Client
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

func (a *App) Start() {
	fmt.Println("starting server on port 5000")
	http.ListenAndServe(":5000", a.Router)
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/", HomeHandler)
	a.Router.HandleFunc("/test", TestHandler)

	a.Router.HandleFunc("/list-languages", a.listLangs).Methods("GET")
}

func (a *App) listLangs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	lang := language.English

	fmt.Println("client: ", a.Client)

	langs, err := a.Client.SupportedLanguages(ctx, lang)
	if err != nil {
		msg := "Failed to get supported languages: " + err.Error()

		utils.RespondWithError(w, http.StatusInternalServerError, msg)
	}
	fmt.Println("langs: ", langs)

	utils.RespondWithJSON(w, http.StatusOK, langs)
	return
}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"home": "en casa!"}
	utils.RespondWithJSON(w, http.StatusOK, response)
	return
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"test": "success"}
	utils.RespondWithJSON(w, http.StatusOK, response)
}


//func TranslateText(text string, language string) (string, error) {
//
//}

//func PickRandomLanguage() string {
//
//}

