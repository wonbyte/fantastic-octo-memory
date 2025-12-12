package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

// PDFService generates bid PDFs
type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

// PDFOptions contains configuration for PDF generation
type PDFOptions struct {
	CompanyInfo   *models.CompanyInfo
	IncludeCover  bool
	IncludeLogo   bool
	LogoPath      string // Path to downloaded logo file if needed
}

// GenerateBidPDF creates a professional bid PDF from bid data
func (s *PDFService) GenerateBidPDF(bid *models.Bid, bidResponse *models.GenerateBidResponse, projectName string) ([]byte, error) {
	return s.GenerateBidPDFWithOptions(bid, bidResponse, projectName, nil)
}

// GenerateBidPDFWithOptions creates a professional bid PDF with custom options
func (s *PDFService) GenerateBidPDFWithOptions(bid *models.Bid, bidResponse *models.GenerateBidResponse, projectName string, options *PDFOptions) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	
	// Add cover page if requested
	if options != nil && options.IncludeCover && options.CompanyInfo != nil {
		s.addCoverPage(pdf, projectName, bid, options.CompanyInfo, options.LogoPath)
	}
	
	pdf.AddPage()

	// Header with company branding
	if options != nil && options.CompanyInfo != nil {
		s.addHeaderWithBranding(pdf, projectName, options.CompanyInfo, options.LogoPath)
	} else {
		s.addHeader(pdf, projectName)
	}

	// Company & Project Info
	pdf.Ln(10)
	s.addSection(pdf, "Project Information")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(40, 6, "Project:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, projectName, "", 0, "L", false, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(40, 6, "Bid ID:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, bid.ID.String()[:8]+"...", "", 0, "L", false, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(40, 6, "Date:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, time.Now().Format("January 2, 2006"), "", 0, "L", false, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(40, 6, "Status:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, string(bid.Status), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Scope of Work
	if bidResponse.ScopeOfWork != "" {
		s.addSection(pdf, "Scope of Work")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, bidResponse.ScopeOfWork, "", "", false)
		pdf.Ln(5)
	}

	// Line Items
	if len(bidResponse.LineItems) > 0 {
		s.addSection(pdf, "Cost Breakdown")
		s.addLineItemsTable(pdf, bidResponse.LineItems)
		pdf.Ln(5)
	}

	// Trade Breakdown
	if len(bidResponse.LineItems) > 0 {
		s.addSection(pdf, "Trade Breakdown")
		s.addTradeBreakdown(pdf, bidResponse.LineItems)
		pdf.Ln(5)
	}

	// Cost Summary
	s.addSection(pdf, "Cost Summary")
	s.addCostSummary(pdf, bidResponse)
	pdf.Ln(5)

	// Inclusions
	if len(bidResponse.Inclusions) > 0 {
		s.addSection(pdf, "Inclusions")
		pdf.SetFont("Arial", "", 10)
		for _, inclusion := range bidResponse.Inclusions {
			pdf.CellFormat(5, 5, "", "", 0, "L", false, 0, "")
			pdf.CellFormat(5, 5, "•", "", 0, "L", false, 0, "")
			pdf.MultiCell(0, 5, inclusion, "", "", false)
		}
		pdf.Ln(3)
	}

	// Exclusions
	if len(bidResponse.Exclusions) > 0 {
		s.addSection(pdf, "Exclusions")
		pdf.SetFont("Arial", "", 10)
		for _, exclusion := range bidResponse.Exclusions {
			pdf.CellFormat(5, 5, "", "", 0, "L", false, 0, "")
			pdf.CellFormat(5, 5, "•", "", 0, "L", false, 0, "")
			pdf.MultiCell(0, 5, exclusion, "", "", false)
		}
		pdf.Ln(3)
	}

	// Schedule
	if len(bidResponse.Schedule) > 0 {
		s.addSection(pdf, "Project Schedule")
		pdf.SetFont("Arial", "", 10)
		for phase, timeline := range bidResponse.Schedule {
			pdf.CellFormat(5, 5, "", "", 0, "L", false, 0, "")
			pdf.CellFormat(80, 5, phase+":", "", 0, "L", false, 0, "")
			pdf.CellFormat(0, 5, timeline, "", 0, "L", false, 0, "")
			pdf.Ln(5)
		}
		pdf.Ln(3)
	}

	// Payment Terms
	if bidResponse.PaymentTerms != "" {
		s.addSection(pdf, "Payment Terms")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, bidResponse.PaymentTerms, "", "", false)
		pdf.Ln(3)
	}

	// Warranty Terms
	if bidResponse.WarrantyTerms != "" {
		s.addSection(pdf, "Warranty")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, bidResponse.WarrantyTerms, "", "", false)
		pdf.Ln(3)
	}

	// Closing Statement
	if bidResponse.ClosingStatement != "" {
		s.addSection(pdf, "Closing")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, bidResponse.ClosingStatement, "", "", false)
	}

	// Footer
	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 10, fmt.Sprintf("Generated on %s | Page %d", time.Now().Format("January 2, 2006"), pdf.PageNo()), "", 0, "C", false, 0, "")

	// Output to buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// addCoverPage creates a professional cover page with company branding
