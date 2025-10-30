package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
)

// generateCSVReport creates a CSV report
func (s *ReportService) generateCSVReport(data *ReportData) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header information
	writer.Write([]string{"Medical Equipment Store - Transaction Report"})
	writer.Write([]string{"Generated", data.GeneratedAt.Format("2006-01-02 15:04:05")})

	if data.StartDate != nil && data.EndDate != nil {
		writer.Write([]string{"Period", fmt.Sprintf("%s to %s", *data.StartDate, *data.EndDate)})
	}
	writer.Write([]string{}) // Empty line

	// Statistics Section
	writer.Write([]string{"Store Statistics"})
	writer.Write([]string{"Metric", "Value"})
	writer.Write([]string{"Total Orders", strconv.FormatInt(data.TotalOrders, 10)})
	writer.Write([]string{"Total Products", strconv.FormatInt(data.TotalProducts, 10)})
	writer.Write([]string{"Total Users", strconv.FormatInt(data.TotalUsers, 10)})
	writer.Write([]string{"Total Revenue", fmt.Sprintf("%.2f", data.TotalRevenue)})
	writer.Write([]string{}) // Empty line

	// Orders by Status
	if len(data.OrdersByStatus) > 0 {
		writer.Write([]string{"Orders by Status"})
		writer.Write([]string{"Status", "Count"})
		for status, count := range data.OrdersByStatus {
			writer.Write([]string{status, strconv.FormatInt(count, 10)})
		}
		writer.Write([]string{}) // Empty line
	}

	// Recent Orders
	if len(data.RecentOrders) > 0 {
		writer.Write([]string{"Recent Orders"})
		writer.Write([]string{"Order ID", "User ID", "Username", "Status", "Total Price", "Payment Method", "Created At"})
		for _, order := range data.RecentOrders {
			writer.Write([]string{
				strconv.FormatUint(uint64(order.OrderID), 10),
				strconv.FormatUint(uint64(order.UserID), 10),
				order.Username,
				order.Status,
				fmt.Sprintf("%.2f", order.TotalPrice),
				order.PaymentMethod,
				order.CreatedAt,
			})
		}
		writer.Write([]string{}) // Empty line
	}

	// Top Products
	if len(data.TopProducts) > 0 {
		writer.Write([]string{"Top Selling Products"})
		writer.Write([]string{"Product ID", "Product Name", "Total Sold", "Total Revenue"})
		for _, product := range data.TopProducts {
			writer.Write([]string{
				strconv.FormatUint(uint64(product.ProductID), 10),
				product.ProductName,
				strconv.FormatInt(product.TotalSold, 10),
				fmt.Sprintf("%.2f", product.TotalRevenue),
			})
		}
		writer.Write([]string{}) // Empty line
	}

	// Top Customers
	if len(data.TopCustomers) > 0 {
		writer.Write([]string{"Top Customers"})
		writer.Write([]string{"User ID", "Username", "Email", "Order Count", "Total Spent"})
		for _, customer := range data.TopCustomers {
			writer.Write([]string{
				strconv.FormatUint(uint64(customer.UserID), 10),
				customer.Username,
				customer.Email,
				strconv.FormatInt(customer.OrderCount, 10),
				fmt.Sprintf("%.2f", customer.TotalSpent),
			})
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to write CSV: %w", err)
	}

	return buf.Bytes(), nil
}
