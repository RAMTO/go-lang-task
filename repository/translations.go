package repository

import "github.com/RAMTO/go-lang-task/model"

type TranlationsRepository interface {
	SaveWord(*model.Word) (id string, err error)
	SaveSentence(*model.Sentence) (id string, err error)
}