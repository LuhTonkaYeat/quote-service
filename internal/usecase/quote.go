package usecase

import (
	"fmt"

	"github.com/LuhTonkaYeat/quote-service/internal/domain"
	"github.com/google/uuid"
)

type QuoteUseCase struct {
	repo domain.QuoteRepository
}

func NewQuoteUseCase(repo domain.QuoteRepository) *QuoteUseCase {
	return &QuoteUseCase{
		repo: repo,
	}
}

func (uc *QuoteUseCase) GetRandomQuote() (*domain.Quote, error) {
	quote, err := uc.repo.GetRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}
	return quote, nil
}

func (uc *QuoteUseCase) GetRandomQuoteByCategory(category string) (*domain.Quote, error) {
	if category == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}

	quote, err := uc.repo.GetRandomByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote by category %s: %w", category, err)
	}
	return quote, nil
}

func (uc *QuoteUseCase) AddQuote(text, author, category string) (*domain.Quote, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}
	if category == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}

	quote := &domain.Quote{
		ID:       uuid.New().String(),
		Text:     text,
		Author:   author,
		Category: category,
	}

	if err := uc.repo.Add(quote); err != nil {
		return nil, fmt.Errorf("failed to add quote: %w", err)
	}

	return quote, nil
}
