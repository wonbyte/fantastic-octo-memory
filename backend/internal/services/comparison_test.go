package services

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

func TestCompareBlueprintRevisions_RoomChanges(t *testing.T) {
	service := NewComparisonService()

	// Create from revision with one room
	fromAnalysis := models.AnalysisResult{
		Rooms: []models.Room{
			{
				Name:       "Living Room",
				Dimensions: "20x15",
				Area:       300.0,
			},
		},
	}
	fromAnalysisJSON, _ := json.Marshal(fromAnalysis)
	fromAnalysisStr := string(fromAnalysisJSON)

	fromRevision := &models.BlueprintRevision{
		ID:           uuid.New(),
		BlueprintID:  uuid.New(),
		Version:      1,
		Filename:     "blueprint_v1.pdf",
		AnalysisData: &fromAnalysisStr,
	}

	// Create to revision with modified and added room
	toAnalysis := models.AnalysisResult{
		Rooms: []models.Room{
			{
				Name:       "Living Room",
				Dimensions: "25x15", // Modified dimensions
				Area:       375.0,
			},
			{
				Name:       "Kitchen", // New room
				Dimensions: "15x12",
				Area:       180.0,
			},
		},
	}
	toAnalysisJSON, _ := json.Marshal(toAnalysis)
	toAnalysisStr := string(toAnalysisJSON)

	toRevision := &models.BlueprintRevision{
		ID:           uuid.New(),
		BlueprintID:  fromRevision.BlueprintID,
		Version:      2,
		Filename:     "blueprint_v2.pdf",
		AnalysisData: &toAnalysisStr,
	}

	// Compare revisions
	comparison, err := service.CompareBlueprintRevisions(fromRevision, toRevision)
	if err != nil {
		t.Fatalf("failed to compare revisions: %v", err)
	}

	// Verify comparison results
	if comparison.FromVersion != 1 {
		t.Errorf("expected from version 1, got %d", comparison.FromVersion)
	}
	if comparison.ToVersion != 2 {
		t.Errorf("expected to version 2, got %d", comparison.ToVersion)
	}

	// Should have 2 changes: 1 modified, 1 added
	if comparison.Summary.TotalChanges != 2 {
		t.Errorf("expected 2 total changes, got %d", comparison.Summary.TotalChanges)
	}
	if comparison.Summary.ModifiedCount != 1 {
		t.Errorf("expected 1 modified change, got %d", comparison.Summary.ModifiedCount)
	}
	if comparison.Summary.AddedCount != 1 {
		t.Errorf("expected 1 added change, got %d", comparison.Summary.AddedCount)
	}

	// Verify changes contain room category
	foundRoomChanges := 0
	for _, change := range comparison.Changes {
		if change.Category == "room" {
			foundRoomChanges++
		}
	}
	if foundRoomChanges != 2 {
		t.Errorf("expected 2 room changes, got %d", foundRoomChanges)
	}
}

func TestCompareBidRevisions_CostChanges(t *testing.T) {
	service := NewComparisonService()

	// Create from revision
	laborCost1 := 5000.0
	materialCost1 := 3000.0
	totalCost1 := 8000.0
	finalPrice1 := 9600.0
	markup1 := 20.0

	bidData1 := models.GenerateBidResponse{
		LaborCost:    laborCost1,
		MaterialCost: materialCost1,
		Subtotal:     totalCost1,
		TotalPrice:   finalPrice1,
		LineItems: []models.LineItem{
			{
				Description: "Framing",
				Trade:       "carpentry",
				Quantity:    100,
				Unit:        "SF",
				UnitCost:    10.0,
				Total:       1000.0,
			},
		},
	}
	bidData1JSON, _ := json.Marshal(bidData1)
	bidData1Str := string(bidData1JSON)

	fromRevision := &models.BidRevision{
		ID:               uuid.New(),
		BidID:            uuid.New(),
		Version:          1,
		LaborCost:        &laborCost1,
		MaterialCost:     &materialCost1,
		TotalCost:        &totalCost1,
		FinalPrice:       &finalPrice1,
		MarkupPercentage: &markup1,
		Status:           models.BidStatusDraft,
		BidData:          &bidData1Str,
	}

	// Create to revision with changed costs
	laborCost2 := 6000.0
	materialCost2 := 3500.0
	totalCost2 := 9500.0
	finalPrice2 := 11400.0
	markup2 := 20.0

	bidData2 := models.GenerateBidResponse{
		LaborCost:    laborCost2,
		MaterialCost: materialCost2,
		Subtotal:     totalCost2,
		TotalPrice:   finalPrice2,
		LineItems: []models.LineItem{
			{
				Description: "Framing",
				Trade:       "carpentry",
				Quantity:    120, // Quantity changed
				Unit:        "SF",
				UnitCost:    10.0,
				Total:       1200.0,
			},
			{
				Description: "Drywall", // New line item
				Trade:       "drywall",
				Quantity:    500,
				Unit:        "SF",
				UnitCost:    5.0,
				Total:       2500.0,
			},
		},
	}
	bidData2JSON, _ := json.Marshal(bidData2)
	bidData2Str := string(bidData2JSON)

	toRevision := &models.BidRevision{
		ID:               uuid.New(),
		BidID:            fromRevision.BidID,
		Version:          2,
		LaborCost:        &laborCost2,
		MaterialCost:     &materialCost2,
		TotalCost:        &totalCost2,
		FinalPrice:       &finalPrice2,
		MarkupPercentage: &markup2,
		Status:           models.BidStatusDraft,
		BidData:          &bidData2Str,
	}

	// Compare revisions
	comparison, err := service.CompareBidRevisions(fromRevision, toRevision)
	if err != nil {
		t.Fatalf("failed to compare revisions: %v", err)
	}

	// Verify comparison results
	if comparison.FromVersion != 1 {
		t.Errorf("expected from version 1, got %d", comparison.FromVersion)
	}
	if comparison.ToVersion != 2 {
		t.Errorf("expected to version 2, got %d", comparison.ToVersion)
	}

	// Should have multiple changes: cost changes + line item changes
	if comparison.Summary.TotalChanges < 4 {
		t.Errorf("expected at least 4 total changes, got %d", comparison.Summary.TotalChanges)
	}

	// Verify we have cost changes
	foundCostChanges := 0
	for _, change := range comparison.Changes {
		if change.Category == "cost" {
			foundCostChanges++
		}
	}
	if foundCostChanges < 3 {
		t.Errorf("expected at least 3 cost changes, got %d", foundCostChanges)
	}

	// Verify we have line item changes
	foundLineItemChanges := 0
	for _, change := range comparison.Changes {
		if change.Category == "line_item" {
			foundLineItemChanges++
		}
	}
	if foundLineItemChanges < 2 {
		t.Errorf("expected at least 2 line item changes, got %d", foundLineItemChanges)
	}
}

