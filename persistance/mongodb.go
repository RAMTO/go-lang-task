package translation

import (
	"context"
	"github.com/RAMTO/go-lang-task/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type TranslationRepository struct {
	db *mongo.Database
}

func (r *TranslationRepository) SaveWord(word *model.Word) (string, error) {
	collection := r.db.Collection("words")
	_, err := collection.InsertOne(context.TODO(), word)
	if err != nil {
		return "", err
	}

	return uuid.New().String(), nil
}

func (r *TranslationRepository) SaveSentence(sentence *model.Sentence) (string, error) {
	collection := r.db.Collection("sentences")
	_, err := collection.InsertOne(context.TODO(), sentence)
	if err != nil {
		return "", err
	}

	return uuid.New().String(), nil
}

func NewTranslationRepository(db *mongo.Database) *TranslationRepository {
	return &TranslationRepository{db: db}
}