package examples

import (
	c "jecolasurdo/go-chaining"
	"jecolasurdo/go-chaining/behavior"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func addOne(value *interface{}) (*interface{}, error) {
	var result interface{} = (*value).(int) + 1
	return &result, nil
}

func multiplyBySix(value *interface{}) (*interface{}, error) {
	var result interface{} = (*value).(int) * 6
	return &result, nil
}

func convertToString(value *interface{}) (*interface{}, error) {
	var result interface{} = strconv.Itoa((*value).(int))
	return &result, nil
}

func Test_ChainItAllTogether(t *testing.T) {
	chain := c.New()
	var initialValue interface{} = 1
	initialBehavior := c.ActionArg{
		Value:    &initialValue,
		Behavior: behavior.InjectSuppliedValue,
	}
	chain.ApplyUnaryIface(addOne, initialBehavior)
	chain.ApplyUnaryIface(multiplyBySix, c.ActionArg{})
	chain.ApplyUnaryIface(convertToString, c.ActionArg{})

	result, err := chain.Flush()
	assert.Equal(t, "12", (*result).(string))
	assert.Nil(t, err)
}
