# Bid Export and Download Guide

This guide explains how to export and download bid proposals in multiple formats from the Construction Estimation & Bidding Platform.

## Table of Contents

- [Overview](#overview)
- [Supported Formats](#supported-formats)
- [PDF Export Features](#pdf-export-features)
- [CSV Export Features](#csv-export-features)
- [Excel Export Features](#excel-export-features)
- [API Endpoints](#api-endpoints)
- [Using the Export Features](#using-the-export-features)
- [Customization Options](#customization-options)
- [Sample Outputs](#sample-outputs)

## Overview

The platform provides professional bid export capabilities in multiple formats to support various business workflows and client preferences. All export formats maintain consistency in data presentation while optimizing for their specific use cases.

## Supported Formats

| Format | Use Case | File Extension | Content Type |
|--------|----------|----------------|--------------|
| **PDF** | Professional proposals for clients | `.pdf` | `application/pdf` |
| **CSV** | Data analysis and import to other systems | `.csv` | `text/csv` |
| **Excel** | Spreadsheet analysis and editing | `.csv` (Excel-compatible) | `text/csv` |

## PDF Export Features

### Professional Template

The PDF export includes a polished, professional format with:

- **Cover Page** (optional)
  - Company logo
  - Company name and branding
  - Project title
  - Bid reference number
  - Preparation date
  - Company contact information (address, phone, email, website)
  - License and insurance information

- **Project Information Section**
  - Project name and details
  - Bid identification
  - Current bid status
  - Date of preparation

- **Scope of Work**
  - Detailed description of work to be performed
  - Clear project objectives

- **Cost Breakdown**
  - Itemized line items table with:
    - Description of work
    - Trade/category
    - Quantity and unit
    - Unit cost
    - Total cost per line item

- **Trade Breakdown**
  - Summary grouped by trade
  - Item count per trade
  - Total cost per trade
  - Grand total

- **Cost Summary**
  - Material costs
  - Labor costs
  - Subtotal
  - Markup amount
  - **Total Price** (prominently displayed)

- **Inclusions**
  - Bulleted list of what's included in the bid
  - Materials and services covered

- **Exclusions**
  - Bulleted list of what's not included
  - Out-of-scope items

- **Project Schedule**
  - Timeline by phase
  - Key milestones

- **Payment Terms**
  - Payment schedule
  - Accepted payment methods
  - Terms and conditions

- **Warranty Terms**
  - Warranty coverage details
  - Duration and scope

- **Closing Statement**
  - Professional closing remarks
  - Next steps
  - Contact for questions

- **Footer**
  - Generation date and timestamp
  - Page numbers
  - Professional formatting

### Branding Options

PDFs can be customized with company branding:

```go
// Example: Generate PDF with company branding
options := &PDFOptions{
    CompanyInfo: &models.CompanyInfo{
        Name:          "Quality Construction Co.",
        Logo:          "/path/to/logo.png",
        Address:       "123 Main St, City, ST 12345",
        Phone:         "(555) 123-4567",
        Email:         "info@qualityconstruction.com",
        Website:       "www.qualityconstruction.com",
        LicenseNumber: "CA-123456",
        InsuranceInfo: "Fully insured and bonded",
    },
    IncludeCover: true,
    IncludeLogo:  true,
}
```

## CSV Export Features

The CSV format is optimized for data portability and analysis:

- **Header Section**
  - Export metadata
  - Project information
  - Bid details

- **Structured Sections**
  - Scope of Work
  - Line Items (with full details)
  - Trade Breakdown
  - Cost Summary
  - Inclusions
  - Exclusions
  - Project Schedule
  - Payment Terms
  - Warranty Terms

- **Data Format**
  - Plain text, comma-separated
  - Compatible with all spreadsheet applications
  - Easy to import into databases
  - Suitable for automated processing

### CSV Structure Example

```csv
Construction Bid Export - CSV Format
Project,Downtown Office Renovation
Bid ID,a1b2c3d4-e5f6-7890-abcd-ef1234567890
Date,2024-01-15
Status,draft

Line Items
Description,Trade,Quantity,Unit,Unit Cost,Total
Framing lumber,Framing,2500.00,BF,2.50,6250.00
Drywall installation,Drywall,1200.00,SF,1.75,2100.00
...
```

## Excel Export Features

The Excel export is a UTF-8 BOM-encoded CSV that opens perfectly in Microsoft Excel:

- All features of CSV export
- UTF-8 BOM encoding for proper character display
- Preserves formatting when opened in Excel
- Compatible with Excel 2007 and later
- Easy to edit and format further in Excel

**Note:** While the file extension is `.csv`, it's optimized for Excel and opens directly with proper formatting.

## API Endpoints

### Generate and Download PDF

```http
GET /bids/{id}/pdf
Authorization: Bearer <token>
```

**Response:**
```json
{
  "pdf_url": "https://s3.amazonaws.com/bucket/bids/project-id/bid-abc12345-20240115.pdf"
}
```

### Download CSV

```http
GET /bids/{id}/csv
Authorization: Bearer <token>
```

**Response:** Binary CSV file download

### Download Excel

```http
GET /bids/{id}/excel
Authorization: Bearer <token>
```

**Response:** Binary CSV file download (Excel-compatible)

## Using the Export Features

### Via Frontend Application

1. Navigate to your project
2. View the generated bid
3. Click on the export option
4. Select desired format (PDF, CSV, or Excel)
5. File downloads automatically

### Via API

#### Download PDF

```typescript
import { bidsApi } from './api/bids';

// Get PDF URL
const { pdf_url } = await bidsApi.getBidPDF(bidId);

// Open or download
window.open(pdf_url, '_blank');
```

#### Download CSV

```typescript
import { bidsApi } from './api/bids';

// Download CSV
const blob = await bidsApi.downloadBidCSV(bidId);

// Create download link
const url = window.URL.createObjectURL(blob);
const link = document.createElement('a');
link.href = url;
link.download = `bid-${bidId}.csv`;
link.click();
```

#### Download Excel

```typescript
import { bidsApi } from './api/bids';

// Download Excel
const blob = await bidsApi.downloadBidExcel(bidId);

// Create download link
const url = window.URL.createObjectURL(blob);
const link = document.createElement('a');
link.href = url;
link.download = `bid-${bidId}.csv`;
link.click();
```

### Via cURL

```bash
# Get PDF URL
curl -X GET "https://api.example.com/bids/{bid-id}/pdf" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Download CSV
curl -X GET "https://api.example.com/bids/{bid-id}/csv" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o bid-export.csv

# Download Excel
curl -X GET "https://api.example.com/bids/{bid-id}/excel" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o bid-export.csv
```

## Customization Options

### Company Branding (Backend)

Update the User model to include branding information:

```go
type User struct {
    CompanyName    *string
    CompanyLogo    *string // S3 URL to logo
    CompanyPhone   *string
    CompanyAddress *string
    LicenseNumber  *string
}
```

### PDF Generation Options

Customize PDF generation in the handler:

```go
pdfService := services.NewPDFService()

// Basic PDF (no branding)
pdfBytes, err := pdfService.GenerateBidPDF(bid, bidResponse, projectName)

// PDF with full branding
options := &services.PDFOptions{
    CompanyInfo:  getCompanyInfo(user),
    IncludeCover: true,
    IncludeLogo:  true,
    LogoPath:     "/path/to/downloaded/logo.png",
}
pdfBytes, err := pdfService.GenerateBidPDFWithOptions(bid, bidResponse, projectName, options)
```

## Sample Outputs

### PDF Sample Structure

```
┌─────────────────────────────────────────┐
│           COVER PAGE                     │
│   [LOGO]                                │
│   Company Name                          │
│   BID PROPOSAL                          │
│   Project Name                          │
│   Date & Reference                      │
│   Contact Information                   │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│   PAGE 2                                │
│   Project Information                   │
│   Scope of Work                         │
│   Cost Breakdown Table                  │
│   Trade Breakdown Summary               │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│   PAGE 3                                │
│   Cost Summary                          │
│   Inclusions & Exclusions               │
│   Schedule                              │
│   Terms & Warranty                      │
└─────────────────────────────────────────┘
```

### CSV Sample Structure

```
Construction Bid Export - CSV Format
Project,Downtown Office Renovation
Bid ID,a1b2c3d4
Date,2024-01-15

Scope of Work
Complete renovation of 5000 sq ft office space...

Line Items
Description,Trade,Quantity,Unit,Unit Cost,Total
Framing lumber,Framing,2500.00,BF,2.50,6250.00
...

Trade Breakdown
Trade,Item Count,Total Cost
Framing,15,45000.00
...

Cost Summary
Material Cost,125000.00
Labor Cost,85000.00
Subtotal,210000.00
Markup Amount,42000.00
Total Price,252000.00
```

## Best Practices

1. **PDF for Client Presentation**
   - Use PDF for formal bid submissions
   - Include company branding for professional appearance
   - Enable cover page for important bids

2. **CSV for Data Analysis**
   - Use CSV for importing into other systems
   - Ideal for cost analysis in other tools
   - Easy to parse programmatically

3. **Excel for Collaboration**
   - Use Excel format when sharing with team members
   - Allows easy editing and formatting
   - Compatible with Microsoft Office

4. **Consistent Data**
   - All formats contain the same core information
   - Choose format based on use case, not content
   - Data integrity maintained across all exports

## Troubleshooting

### PDF Generation Issues

- **Logo not appearing**: Ensure logo path is accessible and in PNG format
- **Missing sections**: Check that bid data includes all required fields
- **Font issues**: The system uses Arial, which is widely available

### CSV/Excel Issues

- **Character encoding**: Excel export includes UTF-8 BOM for proper display
- **Date formats**: Dates are in ISO format (YYYY-MM-DD)
- **Number formats**: Numbers use decimal points (.) not commas

## Support

For additional help or feature requests:
- Open a GitHub issue
- Contact support team
- Refer to main README.md for general platform documentation
