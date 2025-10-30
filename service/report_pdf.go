package service

import (
	"bytes"
	"fmt"

	"github.com/unidoc/unipdf/v3/creator"
)

// generatePDFReport creates a PDF report with improved formatting
func (s *ReportService) generatePDFReport(data *ReportData, req ReportRequest) ([]byte, error) {
	c := creator.New()
	c.SetPageMargins(50, 50, 50, 50)

	// Title Page
	title := c.NewParagraph("Transaction Report")
	title.SetFontSize(24)
	title.SetColor(creator.ColorRGBFrom8bit(0, 51, 102))
	c.Draw(title)

	c.NewPage()

	// Report metadata
	meta := c.NewParagraph(fmt.Sprintf("Generated: %s", data.GeneratedAt.Format("2006-01-02 15:04:05")))
	meta.SetFontSize(10)
	c.Draw(meta)

	if data.StartDate != nil && data.EndDate != nil {
		dateRange := c.NewParagraph(fmt.Sprintf("Period: %s to %s", *data.StartDate, *data.EndDate))
		dateRange.SetFontSize(10)
		c.Draw(dateRange)
	}

	spacer := c.NewParagraph("\n")
	c.Draw(spacer)

	// Statistics Section
	statsTitle := c.NewParagraph("Store Statistics")
	statsTitle.SetFontSize(18)
	statsTitle.SetColor(creator.ColorRGBFrom8bit(0, 51, 102))
	c.Draw(statsTitle)

	spacer2 := c.NewParagraph("\n")
	c.Draw(spacer2)

	// Create statistics table
	statsTable := c.NewTable(2)
	statsTable.SetColumnWidths(0.5, 0.5)

	addTableCell := func(table *creator.Table, text string, isHeader bool) {
		p := c.NewParagraph(text)
		if isHeader {
			p.SetFontSize(12)
			p.SetColor(creator.ColorRGBFrom8bit(255, 255, 255))
		} else {
			p.SetFontSize(11)
		}
		cell := table.NewCell()
		if isHeader {
			cell.SetBackgroundColor(creator.ColorRGBFrom8bit(0, 51, 102))
		} else {
			cell.SetBackgroundColor(creator.ColorRGBFrom8bit(240, 240, 240))
		}
		cell.SetBorder(creator.CellBorderSideAll, creator.CellBorderStyleSingle, 1)
		cell.SetContent(p)
	}

	addTableCell(statsTable, "Metric", true)
	addTableCell(statsTable, "Value", true)

	addTableCell(statsTable, "Total Orders", false)
	addTableCell(statsTable, fmt.Sprintf("%d", data.TotalOrders), false)

	addTableCell(statsTable, "Total Products", false)
	addTableCell(statsTable, fmt.Sprintf("%d", data.TotalProducts), false)

	addTableCell(statsTable, "Total Users", false)
	addTableCell(statsTable, fmt.Sprintf("%d", data.TotalUsers), false)

	addTableCell(statsTable, "Total Revenue", false)
	addTableCell(statsTable, fmt.Sprintf("$%.2f", data.TotalRevenue), false)

	c.Draw(statsTable)

	// Orders by Status
	if len(data.OrdersByStatus) > 0 {
		c.NewPage()

		statusTitle := c.NewParagraph("Orders by Status")
		statusTitle.SetFontSize(18)
		statusTitle.SetColor(creator.ColorRGBFrom8bit(0, 51, 102))
		c.Draw(statusTitle)

		spacer3 := c.NewParagraph("\n")
		c.Draw(spacer3)

		statusTable := c.NewTable(2)
		statusTable.SetColumnWidths(0.5, 0.5)

		addTableCell(statusTable, "Status", true)
		addTableCell(statusTable, "Count", true)

		for status, count := range data.OrdersByStatus {
			addTableCell(statusTable, status, false)
			addTableCell(statusTable, fmt.Sprintf("%d", count), false)
		}

		c.Draw(statusTable)
	}

	// Recent Orders Table
	if len(data.RecentOrders) > 0 {
		c.NewPage()

		ordersTitle := c.NewParagraph("Recent Orders")
		ordersTitle.SetFontSize(18)
		ordersTitle.SetColor(creator.ColorRGBFrom8bit(0, 51, 102))
		c.Draw(ordersTitle)

		spacer4 := c.NewParagraph("\n")
		c.Draw(spacer4)

		ordersTable := c.NewTable(5)
		ordersTable.SetColumnWidths(0.15, 0.25, 0.2, 0.2, 0.2)

		addTableCell(ordersTable, "Order ID", true)
		addTableCell(ordersTable, "Customer", true)
		addTableCell(ordersTable, "Status", true)
		addTableCell(ordersTable, "Amount", true)
		addTableCell(ordersTable, "Date", true)

		for _, order := range data.RecentOrders {
			addTableCell(ordersTable, fmt.Sprintf("%d", order.OrderID), false)
			addTableCell(ordersTable, order.Username, false)
			addTableCell(ordersTable, order.Status, false)
			addTableCell(ordersTable, fmt.Sprintf("$%.2f", order.TotalPrice), false)
			addTableCell(ordersTable, order.CreatedAt, false)
		}

		c.Draw(ordersTable)
	}

	// Top Products
	if len(data.TopProducts) > 0 {
		c.NewPage()

		productsTitle := c.NewParagraph("Top Selling Products")
		productsTitle.SetFontSize(18)
		productsTitle.SetColor(creator.ColorRGBFrom8bit(0, 51, 102))
		c.Draw(productsTitle)

		spacer5 := c.NewParagraph("\n")
		c.Draw(spacer5)

		productsTable := c.NewTable(4)
		productsTable.SetColumnWidths(0.15, 0.35, 0.25, 0.25)

		addTableCell(productsTable, "Product ID", true)
		addTableCell(productsTable, "Product Name", true)
		addTableCell(productsTable, "Units Sold", true)
		addTableCell(productsTable, "Revenue", true)

		for _, product := range data.TopProducts {
			addTableCell(productsTable, fmt.Sprintf("%d", product.ProductID), false)
			addTableCell(productsTable, product.ProductName, false)
			addTableCell(productsTable, fmt.Sprintf("%d", product.TotalSold), false)
			addTableCell(productsTable, fmt.Sprintf("$%.2f", product.TotalRevenue), false)
		}

		c.Draw(productsTable)
	}

	// Top Customers
	if len(data.TopCustomers) > 0 {
		c.NewPage()

		customersTitle := c.NewParagraph("Top Customers")
		customersTitle.SetFontSize(18)
		customersTitle.SetColor(creator.ColorRGBFrom8bit(0, 51, 102))
		c.Draw(customersTitle)

		spacer6 := c.NewParagraph("\n")
		c.Draw(spacer6)

		customersTable := c.NewTable(4)
		customersTable.SetColumnWidths(0.3, 0.3, 0.2, 0.2)

		addTableCell(customersTable, "Username", true)
		addTableCell(customersTable, "Email", true)
		addTableCell(customersTable, "Orders", true)
		addTableCell(customersTable, "Total Spent", true)

		for _, customer := range data.TopCustomers {
			addTableCell(customersTable, customer.Username, false)
			addTableCell(customersTable, customer.Email, false)
			addTableCell(customersTable, fmt.Sprintf("%d", customer.OrderCount), false)
			addTableCell(customersTable, fmt.Sprintf("$%.2f", customer.TotalSpent), false)
		}

		c.Draw(customersTable)
	}

	// Write to buffer
	var buf bytes.Buffer
	err := c.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf.Bytes(), nil
}
