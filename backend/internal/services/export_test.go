package services

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

func TestGenerateBidCSV(t *testing.T) {
	service := NewExportService()

	// Create test data
	bidID := uuid.New()
	projectID := uuid.New()
	totalCost := 100000.0
	laborCost := 60000.0
	materialCost := 40000.0
	markup := 20.0
	finalPrice := 120000.0
	bidName := "Test Bid"

	bid := &models.Bid{
		ID:               bidID,
		ProjectID:        projectID,
		Name:             &bidName,
		TotalCost:        &totalCost,
		LaborCost:        &laborCost,
		MaterialCost:     &materialCost,
		MarkupPercentage: &markup,
		FinalPrice:       &finalPrice,
		Status:           models.BidStatusDraft,
		Version:          1,
		IsLatest:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	bidResponse := &models.GenerateBidResponse{
		BidID:        bidID.String(),
		ProjectID:    projectID.String(),
		Status:       "draft",
		ScopeOfWork:  "Complete office renovation",
		LineItems: []models.LineItem{
			{
				Description: "Framing lumber",
				Trade:       "Framing",
				Quantity:    2500,
				Unit:        "BF",
				UnitCost:    2.50,
				Total:       6250,
			},
			{
				Description: "Drywall installation",
				Trade:       "Drywall",
				Quantity:    1200,
				Unit:        "SF",
				UnitCost:    1.75,
				Total:       2100,
			},
			{
				Description: "Electrical outlets",
				Trade:       "Electrical",
				Quantity:    25,
				Unit:        "EA",
				UnitCost:    125,
				Total:       3125,
			},
		},
		LaborCost:    60000,
		MaterialCost: 40000,
		Subtotal:     100000,
		MarkupAmount: 20000,
		TotalPrice:   120000,
		Inclusions: []string{
			"All materials",
			"Labor",
		},
		Exclusions: []string{
			"Furniture",
		},
		Schedule: map[string]string{
			"Phase 1": "1 week",
		},
		PaymentTerms:     "50% deposit",
		WarrantyTerms:    "1 year",
		ClosingStatement: "Thank you",
	}

	projectName := "Test Project"

	t.Run("generate valid CSV", func(t *testing.T) {
		csvBytes, err := service.GenerateBidCSV(bid, bidResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidCSV() error = %v", err)
			return
		}

		if len(csvBytes) == 0 {
			t.Error("GenerateBidCSV() returned empty CSV")
		}

		// Parse CSV to verify structure - use FieldsPerRecord = -1 for variable fields
		reader := csv.NewReader(strings.NewReader(string(csvBytes)))
		reader.FieldsPerRecord = -1 // Allow variable number of fields per record
		records, err := reader.ReadAll()
		if err != nil {
			t.Errorf("Failed to parse generated CSV: %v", err)
			return
		}

		if len(records) == 0 {
			t.Error("Generated CSV has no records")
		}

		// Check for key sections
		csvContent := string(csvBytes)
		expectedSections := []string{
			"Construction Bid Export",
			"Project",
			"Bid ID",
			"Line Items",
			"Trade Breakdown",
			"Cost Summary",
		}

		for _, section := range expectedSections {
			if !strings.Contains(csvContent, section) {
				t.Errorf("CSV missing expected section: %s", section)
			}
		}
	})

	t.Run("verify line items in CSV", func(t *testing.T) {
		csvBytes, err := service.GenerateBidCSV(bid, bidResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidCSV() error = %v", err)
			return
		}

		csvContent := string(csvBytes)

		// Check that line items are included
		for _, item := range bidResponse.LineItems {
			if !strings.Contains(csvContent, item.Description) {
				t.Errorf("CSV missing line item: %s", item.Description)
			}
			if !strings.Contains(csvContent, item.Trade) {
				t.Errorf("CSV missing trade: %s", item.Trade)
			}
		}
	})

	t.Run("verify cost summary in CSV", func(t *testing.T) {
		csvBytes, err := service.GenerateBidCSV(bid, bidResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidCSV() error = %v", err)
			return
		}

		csvContent := string(csvBytes)

		// Check cost values
		if !strings.Contains(csvContent, "60000.00") { // Labor cost
			t.Error("CSV missing labor cost")
		}
		if !strings.Contains(csvContent, "40000.00") { // Material cost
			t.Error("CSV missing material cost")
		}
		if !strings.Contains(csvContent, "120000.00") { // Total price
			t.Error("CSV missing total price")
		}
	})

	t.Run("generate CSV with empty line items", func(t *testing.T) {
		emptyResponse := &models.GenerateBidResponse{
			BidID:        bidID.String(),
			ProjectID:    projectID.String(),
			Status:       "draft",
			ScopeOfWork:  "Simple project",
			LineItems:    []models.LineItem{},
			LaborCost:    5000,
			MaterialCost: 3000,
			Subtotal:     8000,
			MarkupAmount: 1600,
			TotalPrice:   9600,
		}

		csvBytes, err := service.GenerateBidCSV(bid, emptyResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidCSV() with empty items error = %v", err)
			return
		}

		if len(csvBytes) == 0 {
			t.Error("GenerateBidCSV() with empty items returned empty CSV")
		}
	})
}

func TestGenerateBidExcel(t *testing.T) {
	service := NewExportService()

	bidID := uuid.New()
	projectID := uuid.New()
	totalCost := 100000.0
	laborCost := 60000.0
	materialCost := 40000.0
	markup := 20.0
	finalPrice := 120000.0
	bidName := "Test Bid"

	bid := &models.Bid{
		ID:               bidID,
		ProjectID:        projectID,
		Name:             &bidName,
		TotalCost:        &totalCost,
		LaborCost:        &laborCost,
		MaterialCost:     &materialCost,
		MarkupPercentage: &markup,
		FinalPrice:       &finalPrice,
		Status:           models.BidStatusDraft,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	bidResponse := &models.GenerateBidResponse{
		BidID:        bidID.String(),
		ProjectID:    projectID.String(),
		Status:       "draft",
		ScopeOfWork:  "Test scope",
		LineItems:    []models.LineItem{},
		LaborCost:    60000,
		MaterialCost: 40000,
		Subtotal:     100000,
		MarkupAmount: 20000,
		TotalPrice:   120000,
	}

	projectName := "Test Project"

	t.Run("generate Excel with UTF-8 BOM", func(t *testing.T) {
		excelBytes, err := service.GenerateBidExcel(bid, bidResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidExcel() error = %v", err)
			return
		}

		if len(excelBytes) < 3 {
			t.Error("GenerateBidExcel() returned data too short for BOM")
			return
		}

		// Check for UTF-8 BOM
		bom := []byte{0xEF, 0xBB, 0xBF}
		if excelBytes[0] != bom[0] || excelBytes[1] != bom[1] || excelBytes[2] != bom[2] {
			t.Error("Excel export missing UTF-8 BOM")
		}
	})

	t.Run("Excel content matches CSV", func(t *testing.T) {
		csvBytes, err := service.GenerateBidCSV(bid, bidResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidCSV() error = %v", err)
			return
		}

		excelBytes, err := service.GenerateBidExcel(bid, bidResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidExcel() error = %v", err)
			return
		}

		// Excel should be CSV + 3 bytes for BOM
		if len(excelBytes) != len(csvBytes)+3 {
			t.Errorf("Excel size mismatch. Expected %d, got %d", len(csvBytes)+3, len(excelBytes))
		}

		// Content after BOM should match CSV
		excelContent := string(excelBytes[3:])
		csvContent := string(csvBytes)
		if excelContent != csvContent {
			t.Error("Excel content (after BOM) doesn't match CSV content")
		}
	})
}

func TestGroupByTrade(t *testing.T) {
	service := NewExportService()

	lineItems := []models.LineItem{
		{Description: "Item 1", Trade: "Framing", Total: 1000},
		{Description: "Item 2", Trade: "Framing", Total: 2000},
		{Description: "Item 3", Trade: "Drywall", Total: 1500},
		{Description: "Item 4", Trade: "", Total: 500}, // Empty trade
		{Description: "Item 5", Trade: "Electrical", Total: 3000},
	}

	groups := service.groupByTrade(lineItems)

	t.Run("correct number of trade groups", func(t *testing.T) {
		// Should have Framing, Drywall, General (for empty), and Electrical
		if len(groups) != 4 {
			t.Errorf("Expected 4 trade groups, got %d", len(groups))
		}
	})

	t.Run("framing has 2 items", func(t *testing.T) {
		framingItems, ok := groups["Framing"]
		if !ok {
			t.Error("Framing group not found")
			return
		}
		if len(framingItems) != 2 {
			t.Errorf("Expected 2 framing items, got %d", len(framingItems))
		}
	})

	t.Run("empty trade becomes General", func(t *testing.T) {
		generalItems, ok := groups["General"]
		if !ok {
			t.Error("General group not found for empty trade")
			return
		}
		if len(generalItems) != 1 {
			t.Errorf("Expected 1 general item, got %d", len(generalItems))
		}
	})
}

func TestGenerateCSVFilename(t *testing.T) {
	service := NewExportService()
	projectID := uuid.New()
	bidID := uuid.New()

	filename := service.GenerateCSVFilename(projectID, bidID)

	if len(filename) == 0 {
		t.Error("GenerateCSVFilename() returned empty string")
	}

	// Check filename format
	expectedPrefix := "bids/" + projectID.String() + "/bid-"
	if len(filename) < len(expectedPrefix) || filename[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Filename doesn't start with expected prefix. Got: %s", filename)
	}

	// Check it ends with .csv
	if len(filename) < 4 || filename[len(filename)-4:] != ".csv" {
		t.Error("Filename doesn't end with .csv")
	}
}

func TestGenerateExcelFilename(t *testing.T) {
	service := NewExportService()
	projectID := uuid.New()
	bidID := uuid.New()

	filename := service.GenerateExcelFilename(projectID, bidID)

	if len(filename) == 0 {
		t.Error("GenerateExcelFilename() returned empty string")
	}

	// Check it ends with .xlsx
	if len(filename) < 5 || filename[len(filename)-5:] != ".xlsx" {
		t.Error("Filename doesn't end with .xlsx")
	}
}

func TestExportServiceParseBidDataFromJSON(t *testing.T) {
	service := NewExportService()

	t.Run("parse valid JSON", func(t *testing.T) {
		jsonData := `{
			"bid_id": "test-id",
			"project_id": "project-id",
			"status": "draft",
			"scope_of_work": "Test",
			"line_items": [],
			"labor_cost": 5000,
			"material_cost": 3000,
			"subtotal": 8000,
			"markup_amount": 1600,
			"total_price": 9600
		}`

		bidResponse, err := service.ParseBidDataFromJSON(jsonData)
		if err != nil {
			t.Errorf("ParseBidDataFromJSON() error = %v", err)
			return
		}

		if bidResponse.TotalPrice != 9600 {
			t.Errorf("Expected total price 9600, got %f", bidResponse.TotalPrice)
		}
	})

	t.Run("parse invalid JSON", func(t *testing.T) {
		jsonData := `{invalid}`

		_, err := service.ParseBidDataFromJSON(jsonData)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})
}
