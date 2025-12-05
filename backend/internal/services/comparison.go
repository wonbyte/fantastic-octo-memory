package services

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/wonbyte/fantastic-octo-memory/backend/internal/models"
)

type ComparisonService struct{}

func NewComparisonService() *ComparisonService {
	return &ComparisonService{}
}

// CompareBlueprintRevisions compares two blueprint revisions and returns the differences
func (s *ComparisonService) CompareBlueprintRevisions(from, to *models.BlueprintRevision) (*models.BlueprintComparison, error) {
	comparison := &models.BlueprintComparison{
		FromVersion: from.Version,
		ToVersion:   to.Version,
		Changes:     []models.BlueprintChange{},
		Summary: models.ComparisonSummary{
			ChangesByCategory: make(map[string]int),
		},
	}

	// Parse analysis data from both revisions
	var fromAnalysis, toAnalysis models.AnalysisResult
	if from.AnalysisData != nil {
		if err := json.Unmarshal([]byte(*from.AnalysisData), &fromAnalysis); err != nil {
			return nil, fmt.Errorf("failed to parse from analysis data: %w", err)
		}
	}
	if to.AnalysisData != nil {
		if err := json.Unmarshal([]byte(*to.AnalysisData), &toAnalysis); err != nil {
			return nil, fmt.Errorf("failed to parse to analysis data: %w", err)
		}
	}

	// Compare rooms
	s.compareRooms(&fromAnalysis, &toAnalysis, comparison)

	// Compare openings
	s.compareOpenings(&fromAnalysis, &toAnalysis, comparison)

	// Compare fixtures
	s.compareFixtures(&fromAnalysis, &toAnalysis, comparison)

	// Compare measurements
	s.compareMeasurements(&fromAnalysis, &toAnalysis, comparison)

	// Compare materials
	s.compareMaterials(&fromAnalysis, &toAnalysis, comparison)

	// Calculate summary
	s.calculateSummary(comparison)

	return comparison, nil
}

func (s *ComparisonService) compareRooms(from, to *models.AnalysisResult, comparison *models.BlueprintComparison) {
	fromRooms := make(map[string]models.Room)
	for _, room := range from.Rooms {
		fromRooms[room.Name] = room
	}

	toRooms := make(map[string]models.Room)
	for _, room := range to.Rooms {
		toRooms[room.Name] = room
	}

	// Find added and modified rooms
	for name, toRoom := range toRooms {
		if fromRoom, exists := fromRooms[name]; exists {
			// Check for modifications
			if fromRoom.Area != toRoom.Area || fromRoom.Dimensions != toRoom.Dimensions {
				impact := "Medium"
				// Only check percentage if fromRoom.Area is not zero
				if fromRoom.Area > 0 && math.Abs(fromRoom.Area-toRoom.Area) > fromRoom.Area*0.2 { // >20% change
					impact = "High"
				} else if fromRoom.Area == 0 && toRoom.Area > 0 {
					impact = "High"
				}
				comparison.Changes = append(comparison.Changes, models.BlueprintChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "room",
					Description: fmt.Sprintf("Room '%s' dimensions changed from %s (%.2f SF) to %s (%.2f SF)", name, fromRoom.Dimensions, fromRoom.Area, toRoom.Dimensions, toRoom.Area),
					OldValue:    fromRoom,
					NewValue:    toRoom,
					Impact:      &impact,
				})
			}
		} else {
			// Room added
			impact := "Medium"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeAdded,
				Category:    "room",
				Description: fmt.Sprintf("Room '%s' added with dimensions %s (%.2f SF)", name, toRoom.Dimensions, toRoom.Area),
				NewValue:    toRoom,
				Impact:      &impact,
			})
		}
	}

	// Find removed rooms
	for name, fromRoom := range fromRooms {
		if _, exists := toRooms[name]; !exists {
			impact := "High"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeRemoved,
				Category:    "room",
				Description: fmt.Sprintf("Room '%s' removed (was %s, %.2f SF)", name, fromRoom.Dimensions, fromRoom.Area),
				OldValue:    fromRoom,
				Impact:      &impact,
			})
		}
	}
}

