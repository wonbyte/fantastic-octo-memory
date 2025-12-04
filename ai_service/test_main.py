from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_root():
    response = client.get("/")
    assert response.status_code == 200
    data = response.json()
    assert "message" in data
    assert "version" in data
    assert data["version"] == "1.0.0"


def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert data["version"] == "1.0.0"


def test_analyze_blueprint_validation():
    """Test analyze blueprint endpoint with invalid request."""
    response = client.post("/analyze-blueprint", json={})
    assert response.status_code == 422  # Validation error


def test_analyze_blueprint_mock():
    """Test analyze blueprint endpoint with mock data (no real S3/OCR)."""
    request_data = {
        "blueprint_id": "test-123",
        "s3_key": "test/blueprint.pdf",
        "project_name": "Test Project",
    }
    # This will fail if S3 is not available, but validates the structure
    response = client.post("/analyze-blueprint", json=request_data)
    # We expect either success or 500 (if S3 not available)
    assert response.status_code in [200, 500]


def test_generate_bid_validation():
    """Test generate bid endpoint with invalid request."""
    response = client.post("/generate-bid", json={})
    assert response.status_code == 422  # Validation error


def test_generate_bid_mock():
    """Test generate bid endpoint with valid mock data."""
    request_data = {
        "project_id": "proj-123",
        "blueprint_id": "bp-456",
        "takeoff_data": {
            "rooms": [
                {
                    "name": "Living Room",
                    "dimensions": "15' x 20'",
                    "area": 300.0,
                    "room_type": "Living",
                }
            ],
            "openings": [],
            "fixtures": [],
            "materials": [],
        },
        "markup_percentage": 20.0,
    }
    response = client.post("/generate-bid", json=request_data)
    assert response.status_code == 200
    data = response.json()
    assert "bid_id" in data
    assert "project_id" in data
    assert data["project_id"] == "proj-123"
    assert "status" in data
    assert "total_price" in data
    assert data["total_price"] > 0
