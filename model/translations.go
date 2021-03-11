package model

type Word struct {
	Original 	string 	`json:"original" bson:"original"` 
	Translated 	string 	`json:"translated" bson:"translated"`
}

type Sentence struct {
	Original 	string 	`json:"original" bson:"original"` 
	Translated 	string 	`json:"translated" bson:"translated"`
}