func (s *ComparisonService) compareOpenings(from, to *models.AnalysisResult, comparison *models.BlueprintComparison) {
	fromOpenings := make(map[string]models.Opening)
	for _, opening := range from.Openings {
		key := fmt.Sprintf("%s-%s", opening.OpeningType, opening.Size)
		fromOpenings[key] = opening
	}

	toOpenings := make(map[string]models.Opening)
	for _, opening := range to.Openings {
		key := fmt.Sprintf("%s-%s", opening.OpeningType, opening.Size)
		toOpenings[key] = opening
	}

	// Compare openings
	for key, toOpening := range toOpenings {
		if fromOpening, exists := fromOpenings[key]; exists {
			if fromOpening.Count != toOpening.Count {
				impact := "Medium"
				comparison.Changes = append(comparison.Changes, models.BlueprintChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "opening",
					Description: fmt.Sprintf("%s (%s) count changed from %d to %d", toOpening.OpeningType, toOpening.Size, fromOpening.Count, toOpening.Count),
					OldValue:    fromOpening,
					NewValue:    toOpening,
					Impact:      &impact,
				})
			}
		} else {
			impact := "Low"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeAdded,
				Category:    "opening",
				Description: fmt.Sprintf("%s (%s) added, count: %d", toOpening.OpeningType, toOpening.Size, toOpening.Count),
				NewValue:    toOpening,
				Impact:      &impact,
			})
		}
	}

	for key, fromOpening := range fromOpenings {
		if _, exists := toOpenings[key]; !exists {
			impact := "Low"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeRemoved,
				Category:    "opening",
				Description: fmt.Sprintf("%s (%s) removed, was count: %d", fromOpening.OpeningType, fromOpening.Size, fromOpening.Count),
				OldValue:    fromOpening,
				Impact:      &impact,
			})
		}
	}
}

func (s *ComparisonService) compareFixtures(from, to *models.AnalysisResult, comparison *models.BlueprintComparison) {
	fromFixtures := make(map[string]models.Fixture)
	for _, fixture := range from.Fixtures {
		key := fmt.Sprintf("%s-%s", fixture.Category, fixture.FixtureType)
		fromFixtures[key] = fixture
	}

	toFixtures := make(map[string]models.Fixture)
	for _, fixture := range to.Fixtures {
		key := fmt.Sprintf("%s-%s", fixture.Category, fixture.FixtureType)
		toFixtures[key] = fixture
	}

	// Compare fixtures
	for key, toFixture := range toFixtures {
		if fromFixture, exists := fromFixtures[key]; exists {
			if fromFixture.Count != toFixture.Count {
				impact := "Low"
				comparison.Changes = append(comparison.Changes, models.BlueprintChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "fixture",
					Description: fmt.Sprintf("%s %s count changed from %d to %d", toFixture.Category, toFixture.FixtureType, fromFixture.Count, toFixture.Count),
					OldValue:    fromFixture,
					NewValue:    toFixture,
					Impact:      &impact,
				})
			}
		} else {
			impact := "Low"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeAdded,
				Category:    "fixture",
				Description: fmt.Sprintf("%s %s added, count: %d", toFixture.Category, toFixture.FixtureType, toFixture.Count),
				NewValue:    toFixture,
				Impact:      &impact,
			})
		}
	}

	for key, fromFixture := range fromFixtures {
		if _, exists := toFixtures[key]; !exists {
			impact := "Low"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeRemoved,
				Category:    "fixture",
				Description: fmt.Sprintf("%s %s removed, was count: %d", fromFixture.Category, fromFixture.FixtureType, fromFixture.Count),
				OldValue:    fromFixture,
				Impact:      &impact,
			})
		}
	}
}

