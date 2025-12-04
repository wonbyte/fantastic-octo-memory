package services

import (
	"testing"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

func TestCalculateTakeoffSummary(t *testing.T) {
	service := NewTakeoffService()

	tests := []struct {
		name     string
		analysis *models.AnalysisResult
		wantErr  bool
	}{
		{
			name:     "nil analysis",
			analysis: nil,
			wantErr:  true,
		},
		{
			name: "empty analysis",
			analysis: &models.AnalysisResult{
				BlueprintID: "test-id",
				Status:      "completed",
				Rooms:       []models.Room{},
				Openings:    []models.Opening{},
				Fixtures:    []models.Fixture{},
			},
			wantErr: false,
		},
		{
			name: "analysis with data",
			analysis: &models.AnalysisResult{
				BlueprintID: "test-id",
				Status:      "completed",
				Rooms: []models.Room{
					{
						Name:       "Living Room",
						Dimensions: "15x20",
						Area:       300,
					},
					{
						Name:       "Bedroom",
						Dimensions: "12x14",
						Area:       168,
					},
				},
				Openings: []models.Opening{
					{
						OpeningType: "door",
						Count:       3,
						Size:        "36x80",
					},
					{
						OpeningType: "window",
						Count:       5,
						Size:        "48x60",
					},
				},
				Fixtures: []models.Fixture{
					{
						FixtureType: "outlet",
						Category:    "electrical",
						Count:       12,
					},
					{
						FixtureType: "light_switch",
						Category:    "electrical",
						Count:       8,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary, err := service.CalculateTakeoffSummary(tt.analysis)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateTakeoffSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Validate summary structure
				if summary == nil {
					t.Error("expected summary, got nil")
					return
				}

				if summary.OpeningCounts == nil {
					t.Error("expected opening_counts map, got nil")
				}
				if summary.FixtureCounts == nil {
					t.Error("expected fixture_counts map, got nil")
				}

				// For analysis with data, verify calculations
				if tt.analysis != nil && len(tt.analysis.Rooms) > 0 {
					expectedArea := 0.0
					for _, room := range tt.analysis.Rooms {
						expectedArea += room.Area
					}
					if summary.TotalArea != expectedArea {
						t.Errorf("expected total_area %f, got %f", expectedArea, summary.TotalArea)
					}

					if summary.RoomCount != len(tt.analysis.Rooms) {
						t.Errorf("expected room_count %d, got %d", len(tt.analysis.Rooms), summary.RoomCount)
					}
				}

				// Verify opening counts
				if tt.analysis != nil && len(tt.analysis.Openings) > 0 {
					for _, opening := range tt.analysis.Openings {
						count, exists := summary.OpeningCounts[opening.OpeningType]
						if !exists {
							t.Errorf("expected opening type %s in counts", opening.OpeningType)
						}
						if count != opening.Count {
							t.Errorf("expected count %d for %s, got %d", opening.Count, opening.OpeningType, count)
						}
					}
				}

				// Verify fixture counts
				if tt.analysis != nil && len(tt.analysis.Fixtures) > 0 {
					// Fixture counts are summed by category
					expectedCounts := make(map[string]int)
					for _, fixture := range tt.analysis.Fixtures {
						expectedCounts[fixture.Category] += fixture.Count
					}
					
					for category, expectedCount := range expectedCounts {
						count, exists := summary.FixtureCounts[category]
						if !exists {
							t.Errorf("expected fixture category %s in counts", category)
						}
						if count != expectedCount {
							t.Errorf("expected count %d for category %s, got %d", expectedCount, category, count)
						}
					}
				}
			}
		})
	}
}

func TestParseAnalysisData(t *testing.T) {
	service := NewTakeoffService()

	tests := []struct {
		name        string
		analysisJSON string
		wantErr     bool
	}{
		{
			name:        "empty string",
			analysisJSON: "",
			wantErr:     true,
		},
		{
			name:        "invalid JSON",
			analysisJSON: "not json",
			wantErr:     true,
		},
		{
			name: "valid JSON",
			analysisJSON: `{
				"blueprint_id": "test-id",
				"status": "completed",
				"rooms": [
					{
						"name": "Living Room",
						"dimensions": "15x20",
						"area": 300
					}
				],
				"openings": [],
				"fixtures": [],
				"measurements": [],
				"materials": [],
				"confidence_score": 0.95,
				"processing_time_ms": 1500
			}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ParseAnalysisData(tt.analysisJSON)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAnalysisData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}