func (s *PDFService) addCoverPage(pdf *gofpdf.Fpdf, projectName string, bid *models.Bid, companyInfo *models.CompanyInfo, logoPath string) {
	pdf.AddPage()
	
	// Add logo if available
	if logoPath != "" {
		// Try to add logo - if it fails, continue without it
		pdf.ImageOptions(logoPath, 70, 30, 70, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	}
	
	// Company Name
	pdf.SetY(100)
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 15, companyInfo.Name, "", 0, "C", false, 0, "")
	pdf.Ln(20)
	
	// Title
	pdf.SetFont("Arial", "B", 28)
	pdf.SetTextColor(41, 128, 185) // Professional blue
	pdf.CellFormat(0, 15, "BID PROPOSAL", "", 0, "C", false, 0, "")
	pdf.Ln(20)
	pdf.SetTextColor(0, 0, 0) // Reset to black
	
	// Project Name
	pdf.SetFont("Arial", "B", 18)
	pdf.MultiCell(0, 10, projectName, "", "C", false)
	pdf.Ln(30)
	
	// Date
	pdf.SetFont("Arial", "", 14)
	pdf.CellFormat(0, 8, "Prepared: "+time.Now().Format("January 2, 2006"), "", 0, "C", false, 0, "")
	pdf.Ln(10)
	
	// Bid ID
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(0, 6, "Reference: "+bid.ID.String()[:13], "", 0, "C", false, 0, "")
	pdf.Ln(40)
	
	// Company Contact Information
	pdf.SetFont("Arial", "", 11)
	if companyInfo.Address != nil {
		pdf.CellFormat(0, 6, *companyInfo.Address, "", 0, "C", false, 0, "")
		pdf.Ln(6)
	}
	if companyInfo.Phone != nil {
		pdf.CellFormat(0, 6, "Phone: "+*companyInfo.Phone, "", 0, "C", false, 0, "")
		pdf.Ln(6)
	}
	if companyInfo.Email != nil {
		pdf.CellFormat(0, 6, "Email: "+*companyInfo.Email, "", 0, "C", false, 0, "")
		pdf.Ln(6)
	}
	if companyInfo.Website != nil {
		pdf.CellFormat(0, 6, "Web: "+*companyInfo.Website, "", 0, "C", false, 0, "")
		pdf.Ln(6)
	}
	if companyInfo.LicenseNumber != nil {
		pdf.Ln(6)
		pdf.SetFont("Arial", "I", 9)
		pdf.CellFormat(0, 6, "License: "+*companyInfo.LicenseNumber, "", 0, "C", false, 0, "")
	}
}

// addHeaderWithBranding creates a header with company branding
func (s *PDFService) addHeaderWithBranding(pdf *gofpdf.Fpdf, projectName string, companyInfo *models.CompanyInfo, logoPath string) {
	startY := pdf.GetY()
	
	// Add small logo if available (top right corner)
	if logoPath != "" {
		pdf.ImageOptions(logoPath, 160, startY, 30, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	}
	
	// Company name and title
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, companyInfo.Name, "", 0, "L", false, 0, "")
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 10, "Construction Bid Proposal", "", 0, "L", false, 0, "")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 6, projectName, "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.SetLineWidth(0.5)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
}

func (s *PDFService) addHeader(pdf *gofpdf.Fpdf, projectName string) {
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(0, 10, "Construction Bid Proposal", "", 0, "L", false, 0, "")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 6, projectName, "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.SetLineWidth(0.5)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
}

