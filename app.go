package main

import (
	"cloud.google.com/go/translate"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"translate_go/utils"
)

type TranslateData struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

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
	a.Router.HandleFunc("/", HomeHandler).Methods("GET")
	a.Router.HandleFunc("/test", TestHandler).Methods("GET")

	a.Router.HandleFunc("/list-languages", a.listLangs).Methods("GET")
	a.Router.HandleFunc("/translate", a.translateText).Methods("POST")
}

func (a *App) translateText(w http.ResponseWriter, r *http.Request) {
	fmt.Println("translate: ", r.URL, " ", r.Method)

	var translationData TranslateData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&translationData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	fmt.Println("translationData: ", translationData)

	text := translationData.Text
	fmt.Println("text to translate: ", text)

	lang, err := language.Parse(translationData.Lang)
	if err != nil {
		msg := "Could not parse the target language.  Verify that it is an available option and formatted correctly (ex. 'en' for english) "
		utils.RespondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	ctx := r.Context()
	resp, err := a.Client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func (a *App) listLangs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("list-languages: ", r.URL)
	if r.Method != http.MethodGet {
		msg := "You cannot use that method"
		utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}

	targetLang := r.URL.Query().Get("target")

	if targetLang == "" { // TODO: make target language optional, with default being "en" ??
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
