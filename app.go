package main

import (
	"cloud.google.com/go/translate"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"translate_go/utils"
)

type TranslateData struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

type TranslationResponse struct {
	Source         language.Tag `json:"sourceLanguage"`
	TargetLanguage language.Tag `json:"targetLanguage"`
	TranslatedText string       `json:"translatedText"`
}

type App struct {
	Client *translate.Client
	Router *mux.Router
}

const apiKey = "TRANSLATE_API_KEY"

func (a *App) Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}

	apiKey := os.Getenv(apiKey)

	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	a.Client = client
	defer client.Close()

	a.Router = mux.NewRouter()
	a.initRoutes()
}

func (a *App) Start() {
	port := os.Getenv("PORT")

	fmt.Printf("starting server on port %s \n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port) , a.Router))
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/", HomeHandler).Methods("GET")
	a.Router.HandleFunc("/test", TestHandler).Methods("GET")

	a.Router.HandleFunc("/list-languages", a.listLangs).Methods("GET")
	a.Router.HandleFunc("/translate", a.translateText).Methods("POST")
}

func (a *App) translateText(w http.ResponseWriter, r *http.Request) {
	var translationData TranslateData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&translationData)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	text := translationData.Text

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

	var transRes = make(map[string]TranslationResponse)
	transRes["response"] = TranslationResponse{
		resp[0].Source,
		lang,
		resp[0].Text,
	}

	utils.RespondWithJSON(w, http.StatusOK, transRes)
}

func (a *App) listLangs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		msg := "You cannot use that method. Only the `GET` method is allowed."
		utils.RespondWithError(w, http.StatusMethodNotAllowed, msg)
		return
	}

	targetLang := r.URL.Query().Get("target")

	if targetLang == "" { // TODO: make target language optional, with default being english ("en") ??
		msg := "You must provide a target language (ex. /list-languages?target=pt)"
		utils.RespondWithError(w, http.StatusBadRequest, msg)
		return
	}

	target, err := language.Parse(targetLang)
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

	utils.RespondWithJSON(w, http.StatusOK, langs)
	return
}

func HomeHandler(w http.ResponseWriter, _r *http.Request) {
	routes := map[string]string{
		"/list-languages": "GET",
		"/translate": "POST",
	}

	response := map[string]map[string]string{"routes": routes}

	utils.RespondWithJSON(w, http.StatusOK, response)
	return
}

func TestHandler(w http.ResponseWriter, _r *http.Request) {
	response := map[string]string{"test": "success"}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
