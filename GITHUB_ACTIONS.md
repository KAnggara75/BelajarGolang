# GitHub Actions CI/CD

This project uses GitHub Actions for automated testing and continuous integration.

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Purpose**: Fast feedback on every push and pull request

**Triggers**:
- Push to `main`, `master`, or `develop` branches
- Pull requests to `main`, `master`, or `develop` branches

**Features**:
- **Matrix Testing**: Tests against Go 1.21 and 1.22
- **PostgreSQL Service**: Spins up PostgreSQL 15 for integration tests
- **Race Detection**: Runs tests with `-race` flag to detect race conditions
- **Build Verification**: Ensures the application builds successfully

**Steps**:
1. Checkout code
2. Set up Go (matrix: 1.21, 1.22)
3. Install dependencies
4. Run tests with race detection
5. Build application

### 2. Test Workflow (`.github/workflows/test.yml`)

**Purpose**: Comprehensive testing with coverage and code quality checks

**Triggers**:
- Push to `main`, `master`, or `develop` branches
- Pull requests to `main`, `master`, or `develop` branches

**Features**:
- **Full Test Suite**: Runs all tests with coverage reporting
- **Code Coverage**: Generates coverage report and uploads to Codecov
- **Static Analysis**: Runs `go vet` for potential issues
- **Linting**: Runs golangci-lint for code quality
- **PostgreSQL Service**: Integration testing with real database

**Jobs**:

#### Test Job
1. Checkout code
2. Set up Go 1.21
3. Download and verify dependencies
4. Run `go vet`
5. Run tests with race detection and coverage
6. Generate coverage report
7. Upload coverage to Codecov
8. Build application

#### Lint Job
1. Checkout code
2. Set up Go 1.21
3. Run golangci-lint with custom configuration

## PostgreSQL Service

Both workflows use PostgreSQL as a service container:

```yaml
services:
  postgres:
    image: postgres:15
    env:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: testdb
    options: >-
      --health-cmd pg_isready
      --health-interval 10s
      --health-timeout 5s
      --health-retries 5
    ports:
      - 5432:5432
```

**Connection String**:
```
postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable
```

## Linting Configuration

The project uses golangci-lint with a custom configuration (`.golangci.yml`):

**Enabled Linters**:
- `errcheck` - Check for unchecked errors
- `gosimple` - Simplify code
- `govet` - Vet examines Go source code
- `ineffassign` - Detect ineffectual assignments
- `staticcheck` - Advanced static analysis
- `typecheck` - Type checking
- `unused` - Check for unused code
- `gosec` - Security checks
- `gofmt` - Format checking
- `goimports` - Import organization
- And many more...

**Excluded from Linting**:
- Test files (for certain linters like `gomnd`, `funlen`)
- Generated files (`*.pb.go`)
- Vendor directory
- Tmp directory

## Running Locally

### Run Tests
```bash
# All tests
go test ./...

# With race detection
go test -race ./...

# With coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# View coverage in browser
go tool cover -html=coverage.out
```

### Run Linter
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run with auto-fix
golangci-lint run --fix
```

### Run Static Analysis
```bash
go vet ./...
```

## Badges

Add these badges to your README to show build status:

```markdown
[![CI](https://github.com/KAnggara75/BelajarGolang/actions/workflows/ci.yml/badge.svg)](https://github.com/KAnggara75/BelajarGolang/actions/workflows/ci.yml)
[![Go Tests](https://github.com/KAnggara75/BelajarGolang/actions/workflows/test.yml/badge.svg)](https://github.com/KAnggara75/BelajarGolang/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/KAnggara75/BelajarGolang)](https://goreportcard.com/report/github.com/KAnggara75/BelajarGolang)
```

## Codecov Integration

The test workflow uploads coverage reports to Codecov. To enable:

1. Sign up at [codecov.io](https://codecov.io)
2. Connect your GitHub repository
3. No additional configuration needed - the workflow handles it

**Optional**: Add Codecov badge to README:
```markdown
[![codecov](https://codecov.io/gh/KAnggara75/BelajarGolang/branch/main/graph/badge.svg)](https://codecov.io/gh/KAnggara75/BelajarGolang)
```

## Workflow Status

You can view workflow runs at:
```
https://github.com/KAnggara75/BelajarGolang/actions
```

## Troubleshooting

### Tests Fail in CI but Pass Locally

**Possible causes**:
1. **Database connection**: Ensure `DATABASE_URL` is set correctly in the workflow
2. **Race conditions**: CI runs with `-race` flag which may catch issues not visible locally
3. **Go version**: CI tests against multiple Go versions

**Solution**: Run tests locally with the same flags:
```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable \
  go test -v -race ./...
```

### Linter Fails in CI

**Possible causes**:
1. Different golangci-lint version
2. Configuration differences

**Solution**: Run linter locally with the same configuration:
```bash
golangci-lint run --timeout=5m
```

### PostgreSQL Service Not Ready

The workflow includes health checks to ensure PostgreSQL is ready before running tests:
```yaml
options: >-
  --health-cmd pg_isready
  --health-interval 10s
  --health-timeout 5s
  --health-retries 5
```

If tests still fail, increase the retry count or interval.

## Best Practices

1. **Always run tests locally** before pushing
2. **Check linter warnings** and fix them
3. **Keep coverage high** - aim for >80%
4. **Write meaningful commit messages** for better CI logs
5. **Review failed workflows** and fix issues promptly
6. **Keep dependencies updated** with `go mod tidy`

## Future Enhancements

Potential improvements to the CI/CD pipeline:

- [ ] Add deployment workflow for staging/production
- [ ] Add security scanning (e.g., Snyk, Trivy)
- [ ] Add performance benchmarking
- [ ] Add Docker image building and publishing
- [ ] Add automated release creation
- [ ] Add dependency update automation (Dependabot)
