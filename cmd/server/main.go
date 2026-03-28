package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/LuhTonkaYeat/quote-service/api/proto"
	grpchandler "github.com/LuhTonkaYeat/quote-service/internal/adapter/grpc"
	"github.com/LuhTonkaYeat/quote-service/internal/adapter/repository"
	"github.com/LuhTonkaYeat/quote-service/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Quote Service...")

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "quotes.db"
	}

	repo, err := repository.NewSQLiteQuoteRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	log.Printf("Connected to SQLite database: %s", dbPath)

	quoteUseCase := usecase.NewQuoteUseCase(repo)

	quoteHandler := grpchandler.NewQuoteHandler(quoteUseCase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterQuoteServiceServer(grpcServer, quoteHandler)

	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server listening on :%s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
