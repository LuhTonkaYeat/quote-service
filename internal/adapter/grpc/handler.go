package grpc

import (
	"context"
	"log"

	pb "github.com/LuhTonkaYeat/quote-service/api/proto"
	"github.com/LuhTonkaYeat/quote-service/internal/domain"
	"github.com/LuhTonkaYeat/quote-service/internal/usecase"
)

type QuoteHandler struct {
	pb.UnimplementedQuoteServiceServer
	useCase *usecase.QuoteUseCase
}

func NewQuoteHandler(useCase *usecase.QuoteUseCase) *QuoteHandler {
	return &QuoteHandler{
		useCase: useCase,
	}
}

func (h *QuoteHandler) GetRandom(ctx context.Context, req *pb.Empty) (*pb.Quote, error) {
	log.Println("GetRandom called")

	quote, err := h.useCase.GetRandomQuote()
	if err != nil {
		log.Printf("Error in GetRandom: %v", err)
		return nil, err
	}

	return convertToProto(quote), nil
}

func (h *QuoteHandler) GetByCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.Quote, error) {
	log.Printf("GetByCategory called with category: %s", req.Category)

	quote, err := h.useCase.GetRandomQuoteByCategory(req.Category)
	if err != nil {
		log.Printf("Error in GetByCategory: %v", err)
		return nil, err
	}

	return convertToProto(quote), nil
}

func (h *QuoteHandler) AddQuote(ctx context.Context, req *pb.AddQuoteRequest) (*pb.AddQuoteResponse, error) {
	log.Printf("AddQuote called: text=%s, author=%s, category=%s", req.Text, req.Author, req.Category)

	quote, err := h.useCase.AddQuote(req.Text, req.Author, req.Category)
	if err != nil {
		log.Printf("Error in AddQuote: %v", err)
		return &pb.AddQuoteResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.AddQuoteResponse{
		Id:      quote.ID,
		Success: true,
		Message: "Quote added successfully",
	}, nil
}

func convertToProto(quote *domain.Quote) *pb.Quote {
	return &pb.Quote{
		Id:       quote.ID,
		Text:     quote.Text,
		Author:   quote.Author,
		Category: quote.Category,
	}
}
