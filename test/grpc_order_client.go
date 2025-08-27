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
	orderClient := proto.NewOrderServiceClient(conn)
	reportClient := proto.NewReportServiceClient(conn)

	fmt.Println("=== Testing Order and Payment System ===")

	// Step 1: Register admin user first
	fmt.Println("\n1. Registering admin user...")
	registerAdminResp, err := userClient.Register(context.Background(), &proto.RegisterRequest{
		Name:     "Admin User",
		Email:    "admin@bookstore.com",
		Password: "admin123",
	})
	if err != nil {
		fmt.Printf("Admin registration failed: %v\n", err)
	} else if !registerAdminResp.Success {
		fmt.Printf("Admin registration failed (user may already exist): %s\n", registerAdminResp.Message)
	} else {
		fmt.Printf("Admin registered: %s\n", registerAdminResp.User.Email)
	}

	// Step 2: Login as admin
	fmt.Println("\n2. Login as admin...")
	loginResp, err := userClient.Login(context.Background(), &proto.LoginRequest{
		Email:    "admin@bookstore.com",
		Password: "admin123",
	})
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	if !loginResp.Success {
		fmt.Printf("Login failed: %s\n", loginResp.Message)
		return
	}
	fmt.Printf("Login successful: %s\n", loginResp.Message)
	adminToken := loginResp.Token

	// Step 3: Create a category for testing
	fmt.Println("\n3. Creating test category...")
	categoryName := fmt.Sprintf("TestCategory_%d", time.Now().Unix())
	categoryResp, err := categoryClient.CreateCategory(context.Background(), &proto.CreateCategoryRequest{
		Name:  categoryName,
		Token: adminToken,
	})
	if err != nil {
		fmt.Printf("Category creation failed: %v\n", err)
		return
	}
	if !categoryResp.Success {
		fmt.Printf("Category creation failed: %s\n", categoryResp.Message)
		return
	}
	fmt.Printf("Category created: %s (ID: %d)\n", categoryResp.Category.Name, categoryResp.Category.Id)
	categoryID := categoryResp.Category.Id

	// Step 4: Create test books
	fmt.Println("\n4. Creating test books...")
	books := []struct {
		title  string
		author string
		price  float64
		stock  int32
		year   int32
	}{
		{"Test Book 1", "Author 1", 25.99, 10, 2023},
		{"Test Book 2", "Author 2", 35.50, 15, 2023},
		{"Test Book 3", "Author 3", 19.99, 5, 2023},
	}

	var bookIDs []uint32
	for _, book := range books {
		bookResp, err := bookClient.CreateBook(context.Background(), &proto.CreateBookRequest{
			Title:      book.title,
			Author:     book.author,
			Price:      book.price,
			Stock:      book.stock,
			Year:       book.year,
			CategoryId: categoryID,
			Token:      adminToken,
		})
		if err != nil {
			fmt.Printf("Book creation failed: %v\n", err)
			continue
		}
		if !bookResp.Success {
			fmt.Printf("Book creation failed: %s\n", bookResp.Message)
			continue
		}
		fmt.Printf("Book created: %s (ID: %d, Price: $%.2f, Stock: %d)\n", 
			bookResp.Book.Title, bookResp.Book.Id, bookResp.Book.Price, bookResp.Book.Stock)
		bookIDs = append(bookIDs, bookResp.Book.Id)
	}

	if len(bookIDs) == 0 {
		fmt.Println("No books created, cannot test orders")
		return
	}

	// Step 5: Register a regular user for ordering
	fmt.Println("\n5. Registering regular user...")
	userEmail := fmt.Sprintf("user_%d@test.com", time.Now().Unix())
	registerResp, err := userClient.Register(context.Background(), &proto.RegisterRequest{
		Name:     "Test User",
		Email:    userEmail,
		Password: "password123",
	})
	if err != nil {
		fmt.Printf("User registration failed: %v\n", err)
		return
	}
	if !registerResp.Success {
		fmt.Printf("User registration failed: %s\n", registerResp.Message)
		return
	}
	fmt.Printf("User registered: %s\n", registerResp.User.Email)

	// Step 6: Login as regular user
	fmt.Println("\n6. Login as regular user...")
	userLoginResp, err := userClient.Login(context.Background(), &proto.LoginRequest{
		Email:    userEmail,
		Password: "password123",
	})
	if err != nil {
		fmt.Printf("User login failed: %v\n", err)
		return
	}
	if !userLoginResp.Success {
		fmt.Printf("User login failed: %s\n", userLoginResp.Message)
		return
	}
	fmt.Printf("User login successful: %s\n", userLoginResp.Message)
	userToken := userLoginResp.Token

	// Step 7: Create an order
	fmt.Println("\n7. Creating order...")
	orderItems := []*proto.OrderItemRequest{
		{BookId: bookIDs[0], Quantity: 2},
		{BookId: bookIDs[1], Quantity: 1},
	}

	orderResp, err := orderClient.CreateOrder(context.Background(), &proto.CreateOrderRequest{
		Items: orderItems,
		Token: userToken,
	})
	if err != nil {
		fmt.Printf("Order creation failed: %v\n", err)
		return
	}
	if !orderResp.Success {
		fmt.Printf("Order creation failed: %s\n", orderResp.Message)
		return
	}
	fmt.Printf("Order created successfully!\n")
	fmt.Printf("Order ID: %d\n", orderResp.Order.Id)
	fmt.Printf("Total Price: $%.2f\n", orderResp.Order.TotalPrice)
	fmt.Printf("Status: %s\n", orderResp.Order.Status)
	fmt.Printf("Items count: %d\n", len(orderResp.Order.Items))

	orderID := orderResp.Order.Id

	// Step 8: Get order details
	fmt.Println("\n8. Getting order details...")
	getOrderResp, err := orderClient.GetOrder(context.Background(), &proto.GetOrderRequest{
		Id:    orderID,
		Token: userToken,
	})
	if err != nil {
		fmt.Printf("Get order failed: %v\n", err)
	} else if !getOrderResp.Success {
		fmt.Printf("Get order failed: %s\n", getOrderResp.Message)
	} else {
		fmt.Printf("Order details retrieved successfully\n")
		fmt.Printf("Order ID: %d, Total: $%.2f, Status: %s\n", 
			getOrderResp.Order.Id, getOrderResp.Order.TotalPrice, getOrderResp.Order.Status)
		for i, item := range getOrderResp.Order.Items {
			fmt.Printf("  Item %d: %s x%d @ $%.2f\n", 
				i+1, item.Book.Title, item.Quantity, item.Price)
		}
	}

	// Step 9: Get user orders
	fmt.Println("\n9. Getting user orders...")
	getOrdersResp, err := orderClient.GetOrders(context.Background(), &proto.GetOrdersRequest{
		Token: userToken,
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		fmt.Printf("Get orders failed: %v\n", err)
	} else if !getOrdersResp.Success {
		fmt.Printf("Get orders failed: %s\n", getOrdersResp.Message)
	} else {
		fmt.Printf("User orders retrieved: %d orders found\n", len(getOrdersResp.Orders))
		for i, order := range getOrdersResp.Orders {
			fmt.Printf("  Order %d: ID=%d, Total=$%.2f, Status=%s\n", 
				i+1, order.Id, order.TotalPrice, order.Status)
		}
	}

	// Step 10: Process payment
	fmt.Println("\n10. Processing payment...")
	paymentResp, err := orderClient.ProcessPayment(context.Background(), &proto.ProcessPaymentRequest{
		OrderId: orderID,
		Token:   userToken,
	})
	if err != nil {
		fmt.Printf("Payment processing failed: %v\n", err)
	} else if !paymentResp.Success {
		fmt.Printf("Payment processing failed: %s\n", paymentResp.Message)
	} else {
		fmt.Printf("Payment processed successfully!\n")
		fmt.Printf("Payment URL: %s\n", paymentResp.PaymentUrl)
	}

	// Step 11: Update order status (admin only)
	fmt.Println("\n11. Updating order status (admin)...")
	updateStatusResp, err := orderClient.UpdateOrderStatus(context.Background(), &proto.UpdateOrderStatusRequest{
		Id:     orderID,
		Status: "completed",
		Token:  adminToken,
	})
	if err != nil {
		fmt.Printf("Update order status failed: %v\n", err)
	} else if !updateStatusResp.Success {
		fmt.Printf("Update order status failed: %s\n", updateStatusResp.Message)
	} else {
		fmt.Printf("Order status updated successfully!\n")
		fmt.Printf("New status: %s\n", updateStatusResp.Order.Status)
	}

	// Step 12: Test reporting (admin only)
	fmt.Println("\n12. Testing sales report (admin)...")
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	salesReportResp, err := reportClient.GetSalesReport(context.Background(), &proto.GetSalesReportRequest{
		StartDate: startDate,
		EndDate:   endDate,
		Token:     adminToken,
	})
	if err != nil {
		fmt.Printf("Sales report failed: %v\n", err)
	} else if !salesReportResp.Success {
		fmt.Printf("Sales report failed: %s\n", salesReportResp.Message)
	} else {
		fmt.Printf("Sales report retrieved successfully!\n")
		fmt.Printf("Total Revenue: $%.2f\n", salesReportResp.TotalRevenue)
		fmt.Printf("Report items: %d\n", len(salesReportResp.Report))
	}

	// Step 13: Test top books report (admin only)
	fmt.Println("\n13. Testing top books report (admin)...")
	topBooksResp, err := reportClient.GetTopBooks(context.Background(), &proto.GetTopBooksRequest{
		Limit: 5,
		Token: adminToken,
	})
	if err != nil {
		fmt.Printf("Top books report failed: %v\n", err)
	} else if !topBooksResp.Success {
		fmt.Printf("Top books report failed: %s\n", topBooksResp.Message)
	} else {
		fmt.Printf("Top books report retrieved successfully!\n")
		fmt.Printf("Top books count: %d\n", len(topBooksResp.Books))
		for i, item := range topBooksResp.Books {
			fmt.Printf("  %d. %s - %d sold\n", i+1, item.Book.Title, item.TotalSold)
		}
	}

	fmt.Println("\n=== Order and Payment System Test Completed ===")
}