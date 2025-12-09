# Revision Comparison Feature - Demo Guide

This guide walks through testing the revision comparison feature manually.

## Prerequisites

1. Backend and frontend services running
2. Test user and project created (run `backend/seed.sh`)
3. At least one blueprint uploaded and analyzed

## Quick Start

### 1. Set Up Test Environment

```bash
# Start services
make dev

# In another terminal, seed test data
cd backend && ./seed.sh
```

### 2. Create Initial Blueprint Revision

1. Log in to the app with test credentials:
   - Email: `test@example.com`
   - Password: `testpassword123`

2. Navigate to your test project

3. Upload a blueprint (PDF, PNG, or CAD file)

4. Analyze the blueprint to create analysis data

5. Once analysis is complete, go to the blueprint detail page

6. The system automatically tracks this as Version 1

### 3. Create a Second Revision

**Option A: Using the API (Simulated Update)**

```bash
# Get the blueprint ID from the UI or database
BLUEPRINT_ID="your-blueprint-id-here"

# Create a revision snapshot
curl -X POST http://localhost:8080/blueprints/$BLUEPRINT_ID/revisions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

**Option B: Upload a Modified Blueprint**

1. Upload a new version of the same blueprint with changes
2. The system will track it as Version 2

### 4. Test Revision Comparison

1. **View Revision History**
   - Go to Blueprint Detail screen
   - Click "View Revisions" button
   - You should see a list of all revisions with version numbers and timestamps

2. **Compare Two Versions**
   - Select Version 1 (click on it)
   - Select Version 2 (click on it)
   - Click "Compare" button
   - The comparison view should display:
     - Summary statistics (total changes, added, removed, modified)
     - Detailed list of changes grouped by category
     - Impact levels (High, Medium, Low) with color coding
     - Change icons (+ for added, - for removed, ~ for modified)

3. **Navigate Back**
   - Click "← Back to Revisions" to return to revision history
   - Click "Hide Revisions" to collapse the revision section

## Expected Results

### Revision History View

```
Revision History
┌─────────────────────────────────────────┐
│ Version 2 - Dec 9, 2025 5:10 PM       │
│ Size: 2.5 MB                            │
│ Has Changes: 5 changes from v1          │
├─────────────────────────────────────────┤
│ Version 1 - Dec 9, 2025 4:30 PM       │
│ Size: 2.3 MB                            │
│ Latest                                   │
└─────────────────────────────────────────┘
```

### Comparison View

```
Comparing Version 1 → Version 2
┌─────────────────────────────────────────┐
│ Summary                                 │
│ • Total Changes: 5                      │
│ • Added: 2                               │
│ • Removed: 1                             │
│ • Modified: 2                            │
│ • High Impact: 1                         │
└─────────────────────────────────────────┘

Changes by Category:
┌─────────────────────────────────────────┐
│ Rooms                                   │
│ [~] Living Room dimensions changed      │
│     from 20x15 (300 SF) to 25x15 (375 SF)│
│     Impact: High                         │
├─────────────────────────────────────────┤
│ Materials                               │
│ [+] 2x6 Lumber added: 200.00 LF        │
│     Impact: Medium                       │
│ [-] 2x4 Lumber removed: 150.00 LF      │
│     Impact: Low                          │
└─────────────────────────────────────────┘
```

## Testing API Endpoints Directly

### List Blueprint Revisions

```bash
curl -X GET http://localhost:8080/blueprints/$BLUEPRINT_ID/revisions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Expected Response:
```json
[
  {
    "id": "uuid",
    "blueprint_id": "uuid",
    "version": 2,
    "filename": "blueprint_v2.pdf",
    "file_size": 2500000,
    "created_at": "2025-12-09T17:10:00Z",
    "changes_summary": {
      "total_changes": 5,
      "added_count": 2,
      "removed_count": 1,
      "modified_count": 2
    }
  },
  {
    "id": "uuid",
    "blueprint_id": "uuid",
    "version": 1,
    "filename": "blueprint_v1.pdf",
    "file_size": 2300000,
    "created_at": "2025-12-09T16:30:00Z"
  }
]
```

### Compare Revisions

```bash
curl -X GET "http://localhost:8080/blueprints/$BLUEPRINT_ID/compare?from=1&to=2" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Expected Response:
```json
{
  "from_version": 1,
  "to_version": 2,
  "changes": [
    {
      "change_type": "modified",
      "category": "room",
      "description": "Room 'Living Room' dimensions changed from 20x15 (300.00 SF) to 25x15 (375.00 SF)",
      "old_value": {
        "name": "Living Room",
        "dimensions": "20x15",
        "area": 300.0
      },
      "new_value": {
        "name": "Living Room",
        "dimensions": "25x15",
        "area": 375.0
      },
      "impact": "High"
    }
  ],
  "summary": {
    "total_changes": 5,
    "added_count": 2,
    "removed_count": 1,
    "modified_count": 2,
    "high_impact_count": 1,
    "changes_by_category": {
      "room": 1,
      "material": 2
    }
  }
}
```

## Troubleshooting

### No Revisions Showing

1. Ensure blueprint has been analyzed (analysis_status = 'completed')
2. Check that analysis_data is not null in the database
3. Create a revision snapshot using the POST /blueprints/{id}/revisions endpoint

### Cannot Compare Revisions

1. Verify you have at least 2 revisions
2. Check that both revisions have analysis_data
3. Ensure you're selecting different version numbers

### Changes Not Detected

1. Verify the analysis_data differs between revisions
2. Check backend logs for comparison service errors
3. Ensure the comparison service is properly handling null/missing data

## Database Queries for Debugging

### Check Blueprint Revisions

```sql
SELECT 
  id, 
  version, 
  filename, 
  created_at,
  analysis_data IS NOT NULL as has_analysis
FROM blueprint_revisions
WHERE blueprint_id = 'your-blueprint-id'
ORDER BY version DESC;
```

### Check Blueprint Version Info

```sql
SELECT 
  id, 
  filename,
  version,
  is_latest,
  parent_blueprint_id,
  analysis_data IS NOT NULL as has_analysis
FROM blueprints
WHERE id = 'your-blueprint-id';
```

## Demo Scenarios

### Scenario 1: Room Dimension Change

1. Create blueprint with room: "Kitchen 15x12 (180 SF)"
2. Create revision
3. Update to "Kitchen 18x12 (216 SF)"
4. Compare - should show High impact modification

### Scenario 2: Material Addition

1. Create blueprint with material list
2. Create revision
3. Add new material "2x6 Lumber: 200 LF"
4. Compare - should show Medium impact addition

### Scenario 3: Opening Removal

1. Create blueprint with door/window
2. Create revision
3. Remove the opening
4. Compare - should show removal with impact level

## Next Steps

1. Test with real blueprints (PDF files with actual construction drawings)
2. Integrate with AI analysis service for automatic change detection
3. Add email notifications for high-impact changes
4. Create export to PDF functionality for comparison reports

## Related Documentation

- [REVISION_COMPARISON_IMPLEMENTATION.md](../REVISION_COMPARISON_IMPLEMENTATION.md) - Full technical implementation details
- [E2E Tests](../e2e/revision-comparison.spec.ts) - Automated test scenarios
- [API Documentation](../backend/API_TEST_GUIDE.md) - API testing guide
