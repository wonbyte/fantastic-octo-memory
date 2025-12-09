package services

import (
	"bytes"
	"testing"
)

func TestFileValidator_ValidateContentType(t *testing.T) {
	validator := NewFileValidator()

	tests := []struct {
		name        string
		contentType string
		expectError bool
	}{
		{"Valid PDF", "application/pdf", false},
		{"Valid JPEG", "image/jpeg", false},
		{"Valid PNG", "image/png", false},
		{"Valid GIF", "image/gif", false},
		{"Valid ZIP", "application/zip", false},
		{"Invalid type", "application/javascript", true},
		{"Empty type", "", true},
		{"Case insensitive", "IMAGE/PNG", false},
		{"With whitespace", " image/jpeg ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateContentType(tt.contentType)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestFileValidator_ValidateFileSize(t *testing.T) {
	validator := NewFileValidator()

	tests := []struct {
		name        string
		size        int64
		expectError bool
	}{
		{"Valid small file", 1024, false},
		{"Valid 1MB file", 1024 * 1024, false},
		{"Valid 50MB file", 50 * 1024 * 1024, false},
		{"Valid max file", MaxFileSize, false},
		{"Too large", MaxFileSize + 1, true},
		{"Zero size", 0, true},
		{"Negative size", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFileSize(tt.size)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestFileValidator_ValidateFileType(t *testing.T) {
	validator := NewFileValidator()

	tests := []struct {
		name        string
		contentType string
		fileContent []byte
		expectError bool
	}{
		{
			name:        "Valid PDF",
			contentType: "application/pdf",
			fileContent: []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x34}, // %PDF-1.4
			expectError: false,
		},
		{
			name:        "Valid JPEG",
			contentType: "image/jpeg",
			fileContent: []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46}, // JPEG header
			expectError: false,
		},
		{
			name:        "Valid PNG",
			contentType: "image/png",
			fileContent: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG signature
			expectError: false,
		},
		{
			name:        "Invalid magic bytes for PDF",
			contentType: "application/pdf",
			fileContent: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectError: true,
		},
		{
			name:        "File too small",
			contentType: "application/pdf",
			fileContent: []byte{0x25},
			expectError: true,
		},
		{
			name:        "Empty file",
			contentType: "application/pdf",
			fileContent: []byte{},
			expectError: true,
		},
		{
			name:        "Wrong type for content",
			contentType: "image/jpeg",
			fileContent: []byte{0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x34}, // PDF, not JPEG
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFileType(tt.contentType, tt.fileContent)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestFileValidator_ValidateWebP(t *testing.T) {
	validator := NewFileValidator()

	// Valid WEBP file header
	validWebP := make([]byte, 12)
	copy(validWebP[0:4], []byte{0x52, 0x49, 0x46, 0x46}) // RIFF
	copy(validWebP[8:12], []byte{0x57, 0x45, 0x42, 0x50}) // WEBP

	err := validator.ValidateFileType("image/webp", validWebP)
	if err != nil {
		t.Errorf("Valid WEBP should pass validation: %v", err)
	}

	// Invalid WEBP (missing WEBP signature)
	invalidWebP := make([]byte, 12)
	copy(invalidWebP[0:4], []byte{0x52, 0x49, 0x46, 0x46}) // RIFF
	copy(invalidWebP[8:12], []byte{0x00, 0x00, 0x00, 0x00}) // Not WEBP

	err = validator.ValidateFileType("image/webp", invalidWebP)
	if err == nil {
		t.Error("Invalid WEBP should fail validation")
	}
}

func TestFileValidator_GetAllowedContentTypes(t *testing.T) {
	validator := NewFileValidator()
	types := validator.GetAllowedContentTypes()

	if len(types) == 0 {
		t.Error("Should have allowed content types")
	}

	// Check that some expected types are present
	expectedTypes := []string{"application/pdf", "image/jpeg", "image/png"}
	for _, expected := range expectedTypes {
		found := false
		for _, actual := range types {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected type %s not found in allowed types", expected)
		}
	}
}

func TestFileValidator_GetMaxFileSize(t *testing.T) {
	validator := NewFileValidator()
	maxSize := validator.GetMaxFileSize()

	if maxSize != MaxFileSize {
		t.Errorf("Expected max size %d, got %d", MaxFileSize, maxSize)
	}
}

func TestFileValidator_RealWorldFiles(t *testing.T) {
	validator := NewFileValidator()

	// Test with realistic file headers
	tests := []struct {
		name        string
		contentType string
		fileHeader  []byte
		expectError bool
	}{
		{
			name:        "Real PDF header",
			contentType: "application/pdf",
			fileHeader:  []byte("%PDF-1.7"),
			expectError: false,
		},
		{
			name:        "JPEG with JFIF marker",
			contentType: "image/jpeg",
			fileHeader:  append([]byte{0xFF, 0xD8, 0xFF}, []byte{0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}...),
			expectError: false,
		},
		{
			name:        "ZIP file",
			contentType: "application/zip",
			fileHeader:  []byte{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x00, 0x00},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFileType(tt.contentType, tt.fileHeader)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// Helper function to create a buffer with specific content
func createBuffer(size int, pattern byte) *bytes.Buffer {
	buf := bytes.NewBuffer(make([]byte, 0, size))
	for i := 0; i < size; i++ {
		buf.WriteByte(pattern)
	}
	return buf
}
