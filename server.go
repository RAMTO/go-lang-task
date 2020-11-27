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

	finalMap := make(map[string][]map[string]string)
	finalSlice := make([]map[string]string, 0)

	for key, value := range historyMap {
		mapElement := map[string]string{key: value}
		finalSlice = append(finalSlice, mapElement)
	}

	finalMap["history"] = finalSlice

	json.NewEncoder(w).Encode(finalMap)
}
