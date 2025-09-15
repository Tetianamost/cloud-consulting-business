package shared

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSimpleTestConfig(t *testing.T) {
	config := LoadSimpleTestConfig()

	assert.NotNil(t, config)
	assert.NotEmpty(t, config.DatabaseURL)
	assert.NotEmpty(t, config.RedisURL)
	assert.Equal(t, 30*time.Second, config.Timeout)
}

func TestSetupSimpleTestLogger(t *testing.T) {
	logger := SetupSimpleTestLogger("info")

	assert.NotNil(t, logger)
	assert.Equal(t, logrus.InfoLevel, logger.Level)
}

func TestCreateTestContext(t *testing.T) {
	timeout := 5 * time.Second
	ctx, cancel := CreateTestContext(timeout)
	defer cancel()

	assert.NotNil(t, ctx)

	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, deadline.After(time.Now()))
}

func TestNewTestSuite(t *testing.T) {
	suite, err := NewTestSuite("test-suite")

	// This might fail if database is not available, which is expected
	if err != nil {
		t.Skipf("Skipping test suite creation due to database unavailability: %v", err)
		return
	}

	require.NoError(t, err)
	assert.NotNil(t, suite)
	assert.Equal(t, "test-suite", suite.Name)
	assert.NotNil(t, suite.Config)
	assert.NotNil(t, suite.Logger)
	assert.NotNil(t, suite.Metrics)

	// Cleanup
	suite.Cleanup()
}

func TestTestMetrics(t *testing.T) {
	metrics := NewTestMetrics()

	assert.NotNil(t, metrics)
	assert.False(t, metrics.StartTime.IsZero())

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	metrics.Finish()

	assert.False(t, metrics.EndTime.IsZero())
	assert.True(t, metrics.Duration > 0)
	assert.True(t, metrics.EndTime.After(metrics.StartTime))
}

func TestAssertEventuallyTrue(t *testing.T) {
	counter := 0
	condition := func() bool {
		counter++
		return counter >= 3
	}

	// This should succeed after a few iterations
	AssertEventuallyTrue(t, condition, 1*time.Second, "counter should reach 3")

	assert.GreaterOrEqual(t, counter, 3)
}

func TestSkipIfShort(t *testing.T) {
	// This test will be skipped if running with -short flag
	SkipIfShort(t, "testing skip functionality")

	// If we reach here, we're not in short mode
	t.Log("Not running in short mode")
}
