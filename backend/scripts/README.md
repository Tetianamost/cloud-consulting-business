# Build and Test Scripts

This directory contains scripts for building and testing the Cloud Consulting Backend.

## Build Scripts

### `build.sh` - Simple Build Script
Quick and easy way to build the main application.

```bash
./scripts/build.sh
```

**What it does:**
- Builds the main server application
- Shows binary size and location
- Creates `./server` executable

### `verify_build.sh` - Comprehensive Build Verification
Complete build verification with multiple tests and checks.

```bash
./scripts/verify_build.sh
```

**What it does:**
- Verifies Go module dependencies
- Tests main application build
- Checks internal package compilation
- Identifies common build issues
- Tests cross-compilation
- Measures build performance
- Provides detailed status report

## Test Scripts

### `test_categories.sh` - Test Category Runner
Organize and run different categories of tests.

```bash
# Show available test categories
./scripts/test_categories.sh list

# Run unit tests
./scripts/test_categories.sh unit

# List integration tests
./scripts/test_categories.sh integration

# List email tests
./scripts/test_categories.sh email

# List performance tests
./scripts/test_categories.sh performance

# Run all tests
./scripts/test_categories.sh all
```

## Test Organization

The backend uses a structured approach to testing:

### Unit Tests (`*_test.go`)
- Located alongside source code in `internal/` packages
- Run with: `go test ./internal/...`
- Standard Go testing conventions

### Integration Tests
- Located in `testing/integration/`
- Standalone executables with `package main`
- Run individually: `go run testing/integration/[filename].go`

### Email Tests
- Located in `testing/email/`
- Standalone executables for email system testing
- Run individually: `go run testing/email/[filename].go`

### Performance Tests
- Located in `testing/performance/`
- Standalone executables for performance testing
- Run individually: `go run testing/performance/[filename].go`

### Standalone Tests
- Located in root directory (`test_*.go`)
- Legacy test files that haven't been moved to testing directories
- Each is a standalone executable

## Build Verification Checklist

The `verify_build.sh` script performs these checks:

1. ✅ **Go Module Dependencies** - Verifies all dependencies are clean
2. ✅ **Main Application Build** - Ensures `cmd/server` builds successfully
3. ✅ **Internal Packages** - Verifies all internal packages compile
4. ✅ **Test Compilation** - Checks test files can be compiled
5. ✅ **Build Conflicts** - Identifies multiple main() function issues
6. ✅ **Testing Organization** - Verifies test directory structure
7. ✅ **Dependency Cleanup** - Checks for unused dependencies
8. ✅ **Build Tags** - Tests building with different tags
9. ✅ **Cross-compilation** - Tests building for different platforms
10. ✅ **Performance** - Measures build time and binary size

## Common Build Issues

### Multiple `main()` Functions
- **Issue**: Multiple files with `package main` cause redeclaration errors
- **Solution**: Standalone test files have been moved to `testing/` directories
- **Status**: ✅ Resolved - main application builds cleanly

### Dependency Conflicts
- **Issue**: Conflicting or unused dependencies
- **Solution**: Run `go mod tidy` and `go mod verify`
- **Status**: ✅ Resolved - all dependencies verified

### Test File Conflicts
- **Issue**: Test files interfering with main build
- **Solution**: Organized tests into categories and separate directories
- **Status**: ✅ Resolved - tests are properly isolated

## Usage Examples

### Quick Development Build
```bash
# Build and run the server
./scripts/build.sh
./server
```

### Pre-deployment Verification
```bash
# Run full build verification
./scripts/verify_build.sh

# Check test organization
./scripts/test_categories.sh list

# Run unit tests
./scripts/test_categories.sh unit
```

### Testing Specific Components
```bash
# Test email functionality
go run testing/email/test_email_simple.go

# Test chat functionality  
go run testing/integration/test_chat_simple.go

# Run performance tests
go run testing/performance/test_performance_simple.go
```

## Script Maintenance

### Adding New Scripts
1. Create script in `backend/scripts/`
2. Make executable: `chmod +x scripts/[script_name].sh`
3. Add documentation to this README
4. Test script functionality

### Updating Build Process
1. Modify `verify_build.sh` for new checks
2. Update `build.sh` for new build requirements
3. Update `test_categories.sh` for new test categories
4. Update this documentation

## Troubleshooting

### Build Fails
1. Run `./scripts/verify_build.sh` for detailed diagnostics
2. Check Go version: `go version`
3. Verify dependencies: `go mod verify`
4. Clean and rebuild: `go clean -cache && ./scripts/build.sh`

### Tests Don't Run
1. Check test organization: `./scripts/test_categories.sh list`
2. Verify test files exist in correct directories
3. Run individual tests to isolate issues

### Permission Errors
```bash
# Make scripts executable
chmod +x scripts/*.sh
```

## Integration with CI/CD

These scripts are designed to work in CI/CD environments:

```yaml
# Example GitHub Actions usage
- name: Verify Build
  run: ./scripts/verify_build.sh

- name: Run Unit Tests
  run: ./scripts/test_categories.sh unit
```

The scripts provide appropriate exit codes and colored output for both local development and automated environments.