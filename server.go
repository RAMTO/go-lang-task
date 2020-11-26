package main

import (
	// "fmt"
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"os"

	"github.com/gorilla/mux"
)

type Word struct {
	Value string `json:"english-word"`
}

type TranslatedWord struct {
	Value string `json:"gopher-word"`
}

func main() {
	// Init router
	r := mux.NewRouter()
	port := os.Getenv("PORT");

	// Route handles & endpoints
	r.HandleFunc("/word", handleWordPostRequest).Methods("POST")

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

func handleWordPostRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var word Word
	var translated TranslatedWord

	_ = json.NewDecoder(r.Body).Decode(&word)

	translated.Value = translateWord(word.Value);

	json.NewEncoder(w).Encode(translated)
}
