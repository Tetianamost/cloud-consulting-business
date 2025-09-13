# Build Issues Troubleshooting Guide

This guide helps resolve common build and compilation issues in the Cloud Consulting Platform.

## Table of Contents

1. [Backend Build Issues](#backend-build-issues)
2. [Frontend Build Issues](#frontend-build-issues)
3. [Test Execution Issues](#test-execution-issues)
4. [Docker Build Issues](#docker-build-issues)
5. [Environment Setup Issues](#environment-setup-issues)
6. [File Organization Issues](#file-organization-issues)

## Backend Build Issues

### Issue: Multiple `package main` Redeclaration Errors

**Symptoms:**
```
main redeclared in this block
previous declaration at ./test_something.go:1:1
```

**Root Cause:** Multiple standalone test files with `package main` causing conflicts during build.

**Solution:**
1. **Use Application-Only Build Script:**
   ```bash
   ./backend/scripts/build_app_only.sh
   ```

2. **Manual Build (Application Only):**
   ```bash
   cd backend
   go build -o bin/server ./cmd/server
   ```

3. **Check File Organization:**
   - Standalone test files should be in `testing/` subdirectories
   - Main application code should be in `internal/` and `cmd/`

### Issue: Go Module Dependency Conflicts

**Symptoms:**
```
go: module requires Go 1.21 or later
go: inconsistent vendoring
```

**Solution:**
```bash
cd backend
go mod tidy
go mod verify
go clean -modcache  # If needed
```

### Issue: Missing Dependencies

**Symptoms:**
```
package "some/package" is not in GOROOT
cannot find package
```

**Solution:**
```bash
cd backend
go mod download
go mod tidy
```

### Issue: Build Fails with Test Files

**Symptoms:**
```
build constraints exclude all Go files
no Go files in /path/to/directory
```

**Solution:**
1. **Use the correct build command:**
   ```bash
   # Build application only (excludes test files)
   go build -o bin/server ./cmd/server
   
   # NOT this (includes test files):
   go build ./...
   ```

2. **Use the build script:**
   ```bash
   ./backend/scripts/build_app_only.sh
   ```

## Frontend Build Issues

### Issue: TypeScript Compilation Errors

**Symptoms:**
```
Type 'string' is not assignable to type 'number'
Property 'xyz' does not exist on type
```

**Solution:**
1. **Check TypeScript Configuration:**
   ```bash
   cd frontend
   npx tsc --noEmit  # Check for type errors
   ```

2. **Fix Type Issues:**
   - Add missing type definitions
   - Update component prop types
   - Check interface definitions

3. **Update Dependencies:**
   ```bash
   cd frontend
   npm update
   npm audit fix
   ```

### Issue: Missing Dependencies

**Symptoms:**
```
Module not found: Can't resolve 'some-package'
Cannot resolve dependency
```

**Solution:**
```bash
cd frontend
npm install
# or for clean install:
rm -rf node_modules package-lock.json
npm install
```

### Issue: Build Memory Issues

**Symptoms:**
```
JavaScript heap out of memory
FATAL ERROR: Ineffective mark-compacts near heap limit
```

**Solution:**
```bash
cd frontend
export NODE_OPTIONS="--max-old-space-size=4096"
npm run build
```

## Test Execution Issues

### Issue: Tests Not Found

**Symptoms:**
```
no test files found
no tests to run
```

**Solution:**
1. **List Available Tests:**
   ```bash
   ./backend/scripts/run_test_categories.sh list
   ```

2. **Check Test Organization:**
   - Unit tests: `internal/**/*_test.go`
   - Integration tests: `testing/integration/`
   - Email tests: `testing/email/`
   - Performance tests: `testing/performance/`

3. **Run Specific Category:**
   ```bash
   ./backend/scripts/run_test_categories.sh unit
   ./backend/scripts/run_test_categories.sh integration
   ```

### Issue: Test Compilation Errors

**Symptoms:**
```
test file compilation failed
undefined: SomeFunction
```

**Solution:**
1. **Check Test Dependencies:**
   ```bash
   cd backend
   go mod tidy
   ```

2. **Verify Test File Organization:**
   - Test files should be in correct directories
   - Import paths should be correct
   - Mock dependencies should be available

3. **Run Tests with Verbose Output:**
   ```bash
   ./backend/scripts/run_test_categories.sh unit -v
   ```

### Issue: Test Timeout

**Symptoms:**
```
test timed out after 30s
panic: test timed out
```

**Solution:**
```bash
# Increase timeout
./backend/scripts/run_test_categories.sh integration --timeout=60s

# Or manually:
go test -timeout=60s ./internal/...
```

## Docker Build Issues

### Issue: Docker Build Context Too Large

**Symptoms:**
```
Sending build context to Docker daemon  XXX.XGB
```

**Solution:**
1. **Check .dockerignore:**
   ```bash
   # Add to .dockerignore:
   node_modules
   coverage
   *.log
   .git
   testing/
   ```

2. **Clean Build Context:**
   ```bash
   docker system prune -f
   docker builder prune -f
   ```

### Issue: Docker Registry Connection Issues

**Symptoms:**
```
failed to resolve source metadata
no such host
```

**Solution:**
1. **Use Local Builds:**
   ```bash
   ./start-local-options.sh  # Choose option 1
   ```

2. **Build Offline:**
   ```bash
   ./build-offline.sh
   docker-compose -f docker-compose.offline.yml up
   ```

### Issue: Docker Port Conflicts

**Symptoms:**
```
port is already allocated
bind: address already in use
```

**Solution:**
1. **Check Port Usage:**
   ```bash
   lsof -i :3006  # Frontend port
   lsof -i :8061  # Backend port
   ```

2. **Stop Conflicting Services:**
   ```bash
   ./stop-local.sh
   docker-compose down
   ```

3. **Change Ports (if needed):**
   - Edit `docker-compose.yml`
   - Update `.env` file
   - Update frontend configuration

## Environment Setup Issues

### Issue: Go Version Incompatibility

**Symptoms:**
```
go: module requires Go 1.21 or later
```

**Solution:**
1. **Check Go Version:**
   ```bash
   go version
   ```

2. **Update Go:**
   - Download from https://golang.org/dl/
   - Or use package manager:
     ```bash
     # macOS
     brew install go
     
     # Ubuntu
     sudo apt update && sudo apt install golang-go
     ```

### Issue: Node.js Version Issues

**Symptoms:**
```
error: This version of Node.js requires npm
unsupported engine
```

**Solution:**
1. **Check Node Version:**
   ```bash
   node --version
   npm --version
   ```

2. **Update Node.js:**
   ```bash
   # Using nvm (recommended)
   nvm install 18
   nvm use 18
   
   # Or download from https://nodejs.org/
   ```

### Issue: AWS Credentials Not Working

**Symptoms:**
```
AWS credentials not found
UnauthorizedOperation
```

**Solution:**
1. **Check Environment Variables:**
   ```bash
   echo $AWS_ACCESS_KEY_ID
   echo $AWS_SECRET_ACCESS_KEY
   echo $AWS_REGION
   ```

2. **Update .env File:**
   ```bash
   # Edit backend/.env
   AWS_ACCESS_KEY_ID=your-key-id
   AWS_SECRET_ACCESS_KEY=your-secret-key
   AWS_REGION=us-east-1
   ```

3. **Test AWS Connection:**
   ```bash
   # Run email connectivity test
   cd backend
   go run testing/email/test_ses_connectivity.go
   ```

## File Organization Issues

### Issue: Test Files in Wrong Location

**Symptoms:**
```
package main redeclared
test files interfering with build
```

**Solution:**
1. **Check Current Organization:**
   ```bash
   ./backend/scripts/run_test_categories.sh list
   ```

2. **Move Files to Correct Locations:**
   ```bash
   # Email tests should be in:
   backend/testing/email/
   
   # Integration tests should be in:
   backend/testing/integration/
   
   # Performance tests should be in:
   backend/testing/performance/
   
   # Unit tests should be in:
   backend/internal/*/
   ```

3. **Verify Organization:**
   ```bash
   # Should build without errors:
   ./backend/scripts/build_app_only.sh
   ```

### Issue: Import Path Errors After File Move

**Symptoms:**
```
cannot find package
import path does not exist
```

**Solution:**
1. **Update Import Paths:**
   ```go
   // Update imports in moved test files
   import (
       "github.com/your-org/cloud-consulting/internal/services"
       "github.com/your-org/cloud-consulting/internal/domain"
   )
   ```

2. **Update Package Declarations:**
   ```go
   // For test executables:
   package main
   
   // For test packages:
   package services_test
   ```

3. **Run Tests to Verify:**
   ```bash
   ./backend/scripts/run_test_categories.sh unit
   ```

## Quick Diagnostic Commands

### Check System Status
```bash
# Go environment
go version
go env GOPATH
go env GOROOT

# Node environment
node --version
npm --version

# Docker status
docker --version
docker-compose --version
docker system df

# Port usage
lsof -i :3006
lsof -i :8061
```

### Build Verification
```bash
# Backend build test
./backend/scripts/build_app_only.sh

# Frontend build test
cd frontend && npm run build

# Test organization check
./backend/scripts/run_test_categories.sh list

# Docker build test
docker-compose -f docker-compose.local.yml build
```

### Clean Reset
```bash
# Clean Go cache
cd backend
go clean -cache -modcache -testcache

# Clean Node modules
cd frontend
rm -rf node_modules package-lock.json
npm install

# Clean Docker
docker system prune -f
docker builder prune -f

# Clean build artifacts
rm -rf backend/bin/
rm -rf frontend/build/
```

## Getting Help

If you're still experiencing issues after trying these solutions:

1. **Check the logs:**
   ```bash
   # Application logs
   docker-compose logs backend
   docker-compose logs frontend
   
   # Build logs
   ./backend/scripts/build_app_only.sh 2>&1 | tee build.log
   ```

2. **Run diagnostic commands:**
   ```bash
   # System information
   go version
   node --version
   docker --version
   
   # Test organization
   ./backend/scripts/run_test_categories.sh list
   ```

3. **Create a minimal reproduction:**
   - Start with a clean environment
   - Follow the quick start guide
   - Document the exact steps that fail

4. **Check for known issues:**
   - Review recent commits
   - Check for environment-specific problems
   - Verify all dependencies are installed

## Prevention Tips

1. **Use the provided build scripts** instead of manual commands
2. **Keep test files organized** in the correct directories
3. **Run builds regularly** to catch issues early
4. **Use Docker for consistent environments** when possible
5. **Keep dependencies updated** but test after updates
6. **Document any custom setup steps** for your environment

This troubleshooting guide should help resolve most build and compilation issues. Keep it updated as new issues are discovered and resolved.