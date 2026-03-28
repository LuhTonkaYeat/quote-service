package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/LuhTonkaYeat/quote-service/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewQuoteServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	command := os.Args[1]

	switch command {
	case "random":
		getRandomQuote(ctx, client)

	case "category":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run client/main.go category <category>")
			fmt.Println("Example: go run client/main.go category motivation")
			os.Exit(1)
		}
		category := os.Args[2]
		getQuoteByCategory(ctx, client, category)

	case "add":
		if len(os.Args) < 5 {
			fmt.Println("Usage: go run client/main.go add <text> <author> <category>")
			fmt.Println("Example: go run client/main.go add \"Keep going\" \"Unknown\" \"motivation\"")
			os.Exit(1)
		}
		text := os.Args[2]
		author := os.Args[3]
		category := os.Args[4]
		addQuote(ctx, client, text, author, category)

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Quote Service CLI Client")
	fmt.Println("\nUsage:")
	fmt.Println("  go run client/main.go random")
	fmt.Println("    Get a random quote")
	fmt.Println()
	fmt.Println("  go run client/main.go category <category>")
	fmt.Println("    Get a random quote from specific category")
	fmt.Println("    Categories: motivation, life, wisdom, coding, etc.")
	fmt.Println()
	fmt.Println("  go run client/main.go add <text> <author> <category>")
	fmt.Println("    Add a new quote")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run client/main.go random")
	fmt.Println("  go run client/main.go category motivation")
	fmt.Println("  go run client/main.go add \"Code is poetry\" \"Developer\" \"coding\"")
}

func getRandomQuote(ctx context.Context, client pb.QuoteServiceClient) {
	resp, err := client.GetRandom(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Error getting random quote: %v", err)
	}

	printQuote(resp)
}

func getQuoteByCategory(ctx context.Context, client pb.QuoteServiceClient, category string) {
	resp, err := client.GetByCategory(ctx, &pb.CategoryRequest{Category: category})
	if err != nil {
		log.Fatalf("Error getting quote by category: %v", err)
	}

	printQuote(resp)
}

func addQuote(ctx context.Context, client pb.QuoteServiceClient, text, author, category string) {
	resp, err := client.AddQuote(ctx, &pb.AddQuoteRequest{
		Text:     text,
		Author:   author,
		Category: category,
	})
	if err != nil {
		log.Fatalf("Error adding quote: %v", err)
	}

	if resp.Success {
		fmt.Printf("✅ Quote added successfully!\n")
		fmt.Printf("   ID: %s\n", resp.Id)
		fmt.Printf("   Message: %s\n", resp.Message)
	} else {
		fmt.Printf("❌ Failed to add quote: %s\n", resp.Message)
	}
}

func printQuote(quote *pb.Quote) {
	fmt.Println("\n" + strings.Repeat("─", 60))
	fmt.Printf("📝 Quote ID: %s\n", quote.Id)
	fmt.Printf("💬 \"%s\"\n", quote.Text)
	fmt.Printf("✍️  Author: %s\n", quote.Author)
	fmt.Printf("🏷️  Category: %s\n", quote.Category)
	fmt.Println(strings.Repeat("─", 60) + "\n")
}
