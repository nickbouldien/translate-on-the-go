## Translate on the Go

### Inspiration
I'm always looking up words/phrases (usually en español ou em português).

(I know this is pointless since a thing called Google Translate exists, but this is a way for me to mess around with Go on something I might use at one point or another.)


### Setup:
- go install
- go build
- ./travel_go
- test it out on port 5000 (or whatever endpoint/port you decide to use) with one of the routes


### Routes:
- /list-languages - lists all possible languages for
- /test - responds with a 200 status (used to verify server is up and running)


### Packages:


### Resources:
- https://cloud.google.com/translate/
- https://thewhitetulip.gitbooks.io/webapp-with-golang-anti-textbook/content/


### TODOs:
- some type of caching - see if requested text has already been translated.  if so, return that text (https://github.com/patrickmn/go-cache ???)
