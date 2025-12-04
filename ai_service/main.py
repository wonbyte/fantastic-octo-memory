import os

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(
    title="Construction Estimation AI Service",
    description="AI/ML service for construction cost estimation and bidding automation",
    version="1.0.0"
)

# CORS configuration
origins = [
    "http://localhost:3000",
    "http://localhost:8080",
    "http://localhost:19006",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {
        "message": "Construction Estimation AI Service",
        "version": "1.0.0",
        "status": "running"
    }

@app.get("/health")
async def health():
    return {
        "status": "healthy",
        "service": "ai_service"
    }

@app.post("/estimate")
async def estimate(data: dict):
    """
    Endpoint for cost estimation.
    This is a placeholder for actual ML model inference.
    """
    return {
        "estimated_cost": 100000.0,
        "confidence": 0.85,
        "breakdown": {
            "materials": 60000.0,
            "labor": 30000.0,
            "equipment": 10000.0
        }
    }

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(app, host="0.0.0.0", port=port)
