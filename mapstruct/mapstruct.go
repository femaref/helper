package mapstruct
import (
	"errors"
	"fmt"
	"reflect"
	"encoding/json"
)

var InvalidTypeError = errors.New("Provided value type didn't match obj field type")

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
        if !val.Type().ConvertibleTo(structFieldType) {
    		return InvalidTypeError
        }

        val = val.Convert(structFieldType)
	}

	structFieldValue.Set(val)
	return nil
}

func MapToStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func StructToMap(s interface{}) (map[string]interface{}, error) {
    marsh, err := json.Marshal(s)

    if err != nil {
        return nil, err
    }

    var unmarsh map[string]interface{}

    err = json.Unmarshal(marsh, &unmarsh)

    return unmarsh, err
}

func MustStructToMap(s interface{}) map[string]interface{} {
    m, err := StructToMap(s)

    if err != nil {
        panic(err)
    }

    return m
}