func TestComparisonService_EmptyRevisions(t *testing.T) {
	service := NewComparisonService()

	// Create empty revisions
	emptyAnalysis := models.AnalysisResult{}
	emptyAnalysisJSON, _ := json.Marshal(emptyAnalysis)
	emptyAnalysisStr := string(emptyAnalysisJSON)

	fromRevision := &models.BlueprintRevision{
		ID:           uuid.New(),
		BlueprintID:  uuid.New(),
		Version:      1,
		Filename:     "blueprint_v1.pdf",
		AnalysisData: &emptyAnalysisStr,
	}

	toRevision := &models.BlueprintRevision{
		ID:           uuid.New(),
		BlueprintID:  fromRevision.BlueprintID,
		Version:      2,
		Filename:     "blueprint_v2.pdf",
		AnalysisData: &emptyAnalysisStr,
	}

	// Compare empty revisions
	comparison, err := service.CompareBlueprintRevisions(fromRevision, toRevision)
	if err != nil {
		t.Fatalf("failed to compare revisions: %v", err)
	}

	// Should have no changes
	if comparison.Summary.TotalChanges != 0 {
		t.Errorf("expected 0 changes, got %d", comparison.Summary.TotalChanges)
	}
}

func TestComparisonService_MaterialChanges(t *testing.T) {
	service := NewComparisonService()

	// Create from revision with materials
	fromAnalysis := models.AnalysisResult{
		Materials: []models.Material{
			{
				MaterialName: "2x4 Lumber",
				Quantity:     100,
				Unit:         "LF",
			},
		},
	}
	fromAnalysisJSON, _ := json.Marshal(fromAnalysis)
	fromAnalysisStr := string(fromAnalysisJSON)

	fromRevision := &models.BlueprintRevision{
		ID:           uuid.New(),
		BlueprintID:  uuid.New(),
		Version:      1,
		Filename:     "blueprint_v1.pdf",
		AnalysisData: &fromAnalysisStr,
	}

	// Create to revision with modified material
	toAnalysis := models.AnalysisResult{
		Materials: []models.Material{
			{
				MaterialName: "2x4 Lumber",
				Quantity:     150, // Increased quantity
				Unit:         "LF",
			},
		},
	}
	toAnalysisJSON, _ := json.Marshal(toAnalysis)
	toAnalysisStr := string(toAnalysisJSON)

	toRevision := &models.BlueprintRevision{
		ID:           uuid.New(),
		BlueprintID:  fromRevision.BlueprintID,
		Version:      2,
		Filename:     "blueprint_v2.pdf",
		AnalysisData: &toAnalysisStr,
	}

	// Compare revisions
	comparison, err := service.CompareBlueprintRevisions(fromRevision, toRevision)
	if err != nil {
		t.Fatalf("failed to compare revisions: %v", err)
	}

	// Should have 1 modified material change
	if comparison.Summary.TotalChanges != 1 {
		t.Errorf("expected 1 change, got %d", comparison.Summary.TotalChanges)
	}
	if comparison.Summary.ModifiedCount != 1 {
		t.Errorf("expected 1 modified change, got %d", comparison.Summary.ModifiedCount)
	}

	// Verify it's a material change with high impact (>20% change)
	if len(comparison.Changes) > 0 {
		change := comparison.Changes[0]
		if change.Category != "material" {
			t.Errorf("expected material category, got %s", change.Category)
		}
		if change.Impact == nil || *change.Impact != "High" {
			t.Errorf("expected High impact for 50%% quantity change")
		}
	}
}
