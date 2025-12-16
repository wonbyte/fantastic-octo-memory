# E2E Test Fixtures

This directory contains test fixtures used by Playwright E2E tests.

## Test Blueprint Files

For testing blueprint upload functionality, add sample PDF files here:

### Recommended Test Files

1. **test-blueprint.pdf** - A small, simple architectural blueprint (1-2 MB)
   - Should contain basic floor plans
   - Used for standard workflow tests

2. **test-blueprint-large.pdf** - A larger blueprint file (8-10 MB)
   - Used for testing file size limits
   - Tests upload progress indicators

3. **test-blueprint-multipage.pdf** - Multi-page blueprint
   - Used for testing multi-page document handling

## Creating Test Fixtures

You can create test PDF files using:

1. **Sample architectural plans** from public domain sources
2. **Generated PDFs** using tools like:
   - PDFKit (Node.js)
   - ReportLab (Python)
   - Online PDF generators

### Example: Generate Simple Test PDF

```javascript
// Using PDFKit in Node.js
const PDFDocument = require('pdfkit');
const fs = require('fs');

const doc = new PDFDocument();
doc.pipe(fs.createWriteStream('test-blueprint.pdf'));

doc.fontSize(25).text('Test Blueprint', 100, 100);
doc.text('This is a test architectural drawing', 100, 150);
doc.text('Room 1: 20\' x 15\' = 300 sq ft', 100, 200);
doc.text('Room 2: 15\' x 12\' = 180 sq ft', 100, 250);

doc.end();
```

## Usage in Tests

```typescript
import path from 'path';

// In your test
const testFilePath = path.join(__dirname, 'fixtures', 'test-blueprint.pdf');
await fileInput.setInputFiles(testFilePath);
```

## Note

Test fixtures are not committed to the repository if they are large binary files. Add them locally for testing or generate them as needed.