func (s *PDFService) addSection(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, title, "", 0, "L", false, 0, "")
	pdf.Ln(8)
}

func (s *PDFService) addLineItemsTable(pdf *gofpdf.Fpdf, items []models.LineItem) {
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	
	// Header
	pdf.CellFormat(80, 6, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 6, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 6, "Unit", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 6, "Unit Cost", "1", 0, "R", true, 0, "")
	pdf.CellFormat(25, 6, "Total", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// Items
	pdf.SetFont("Arial", "", 9)
	for _, item := range items {
		pdf.CellFormat(80, 6, item.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%.1f", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 6, item.Unit, "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("$%.2f", item.UnitCost), "1", 0, "R", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("$%.2f", item.Total), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
}

// addTradeBreakdown groups line items by trade and shows totals
func (s *PDFService) addTradeBreakdown(pdf *gofpdf.Fpdf, items []models.LineItem) {
	// Group items by trade
	tradeGroups := make(map[string][]models.LineItem)
	tradeTotals := make(map[string]float64)
	
	for _, item := range items {
		trade := item.Trade
		if trade == "" {
			trade = "General"
		}
		tradeGroups[trade] = append(tradeGroups[trade], item)
		tradeTotals[trade] += item.Total
	}
	
	// Display trade summary table
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(240, 240, 240)
	
	// Header
	pdf.CellFormat(120, 6, "Trade", "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 6, "Items", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 6, "Total", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)
	
	// Trade rows
	pdf.SetFont("Arial", "", 9)
	var grandTotal float64
	for trade, items := range tradeGroups {
		total := tradeTotals[trade]
		grandTotal += total
		
		pdf.CellFormat(120, 6, trade, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("%d", len(items)), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("$%.2f", total), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	
	// Grand total
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(120, 6, "Total", "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 6, "", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 6, fmt.Sprintf("$%.2f", grandTotal), "1", 0, "R", true, 0, "")
	pdf.Ln(-1)
}

func (s *PDFService) addCostSummary(pdf *gofpdf.Fpdf, bidResponse *models.GenerateBidResponse) {
	pdf.SetFont("Arial", "", 10)
	
	// Right-align summary
	x := 120.0
	
	pdf.SetX(x)
	pdf.CellFormat(40, 6, "Material Cost:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, fmt.Sprintf("$%.2f", bidResponse.MaterialCost), "", 0, "R", false, 0, "")
	pdf.Ln(6)
	
	pdf.SetX(x)
	pdf.CellFormat(40, 6, "Labor Cost:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, fmt.Sprintf("$%.2f", bidResponse.LaborCost), "", 0, "R", false, 0, "")
	pdf.Ln(6)
	
	pdf.SetX(x)
	pdf.CellFormat(40, 6, "Subtotal:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, fmt.Sprintf("$%.2f", bidResponse.Subtotal), "", 0, "R", false, 0, "")
	pdf.Ln(6)
	
	pdf.SetX(x)
	pdf.CellFormat(40, 6, "Markup:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 6, fmt.Sprintf("$%.2f", bidResponse.MarkupAmount), "", 0, "R", false, 0, "")
	pdf.Ln(6)
	
	// Total with emphasis
	pdf.SetFont("Arial", "B", 12)
	pdf.SetX(x)
	pdf.CellFormat(40, 8, "Total Price:", "", 0, "L", false, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", bidResponse.TotalPrice), "", 0, "R", false, 0, "")
	pdf.Ln(8)
}

// ParseBidDataFromJSON parses bid_data JSONB field into GenerateBidResponse
func (s *PDFService) ParseBidDataFromJSON(bidData string) (*models.GenerateBidResponse, error) {
	var bidResponse models.GenerateBidResponse
	if err := json.Unmarshal([]byte(bidData), &bidResponse); err != nil {
		return nil, fmt.Errorf("failed to parse bid data: %w", err)
	}
	return &bidResponse, nil
}

// GeneratePDFFilename creates a unique filename for the bid PDF
func (s *PDFService) GeneratePDFFilename(projectID uuid.UUID, bidID uuid.UUID) string {
	timestamp := time.Now().Format("20060102-150405")
	return fmt.Sprintf("bids/%s/bid-%s-%s.pdf", projectID.String(), bidID.String()[:8], timestamp)
}