func (s *ComparisonService) compareMeasurements(from, to *models.AnalysisResult, comparison *models.BlueprintComparison) {
	fromMeasurements := make(map[string]models.Measurement)
	for _, measurement := range from.Measurements {
		key := measurement.MeasurementType
		if measurement.Location != nil {
			key = fmt.Sprintf("%s-%s", measurement.MeasurementType, *measurement.Location)
		}
		fromMeasurements[key] = measurement
	}

	toMeasurements := make(map[string]models.Measurement)
	for _, measurement := range to.Measurements {
		key := measurement.MeasurementType
		if measurement.Location != nil {
			key = fmt.Sprintf("%s-%s", measurement.MeasurementType, *measurement.Location)
		}
		toMeasurements[key] = measurement
	}

	// Compare measurements
	for key, toMeasurement := range toMeasurements {
		if fromMeasurement, exists := fromMeasurements[key]; exists {
			if fromMeasurement.Value != toMeasurement.Value {
				impact := "Medium"
				// Only check percentage if fromMeasurement.Value is not zero
				if fromMeasurement.Value > 0 && math.Abs(fromMeasurement.Value-toMeasurement.Value) > fromMeasurement.Value*0.2 {
					impact = "High"
				} else if fromMeasurement.Value == 0 && toMeasurement.Value > 0 {
					impact = "High"
				}
				comparison.Changes = append(comparison.Changes, models.BlueprintChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "measurement",
					Description: fmt.Sprintf("%s changed from %.2f %s to %.2f %s", toMeasurement.MeasurementType, fromMeasurement.Value, fromMeasurement.Unit, toMeasurement.Value, toMeasurement.Unit),
					OldValue:    fromMeasurement,
					NewValue:    toMeasurement,
					Impact:      &impact,
				})
			}
		} else {
			impact := "Low"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeAdded,
				Category:    "measurement",
				Description: fmt.Sprintf("%s added: %.2f %s", toMeasurement.MeasurementType, toMeasurement.Value, toMeasurement.Unit),
				NewValue:    toMeasurement,
				Impact:      &impact,
			})
		}
	}

	for key, fromMeasurement := range fromMeasurements {
		if _, exists := toMeasurements[key]; !exists {
			impact := "Medium"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeRemoved,
				Category:    "measurement",
				Description: fmt.Sprintf("%s removed, was: %.2f %s", fromMeasurement.MeasurementType, fromMeasurement.Value, fromMeasurement.Unit),
				OldValue:    fromMeasurement,
				Impact:      &impact,
			})
		}
	}
}

func (s *ComparisonService) compareMaterials(from, to *models.AnalysisResult, comparison *models.BlueprintComparison) {
	fromMaterials := make(map[string]models.Material)
	for _, material := range from.Materials {
		fromMaterials[material.MaterialName] = material
	}

	toMaterials := make(map[string]models.Material)
	for _, material := range to.Materials {
		toMaterials[material.MaterialName] = material
	}

	// Compare materials
	for name, toMaterial := range toMaterials {
		if fromMaterial, exists := fromMaterials[name]; exists {
			if fromMaterial.Quantity != toMaterial.Quantity {
				impact := "Medium"
				// Only check percentage if fromMaterial.Quantity is not zero
				if fromMaterial.Quantity > 0 && math.Abs(fromMaterial.Quantity-toMaterial.Quantity) > fromMaterial.Quantity*0.2 {
					impact = "High"
				} else if fromMaterial.Quantity == 0 && toMaterial.Quantity > 0 {
					impact = "High"
				}
				comparison.Changes = append(comparison.Changes, models.BlueprintChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "material",
					Description: fmt.Sprintf("%s quantity changed from %.2f %s to %.2f %s", name, fromMaterial.Quantity, fromMaterial.Unit, toMaterial.Quantity, toMaterial.Unit),
					OldValue:    fromMaterial,
					NewValue:    toMaterial,
					Impact:      &impact,
				})
			}
		} else {
			impact := "Medium"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeAdded,
				Category:    "material",
				Description: fmt.Sprintf("%s added: %.2f %s", name, toMaterial.Quantity, toMaterial.Unit),
				NewValue:    toMaterial,
				Impact:      &impact,
			})
		}
	}

	for name, fromMaterial := range fromMaterials {
		if _, exists := toMaterials[name]; !exists {
			impact := "Medium"
			comparison.Changes = append(comparison.Changes, models.BlueprintChange{
				ChangeType:  models.ChangeTypeRemoved,
				Category:    "material",
				Description: fmt.Sprintf("%s removed, was: %.2f %s", name, fromMaterial.Quantity, fromMaterial.Unit),
				OldValue:    fromMaterial,
				Impact:      &impact,
			})
		}
	}
}

