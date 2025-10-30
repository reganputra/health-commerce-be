package service

import (
	"fmt"
	"health-store/repositories"
	"log"
	"time"
)

// ReportService handles business logic for reports
type ReportService struct {
	orderRepo   repositories.OrderRepositoryInterface
	productRepo repositories.ProductRepositoryInterface
	userRepo    repositories.UserRepositoryInterface
}

// NewReportService creates a new report service
func NewReportService(
	orderRepo repositories.OrderRepositoryInterface,
	productRepo repositories.ProductRepositoryInterface,
	userRepo repositories.UserRepositoryInterface,
) *ReportService {
	return &ReportService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

// ReportRequest defines the parameters for report generation
type ReportRequest struct {
	ReportType      string   // "summary", "detailed", "financial"
	Format          string   // "pdf", "csv"
	StartDate       *string  // Optional date range filter
	EndDate         *string  // Optional date range filter
	IncludeSections []string // ["statistics", "orders", "products", "revenue", "customers"]
	Limit           int      // Limit for items (default 10)
}

// ReportData holds all the data needed for report generation
type ReportData struct {
	TotalOrders    int64
	TotalProducts  int64
	TotalUsers     int64
	TotalRevenue   float64
	OrdersByStatus map[string]int64
	RecentOrders   []OrderSummary
	TopProducts    []ProductSummary
	TopCustomers   []CustomerSummary
	GeneratedAt    time.Time
	StartDate      *string
	EndDate        *string
}

// OrderSummary represents a summary of an order for reports
type OrderSummary struct {
	OrderID       uint
	UserID        uint
	Username      string
	Status        string
	TotalPrice    float64
	PaymentMethod string
	CreatedAt     string
}

// ProductSummary represents a summary of a product for reports
type ProductSummary struct {
	ProductID    uint
	ProductName  string
	TotalSold    int64
	TotalRevenue float64
}

// CustomerSummary represents a summary of a customer for reports
type CustomerSummary struct {
	UserID     uint
	Username   string
	Email      string
	OrderCount int64
	TotalSpent float64
}

// GenerateReport generates a report based on the request parameters
func (s *ReportService) GenerateReport(req ReportRequest) ([]byte, error) {
	// Set defaults
	if req.Limit == 0 {
		req.Limit = 10
	}

	// Gather report data
	data, err := s.gatherReportData(req)
	if err != nil {
		log.Printf("Error gathering report data: %v", err)
		return nil, fmt.Errorf("failed to gather report data: %w", err)
	}

	// Generate report based on format
	switch req.Format {
	case "csv":
		return s.generateCSVReport(data)
	case "pdf":
		fallthrough
	default:
		return s.generatePDFReport(data, req)
	}
}

// GenerateTransactionReport generates a PDF transaction report (backward compatibility)
func (s *ReportService) GenerateTransactionReport() ([]byte, error) {
	req := ReportRequest{
		ReportType:      "summary",
		Format:          "pdf",
		Limit:           10,
		IncludeSections: []string{"statistics", "orders"},
	}
	return s.GenerateReport(req)
}

// gatherReportData collects all necessary data for the report
func (s *ReportService) gatherReportData(req ReportRequest) (*ReportData, error) {
	data := &ReportData{
		GeneratedAt: time.Now(),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	var err error

	// Get basic statistics
	data.TotalOrders, err = s.orderRepo.GetOrderStatistics()
	if err != nil {
		log.Printf("Warning: Failed to get order statistics: %v", err)
	}

	data.TotalProducts, err = s.productRepo.GetProductCount()
	if err != nil {
		log.Printf("Warning: Failed to get product count: %v", err)
	}

	data.TotalUsers, err = s.userRepo.GetUserCount()
	if err != nil {
		log.Printf("Warning: Failed to get user count: %v", err)
	}

	// Get revenue
	if req.StartDate != nil && req.EndDate != nil {
		data.TotalRevenue, err = s.orderRepo.GetRevenueByDateRange(*req.StartDate, *req.EndDate)
	} else {
		data.TotalRevenue, err = s.orderRepo.GetTotalRevenue()
	}
	if err != nil {
		log.Printf("Warning: Failed to get revenue: %v", err)
	}

	// Get orders by status
	data.OrdersByStatus, err = s.orderRepo.GetOrdersByStatus()
	if err != nil {
		log.Printf("Warning: Failed to get orders by status: %v", err)
	}

	// Get recent orders
	var orders []OrderSummary
	if req.StartDate != nil && req.EndDate != nil {
		dbOrders, err := s.orderRepo.GetOrdersByDateRange(*req.StartDate, *req.EndDate)
		if err != nil {
			log.Printf("Warning: Failed to get orders by date range: %v", err)
		} else {
			for _, order := range dbOrders {
				orders = append(orders, OrderSummary{
					OrderID:       order.ID,
					UserID:        order.UserID,
					Username:      order.User.Username,
					Status:        order.Status,
					TotalPrice:    order.TotalPrice,
					PaymentMethod: order.PaymentMethod,
					CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		}
	} else {
		dbOrders, err := s.orderRepo.GetRecentOrders(req.Limit)
		if err != nil {
			log.Printf("Warning: Failed to get recent orders: %v", err)
		} else {
			for _, order := range dbOrders {
				orders = append(orders, OrderSummary{
					OrderID:       order.ID,
					UserID:        order.UserID,
					Username:      order.User.Username,
					Status:        order.Status,
					TotalPrice:    order.TotalPrice,
					PaymentMethod: order.PaymentMethod,
					CreatedAt:     order.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		}
	}
	data.RecentOrders = orders

	// Get top products
	topProducts, err := s.productRepo.GetTopSellingProducts(req.Limit)
	if err != nil {
		log.Printf("Warning: Failed to get top products: %v", err)
	} else {
		for _, product := range topProducts {
			data.TopProducts = append(data.TopProducts, ProductSummary{
				ProductID:    product.ProductID,
				ProductName:  product.ProductName,
				TotalSold:    product.TotalSold,
				TotalRevenue: product.TotalRevenue,
			})
		}
	}

	// Get top customers
	topCustomers, err := s.orderRepo.GetTopCustomers(req.Limit)
	if err != nil {
		log.Printf("Warning: Failed to get top customers: %v", err)
	} else {
		for _, customer := range topCustomers {
			data.TopCustomers = append(data.TopCustomers, CustomerSummary{
				UserID:     customer.UserID,
				Username:   customer.Username,
				Email:      customer.Email,
				OrderCount: customer.OrderCount,
				TotalSpent: customer.TotalSpent,
			})
		}
	}

	return data, nil
}
