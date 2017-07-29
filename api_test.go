package chaining_test

import (
	"errors"
	"testing"

	"jecolasurdo/go-chaining"
	"jecolasurdo/go-chaining/behavior"

	"github.com/stretchr/testify/assert"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// Atomic Function Tests
/// All chaining methods are ultimately stated in terms of a single atomic function.
/// Base functionality is guaranteed via the atomic function.
/// Since the default atomic function is not public, we test it through its public wrapper "ApplyUnaryInterface"
///

func Test_AtomicFunction_PreviousError_IgnoresAction(t *testing.T) {
	d := chaining.New()
	timesActionWasCalled := 0
	action := func(*interface{}) (*interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	d.LocalError = errors.New("test error")
	d.ApplyUnaryIface(action, chaining.ActionArg{})

	assert.Equal(t, 0, timesActionWasCalled)
}

func Test_AtomicFunction_NoPreviousError_ExecutesAction(t *testing.T) {
	d := chaining.New()
	timesActionWasCalled := 0
	action := func(*interface{}) (*interface{}, error) {
		timesActionWasCalled++
		return nil, nil
	}

	d.ApplyUnaryIface(action, chaining.ActionArg{})

	assert.Equal(t, 1, timesActionWasCalled)
}

func Test_AtomicFunction_NoPreviousError_BehaviorIsNotSpecified_InjectsPreviousValue(t *testing.T) {
	d := chaining.New()
	var injectedValue *interface{}
	action := func(value *interface{}) (*interface{}, error) {
		injectedValue = value
		return nil, nil
	}
	argWithBehaviorNotSpecified := chaining.ActionArg{}
	var simulatedValueOfPreviousActionInChain interface{} = "somevalue"
	d.PreviousActionResult = &simulatedValueOfPreviousActionInChain

	d.ApplyUnaryIface(action, argWithBehaviorNotSpecified)

	assert.Equal(t, &simulatedValueOfPreviousActionInChain, injectedValue)
}

func Test_AtomicFunction_NoPreviousError_BehaviorIsUsePrevious_InjectsPreviousValue(t *testing.T) {
	d := chaining.New()
	var injectedValue interface{}
	action := func(value *interface{}) (*interface{}, error) {
		injectedValue = value
		return nil, nil
	}
	argWithSpecifiedBehavior := chaining.ActionArg{
		Behavior: behavior.InjectPreviousResult,
	}
	var simulatedValueOfPreviousActionInChain interface{} = "somevalue"
	d.PreviousActionResult = &simulatedValueOfPreviousActionInChain

	d.ApplyUnaryIface(action, argWithSpecifiedBehavior)

	assert.Equal(t, &simulatedValueOfPreviousActionInChain, injectedValue)
}

func Test_AtomicFunction_NoPreviousError_BehaviorIsOverridePrevious_InjectsSuppliedValue(t *testing.T) {
	d := chaining.New()
	injectedValue := ""
	action := func(value *interface{}) (*interface{}, error) {
		injectedValue = (*value).(string)
		return nil, nil
	}
	var valueSubmittedThroughArg interface{} = "valueFromArg"
	argWithSpecifiedBehavior := chaining.ActionArg{
		Behavior: behavior.InjectSuppliedValue,
		Value:    &valueSubmittedThroughArg,
	}
	var simulatedValueOfPreviousActionInChain interface{} = "previousValue"
	d.PreviousActionResult = &simulatedValueOfPreviousActionInChain

	d.ApplyUnaryIface(action, argWithSpecifiedBehavior)

	assert.Equal(t, valueSubmittedThroughArg, injectedValue)
}

func Test_AtomicFunction_NoPreviousError_ForAnySpecifiedBehavior_SetsPreviousActionResult(t *testing.T) {
	d := chaining.New()
	var expectedReturnValue interface{} = "expectedReturnValue"
	action := func(value *interface{}) (*interface{}, error) { return &expectedReturnValue, nil }
	var valueFromArg interface{} = "valueFromArg"
	arg := chaining.ActionArg{
		Value: &valueFromArg,
	}

	d.PreviousActionResult = nil
	arg.Behavior = behavior.InjectSuppliedValue
	d.ApplyUnaryIface(action, arg)
	assert.Equal(t, &expectedReturnValue, d.PreviousActionResult)

	d.PreviousActionResult = nil
	arg.Behavior = behavior.InjectPreviousResult
	d.ApplyUnaryIface(action, arg)
	assert.Equal(t, &expectedReturnValue, d.PreviousActionResult)

	d.PreviousActionResult = nil
	arg.Behavior = behavior.NotSpecified
	d.ApplyUnaryIface(action, arg)
	assert.Equal(t, &expectedReturnValue, d.PreviousActionResult)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// Interface Function Tests
/// Since all capabilities for interface functions are guaranteed by the atomic function,
/// the goal of these tests is simply to ensure that all chaining methods call the atomic function.
///

var numberOfTimesAtomicCalled int
var suppliedContext *chaining.Context
var suppliedAction func(interface{}) (interface{}, error)
var suppliedArg chaining.ActionArg
var mockContext *chaining.Context

func ResetTestParameters() {
	numberOfTimesAtomicCalled = 0
	mockContext = &chaining.Context{
		AtomicFunc: func(c *chaining.Context, action func(*interface{}) (*interface{}, error), arg chaining.ActionArg) {
			numberOfTimesAtomicCalled++
			result, _ := action(arg.Value)
			c.PreviousActionResult = result
		},
	}
}

func Test_ApplyNullary_Normally_CallsAtomic(t *testing.T) {
	ResetTestParameters()
	action := func() error { return nil }
	mockContext.ApplyNullary(action, behavior.NotSpecified)
	assert.Equal(t, 1, numberOfTimesAtomicCalled)
}

func Test_ApplyNullaryIface_Normally_CallsAtomic(t *testing.T) {
	ResetTestParameters()
	action := func() (*interface{}, error) { return nil, nil }
	mockContext.ApplyNullaryIface(action, behavior.NotSpecified)
	assert.Equal(t, 1, numberOfTimesAtomicCalled)
}

func Test_ApplyUnary_Normally_CallsAtomic(t *testing.T) {
	ResetTestParameters()
	action := func(*interface{}) error { return nil }
	mockContext.ApplyUnary(action, chaining.ActionArg{})
	assert.Equal(t, 1, numberOfTimesAtomicCalled)
}

func Test_ApplyUnaryIface_Normally_CallsAtomic(t *testing.T) {
	ResetTestParameters()
	action := func(*interface{}) (*interface{}, error) { return nil, nil }
	mockContext.ApplyUnaryIface(action, chaining.ActionArg{})
	assert.Equal(t, 1, numberOfTimesAtomicCalled)
}

func Test_ApplyNullaryBool_Normally_CallsAtomic(t *testing.T) {
	ResetTestParameters()
	False := false
	action := func() (*bool, error) { return &False, nil }
	mockContext.ApplyNullaryBool(action, behavior.NotSpecified)
	assert.Equal(t, 1, numberOfTimesAtomicCalled)
}

func Test_ApplyUnaryBool_Normally_CallsAtomic(t *testing.T) {
	ResetTestParameters()
	False := false
	action := func(*interface{}) (*bool, error) { return &False, nil }
	mockContext.ApplyUnaryBool(action, chaining.ActionArg{})
	assert.Equal(t, 1, numberOfTimesAtomicCalled)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// Special Nullary Boolean Function Tests
///

func Test_ApplyNullaryBool_ActionTrueNoError_ReturnsTrue(t *testing.T) {
	d := chaining.New()
	True := true
	trueAction := func() (*bool, error) { return &True, nil }
	result := d.ApplyNullaryBool(trueAction, behavior.NotSpecified)
	assert.True(t, result)
}

func Test_ApplyNullaryBool_ActionFalseNoError_ReturnsFalse(t *testing.T) {
	d := chaining.New()
	False := false
	falseAction := func() (*bool, error) { return &False, nil }
	result := d.ApplyNullaryBool(falseAction, behavior.NotSpecified)
	assert.False(t, result)
}

func Test_ApplyNullaryBool_PreviousError_ReturnsFalse(t *testing.T) {
	d := chaining.New()
	trueAction := func() (*bool, error) {
		True := true
		return &True, nil
	}

	d.LocalError = errors.New("test error")
	result := d.ApplyNullaryBool(trueAction, behavior.NotSpecified)

	assert.False(t, result)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// Special Unary Boolean Function Tests
///

func Test_ApplyUnaryBool_ActionTrueNoError_ReturnsTrue(t *testing.T) {
	d := chaining.New()
	True := true
	trueAction := func(*interface{}) (*bool, error) { return &True, nil }
	result := d.ApplyUnaryBool(trueAction, chaining.ActionArg{})
	assert.True(t, result)
}

func Test_ApplyUnaryBool_ActionFalseNoError_ReturnsFalse(t *testing.T) {
	d := chaining.New()
	False := false
	falseAction := func(*interface{}) (*bool, error) { return &False, nil }
	result := d.ApplyUnaryBool(falseAction, chaining.ActionArg{})
	assert.False(t, result)
}

func Test_ApplyUnaryBool_PreviousError_ReturnsFalse(t *testing.T) {
	d := chaining.New()
	trueAction := func(*interface{}) (*bool, error) {
		True := true
		return &True, nil
	}
	d.LocalError = errors.New("test error")

	result := d.ApplyUnaryBool(trueAction, chaining.ActionArg{})

	assert.False(t, result)
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// /// Flush Tests
// ///

// func Test_Flush_Normally_ResetsContext(t *testing.T) {
// 	d := chaining.New()
// 	d.LocalError = errors.New("test error")

// 	var notNilValue interface{} = "Not nil"
// 	d.PreviousActionResult = &notNilValue

// 	d.Flush()

// 	assert.Nil(t, d.LocalError)
// 	assert.Nil(t, d.PreviousActionResult)
// }

// func Test_Flush_Normally_ReturnsErrorAndFinalResult(t *testing.T) {
// 	d := chaining.New()
// 	d.LocalError = errors.New("test error")
// 	var expectedFinalResult interface{} = "FinalResult"
// 	d.PreviousActionResult = &expectedFinalResult

// 	result, err := d.Flush()

// 	assert.NotNil(t, err)
// 	assert.Equal(t, &expectedFinalResult, result)
// }
