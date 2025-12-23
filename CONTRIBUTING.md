# Contributing to Construction Estimation & Bidding Platform

Thank you for your interest in contributing! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Be respectful, inclusive, and professional in all interactions.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/fantastic-octo-memory.git
   cd fantastic-octo-memory
   ```
3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/wonbyte/fantastic-octo-memory.git
   ```
4. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- **Docker** with Docker Compose V2
- **Make** (for convenience commands)
- **Node.js** 22 LTS (for local frontend development)
- **Go** 1.25+ (for local backend development)
- **Python** 3.12+ (for local AI service development)

### Quick Setup

```bash
# Copy environment files
cp .env.example .env

# Start all services in development mode
make dev
```

### Service-Specific Setup

#### Backend (Go)
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### AI Service (Python)
```bash
cd ai_service
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install -r requirements.txt -r requirements-dev.txt
uvicorn app.main:app --reload
```

#### Frontend (React Native)
```bash
cd app
npm install
npm start
```

## How to Contribute

### Types of Contributions

- 🐛 **Bug Reports**: Use GitHub Issues with the bug template
- ✨ **Feature Requests**: Use GitHub Issues with the feature template
- 📝 **Documentation**: Improvements to docs, README, code comments
- 🧪 **Tests**: Adding or improving test coverage
- 🔧 **Code**: Bug fixes, features, refactoring

### Before You Start

1. **Check existing issues** to avoid duplicates
2. **Discuss major changes** by opening an issue first
3. **Keep changes focused** - one feature/fix per PR
4. **Update documentation** for any user-facing changes

## Coding Standards

### General Principles

- Write clean, readable, maintainable code
- Follow the existing code style and patterns
- Keep functions small and focused (< 50 lines ideally)
- Use meaningful variable and function names
- Add comments for complex logic only
- Avoid premature optimization

### Backend (Go)

**Style Guide:**
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `go vet` before committing
- Handle all errors explicitly
- Use structured logging (`slog`)

**Example:**
```go
// Good: Clear function with error handling
func GetProject(ctx context.Context, id string) (*Project, error) {
    if id == "" {
        return nil, ErrInvalidID
    }
    
    project, err := repo.Find(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("finding project: %w", err)
    }
    
    return project, nil
}
```

**Linting:**
```bash
cd backend
go vet ./...
go test ./...
```

### AI Service (Python)

