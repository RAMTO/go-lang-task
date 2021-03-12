package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"os"
	"context"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	translationsMongo "github.com/RAMTO/go-lang-task/persistance"
	translationsModel "github.com/RAMTO/go-lang-task/model"
)

// Word is...
type Word struct {
	Value string `json:"english-word" bson:"english-word"`
}

// Sentence is...
type Sentence struct {
	Value string `json:"english-sentence" bson:"english-sentence"`
}

// TranslatedWord is...
type TranslatedWord struct {
	Value string `json:"gopher-word" bson:"gopher-word"`
}

// TranslatedSentence is...
type TranslatedSentence struct {
	Value string `json:"gopher-sentence" bson:"gopher-sentence"`
}

var historyMap = make(map[string]string)

var client *mongo.Client
var db *mongo.Database
var translationsRepo *translationsMongo.TranslationRepository;

func main() {
	// Connect to Mongo
	client, db = connectToDb()
	err := client.Ping(context.Background(), readpref.Primary())

	defer client.Disconnect(context.Background())

	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	translationsRepo = translationsMongo.NewTranslationRepository(db)

	// Init router
	r := mux.NewRouter()
	var port string

	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	} else {
		port = os.Args[1]
	}

	// Route handles & endpoints
	r.HandleFunc("/word", handleWordPostRequest).Methods("POST")
	r.HandleFunc("/sentence", handleSentencePostRequest).Methods("POST")
	r.HandleFunc("/history", handleHistoryGetRequest).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":" + port, r))
}

func connectToDb() (*mongo.Client, *mongo.Database) {
	clientOptions := options.Client().ApplyURI("mongodb+srv://test:test123@cluster0.3hy86.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db := client.Database("translations")

	return client, db
}

// Helper functions
func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func translateWord(word string) string {
	translated := word
	vowels := [6]string{"a", "e", "i", "o", "u", "y"}
	consonants := [20]string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "z"}
	consonantLetters := "xr"
	
	prexif := "g"
	prexifConsonant := "ge"
	suffix := "ogo"
	additionalCheck := "qu"
	
	fistChar := word[0:1]
	firstCharRemoved := strings.Replace(word, fistChar, "", -1)

	// Check if word starts with vowel
	if itemExists(vowels, fistChar) { 
		translated = prexif + word
	} else {
		// Check for spesific consonant
		if strings.HasPrefix(word, consonantLetters) { 
			translated = prexifConsonant + word
		} else { 
			// Additional check for "qu"
			if(strings.HasPrefix(firstCharRemoved, additionalCheck)) { 
				replaced := strings.Replace(firstCharRemoved, additionalCheck, "", -1)
				translated = replaced + fistChar + additionalCheck + suffix
			} else {
				consonantsToBeReplacedSlice := []string{}
	
				for _, value := range word {
					if(itemExists(consonants, string(value))) {
						consonantsToBeReplacedSlice = append(consonantsToBeReplacedSlice, string(value))
					}else {
						break
					}
				}
				
				consonantsToBeReplaced := strings.Join(consonantsToBeReplacedSlice[:], "")
	
				replaced := strings.Replace(word, consonantsToBeReplaced, "", -1)
				translated = replaced + consonantsToBeReplaced + suffix
			}
		} 
	}
		
	return translated
}

func translateSentence(sentence string) string {
	words := strings.Fields(sentence)
	translatedWords := []string{};

	for _, word := range words {
		translatedWords = append(translatedWords, translateWord(word))
	}

	translated := strings.Join(translatedWords[:], " ")

	return translated
}

// Route handlers
func handleWordPostRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var word Word
	var translated TranslatedWord

	_ = json.NewDecoder(r.Body).Decode(&word)

	translated.Value = translateWord(word.Value);

	// Log in history
	historyMap[word.Value] = translated.Value

	wordObj := &translationsModel.Word{Original: word.Value, Translated: translated.Value}
	translationsRepo.SaveWord(wordObj)

	json.NewEncoder(w).Encode(translated)
}

func handleSentencePostRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var sentence Sentence
	var translatedSentence TranslatedSentence

	_ = json.NewDecoder(r.Body).Decode(&sentence)
	
	translatedSentence.Value = translateSentence(sentence.Value);
	
	// Log in history
	historyMap[sentence.Value] = translatedSentence.Value
	
	sentenceObj := &translationsModel.Sentence{Original: sentence.Value, Translated: translatedSentence.Value}
	translationsRepo.SaveSentence(sentenceObj)

	json.NewEncoder(w).Encode(translatedSentence)
}

func handleHistoryGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	finalMap := make(map[string][]map[string]string)
	finalSlice := make([]map[string]string, 0)

	for key, value := range historyMap {
		mapElement := map[string]string{key: value}
		finalSlice = append(finalSlice, mapElement)
	}

	finalMap["history"] = finalSlice

	json.NewEncoder(w).Encode(finalMap)
}
