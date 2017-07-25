package errorhandler_test

import (
	"testing"

	"jecolasurdo/go-deferrederrors/errorhandler"

	"github.com/stretchr/testify/assert"
)

func Test_TryBool_ActionTrueNoError_ReturnsTrue(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	action := func() (bool, error) { return true, nil }
	result := d.TryBool(action)
	assert.True(t, result)
}
