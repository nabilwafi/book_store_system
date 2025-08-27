package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nabil/book-store-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create clients
	userClient := proto.NewUserServiceClient(conn)
	categoryClient := proto.NewCategoryServiceClient(conn)
	bookClient := proto.NewBookServiceClient(conn)

	fmt.Println("=== Testing gRPC CRUD Operations ===")

	// First, register and login to get token
	token := setupUserAndGetToken(userClient)
	if token == "" {
		return
	}

	// Test Category CRUD
	categoryID := testCategoryCRUD(categoryClient, token)

	// Test Book CRUD
	testBookCRUD(bookClient, categoryClient, token, categoryID)

	fmt.Println("\n=== CRUD Testing Complete ===")
}

func setupUserAndGetToken(client proto.UserServiceClient) string {
	fmt.Println("\n1. Setting up user and getting token...")

	// Register user
	registerReq := &proto.RegisterRequest{
		Name:     "Admin User",
		Email:    "admin@bookstore.com",
		Password: "admin123",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Register(ctx, registerReq)
	if err != nil {
		// User might already exist, try to login
		fmt.Printf("Register failed (user might exist): %v\n", err)
	}

	// Login to get token
	loginReq := &proto.LoginRequest{
		Email:    "admin@bookstore.com",
		Password: "admin123",
	}

	loginResp, err := client.Login(ctx, loginReq)
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		return ""
	}

	if loginResp.Success {
		fmt.Printf("✅ Login successful! Token obtained.\n")
		return loginResp.Token
	} else {
		fmt.Printf("❌ Login failed: %s\n", loginResp.Message)
		return ""
	}
}

func testCategoryCRUD(client proto.CategoryServiceClient, token string) uint32 {
	fmt.Println("\n2. Testing Category CRUD Operations")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create Category
	fmt.Println("\n2.1 Creating Category...")
	createReq := &proto.CreateCategoryRequest{
		Name:  fmt.Sprintf("Fiction_%d", time.Now().Unix()),
		Token: token,
	}

	createResp, err := client.CreateCategory(ctx, createReq)
	if err != nil {
		fmt.Printf("❌ Create category failed: %v\n", err)
		return 0
	}

	if createResp.Success {
		fmt.Printf("✅ Category created successfully!\n")
		fmt.Printf("   ID: %d, Name: %s\n", createResp.Category.Id, createResp.Category.Name)
	} else {
		fmt.Printf("❌ Create category failed: %s\n", createResp.Message)
		return 0
	}

	categoryID := createResp.Category.Id

	// Get Categories
	fmt.Println("\n2.2 Getting Categories...")
	getCategoriesReq := &proto.GetCategoriesRequest{
		Page:  1,
		Limit: 10,
	}

	getCategoriesResp, err := client.GetCategories(ctx, getCategoriesReq)
	if err != nil {
		fmt.Printf("❌ Get categories failed: %v\n", err)
	} else if getCategoriesResp.Success {
		fmt.Printf("✅ Categories retrieved successfully!\n")
		fmt.Printf("   Total: %d categories\n", getCategoriesResp.Total)
		for _, cat := range getCategoriesResp.Categories {
			fmt.Printf("   - ID: %d, Name: %s\n", cat.Id, cat.Name)
		}
	}

	// Get Single Category
	fmt.Println("\n2.3 Getting Single Category...")
	getCategoryReq := &proto.GetCategoryRequest{
		Id: categoryID,
	}

	getCategoryResp, err := client.GetCategory(ctx, getCategoryReq)
	if err != nil {
		fmt.Printf("❌ Get category failed: %v\n", err)
	} else if getCategoryResp.Success {
		fmt.Printf("✅ Category retrieved successfully!\n")
		fmt.Printf("   ID: %d, Name: %s\n", getCategoryResp.Category.Id, getCategoryResp.Category.Name)
	}

	// Update Category
	fmt.Println("\n2.4 Updating Category...")
	updateReq := &proto.UpdateCategoryRequest{
		Id:    categoryID,
		Name:  "Science Fiction",
		Token: token,
	}

	updateResp, err := client.UpdateCategory(ctx, updateReq)
	if err != nil {
		fmt.Printf("❌ Update category failed: %v\n", err)
	} else if updateResp.Success {
		fmt.Printf("✅ Category updated successfully!\n")
		fmt.Printf("   ID: %d, Name: %s\n", updateResp.Category.Id, updateResp.Category.Name)
	}

	return categoryID
}

