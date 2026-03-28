package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/LuhTonkaYeat/quote-service/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteQuoteRepository struct {
	db *sql.DB
}

func NewSQLiteQuoteRepository(dbPath string) (*SQLiteQuoteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS quotes (
		id TEXT PRIMARY KEY,
		text TEXT NOT NULL,
		author TEXT NOT NULL,
		category TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_category ON quotes(category);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	repo := &SQLiteQuoteRepository{db: db}
	repo.seedInitialQuotes()

	return repo, nil
}

func (r *SQLiteQuoteRepository) seedInitialQuotes() {
	quotes := []*domain.Quote{
		{
			ID:       "1",
			Text:     "The only way to do great work is to love what you do.",
			Author:   "Steve Jobs",
			Category: "motivation",
		},
		{
			ID:       "2",
			Text:     "Life is what happens when you're busy making other plans.",
			Author:   "John Lennon",
			Category: "life",
		},
		{
			ID:       "3",
			Text:     "Stay hungry, stay foolish.",
			Author:   "Steve Jobs",
			Category: "motivation",
		},
		{
			ID:       "4",
			Text:     "The journey of a thousand miles begins with one step.",
			Author:   "Lao Tzu",
			Category: "wisdom",
		},
	}

	for _, quote := range quotes {
		var count int
		r.db.QueryRow("SELECT COUNT(*) FROM quotes WHERE id = ?", quote.ID).Scan(&count)
		if count == 0 {
			_, err := r.db.Exec(
				"INSERT INTO quotes (id, text, author, category) VALUES (?, ?, ?, ?)",
				quote.ID, quote.Text, quote.Author, quote.Category,
			)
			if err != nil {
				log.Printf("Warning: failed to seed quote: %v", err)
			}
		}
	}
}

func (r *SQLiteQuoteRepository) GetRandom() (*domain.Quote, error) {
	var quote domain.Quote
	query := `SELECT id, text, author, category, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`

	err := r.db.QueryRow(query).Scan(
		&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no quotes found")
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return &quote, nil
}

func (r *SQLiteQuoteRepository) GetRandomByCategory(category string) (*domain.Quote, error) {
	var quote domain.Quote
	query := `SELECT id, text, author, category, created_at 
	          FROM quotes 
	          WHERE category = ? 
	          ORDER BY RANDOM() 
	          LIMIT 1`

	err := r.db.QueryRow(query, category).Scan(
		&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no quotes found in category: %s", category)
		}
		return nil, fmt.Errorf("failed to get random quote by category: %w", err)
	}

	return &quote, nil
}

func (r *SQLiteQuoteRepository) Add(quote *domain.Quote) error {
	if err := quote.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO quotes (id, text, author, category) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, quote.ID, quote.Text, quote.Author, quote.Category)
	if err != nil {
		return fmt.Errorf("failed to add quote: %w", err)
	}

	return nil
}

func (r *SQLiteQuoteRepository) GetAll() ([]*domain.Quote, error) {
	rows, err := r.db.Query("SELECT id, text, author, category, created_at FROM quotes")
	if err != nil {
		return nil, fmt.Errorf("failed to get all quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*domain.Quote
	for rows.Next() {
		var quote domain.Quote
		err := rows.Scan(&quote.ID, &quote.Text, &quote.Author, &quote.Category, &quote.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, &quote)
	}

	return quotes, nil
}

func (r *SQLiteQuoteRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
