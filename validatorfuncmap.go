package validator

import (
	"fmt"
)

var validatorFuncMap = map[string]map[string]interface{}{
	"int": {
		"nonzero": func(v int, opts string) error {
			var zeroValue int
			if v == zeroValue {
				return fmt.Errorf("Expected non-zero value")
			}
			return nil
		},
	},
}
