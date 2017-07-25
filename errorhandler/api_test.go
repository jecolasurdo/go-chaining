package errorhandler_test

import (
	"errors"
	"testing"

	"jecolasurdo/go-deferrederrors/errorhandler"

	"github.com/stretchr/testify/assert"
)

func Test_TryBool_ActionTrueNoError_ReturnsTrue(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	trueAction := func() (bool, error) { return true, nil }
	result := d.TryBool(trueAction)
	assert.True(t, result)
}

func Test_TryBool_ActionFalseNoError_ReturnsFalse(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	falseAction := func() (bool, error) { return false, nil }
	result := d.TryBool(falseAction)
	assert.False(t, result)
}

func Test_TryBool_PreviousError_IgnoresAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func() (bool, error) {
		timesActionWasCalled++
		return true, nil
	}

	d.LocalError = errors.New("test error")
	d.TryBool(action)

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_TryBool_PreviousError_ReturnsFalse(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	trueAction := func() (bool, error) {
		return true, nil
	}

	d.LocalError = errors.New("test error")
	result := d.TryBool(trueAction)

	assert.False(t, result)
}

func Test_TryBool_ActionErrors_SetsLocalError(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	errorAction := func() (bool, error) { return false, errors.New("test error") }

	d.TryBool(errorAction)

	assert.NotNil(t, d.LocalError)
}

func Test_TryVoid_NoPreviousError_ExecutesAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func() error {
		timesActionWasCalled++
		return nil
	}

	d.TryVoid(action)

	assert.Equal(t, 1, timesActionWasCalled)
}

func Test_TryVoid_PreviousError_IgnoresAction(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	timesActionWasCalled := 0
	action := func() error {
		timesActionWasCalled++
		return nil
	}

	d.LocalError = errors.New("test error")
	d.TryVoid(action)

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_TryVoid_ActionErrors_SetsLocalError(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	errorAction := func() error { return errors.New("test error") }

	d.TryVoid(errorAction)

	assert.NotNil(t, d.LocalError)
}

func Test_FlushError_Normally_ResetsErrorToNil(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	d.LocalError = errors.New("test error")
	d.FlushError()
	assert.Nil(t, d.LocalError)
}

func Test_FlushError_Normally_ReturnsAnyEror(t *testing.T) {
	d := new(errorhandler.DeferredErrorContext)
	d.LocalError = errors.New("test error")
	err := d.FlushError()
	assert.NotNil(t, err)
}
