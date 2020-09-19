## Translate on the Go

### Inspiration
I'm always looking up words/phrases (usually en español ou em português).

(I know this is pointless since a thing called Google Translate exists, but this is a way for me to mess around with Go on something I might use at one point or another.)


### Setup:
- `go install`
- `go build`
- `./translate-on-the-go`
- test it out on port 5000 (or whatever endpoint/port you decide to use) with one of the routes


### Routes:
- GET `/list-languages?target={target_language_code}` - lists all possible languages to translate to
- POST `/translate` - translate text to a target language
- GET `/` - responds with the available routes

target language code info: https://cloud.google.com/translate/docs/languages

### Sample requests:

1 - list languages with a target of english (`en`)

request:
```bash
curl --request GET \
  --url 'http://localhost:5000/list-languages?target=en'
```

response (truncated here for display purposes):
```json
[
  {
    "Name": "Afrikaans",
    "Tag": "af"
  },
  {
    "Name": "Albanian",
    "Tag": "sq"
  },
  {
    "Name": "Amharic",
    "Tag": "am"
  },
  {
    "Name": "Arabic",
    "Tag": "ar"
  }
]
```


2 - translate the word "hello" from english to português

request:
```bash
curl --request POST \
  --url http://localhost:5000/translate \
  --header 'content-type: application/json' \
  --data '{
	"lang": "pt",
	"text": "hello"
}'
```

response:

```json
{
  "response": {
    "sourceLanguage": "en",
    "targetLanguage": "pt",
    "translatedText": "Olá"
  }
}
```

3 - the base/home route to list the available routes

request:
```bash
curl --request GET --url http://localhost:5000/
```

response:

```json
{
  "routes": {
    "/list-languages": "GET",
    "/translate": "POST"
  }
}
```


### Packages:
- [go-redis](https://github.com/go-redis/redis)
- [gorilla mux router](https://github.com/gorilla/mux)
- [godotenv](https://github.com/joho/godotenv)


### Resources:
- [Google Translate](https://cloud.google.com/translate/)
- [Google Cloud Translation API](https://cloud.google.com/translate/docs/) - responsible for all translations
- [webapp-with-golang-anti-textbook](https://thewhitetulip.gitbooks.io/webapp-with-golang-anti-textbook/content/)
- [Getting Started with Redis and Go - Tutorial](https://tutorialedge.net/golang/go-redis-tutorial/)
- [just for func](https://www.youtube.com/channel/UC_BzFbxG2za3bp5NRRRXJSw)
