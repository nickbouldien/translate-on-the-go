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

// TODO: make these structs with "abbreviation" ("en") and "display name" ("english") fields??
const (
	EN  = "en" // english
	ES  = "es" // español
	PT  = "pt" // português
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
	a.Router.HandleFunc("/test", TestHandler).Methods("GET")

	a.Router.HandleFunc("/list-languages", a.listLangs).Queries("target", "{target}").Methods("GET")
	a.Router.HandleFunc("/translate", a.translateHandler).Methods("GET", "POST")
}

func (a *App) translateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)

	lg := r.FormValue("key1")

	fmt.Println("translate handler: ", lg)
	ctx := r.Context()

	lang := language.Spanish
	//target, err := language.Parse("ru")
	//if err != nil {
	//	log.Fatalf("Failed to parse target language: %v", err)
	//}
	fmt.Println("lang: ", lang)

	text := "hola, mundo"
	fmt.Println("text to translate: ", text)

	resp, err := a.Client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

//func (a *App) translateText(ctx context.Context , text string, language string) (string, error) {
//	translation, err := a.Client.Translate(ctx, text, language )
//	if err != nil {
//		return nil, err
//	}
//
//	return translation, nil
//}

func (a *App) listLangs(w http.ResponseWriter, r *http.Request) {
	vals := mux.Vars(r)
	targetLang := vals["target"]
	fmt.Println("targetLang: ", targetLang)

	fmt.Println("method: ", r.Method)
	if r.Method != http.MethodGet {
		msg := "You cannot use that method "
		utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}

	ctx := r.Context()
	fmt.Println("listlangs")

	lang := language.English

	langs, err := a.Client.SupportedLanguages(ctx, lang)
	if err != nil {
		msg := "Failed to get supported languages: " + err.Error()

		utils.RespondWithError(w, http.StatusInternalServerError, msg)
	}
	//fmt.Println("langs: ", langs)
	//for _, lang := range langs {
	//	fmt.Fprintf(w, "%q: %s\n", lang.Tag, lang.Name)
	//}

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

//func PickRandomLanguage() string {
//
//}

