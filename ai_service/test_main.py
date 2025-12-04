from fastapi.testclient import TestClient
from main import app

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
    assert data["status"] == "healthy"
    assert data["service"] == "ai_service"

def test_estimate():
    response = client.post("/estimate", json={})
    assert response.status_code == 200
    data = response.json()
    assert "estimated_cost" in data
    assert "confidence" in data
    assert "breakdown" in data
