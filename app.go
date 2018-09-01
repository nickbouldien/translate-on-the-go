package main

import (
	"cloud.google.com/go/translate"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"log"
	"net/http"

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

	a.Router.HandleFunc("/list-languages", a.listLangs).Methods("GET")
	a.Router.HandleFunc("/translate", a.translateHandler).Methods("GET", "POST")
}

func (a *App) translateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)

	lg := r.FormValue("key1")
	//txt := r.FormValue("text")

	fmt.Println("translate handler: ", lg)
	//fmt.Println("translate handler: ", lg)

	ctx := r.Context()

	lang := language.BrazilianPortuguese
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
	fmt.Println("listlangs")
	fmt.Println("method: ", r.Method)
	if r.Method != http.MethodGet {
		msg := "You cannot use that method"
		utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}

	targetLang := r.URL.Query().Get("target")
	fmt.Println("targetLang: ", targetLang)

	if targetLang == "" { // TODO: make target language optional, with default being "en"
		msg := "You must provide a target language (ex. /list-languages?target=pt)"
		utils.RespondWithError(w, http.StatusBadRequest, msg)
		return
	}

	target, err := language.Parse(targetLang)
	fmt.Println("parsed Target: ", target)
	if err != nil {
		msg := "Could not parse the target language.  Verify that it is an available option and formatted correctly (ex. 'en' for english) "
		utils.RespondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	ctx := r.Context()

	langs, err := a.Client.SupportedLanguages(ctx, target)
	if err != nil {
		msg := "Failed to get supported languages: " + err.Error()
		utils.RespondWithError(w, http.StatusInternalServerError, msg)
		return
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

