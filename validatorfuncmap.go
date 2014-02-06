package validator

import (
	"fmt"
	"reflect"
)

var validatorFuncMap = map[string]map[string]ValidatorFunc{
	"int": {
		"nonzero": ValidatorFunc{func(v int, opts string) error {
			var zeroValue int
			if v == zeroValue {
				return fmt.Errorf("Expected non-zero value")
			}
			return nil
		}},
	},
	"string": {
		"nonzero": ValidatorFunc{func(v string, opts string) error {
			var zeroValue string
			if v == zeroValue {
				return fmt.Errorf("Expected non-zero value")
			}
			return nil
		}},
	},
}

type ValidatorFunc struct {
	F interface{}
}

func (vf ValidatorFunc) Call(v reflect.Value, opts string) error {
	vfv := reflect.ValueOf(vf.F)
	if !v.Type().AssignableTo(vfv.Type().In(0)) {
		return fmt.Errorf("Wrong type")
	}
	err := vfv.Call([]reflect.Value{v, reflect.ValueOf(opts)})[0]
	if err.Interface() == nil {
		return nil
	}
	return err.Interface().(error)
}
