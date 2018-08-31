package main

import (
	"cloud.google.com/go/translate"
	"golang.org/x/net/context"
	"golang.org/x/text/language"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
)

// TODO: make these structs with "abbreviation" ("en") and "display name" ("english") fields??
const (
	EN  = "en" // english
	ES  = "es" // español
	PT  = "pt" // português
)

func main() {
	ctx := context.Background()

	// Creates a client.
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the text to translate.

	selectedLang := EN
	text := "Hello, world!"
	// Sets the target language.

	target, err := language.Parse(selectedLang)
	if err != nil {
		log.Fatalf("Failed to parse target language: %v", err)
	}

	langs, err := client.SupportedLanguages(ctx, target)
	if err != nil {
		log.Fatalf("Failed to get supported languages: %v", err)
	}

	fmt.Println("langs: ", langs)

	// Translates the text into Russian.
	translations, err := client.Translate(ctx, []string{text}, target, nil)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	fmt.Printf("Text: %v\n", text)
	fmt.Printf("Translation: %v\n", translations[0].Text)


	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/test", TestHandler)
	http.Handle("/", r)
	fmt.Println("starting server on port 5000")
	http.ListenAndServe(":5000", r)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	return
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)

	response := map[string]string{"test": "success"}
	RespondWithJSON(w, http.StatusOK, response)
}


//func TranslateText(text string, language string) (string, error) {
//
//}

//func PickRandomLanguage() string {
//
//}
