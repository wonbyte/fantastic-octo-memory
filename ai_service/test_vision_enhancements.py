"""Tests for enhanced vision service features."""

import pytest

from app.services.vision_service import VisionService


@pytest.fixture
def vision_service():
    """Create a vision service instance."""
    return VisionService()


class TestScaleDetection:
    """Test scale detection functionality."""

    def test_detect_scale_quarter_inch(self, vision_service):
        """Test detection of 1/4\" = 1'-0\" scale."""
        ocr_text = "Blueprint Floor Plan\nSCALE: 1/4\" = 1'-0\"\nLiving Room"
        scale_info = vision_service._detect_scale(ocr_text)

        assert scale_info is not None
        assert "1/4" in scale_info["scale_string"]
        assert scale_info["confidence"] > 0

    def test_detect_scale_eighth_inch(self, vision_service):
        """Test detection of 1/8\" = 1'-0\" scale."""
        ocr_text = "SCALE 1/8\" = 1'-0\""
        scale_info = vision_service._detect_scale(ocr_text)

        assert scale_info is not None
        assert "1/8" in scale_info["scale_string"]

    def test_detect_scale_metric(self, vision_service):
        """Test detection of metric scale 1:100."""
        ocr_text = "SCALE: 1:100"
        scale_info = vision_service._detect_scale(ocr_text)

        assert scale_info is not None
        assert "1:100" in scale_info["scale_string"]

    def test_no_scale_detected(self, vision_service):
        """Test when no scale is present."""
        ocr_text = "Living Room\nBedroom\nKitchen"
        scale_info = vision_service._detect_scale(ocr_text)

        assert scale_info is None

    def test_detect_scale_empty_text(self, vision_service):
        """Test scale detection with empty text."""
        scale_info = vision_service._detect_scale("")
        assert scale_info is None


class TestTradeTypeDetection:
    """Test trade type detection functionality."""

    def test_detect_electrical_trade(self, vision_service):
        """Test detection of electrical trade type."""
        ocr_text = "Electrical Plan\nPanel Schedule\n120V outlets\nLighting fixtures\nCircuit breakers"
        trade_type = vision_service._detect_trade_type(ocr_text, None)

        assert trade_type == "electrical"

    def test_detect_plumbing_trade(self, vision_service):
        """Test detection of plumbing trade type."""
        ocr_text = "Plumbing Plan\nWater supply lines\nDrainage system\nFixture schedule"
        trade_type = vision_service._detect_trade_type(ocr_text, None)

        assert trade_type == "plumbing"

    def test_detect_hvac_trade(self, vision_service):
        """Test detection of HVAC trade type."""
        ocr_text = "Mechanical Plan\nHVAC Layout\nDuctwork\nFurnace location\nAir conditioning"
        trade_type = vision_service._detect_trade_type(ocr_text, None)

        assert trade_type == "hvac"

    def test_detect_structural_trade(self, vision_service):
        """Test detection of structural trade type."""
        ocr_text = "Structural Plan\nFoundation details\nBeam schedule\nColumn layout\nFooting"
        trade_type = vision_service._detect_trade_type(ocr_text, None)

        assert trade_type == "structural"

    def test_detect_general_trade(self, vision_service):
        """Test detection defaults to general when no specific trade found."""
        ocr_text = "Floor Plan\nLiving Room\nBedroom\nKitchen"
        trade_type = vision_service._detect_trade_type(ocr_text, None)

        assert trade_type == "general"

    def test_trade_type_from_context(self, vision_service):
        """Test trade type specified in context takes precedence."""
        ocr_text = "Some random text"
        context = {"trade_type": "electrical"}
        trade_type = vision_service._detect_trade_type(ocr_text, context)

        assert trade_type == "electrical"

    def test_invalid_trade_type_in_context(self, vision_service):
        """Test invalid trade type in context is ignored."""
        ocr_text = "Electrical plan with outlets"
        context = {"trade_type": "invalid_trade"}
        trade_type = vision_service._detect_trade_type(ocr_text, context)

        # Should detect from OCR text instead
        assert trade_type == "electrical"


class TestTradePromptSelection:
    """Test trade-specific prompt and system message selection."""

    def test_get_electrical_prompt(self, vision_service):
        """Test getting electrical trade prompt."""
        prompt, system = vision_service._get_trade_prompt_and_system("electrical")

        assert "electrician" in system.lower()
        assert "outlets" in prompt.lower()
        assert "switches" in prompt.lower()

    def test_get_plumbing_prompt(self, vision_service):
        """Test getting plumbing trade prompt."""
        prompt, system = vision_service._get_trade_prompt_and_system("plumbing")

        assert "plumber" in system.lower()
        assert "fixtures" in prompt.lower()
        assert "drainage" in prompt.lower()

    def test_get_hvac_prompt(self, vision_service):
        """Test getting HVAC trade prompt."""
        prompt, system = vision_service._get_trade_prompt_and_system("hvac")

        assert "hvac" in system.lower()
        assert "ductwork" in prompt.lower()

    def test_get_structural_prompt(self, vision_service):
        """Test getting structural trade prompt."""
        prompt, system = vision_service._get_trade_prompt_and_system("structural")

        assert "structural engineer" in system.lower()
        assert "beam" in prompt.lower()

    def test_get_general_prompt(self, vision_service):
        """Test getting general trade prompt."""
        prompt, system = vision_service._get_trade_prompt_and_system("general")

        assert "construction estimator" in system.lower()

    def test_get_unknown_trade_defaults_to_general(self, vision_service):
        """Test unknown trade type defaults to general."""
        prompt, system = vision_service._get_trade_prompt_and_system("unknown")

        assert "construction estimator" in system.lower()


