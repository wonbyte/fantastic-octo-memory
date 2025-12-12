package services

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

// ExportService handles exporting bid data to various formats
type ExportService struct{}

func NewExportService() *ExportService {
	return &ExportService{}
}

// ExportFormat represents the export file format
type ExportFormat string

const (
	ExportFormatCSV   ExportFormat = "csv"
	ExportFormatExcel ExportFormat = "xlsx"
)

// GenerateBidCSV exports bid data to CSV format
func (s *ExportService) GenerateBidCSV(bid *models.Bid, bidResponse *models.GenerateBidResponse, projectName string) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header section
	writer.Write([]string{"Construction Bid Export - CSV Format"})
	writer.Write([]string{"Project", projectName})
	writer.Write([]string{"Bid ID", bid.ID.String()})
	writer.Write([]string{"Date", time.Now().Format("2006-01-02")})
	writer.Write([]string{"Status", string(bid.Status)})
	writer.Write([]string{}) // Empty row

	// Scope of Work
	if bidResponse.ScopeOfWork != "" {
		writer.Write([]string{"Scope of Work"})
		writer.Write([]string{bidResponse.ScopeOfWork})
		writer.Write([]string{}) // Empty row
	}

	// Line Items
	if len(bidResponse.LineItems) > 0 {
		writer.Write([]string{"Line Items"})
		writer.Write([]string{"Description", "Trade", "Quantity", "Unit", "Unit Cost", "Total"})
		
		for _, item := range bidResponse.LineItems {
			writer.Write([]string{
				item.Description,
				item.Trade,
				fmt.Sprintf("%.2f", item.Quantity),
				item.Unit,
				fmt.Sprintf("%.2f", item.UnitCost),
				fmt.Sprintf("%.2f", item.Total),
			})
		}
		writer.Write([]string{}) // Empty row
	}

	// Trade Breakdown
	if len(bidResponse.LineItems) > 0 {
		writer.Write([]string{"Trade Breakdown"})
		writer.Write([]string{"Trade", "Item Count", "Total Cost"})
		
		tradeGroups := s.groupByTrade(bidResponse.LineItems)
		for trade, items := range tradeGroups {
			total := 0.0
			for _, item := range items {
				total += item.Total
			}
			writer.Write([]string{
				trade,
				strconv.Itoa(len(items)),
				fmt.Sprintf("%.2f", total),
			})
		}
		writer.Write([]string{}) // Empty row
	}

	// Cost Summary
	writer.Write([]string{"Cost Summary"})
	writer.Write([]string{"Material Cost", fmt.Sprintf("%.2f", bidResponse.MaterialCost)})
	writer.Write([]string{"Labor Cost", fmt.Sprintf("%.2f", bidResponse.LaborCost)})
	writer.Write([]string{"Subtotal", fmt.Sprintf("%.2f", bidResponse.Subtotal)})
	writer.Write([]string{"Markup Amount", fmt.Sprintf("%.2f", bidResponse.MarkupAmount)})
	writer.Write([]string{"Total Price", fmt.Sprintf("%.2f", bidResponse.TotalPrice)})
	writer.Write([]string{}) // Empty row

	// Inclusions
	if len(bidResponse.Inclusions) > 0 {
		writer.Write([]string{"Inclusions"})
		for _, inclusion := range bidResponse.Inclusions {
			writer.Write([]string{inclusion})
		}
		writer.Write([]string{}) // Empty row
	}

	// Exclusions
	if len(bidResponse.Exclusions) > 0 {
		writer.Write([]string{"Exclusions"})
		for _, exclusion := range bidResponse.Exclusions {
			writer.Write([]string{exclusion})
		}
		writer.Write([]string{}) // Empty row
	}

	// Schedule
	if len(bidResponse.Schedule) > 0 {
		writer.Write([]string{"Project Schedule"})
		writer.Write([]string{"Phase", "Timeline"})
		for phase, timeline := range bidResponse.Schedule {
			writer.Write([]string{phase, timeline})
		}
		writer.Write([]string{}) // Empty row
	}

	// Payment Terms
	if bidResponse.PaymentTerms != "" {
		writer.Write([]string{"Payment Terms"})
		writer.Write([]string{bidResponse.PaymentTerms})
		writer.Write([]string{}) // Empty row
	}

	// Warranty Terms
	if bidResponse.WarrantyTerms != "" {
		writer.Write([]string{"Warranty Terms"})
		writer.Write([]string{bidResponse.WarrantyTerms})
		writer.Write([]string{}) // Empty row
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to write CSV: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateBidExcel exports bid data to Excel-compatible CSV format (with UTF-8 BOM)
// Note: This generates a CSV that Excel can open properly. For true .xlsx format,
// we would need to add the excelize library. This approach keeps dependencies minimal
// while maintaining Excel compatibility.
func (s *ExportService) GenerateBidExcel(bid *models.Bid, bidResponse *models.GenerateBidResponse, projectName string) ([]byte, error) {
	csvData, err := s.GenerateBidCSV(bid, bidResponse, projectName)
	if err != nil {
		return nil, err
	}

	// Add UTF-8 BOM for Excel compatibility
	bom := []byte{0xEF, 0xBB, 0xBF}
	excelData := append(bom, csvData...)
	
	return excelData, nil
}

// groupByTrade groups line items by their trade
func (s *ExportService) groupByTrade(items []models.LineItem) map[string][]models.LineItem {
	groups := make(map[string][]models.LineItem)
	for _, item := range items {
		trade := item.Trade
		if trade == "" {
			trade = "General"
		}
		groups[trade] = append(groups[trade], item)
	}
	return groups
}

// GenerateCSVFilename creates a unique filename for the bid CSV
func (s *ExportService) GenerateCSVFilename(projectID uuid.UUID, bidID uuid.UUID) string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("bids/%s/bid-%s-%s.csv", projectID.String(), bidID.String()[:8], timestamp)
}

// GenerateExcelFilename creates a unique filename for the bid Excel file
func (s *ExportService) GenerateExcelFilename(projectID uuid.UUID, bidID uuid.UUID) string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("bids/%s/bid-%s-%s.xlsx", projectID.String(), bidID.String()[:8], timestamp)
}

// ParseBidDataFromJSON parses bid_data JSONB field into GenerateBidResponse
func (s *ExportService) ParseBidDataFromJSON(bidData string) (*models.GenerateBidResponse, error) {
	var bidResponse models.GenerateBidResponse
	if err := json.Unmarshal([]byte(bidData), &bidResponse); err != nil {
		return nil, fmt.Errorf("failed to parse bid data: %w", err)
	}
	return &bidResponse, nil
}
