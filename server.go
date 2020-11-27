package main

import (
	"fmt"
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

	// Start server
	log.Fatal(http.ListenAndServe(":" + port, r))
}

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
	vowels := [6]string{"a", "e", "i", "o", "u", "y"}
	consonantLetters := "xr"
	fistChar := word[0:1]
	prexif := "g"
	prexifConsonant := "ge"
	
	// Check for vowels
	if itemExists(vowels, fistChar) {
		word = prexif + word
	} 
		
	// Check for consonant
	if strings.HasPrefix(word, consonantLetters) {
		word = prexifConsonant + word
	}

	return word
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

func saveTranslations(translation string) {
	// update translation history array
}

func handleWordPostRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var word Word
	var translated TranslatedWord

	_ = json.NewDecoder(r.Body).Decode(&word)

	translated.Value = translateWord(word.Value);

	json.NewEncoder(w).Encode(translated)
}

func handleSentencePostRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var sentence Sentence
	var translatedSentence TranslatedSentence

	_ = json.NewDecoder(r.Body).Decode(&sentence)
	
	translatedSentence.Value = translateSentence(sentence.Value);
	
	fmt.Print(translatedSentence.Value)

	json.NewEncoder(w).Encode(translatedSentence)
}
