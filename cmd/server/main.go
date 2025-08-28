package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/nabil/book-store-system/config"
	"github.com/nabil/book-store-system/internal/repository"
	"github.com/nabil/book-store-system/internal/service"
	"github.com/nabil/book-store-system/internal/transport/grpc"
	"github.com/nabil/book-store-system/pkg/database"
	"github.com/nabil/book-store-system/pkg/helpers"
	"github.com/nabil/book-store-system/pkg/logger"
	"github.com/nabil/book-store-system/proto"
	grpcServer "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Info("Starting Book Store gRPC Server...")

	// Load configuration
	cfg := config.LoadConfig()
	logger.Infof("Configuration loaded: gRPC Port=%d", cfg.GRPCPort)

	// Initialize database
	database.Connect()
	db := database.DB
	logger.Info("Database connected successfully")

	// Run database migrations
	database.Migrate()
	logger.Info("Database migration completed")

	// Initialize Midtrans
	helpers.InitMidtrans()
	logger.Info("Midtrans payment gateway initialized")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	bookRepo := repository.NewBookRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reportRepo := repository.NewReportRepository(db)
	txRepo := repository.NewTransactionRepository(db)
	logger.Info("Repositories initialized")

	// Initialize services
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo, userRepo)
	bookService := service.NewBookService(bookRepo, categoryRepo, userRepo)
	orderService := service.NewOrderService(orderRepo, bookRepo, userRepo, txRepo)
	reportService := service.NewReportService(reportRepo, userRepo)
	logger.Info("Services initialized")

	// Initialize gRPC handlers
	userHandler := grpc.NewUserHandler(userService)
	categoryHandler := grpc.NewCategoryHandler(categoryService)
	bookHandler := grpc.NewBookHandler(bookService)
	orderHandler := grpc.NewOrderHandler(orderService)
	reportHandler := grpc.NewReportHandler(reportService)
	logger.Info("gRPC handlers initialized")

	// Create gRPC server
	grpcSrv := grpcServer.NewServer()

	// Register services
	proto.RegisterUserServiceServer(grpcSrv, userHandler)
	proto.RegisterCategoryServiceServer(grpcSrv, categoryHandler)
	proto.RegisterBookServiceServer(grpcSrv, bookHandler)
	proto.RegisterOrderServiceServer(grpcSrv, orderHandler)
	proto.RegisterReportServiceServer(grpcSrv, reportHandler)

	reflection.Register(grpcSrv)
	logger.Info("gRPC services registered")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.GRPCPort))
	if err != nil {
		logger.Errorf("Failed to listen on port %d: %v", cfg.GRPCPort, err)
		log.Fatal(err)
	}

	go func() {
		logger.Infof("gRPC server starting on port %d", cfg.GRPCPort)
		if err := grpcSrv.Serve(lis); err != nil {
			logger.Errorf("Failed to serve gRPC server: %v", err)
			log.Fatal(err)
		}
	}()

	logger.Info("Book Store gRPC Server started successfully")
	fmt.Printf("gRPC Server is running on port %d\n", cfg.GRPCPort)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logger.Info("Shutting down gRPC server...")

	grpcSrv.GracefulStop()
	logger.Info("gRPC server stopped")

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
		logger.Info("Database connection closed")
	}

	logger.Info("Server shutdown complete")
	fmt.Println("Server shutdown complete")
}
