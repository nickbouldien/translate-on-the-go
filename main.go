package main

func main() {
	app := App{}
	app.Init()

	app.Start()
}

//ctx := context.Background()
//
//// Creates a client.
//client, err := translate.NewClient(ctx)
//if err != nil {
//	log.Fatalf("Failed to create client: %v", err)
//}
//defer client.Close()
//
//// Sets the text to translate.
//
//selectedLang := EN
//text := "Hello, world!"
//// Sets the target language.
//
//target, err := language.Parse(selectedLang)
//if err != nil {
//	log.Fatalf("Failed to parse target language: %v", err)
//}
//
//langs, err := client.SupportedLanguages(ctx, target)
//if err != nil {
//	log.Fatalf("Failed to get supported languages: %v", err)
//}
//
//fmt.Println("langs: ", langs)
//
//// Translates the text into Russian.
//translations, err := client.Translate(ctx, []string{text}, target, nil)
//if err != nil {
//	log.Fatalf("Failed to translate text: %v", err)
//}
//
//fmt.Printf("Text: %v\n", text)
//fmt.Printf("Translation: %v\n", translations[0].Text)
