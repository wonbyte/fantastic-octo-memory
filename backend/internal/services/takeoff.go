package services

import (
	"encoding/json"
	"fmt"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type TakeoffService struct{}

func NewTakeoffService() *TakeoffService {
	return &TakeoffService{}
}

// CalculateTakeoffSummary computes deterministic takeoff summary from analysis data
func (s *TakeoffService) CalculateTakeoffSummary(analysis *models.AnalysisResult) (*models.TakeoffSummary, error) {
	if analysis == nil {
		return nil, fmt.Errorf("analysis result is nil")
	}

	summary := &models.TakeoffSummary{
		OpeningCounts:    make(map[string]int),
		FixtureCounts:    make(map[string]int),
		RoomBreakdown:    make([]models.RoomSummary, 0),
		OpeningBreakdown: make([]models.OpeningSummary, 0),
		FixtureBreakdown: make([]models.FixtureSummary, 0),
	}

	// Calculate room totals
	for _, room := range analysis.Rooms {
		summary.TotalArea += room.Area
		summary.RoomCount++

		// Parse dimensions to calculate perimeter if possible
		// Assuming dimensions are in format "WxL" or similar
		// For now, we estimate perimeter as 2*(sqrt(area)*2) if dimensions not parseable
		// In a production system, you'd parse dimensions more robustly
		perimeter := estimatePerimeter(room.Area, room.Dimensions)
		summary.TotalPerimeter += perimeter

		summary.RoomBreakdown = append(summary.RoomBreakdown, models.RoomSummary{
			Name:       room.Name,
			RoomType:   room.RoomType,
			Area:       room.Area,
			Dimensions: room.Dimensions,
		})
	}

	// Count openings by type
	for _, opening := range analysis.Openings {
		summary.OpeningCounts[opening.OpeningType] += opening.Count

		summary.OpeningBreakdown = append(summary.OpeningBreakdown, models.OpeningSummary{
			OpeningType: opening.OpeningType,
			Count:       opening.Count,
			Size:        opening.Size,
		})
	}

	// Count fixtures by category
	for _, fixture := range analysis.Fixtures {
		summary.FixtureCounts[fixture.Category] += fixture.Count

		summary.FixtureBreakdown = append(summary.FixtureBreakdown, models.FixtureSummary{
			FixtureType: fixture.FixtureType,
			Category:    fixture.Category,
			Count:       fixture.Count,
		})
	}

	return summary, nil
}

// estimatePerimeter calculates perimeter from area and dimensions string
// This is a simplified implementation - in production, parse actual dimensions
func estimatePerimeter(area float64, dimensions string) float64 {
	// Try to parse dimensions like "10x12" or "10' x 12'"
	// For now, use simple approximation: assume square room
	// perimeter = 4 * sqrt(area)
	if area <= 0 {
		return 0
	}

	// Simple approximation for square room
	// For rectangular room, approximate as 2*(W+L) where W*L = area
	// Use golden ratio approximation: W = sqrt(area/1.618), L = sqrt(area*1.618)
	// This gives reasonable perimeter estimates
	return 4.0 * (area / 10.0) // Simplified: assume average room is ~10ft on a side per 100 sq ft
}

// ParseAnalysisData parses JSONB string into AnalysisResult
func (s *TakeoffService) ParseAnalysisData(analysisJSON string) (*models.AnalysisResult, error) {
	if analysisJSON == "" {
		return nil, fmt.Errorf("analysis data is empty")
	}

	var analysis models.AnalysisResult
	if err := json.Unmarshal([]byte(analysisJSON), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse analysis data: %w", err)
	}

	return &analysis, nil
}
