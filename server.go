package main

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"os"

	"github.com/gorilla/mux"
)

// Word is...
type Word struct {
	Value string `json:"english-word"`
}

// Sentence is...
type Sentence struct {
	Value string `json:"english-sentence"`
}

// TranslatedWord is...
type TranslatedWord struct {
	Value string `json:"gopher-word"`
}

// TranslatedSentence is...
type TranslatedSentence struct {
	Value string `json:"gopher-sentence"`
}

var historyMap = make(map[string]string)

func main() {
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
	consonants := [21]string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "z"}
	consonantLetters := "xr"
	
	fistChar := word[0:1]
	
	prexif := "g"
	prexifConsonant := "ge"
	suffix := "ogo"

	consonantsToBeReplacedSlice := []string{}

	for _, value := range word {
		if(itemExists(consonants, string(value))) {
			consonantsToBeReplacedSlice = append(consonantsToBeReplacedSlice, string(value))
		}else {
			break
		}
	}
	
	consonantsToBeReplaced := strings.Join(consonantsToBeReplacedSlice[:], "")
	
	if itemExists(vowels, fistChar) { // Check for vowels
		translated = prexif + word
	} else if strings.HasPrefix(word, consonantLetters) { // Check for consonant
		translated = prexifConsonant + word
	} else if itemExists(consonants, fistChar) { // Check for consonants
		replaced := strings.Replace(word, consonantsToBeReplaced, "", -1)
		translated = replaced + consonantsToBeReplaced + suffix
	} 
		
	return translated
}

func translateSentence(sentence string) string {
	// Translate sentence using translateWord function
	words := strings.Fields(sentence)
	translatedWords := [] string {};

	for i :=0; i < len(words); i++ {
		translatedWords = append(translatedWords, translateWord(words[i]))
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

	json.NewEncoder(w).Encode(translatedSentence)
}

func handleHistoryGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// fmt.Print(historyMap)
	finalMap := make(map[string]map[string]string)
	finalMap["history"] = historyMap
	json.NewEncoder(w).Encode(finalMap)
}