func (s *ComparisonService) calculateSummary(comparison *models.BlueprintComparison) {
	comparison.Summary.TotalChanges = len(comparison.Changes)

	for _, change := range comparison.Changes {
		switch change.ChangeType {
		case models.ChangeTypeAdded:
			comparison.Summary.AddedCount++
		case models.ChangeTypeRemoved:
			comparison.Summary.RemovedCount++
		case models.ChangeTypeModified:
			comparison.Summary.ModifiedCount++
		}

		if change.Impact != nil && *change.Impact == "High" {
			comparison.Summary.HighImpactCount++
		}

		comparison.Summary.ChangesByCategory[change.Category]++
	}
}

// CompareBidRevisions compares two bid revisions and returns the differences
func (s *ComparisonService) CompareBidRevisions(from, to *models.BidRevision) (*models.BidComparison, error) {
	comparison := &models.BidComparison{
		FromVersion: from.Version,
		ToVersion:   to.Version,
		Changes:     []models.BidChange{},
		Summary: models.ComparisonSummary{
			ChangesByCategory: make(map[string]int),
		},
	}

	// Compare basic costs
	s.compareBidCosts(from, to, comparison)

	// Compare bid data if available
	if from.BidData != nil && to.BidData != nil {
		var fromBidData, toBidData models.GenerateBidResponse
		if err := json.Unmarshal([]byte(*from.BidData), &fromBidData); err == nil {
			if err := json.Unmarshal([]byte(*to.BidData), &toBidData); err == nil {
				s.compareBidLineItems(&fromBidData, &toBidData, comparison)
				s.compareBidTerms(&fromBidData, &toBidData, comparison)
			}
		}
	}

	// Calculate summary
	s.calculateBidSummary(comparison)

	return comparison, nil
}

func (s *ComparisonService) compareBidCosts(from, to *models.BidRevision, comparison *models.BidComparison) {
	// Compare total cost
	if from.TotalCost != nil && to.TotalCost != nil && *from.TotalCost != *to.TotalCost {
		impact := "High"
		diff := *to.TotalCost - *from.TotalCost
		var description string
		if *from.TotalCost > 0 {
			percentChange := (diff / *from.TotalCost) * 100
			description = fmt.Sprintf("Total cost changed from $%.2f to $%.2f (%.2f%%)", *from.TotalCost, *to.TotalCost, percentChange)
		} else {
			description = fmt.Sprintf("Total cost changed from $%.2f to $%.2f", *from.TotalCost, *to.TotalCost)
		}
		comparison.Changes = append(comparison.Changes, models.BidChange{
			ChangeType:  models.ChangeTypeModified,
			Category:    "cost",
			Description: description,
			OldValue:    *from.TotalCost,
			NewValue:    *to.TotalCost,
			Impact:      &impact,
		})
	}

	// Compare labor cost
	if from.LaborCost != nil && to.LaborCost != nil && *from.LaborCost != *to.LaborCost {
		impact := "Medium"
		diff := *to.LaborCost - *from.LaborCost
		var description string
		if *from.LaborCost > 0 {
			percentChange := (diff / *from.LaborCost) * 100
			description = fmt.Sprintf("Labor cost changed from $%.2f to $%.2f (%.2f%%)", *from.LaborCost, *to.LaborCost, percentChange)
		} else {
			description = fmt.Sprintf("Labor cost changed from $%.2f to $%.2f", *from.LaborCost, *to.LaborCost)
		}
		comparison.Changes = append(comparison.Changes, models.BidChange{
			ChangeType:  models.ChangeTypeModified,
			Category:    "cost",
			Description: description,
			OldValue:    *from.LaborCost,
			NewValue:    *to.LaborCost,
			Impact:      &impact,
		})
	}

	// Compare material cost
	if from.MaterialCost != nil && to.MaterialCost != nil && *from.MaterialCost != *to.MaterialCost {
		impact := "Medium"
		diff := *to.MaterialCost - *from.MaterialCost
		var description string
		if *from.MaterialCost > 0 {
			percentChange := (diff / *from.MaterialCost) * 100
			description = fmt.Sprintf("Material cost changed from $%.2f to $%.2f (%.2f%%)", *from.MaterialCost, *to.MaterialCost, percentChange)
		} else {
			description = fmt.Sprintf("Material cost changed from $%.2f to $%.2f", *from.MaterialCost, *to.MaterialCost)
		}
		comparison.Changes = append(comparison.Changes, models.BidChange{
			ChangeType:  models.ChangeTypeModified,
			Category:    "cost",
			Description: description,
			OldValue:    *from.MaterialCost,
			NewValue:    *to.MaterialCost,
			Impact:      &impact,
		})
	}

	// Compare markup percentage
	if from.MarkupPercentage != nil && to.MarkupPercentage != nil && *from.MarkupPercentage != *to.MarkupPercentage {
		impact := "Medium"
		comparison.Changes = append(comparison.Changes, models.BidChange{
			ChangeType:  models.ChangeTypeModified,
			Category:    "terms",
			Description: fmt.Sprintf("Markup percentage changed from %.2f%% to %.2f%%", *from.MarkupPercentage, *to.MarkupPercentage),
			OldValue:    *from.MarkupPercentage,
			NewValue:    *to.MarkupPercentage,
			Impact:      &impact,
		})
	}

	// Compare final price
	if from.FinalPrice != nil && to.FinalPrice != nil && *from.FinalPrice != *to.FinalPrice {
		impact := "High"
		diff := *to.FinalPrice - *from.FinalPrice
		var description string
		if *from.FinalPrice > 0 {
			percentChange := (diff / *from.FinalPrice) * 100
			description = fmt.Sprintf("Final price changed from $%.2f to $%.2f (%.2f%%)", *from.FinalPrice, *to.FinalPrice, percentChange)
		} else {
			description = fmt.Sprintf("Final price changed from $%.2f to $%.2f", *from.FinalPrice, *to.FinalPrice)
		}
		comparison.Changes = append(comparison.Changes, models.BidChange{
			ChangeType:  models.ChangeTypeModified,
			Category:    "cost",
			Description: description,
			OldValue:    *from.FinalPrice,
			NewValue:    *to.FinalPrice,
			Impact:      &impact,
		})
	}
}

