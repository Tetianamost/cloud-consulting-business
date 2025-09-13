package shared

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigManager(t *testing.T) {
	logger := logrus.New()
	cm := NewConfigManager(logger)

	assert.NotNil(t, cm)
	assert.NotNil(t, cm.configs)
	assert.Equal(t, logger, cm.logger)
}

func TestCreateDefaultConfig(t *testing.T) {
	logger := logrus.New()
	cm := NewConfigManager(logger)

	config := cm.createDefaultConfig(TestEnvUnit)

	assert.Equal(t, TestEnvUnit, config.Environment)
	assert.Equal(t, "sqlite://memory", config.Database.URL)
	assert.True(t, config.Parallel)
	assert.True(t, config.SkipSlow)
	assert.False(t, config.CleanupDB)
	assert.True(t, config.Bedrock.UseMock)
}

func TestCreateDefaultConfigIntegration(t *testing.T) {
	logger := logrus.New()
	cm := NewConfigManager(logger)

	config := cm.createDefaultConfig(TestEnvIntegration)

	assert.Equal(t, TestEnvIntegration, config.Environment)
	assert.True(t, config.SeedData)
	assert.False(t, config.Bedrock.UseMock)
	assert.False(t, config.SES.UseMock)
	assert.False(t, config.AWS.UseMocks)
}

func TestOverrideWithEnvVars(t *testing.T) {
	logger := logrus.New()
	cm := NewConfigManager(logger)

	// Set environment variables
	os.Setenv("TEST_DATABASE_URL", "postgres://test:test@localhost/override_db")
	os.Setenv("TEST_REDIS_URL", "redis://localhost:6379/9")
	os.Setenv("BEDROCK_USE_MOCK", "false")
	os.Setenv("TEST_PARALLEL", "false")
	defer func() {
		os.Unsetenv("TEST_DATABASE_URL")
		os.Unsetenv("TEST_REDIS_URL")
		os.Unsetenv("BEDROCK_USE_MOCK")
		os.Unsetenv("TEST_PARALLEL")
	}()

	config := cm.createDefaultConfig(TestEnvUnit)
	cm.overrideWithEnvVars(config)

	assert.Equal(t, "postgres://test:test@localhost/override_db", config.Database.URL)
	assert.Equal(t, "redis://localhost:6379/9", config.Redis.URL)
	assert.False(t, config.Bedrock.UseMock)
	assert.False(t, config.Parallel)
}

func TestValidateConfig(t *testing.T) {
	logger := logrus.New()
	cm := NewConfigManager(logger)

	// Valid configuration
	config := cm.createDefaultConfig(TestEnvUnit)
	err := cm.validateConfig(config)
	assert.NoError(t, err)

	// Invalid configuration - empty database URL
	config.Database.URL = ""
	err = cm.validateConfig(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database URL is required")
}

func TestListEnvironments(t *testing.T) {
	logger := logrus.New()
	cm := NewConfigManager(logger)

	envs := cm.ListEnvironments()

	assert.Len(t, envs, 6)
	assert.Contains(t, envs, TestEnvUnit)
	assert.Contains(t, envs, TestEnvIntegration)
	assert.Contains(t, envs, TestEnvPerformance)
	assert.Contains(t, envs, TestEnvE2E)
	assert.Contains(t, envs, TestEnvLocal)
	assert.Contains(t, envs, TestEnvCI)
}

func TestConfigUtils(t *testing.T) {
	logger := logrus.New()
	cu := NewConfigUtils(logger)

	assert.NotNil(t, cu)
	assert.NotNil(t, cu.configManager)
	assert.Equal(t, logger, cu.logger)
}

func TestDetectEnvironment(t *testing.T) {
	logger := logrus.New()
	cu := NewConfigUtils(logger)

	// Test explicit environment variable
	os.Setenv("TEST_ENVIRONMENT", "integration")
	defer os.Unsetenv("TEST_ENVIRONMENT")

	env := cu.DetectEnvironment()
	assert.Equal(t, TestEnvIntegration, env)
}

func TestDetectEnvironmentCI(t *testing.T) {
	logger := logrus.New()
	cu := NewConfigUtils(logger)

	// Test CI detection
	os.Setenv("CI", "true")
	defer os.Unsetenv("CI")

	env := cu.DetectEnvironment()
	assert.Equal(t, TestEnvCI, env)
}

func TestDetectEnvironmentShort(t *testing.T) {
	logger := logrus.New()
	cu := NewConfigUtils(logger)

	// Test short test detection
	os.Setenv("TEST_SHORT", "true")
	defer os.Unsetenv("TEST_SHORT")

	env := cu.DetectEnvironment()
	assert.Equal(t, TestEnvUnit, env)
}

func TestDetermineEnvironmentFromTestName(t *testing.T) {
	logger := logrus.New()
	cu := NewConfigUtils(logger)

	testCases := []struct {
		testName string
		expected TestEnvironment
	}{
		{"TestUnitFunction", TestEnvUnit},
		{"TestMockService", TestEnvUnit},
		{"TestIntegrationAPI", TestEnvIntegration},
		{"TestAPIEndpoint", TestEnvIntegration},
		{"TestPerformanceLoad", TestEnvPerformance},
		{"BenchmarkFunction", TestEnvPerformance},
		{"TestE2EFlow", TestEnvE2E},
		{"TestEndToEnd", TestEnvE2E},
		{"TestSomethingElse", TestEnvUnit}, // Default
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			env := cu.determineEnvironmentFromTestName(tc.testName)
			assert.Equal(t, tc.expected, env)
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test ShouldSkipSlow
	os.Setenv("TEST_SKIP_SLOW", "true")
	defer os.Unsetenv("TEST_SKIP_SLOW")

	assert.True(t, ShouldSkipSlow())

	// Test ShouldRunInParallel
	os.Setenv("TEST_PARALLEL", "true")
	defer os.Unsetenv("TEST_PARALLEL")

	assert.True(t, ShouldRunInParallel())

	// Test IsTestEnvironment
	os.Setenv("TEST_ENVIRONMENT", "unit")
	defer os.Unsetenv("TEST_ENVIRONMENT")

	assert.True(t, IsTestEnvironment(TestEnvUnit))
	assert.False(t, IsTestEnvironment(TestEnvIntegration))
}