func testBookCRUD(bookClient proto.BookServiceClient, categoryClient proto.CategoryServiceClient, token string, categoryID uint32) {
	fmt.Println("\n3. Testing Book CRUD Operations")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create Book
	fmt.Println("\n3.1 Creating Book...")
	createBookReq := &proto.CreateBookRequest{
		Title:      "The Hitchhiker's Guide to the Galaxy",
		Author:     "Douglas Adams",
		Price:      15.99,
		Stock:      50,
		Year:       1979,
		CategoryId: categoryID,
		Token:      token,
	}

	createBookResp, err := bookClient.CreateBook(ctx, createBookReq)
	if err != nil {
		fmt.Printf("❌ Create book failed: %v\n", err)
		return
	}

	var bookID uint32
	if createBookResp.Success {
		fmt.Printf("✅ Book created successfully!\n")
		fmt.Printf("   ID: %d, Title: %s, Author: %s, Price: $%.2f\n",
			createBookResp.Book.Id, createBookResp.Book.Title,
			createBookResp.Book.Author, createBookResp.Book.Price)
		bookID = createBookResp.Book.Id
	} else {
		fmt.Printf("❌ Create book failed: %s\n", createBookResp.Message)
		return
	}

	// Get Books
	fmt.Println("\n3.2 Getting Books...")
	getBooksReq := &proto.GetBooksRequest{
		Page:  1,
		Limit: 10,
	}

	getBooksResp, err := bookClient.GetBooks(ctx, getBooksReq)
	if err != nil {
		fmt.Printf("❌ Get books failed: %v\n", err)
	} else if getBooksResp.Success {
		fmt.Printf("✅ Books retrieved successfully!\n")
		fmt.Printf("   Total: %d books\n", getBooksResp.Total)
		for _, book := range getBooksResp.Books {
			fmt.Printf("   - ID: %d, Title: %s, Author: %s, Price: $%.2f\n",
				book.Id, book.Title, book.Author, book.Price)
		}
	}

	// Get Single Book
	fmt.Println("\n3.3 Getting Single Book...")
	getBookReq := &proto.GetBookRequest{
		Id: bookID,
	}

	getBookResp, err := bookClient.GetBook(ctx, getBookReq)
	if err != nil {
		fmt.Printf("❌ Get book failed: %v\n", err)
	} else if getBookResp.Success {
		fmt.Printf("✅ Book retrieved successfully!\n")
		fmt.Printf("   ID: %d, Title: %s, Author: %s, Price: $%.2f, Stock: %d\n",
			getBookResp.Book.Id, getBookResp.Book.Title, getBookResp.Book.Author,
			getBookResp.Book.Price, getBookResp.Book.Stock)
	}

	// Update Book
	fmt.Println("\n3.4 Updating Book...")
	updateBookReq := &proto.UpdateBookRequest{
		Id:         bookID,
		Title:      "The Hitchhiker's Guide to the Galaxy (Updated)",
		Author:     "Douglas Adams",
		Price:      18.99,
		Stock:      45,
		Year:       1979,
		CategoryId: categoryID,
		Token:      token,
	}

	updateBookResp, err := bookClient.UpdateBook(ctx, updateBookReq)
	if err != nil {
		fmt.Printf("❌ Update book failed: %v\n", err)
	} else if updateBookResp.Success {
		fmt.Printf("✅ Book updated successfully!\n")
		fmt.Printf("   ID: %d, Title: %s, Price: $%.2f, Stock: %d\n",
			updateBookResp.Book.Id, updateBookResp.Book.Title,
			updateBookResp.Book.Price, updateBookResp.Book.Stock)
	}

	// Get Books by Category
	fmt.Println("\n3.5 Getting Books by Category...")
	getBooksByCategoryReq := &proto.GetBooksByCategoryRequest{
		CategoryId: categoryID,
		Page:       1,
		Limit:      10,
	}

	getBooksByCategoryResp, err := bookClient.GetBooksByCategory(ctx, getBooksByCategoryReq)
	if err != nil {
		fmt.Printf("❌ Get books by category failed: %v\n", err)
	} else if getBooksByCategoryResp.Success {
		fmt.Printf("✅ Books by category retrieved successfully!\n")
		fmt.Printf("   Total: %d books in category\n", getBooksByCategoryResp.Total)
		for _, book := range getBooksByCategoryResp.Books {
			fmt.Printf("   - ID: %d, Title: %s, Author: %s\n",
				book.Id, book.Title, book.Author)
		}
	}

	// Test Delete Book
	fmt.Println("\n3.6 Deleting Book...")
	deleteBookReq := &proto.DeleteBookRequest{
		Id:    bookID,
		Token: token,
	}

	deleteBookResp, err := bookClient.DeleteBook(ctx, deleteBookReq)
	if err != nil {
		fmt.Printf("❌ Delete book failed: %v\n", err)
	} else if deleteBookResp.Success {
		fmt.Printf("✅ Book deleted successfully!\n")
	} else {
		fmt.Printf("❌ Delete book failed: %s\n", deleteBookResp.Message)
	}

	// Test Delete Category
	fmt.Println("\n3.7 Deleting Category...")
	deleteCategoryReq := &proto.DeleteCategoryRequest{
		Id:    categoryID,
		Token: token,
	}

	deleteCategoryResp, err := categoryClient.DeleteCategory(ctx, deleteCategoryReq)
	if err != nil {
		fmt.Printf("❌ Delete category failed: %v\n", err)
	} else if deleteCategoryResp.Success {
		fmt.Printf("✅ Category deleted successfully!\n")
	} else {
		fmt.Printf("❌ Delete category failed: %s\n", deleteCategoryResp.Message)
	}
}