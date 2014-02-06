package validator

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
)

func Validate(x interface{}) error {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("Only structs or pointer to structs can be validated")
	}

	for i := 0; i < v.Type().NumField(); i++ {
		if !isExportedField(v.Type().Field(i)) {
			continue
		}
		if err := validate(v.Field(i), v.Type().Field(i).Tag.Get("validator")); err != nil {
			return fmt.Errorf("Field %s is invalid: %s", v.Type().Field(i).Name, err)
		}
	}
	return nil
}

func validate(v reflect.Value, options string) error {
	name := v.Type().Name()
	validatorFuncs, ok := validatorFuncMap[name]
	if !ok {
		name = v.Kind().String()
		validatorFuncs, ok = validatorFuncMap[v.Kind().String()]
		if !ok {
			return fmt.Errorf("No validators found for either %s nor %s", v.Type().Name(), v.Kind().String())
		}
	}
	optionsMap, err := parseOptions(options)
	if err != nil {
		return err
	}
	for option, args := range optionsMap {

		validatorFunc, ok := validatorFuncs[options]
		if !ok {
			return fmt.Errorf("Unknown option \"%s\" for type %s", option, name)
		}
		if err := validatorFunc.Call(v, args); err != nil {
			return err
		}
	}

	return nil
}

func isExportedField(f reflect.StructField) bool {
	return strings.ToUpper(f.Name[0:1]) == f.Name[0:1]
}

func parseOptions(options string) (map[string]string, error) {
	optionsMap := map[string]string{}

	s := bufio.NewScanner(strings.NewReader(options))

	stringMode := false
	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := range data {
			if data[i] == '\'' && !stringMode {
				stringMode = true
			} else if data[i] == '\'' && stringMode {
				stringMode = false
			} else if data[i] == ',' && !stringMode {
				return i + 1, data[0:i], nil
			}
		}
		if stringMode && atEOF {
			return 0, nil, fmt.Errorf("Untermindated string")
		}
		if !stringMode && atEOF {
			return len(data), data, nil
		}

		// Slice doesn't contain a complete token, try again with more data
		return 0, nil, nil
	})

	for s.Scan() {
		token := s.Text()
		idx := strings.IndexByte(token, '=')
		if idx == -1 {
			token = token + "="
			idx = len(token) - 1
		}
		optionsMap[strings.TrimSpace(token[0:idx])] = strings.Trim(token[idx+1:], " '")
	}
	return optionsMap, s.Err()
}
