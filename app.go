package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	"translate-on-the-go/cache"
	"translate-on-the-go/utils"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	Cache  *cache.Cache
	Client *translate.Client
	Router *mux.Router
}

const FIFTEEN_DAYS = time.Second * 86400 * 15

// Returns the available routes
func IndexHandler(w http.ResponseWriter, _r *http.Request) {
	routes := map[string]string{
		"/list-languages": "GET",
		"/translate":      "POST",
	}

	response := map[string]map[string]string{"routes": routes}

	utils.RespondWithJSON(w, http.StatusOK, response)
	return
}

func (a *App) Init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading the .env file")
	}

	apiKey := os.Getenv("TRANSLATE_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing api key")
	}

	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create translation client: %v", err)
	}
	a.Client = client
	defer client.Close()

	a.Cache = cache.NewCache()
	a.Router = mux.NewRouter()

	a.initRoutes()
}

func (a *App) Start() {
	port := os.Getenv("PORT")

	fmt.Printf("starting server on port %s \n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), a.Router))
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/", IndexHandler).Methods("GET")

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

	// the decoded translation text
	text := translationData.Text

	lang, err := language.Parse(translationData.Lang)
	if err != nil {
		msg := "Could not parse the target language.  Verify that it is an available option and formatted correctly (ex. 'en' for english) "
		utils.RespondWithError(w, http.StatusBadRequest, msg)
		return
	}

	// TODO - might have to hash/encode this, or only do it for short strings
	// the key includes the target language code and the text (i.e. "en-hola" for the case of wanting to translate "hola" to english)
	key := fmt.Sprintf("%s-%s", lang.String(), text)

	// check to see if the translation is in the cache
	val, err := a.Cache.Get(key)
	if err != nil {
		if err == redis.Nil {
			fmt.Println("the key does not exist in the cache")
		} else {
			fmt.Println("cache error: ", err.Error())
			fmt.Println("there was a problem getting the value from the cache. Continuing...")
		}
	}

	if len(val) > 0 {
		// return the value from the cache
		var cachedTranslation TranslationResponse
		err := json.Unmarshal(val, &cachedTranslation)
		if err != nil {
			fmt.Println("error marshalling the cached data to JSON ", err)
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, cachedTranslation)
		return
	}

	ctx := r.Context()
	resp, err := a.Client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(resp) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("Translate returned empty response to text: %s", text).Error())
		return
	}

	translationObject := TranslationResponse{
		resp[0].Source,
		lang,
		resp[0].Text,
	}

	// set the translation to the cache
	err = a.Cache.Set(key, translationObject, FIFTEEN_DAYS)
	if err != nil {
		fmt.Println("cache set error: ", err.Error())
		fmt.Println("there was a problem setting the value to the cache.")
	}

	transRes := map[string]TranslationResponse{
		"response": translationObject,
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
