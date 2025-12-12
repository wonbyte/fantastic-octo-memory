package services

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

func TestGenerateBidPDF(t *testing.T) {
	service := NewPDFService()

	// Create test bid
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

	// Create test bid response
	bidResponse := &models.GenerateBidResponse{
		BidID:        bidID.String(),
		ProjectID:    projectID.String(),
		Status:       "draft",
		ScopeOfWork:  "Complete office renovation including framing, drywall, electrical, and plumbing.",
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
			"All materials specified",
			"Labor for installation",
			"Job site cleanup",
		},
		Exclusions: []string{
			"Furniture and equipment",
			"IT infrastructure",
		},
		Schedule: map[string]string{
			"Demolition":     "1 week",
			"Framing":        "2 weeks",
			"Finish work":    "3 weeks",
		},
		PaymentTerms:     "50% deposit, 50% on completion",
		WarrantyTerms:    "1-year workmanship warranty",
		ClosingStatement: "Thank you for considering our proposal.",
	}

	projectName := "Downtown Office Renovation"

	t.Run("generate basic PDF without branding", func(t *testing.T) {
		pdfBytes, err := service.GenerateBidPDF(bid, bidResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidPDF() error = %v", err)
			return
		}

		if len(pdfBytes) == 0 {
			t.Error("GenerateBidPDF() returned empty PDF")
		}

		// Check PDF magic number
		if len(pdfBytes) < 4 || string(pdfBytes[:4]) != "%PDF" {
			t.Error("Generated file does not appear to be a valid PDF")
		}
	})

	t.Run("generate PDF with company info", func(t *testing.T) {
		companyAddress := "123 Main St, City, ST 12345"
		companyPhone := "(555) 123-4567"
		companyEmail := "info@example.com"
		companyWebsite := "www.example.com"
		license := "CA-123456"
		insurance := "Fully insured"

		options := &PDFOptions{
			CompanyInfo: &models.CompanyInfo{
				Name:          "Quality Construction Co.",
				Address:       &companyAddress,
				Phone:         &companyPhone,
				Email:         &companyEmail,
				Website:       &companyWebsite,
				LicenseNumber: &license,
				InsuranceInfo: &insurance,
			},
			IncludeCover: true,
			IncludeLogo:  false,
		}

		pdfBytes, err := service.GenerateBidPDFWithOptions(bid, bidResponse, projectName, options)
		if err != nil {
			t.Errorf("GenerateBidPDFWithOptions() error = %v", err)
			return
		}

		if len(pdfBytes) == 0 {
			t.Error("GenerateBidPDFWithOptions() returned empty PDF")
		}

		// Check PDF magic number
		if len(pdfBytes) < 4 || string(pdfBytes[:4]) != "%PDF" {
			t.Error("Generated file does not appear to be a valid PDF")
		}
	})

	t.Run("generate PDF with empty line items", func(t *testing.T) {
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

		pdfBytes, err := service.GenerateBidPDF(bid, emptyResponse, projectName)
		if err != nil {
			t.Errorf("GenerateBidPDF() with empty items error = %v", err)
			return
		}

		if len(pdfBytes) == 0 {
			t.Error("GenerateBidPDF() with empty items returned empty PDF")
		}
	})
}

func TestParseBidDataFromJSON(t *testing.T) {
	service := NewPDFService()

	t.Run("parse valid JSON", func(t *testing.T) {
		jsonData := `{
			"bid_id": "test-id",
			"project_id": "project-id",
			"status": "draft",
			"scope_of_work": "Test scope",
			"line_items": [
				{
					"description": "Test item",
					"trade": "General",
					"quantity": 10,
					"unit": "EA",
					"unit_cost": 100,
					"total": 1000
				}
			],
			"labor_cost": 5000,
			"material_cost": 3000,
			"subtotal": 8000,
			"markup_amount": 1600,
			"total_price": 9600,
			"inclusions": ["Item 1"],
			"exclusions": ["Item 2"],
			"schedule": {"Phase 1": "1 week"},
			"payment_terms": "50% deposit",
			"warranty_terms": "1 year",
			"closing_statement": "Thank you"
		}`

		bidResponse, err := service.ParseBidDataFromJSON(jsonData)
		if err != nil {
			t.Errorf("ParseBidDataFromJSON() error = %v", err)
			return
		}

		if bidResponse.BidID != "test-id" {
			t.Errorf("Expected bid_id 'test-id', got '%s'", bidResponse.BidID)
		}

		if len(bidResponse.LineItems) != 1 {
			t.Errorf("Expected 1 line item, got %d", len(bidResponse.LineItems))
		}

		if bidResponse.TotalPrice != 9600 {
			t.Errorf("Expected total price 9600, got %f", bidResponse.TotalPrice)
		}
	})

	t.Run("parse invalid JSON", func(t *testing.T) {
		jsonData := `{invalid json}`

		_, err := service.ParseBidDataFromJSON(jsonData)
		if err == nil {
			t.Error("Expected error for invalid JSON, got nil")
		}
	})

	t.Run("parse empty JSON", func(t *testing.T) {
		jsonData := ``

		_, err := service.ParseBidDataFromJSON(jsonData)
		if err == nil {
			t.Error("Expected error for empty JSON, got nil")
		}
	})
}

func TestGeneratePDFFilename(t *testing.T) {
	service := NewPDFService()
	projectID := uuid.New()
	bidID := uuid.New()

	filename := service.GeneratePDFFilename(projectID, bidID)

	// Check filename format
	if len(filename) == 0 {
		t.Error("GeneratePDFFilename() returned empty string")
	}

	// Check it contains expected parts
	expectedPrefix := "bids/" + projectID.String() + "/bid-"
	if len(filename) < len(expectedPrefix) || filename[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Filename doesn't start with expected prefix. Got: %s", filename)
	}

	// Check it ends with .pdf
	if len(filename) < 4 || filename[len(filename)-4:] != ".pdf" {
		t.Error("Filename doesn't end with .pdf")
	}
}