class TestSymbolLibrary:
    """Test symbol library retrieval."""

    def test_get_electrical_symbols(self, vision_service):
        """Test getting electrical symbol library."""
        symbols = vision_service._get_symbol_library("electrical")

        assert "outlets" in symbols
        assert "switches" in symbols
        assert "lighting" in symbols
        assert len(symbols["outlets"]) > 0

    def test_get_plumbing_symbols(self, vision_service):
        """Test getting plumbing symbol library."""
        symbols = vision_service._get_symbol_library("plumbing")

        assert "fixtures" in symbols
        assert "supply" in symbols
        assert "drainage" in symbols

    def test_get_hvac_symbols(self, vision_service):
        """Test getting HVAC symbol library."""
        symbols = vision_service._get_symbol_library("hvac")

        assert "equipment" in symbols
        assert "ductwork" in symbols
        assert "ventilation" in symbols

    def test_get_structural_symbols(self, vision_service):
        """Test getting structural symbol library."""
        symbols = vision_service._get_symbol_library("structural")

        assert "foundation" in symbols
        assert "framing" in symbols
        assert "vertical" in symbols

    def test_get_general_symbols_empty(self, vision_service):
        """Test general trade returns empty symbol library."""
        symbols = vision_service._get_symbol_library("general")

        assert symbols == {}


class TestMockResponses:
    """Test mock response generation for different trade types."""

    def test_mock_response_electrical(self, vision_service):
        """Test mock response for electrical trade."""
        response = vision_service._mock_vision_response("electrical")

        assert "fixtures" in response
        assert len(response["fixtures"]) > 0
        assert response["trade_type"] == "electrical"
        assert any("outlet" in f["fixture_type"].lower() for f in response["fixtures"])

    def test_mock_response_plumbing(self, vision_service):
        """Test mock response for plumbing trade."""
        response = vision_service._mock_vision_response("plumbing")

        assert "fixtures" in response
        assert response["trade_type"] == "plumbing"
        assert any("plumbing" in f["category"].lower() for f in response["fixtures"])

    def test_mock_response_hvac(self, vision_service):
        """Test mock response for HVAC trade."""
        response = vision_service._mock_vision_response("hvac")

        assert "fixtures" in response
        assert response["trade_type"] == "hvac"
        assert any("hvac" in f["category"].lower() for f in response["fixtures"])

    def test_mock_response_structural(self, vision_service):
        """Test mock response for structural trade."""
        response = vision_service._mock_vision_response("structural")

        assert response["trade_type"] == "structural"
        assert len(response["measurements"]) > 2  # Should have extra structural measurements

    def test_mock_response_general(self, vision_service):
        """Test mock response for general trade."""
        response = vision_service._mock_vision_response("general")

        assert "rooms" in response
        assert "fixtures" in response
        assert "confidence_score" in response
        assert response["confidence_score"] > 0

    def test_mock_response_includes_scale_info(self, vision_service):
        """Test mock response includes scale info."""
        response = vision_service._mock_vision_response("general")

        assert "scale_info" in response
        assert response["scale_info"] is not None


@pytest.mark.asyncio
class TestBlueprintAnalysis:
    """Test blueprint analysis with enhancements."""

    async def test_analyze_blueprint_with_scale(self, vision_service):
        """Test blueprint analysis detects scale information."""
        image_bytes = b"fake_image_data"
        ocr_text = "Floor Plan\nSCALE: 1/4\" = 1'-0\"\nLiving Room"

        analysis = await vision_service.analyze_blueprint(image_bytes, ocr_text)

        assert analysis.scale_info is not None
        assert analysis.confidence_score > 0

    async def test_analyze_blueprint_electrical_trade(self, vision_service):
        """Test blueprint analysis with electrical trade type."""
        image_bytes = b"fake_image_data"
        ocr_text = "Electrical Plan\nPanel Schedule\nOutlets and switches"

        analysis = await vision_service.analyze_blueprint(image_bytes, ocr_text)

        assert analysis.trade_type == "electrical"
        assert len(analysis.fixtures) > 0

    async def test_analyze_blueprint_with_context_trade(self, vision_service):
        """Test blueprint analysis with trade type in context."""
        image_bytes = b"fake_image_data"
        ocr_text = "Some text"
        context = {"trade_type": "plumbing"}

        analysis = await vision_service.analyze_blueprint(image_bytes, ocr_text, context)

        assert analysis.trade_type == "plumbing"

    async def test_multi_page_analysis(self, vision_service):
        """Test multi-page blueprint analysis."""
        pages = [
            (b"page1_data", "Floor Plan\nLiving Room"),
            (b"page2_data", "Electrical Plan\nPanel Schedule"),
            (b"page3_data", "Plumbing Plan\nFixture Layout"),
        ]

        analysis = await vision_service.analyze_multi_page_blueprint(pages)

        assert analysis.confidence_score > 0
        assert len(analysis.rooms) > 0
        # Should aggregate results from all pages

    async def test_multi_page_analysis_captures_first_page_info(self, vision_service):
        """Test multi-page analysis captures scale and trade from first page."""
        pages = [
            (b"page1_data", "SCALE: 1/4\" = 1'-0\"\nElectrical Plan with outlets and switches"),
            (b"page2_data", "Additional details"),
        ]

        analysis = await vision_service.analyze_multi_page_blueprint(pages)

        assert analysis.scale_info is not None
        # Trade type should be detected if keywords are present
        assert analysis.trade_type is not None or analysis.confidence_score > 0
