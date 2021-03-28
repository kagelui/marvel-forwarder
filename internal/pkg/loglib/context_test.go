package loglib

import (
	"context"
	"testing"

	"github.com/kagelui/marvel-forwarder/internal/pkg/testutil"
)

func TestContext(t *testing.T) {
	// Given:
	ctx := context.Background()

	// When: GetLogger
	result := GetLogger(ctx)

	// Then:
	testutil.Asserts(t, result != nil, "should have logger when not set")

	// Given:
	given := NewLogger(nil)

	// When: SetLogger
	ctx = SetLogger(ctx, given)

	// When: HasLogger
	resultNow, ok := HasLogger(ctx)

	// Then:
	testutil.Asserts(t, ok, "should have logger")
	testutil.Equals(t, given, resultNow)

	// When: GetLogger
	result = GetLogger(ctx)

	// Then:
	testutil.Equals(t, given, resultNow)
}