**Style Guide:**
- Follow [PEP 8](https://pep8.org/)
- Use type hints for function signatures
- Docstrings for public functions
- Line length: 100 characters max
- Use Ruff for linting and formatting

**Example:**
```python
async def analyze_blueprint(
    blueprint_id: str,
    options: AnalysisOptions
) -> AnalysisResult:
    """Analyze a blueprint using AI vision models.
    
    Args:
        blueprint_id: Unique identifier for the blueprint
        options: Configuration options for analysis
        
    Returns:
        Analysis results with detected elements
        
    Raises:
        ValueError: If blueprint_id is invalid
        AnalysisError: If analysis fails
    """
    if not blueprint_id:
        raise ValueError("blueprint_id is required")
    
    # ... implementation
```

**Linting:**
```bash
cd ai_service
ruff check .
ruff format .
pytest
```

### Frontend (React Native)

**Style Guide:**
- Use TypeScript for type safety
- Functional components with hooks
- Props interface for all components
- ESLint configuration enforced
- Consistent file naming (PascalCase for components)

**Example:**
```typescript
interface ProjectCardProps {
  project: Project;
  onPress: (id: string) => void;
}

export function ProjectCard({ project, onPress }: ProjectCardProps) {
  const handlePress = () => {
    onPress(project.id);
  };

  return (
    <Pressable onPress={handlePress}>
      <Text>{project.name}</Text>
    </Pressable>
  );
}
```

**Linting:**
```bash
cd app
npm run lint
npm run type-check
npm test
```

## Commit Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic changes)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates
- `perf`: Performance improvements
- `ci`: CI/CD changes

### Examples

```bash
feat(backend): add bid export to CSV format

Implements CSV export functionality for bids with all
pricing details and line items.

Closes #123

---

fix(ai-service): handle empty blueprint uploads

Added validation to check for empty files before processing.
Returns 400 Bad Request with helpful error message.

---

docs(readme): update installation instructions

- Added Node.js version requirement
- Clarified Docker Compose V2 usage
- Fixed broken links
```

### Rules

- Use present tense ("add feature" not "added feature")
- Use imperative mood ("move cursor" not "moves cursor")
- Capitalize first letter of subject
- No period at the end of subject line
- Limit subject line to 72 characters
- Separate subject from body with blank line
- Reference issues and PRs in footer

## Pull Request Process

### Before Submitting

1. **Ensure all tests pass:**
   ```bash
   make test
   ```

2. **Run linters:**
   ```bash
   # Backend
   cd backend && go vet ./...
   
   # AI Service
   cd ai_service && ruff check .
   
   # Frontend
   cd app && npm run lint && npm run type-check
   ```

3. **Update documentation** if needed

4. **Add tests** for new functionality

5. **Rebase on latest main:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### PR Guidelines

1. **Use the PR template** when creating your PR

2. **Link related issues** using keywords:
   - `Closes #123`
   - `Fixes #456`
   - `Relates to #789`

3. **Provide clear description:**
   - What changes were made
   - Why the changes were necessary
   - How to test the changes

4. **Keep PRs focused:**
   - One feature or fix per PR
   - Avoid mixing refactoring with features
   - Split large changes into multiple PRs

5. **Request reviews** from relevant maintainers

6. **Respond to feedback** promptly and professionally

### PR Review Process

1. Automated checks must pass (CI/CD)
2. Code review from at least one maintainer
3. All conversations must be resolved
4. Final approval from maintainer
5. Squash and merge to main

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] All tests pass locally
- [ ] New tests added for new functionality
- [ ] Documentation updated
- [ ] Commit messages follow convention
- [ ] No merge conflicts with main
- [ ] PR description is clear and complete

## Testing

### Test Coverage Requirements

- **Minimum 80% coverage** for new code
- **All new features** must have tests
- **Bug fixes** should include regression tests

### Running Tests

```bash
# All services
make test

# Backend only
cd backend && go test -v -race -coverprofile=coverage.out ./...

# AI Service only
cd ai_service && pytest --cov=. --cov-report=term-missing

# Frontend only
cd app && npm test -- --coverage
```

### Writing Tests

#### Backend (Go)
```go
func TestGetProject(t *testing.T) {
    // Arrange
    repo := &MockRepository{}
    service := NewProjectService(repo)
    
    // Act
    project, err := service.GetProject(context.Background(), "123")
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, project)
}
```

#### AI Service (Python)
```python
@pytest.mark.asyncio
async def test_analyze_blueprint():
    # Arrange
    service = VisionService()
    
    # Act
    result = await service.analyze_blueprint("test-id")
    
    # Assert
    assert result is not None
    assert result.elements_count > 0
```

#### Frontend (React Native)
```typescript
describe('ProjectCard', () => {
  it('renders project name', () => {
    const project = { id: '1', name: 'Test Project' };
    const { getByText } = render(<ProjectCard project={project} onPress={jest.fn()} />);
    
    expect(getByText('Test Project')).toBeTruthy();
  });
});
```

## Documentation

### Code Documentation

- **Public APIs**: Must have documentation
- **Complex logic**: Add explanatory comments
- **Exported functions**: Include docstrings/comments
- **Configuration**: Document all options

### README Updates

Update relevant README files for:
- New features
- Changed behavior
- New dependencies
- Setup instructions

### API Documentation

- Update API docs for endpoint changes
- Include request/response examples
- Document error responses
- Note breaking changes

## Questions?

- 📫 Open an issue for questions
- 💬 Join our discussions (if applicable)
- 📖 Check existing documentation

## Recognition

Contributors will be recognized in:
- CHANGELOG.md for significant contributions
- GitHub contributors list
- Release notes

Thank you for contributing! 🎉
