package domain

import (
	"errors"
	"time"
)

type Quote struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type QuoteRepository interface {
	GetRandom() (*Quote, error)

	GetRandomByCategory(category string) (*Quote, error)

	Add(quote *Quote) error

	GetAll() ([]*Quote, error)
}

func (q *Quote) Validate() error {
	if q.Text == "" {
		return errors.New("text cannot be empty")
	}
	if q.Author == "" {
		return errors.New("author cannot be empty")
	}
	if q.Category == "" {
		return errors.New("category cannot be empty")
	}
	return nil
}