func (s *ComparisonService) compareBidLineItems(from, to *models.GenerateBidResponse, comparison *models.BidComparison) {
	fromItems := make(map[string]models.LineItem)
	for _, item := range from.LineItems {
		key := fmt.Sprintf("%s-%s", item.Trade, item.Description)
		fromItems[key] = item
	}

	toItems := make(map[string]models.LineItem)
	for _, item := range to.LineItems {
		key := fmt.Sprintf("%s-%s", item.Trade, item.Description)
		toItems[key] = item
	}

	// Compare line items
	for key, toItem := range toItems {
		trade := toItem.Trade
		if fromItem, exists := fromItems[key]; exists {
			// Check for quantity changes
			if fromItem.Quantity != toItem.Quantity {
				impact := "Medium"
				comparison.Changes = append(comparison.Changes, models.BidChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "quantity",
					Trade:       &trade,
					Description: fmt.Sprintf("%s - %s: quantity changed from %.2f to %.2f %s", toItem.Trade, toItem.Description, fromItem.Quantity, toItem.Quantity, toItem.Unit),
					OldValue:    fromItem.Quantity,
					NewValue:    toItem.Quantity,
					Impact:      &impact,
				})
			}
			// Check for unit cost changes
			if fromItem.UnitCost != toItem.UnitCost {
				impact := "Low"
				comparison.Changes = append(comparison.Changes, models.BidChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "cost",
					Trade:       &trade,
					Description: fmt.Sprintf("%s - %s: unit cost changed from $%.2f to $%.2f", toItem.Trade, toItem.Description, fromItem.UnitCost, toItem.UnitCost),
					OldValue:    fromItem.UnitCost,
					NewValue:    toItem.UnitCost,
					Impact:      &impact,
				})
			}
			// Check for total changes
			if fromItem.Total != toItem.Total {
				impact := "Medium"
				comparison.Changes = append(comparison.Changes, models.BidChange{
					ChangeType:  models.ChangeTypeModified,
					Category:    "line_item",
					Trade:       &trade,
					Description: fmt.Sprintf("%s - %s: total changed from $%.2f to $%.2f", toItem.Trade, toItem.Description, fromItem.Total, toItem.Total),
					OldValue:    fromItem.Total,
					NewValue:    toItem.Total,
					Impact:      &impact,
				})
			}
		} else {
			impact := "Medium"
			comparison.Changes = append(comparison.Changes, models.BidChange{
				ChangeType:  models.ChangeTypeAdded,
				Category:    "line_item",
				Trade:       &trade,
				Description: fmt.Sprintf("%s - %s added: %.2f %s @ $%.2f = $%.2f", toItem.Trade, toItem.Description, toItem.Quantity, toItem.Unit, toItem.UnitCost, toItem.Total),
				NewValue:    toItem,
				Impact:      &impact,
			})
		}
	}

	for key, fromItem := range fromItems {
		if _, exists := toItems[key]; !exists {
			trade := fromItem.Trade
			impact := "High"
			comparison.Changes = append(comparison.Changes, models.BidChange{
				ChangeType:  models.ChangeTypeRemoved,
				Category:    "line_item",
				Trade:       &trade,
				Description: fmt.Sprintf("%s - %s removed: was %.2f %s @ $%.2f = $%.2f", fromItem.Trade, fromItem.Description, fromItem.Quantity, fromItem.Unit, fromItem.UnitCost, fromItem.Total),
				OldValue:    fromItem,
				Impact:      &impact,
			})
		}
	}
}

