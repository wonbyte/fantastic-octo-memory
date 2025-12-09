package services

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	// Max file size: 100MB
	MaxFileSize = 100 * 1024 * 1024
)

// FileValidator provides file validation functionality
type FileValidator struct {
	maxSize       int64
	allowedTypes  map[string][]byte
}

// NewFileValidator creates a new file validator
func NewFileValidator() *FileValidator {
	return &FileValidator{
		maxSize: MaxFileSize,
		allowedTypes: map[string][]byte{
			// PDF files
			"application/pdf": {0x25, 0x50, 0x44, 0x46}, // %PDF
			// Image files
			"image/jpeg": {0xFF, 0xD8, 0xFF},         // JPEG
			"image/png":  {0x89, 0x50, 0x4E, 0x47},   // PNG
			"image/gif":  {0x47, 0x49, 0x46, 0x38},   // GIF
			"image/bmp":  {0x42, 0x4D},               // BMP
			"image/webp": {0x52, 0x49, 0x46, 0x46},   // WEBP (RIFF)
			// CAD files
			"application/acad":            {0x41, 0x43, 0x31, 0x30}, // DWG (AutoCAD)
			"application/x-autocad":       {0x41, 0x43, 0x31, 0x30}, // DWG (AutoCAD)
			"application/dxf":             {0x30, 0x0D, 0x0A},       // DXF (ASCII)
			"image/vnd.dwg":              {0x41, 0x43, 0x31, 0x30}, // DWG
			// ZIP-based formats (might contain CAD files)
			"application/zip":             {0x50, 0x4B, 0x03, 0x04}, // ZIP
			"application/x-zip-compressed": {0x50, 0x4B, 0x03, 0x04}, // ZIP
		},
	}
}

// ValidateFileType validates a file based on its magic bytes (file signature)
func (fv *FileValidator) ValidateFileType(contentType string, fileContent []byte) error {
	if len(fileContent) == 0 {
		return fmt.Errorf("file content is empty")
	}

	// Normalize content type
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	// Get expected magic bytes for this content type
	expectedMagic, exists := fv.allowedTypes[contentType]
	if !exists {
		return fmt.Errorf("content type '%s' is not allowed", contentType)
	}

	// Check if file has enough bytes to match the magic signature
	if len(fileContent) < len(expectedMagic) {
		return fmt.Errorf("file is too small to validate")
	}

	// Compare magic bytes
	actualMagic := fileContent[:len(expectedMagic)]
	
	// Special handling for WEBP - need to check for WEBP in the file header
	if contentType == "image/webp" {
		if len(fileContent) >= 12 {
			// WEBP format: RIFF....WEBP
			if bytes.Equal(fileContent[0:4], []byte{0x52, 0x49, 0x46, 0x46}) &&
			   bytes.Equal(fileContent[8:12], []byte{0x57, 0x45, 0x42, 0x50}) {
				return nil
			}
		}
		return fmt.Errorf("file does not match expected WEBP format")
	}

	// Standard magic bytes comparison
	if !bytes.Equal(actualMagic, expectedMagic) {
		return fmt.Errorf("file type mismatch: content type is '%s' but file signature is %s", 
			contentType, hex.EncodeToString(actualMagic))
	}

	return nil
}

// ValidateFileSize validates file size
func (fv *FileValidator) ValidateFileSize(size int64) error {
	if size <= 0 {
		return fmt.Errorf("file size must be greater than 0")
	}

	if size > fv.maxSize {
		return fmt.Errorf("file size (%d bytes) exceeds maximum allowed size (%d bytes)", size, fv.maxSize)
	}

	return nil
}

// ValidateContentType validates the content type string
func (fv *FileValidator) ValidateContentType(contentType string) error {
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	
	if contentType == "" {
		return fmt.Errorf("content type is required")
	}

	if _, exists := fv.allowedTypes[contentType]; !exists {
		return fmt.Errorf("content type '%s' is not allowed", contentType)
	}

	return nil
}

// GetAllowedContentTypes returns a list of allowed content types
func (fv *FileValidator) GetAllowedContentTypes() []string {
	types := make([]string, 0, len(fv.allowedTypes))
	for contentType := range fv.allowedTypes {
		types = append(types, contentType)
	}
	return types
}

// GetMaxFileSize returns the maximum allowed file size
func (fv *FileValidator) GetMaxFileSize() int64 {
	return fv.maxSize
}