func (s *ComparisonService) compareBidTerms(from, to *models.GenerateBidResponse, comparison *models.BidComparison) {
	// Compare payment terms
	if from.PaymentTerms != to.PaymentTerms {
		impact := "Medium"
		comparison.Changes = append(comparison.Changes, models.BidChange{
			ChangeType:  models.ChangeTypeModified,
			Category:    "terms",
			Description: "Payment terms changed",
			OldValue:    from.PaymentTerms,
			NewValue:    to.PaymentTerms,
			Impact:      &impact,
		})
	}

	// Compare warranty terms
	if from.WarrantyTerms != to.WarrantyTerms {
		impact := "Low"
		comparison.Changes = append(comparison.Changes, models.BidChange{
			ChangeType:  models.ChangeTypeModified,
			Category:    "terms",
			Description: "Warranty terms changed",
			OldValue:    from.WarrantyTerms,
			NewValue:    to.WarrantyTerms,
			Impact:      &impact,
		})
	}

	// Compare scope changes (inclusions/exclusions)
	fromInclusions := make(map[string]bool)
	for _, inc := range from.Inclusions {
		fromInclusions[inc] = true
	}
	toInclusions := make(map[string]bool)
	for _, inc := range to.Inclusions {
		toInclusions[inc] = true
	}

	for inc := range toInclusions {
		if !fromInclusions[inc] {
			impact := "Low"
			comparison.Changes = append(comparison.Changes, models.BidChange{
				ChangeType:  models.ChangeTypeAdded,
				Category:    "scope",
				Description: fmt.Sprintf("Inclusion added: %s", inc),
				NewValue:    inc,
				Impact:      &impact,
			})
		}
	}

	for inc := range fromInclusions {
		if !toInclusions[inc] {
			impact := "Medium"
			comparison.Changes = append(comparison.Changes, models.BidChange{
				ChangeType:  models.ChangeTypeRemoved,
				Category:    "scope",
				Description: fmt.Sprintf("Inclusion removed: %s", inc),
				OldValue:    inc,
				Impact:      &impact,
			})
		}
	}
}

func (s *ComparisonService) calculateBidSummary(comparison *models.BidComparison) {
	comparison.Summary.TotalChanges = len(comparison.Changes)

	for _, change := range comparison.Changes {
		switch change.ChangeType {
		case models.ChangeTypeAdded:
			comparison.Summary.AddedCount++
		case models.ChangeTypeRemoved:
			comparison.Summary.RemovedCount++
		case models.ChangeTypeModified:
			comparison.Summary.ModifiedCount++
		}

		if change.Impact != nil && *change.Impact == "High" {
			comparison.Summary.HighImpactCount++
		}

		comparison.Summary.ChangesByCategory[change.Category]++
	}
}
