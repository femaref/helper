package mapstruct
import (
	"errors"
	"fmt"
	"reflect"
	"encoding/json"
	"github.com/chuckpreslar/inflect"
	"strings"
)

var InvalidTypeError = errors.New("Provided value type didn't match obj field type")

// set a single map field based on the name in the map
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

// maps a map to a struct by iterating over the map keys
func MapToStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func name(t reflect.StructField) string {
	tag := t.Tag

	value := tag.Get("json")

	if value != "" {
	    arr := strings.Split(value, ",")
	    if arr[0] != "" {
    		return value
    	}
	}

	return inflect.Underscore(t.Name)
}


// maps a map to a struct by going over the struct fields
// this works recursively, as well generating names for exported fields etc.
func MapToStructv2(s interface{}, m map[string]interface{}) error {
    // get the type in the pointer
    structValue := reflect.ValueOf(s).Elem()
    t := structValue.Type()
    // go over all fields
    for i := 0; i < t.NumField(); i++ {
        // get the field Type and the field Value
        field_t := t.Field(i)
        field_v := structValue.Field(i)

        // we can't assign to it, skip it (unexported name)
        if !field_v.CanSet() {
            continue
        }

        // generate the name of the field in the map
        field_in_map := name(field_t)


        // get the value
        v, ok := m[field_in_map]

        // couldn't find it, skip
        if !ok {
            continue
        }

        // check if the map contains a slice of map[string]interface{}
        // in this case, we expect to find a struct slice in the target struct
        slice_of_maps, slice_of_maps_ok := v.([]map[string]interface{})

        // check if we have a struct field
        if field_t.Type.Kind() == reflect.Struct {
            // check if we have a map[string]interface{} in the map to map to the struct
            st, ok := v.(map[string]interface{})

            if !ok {
                return fmt.Errorf("field %s is a struct type, map didn't contain nested map", field_t.Name)
            }

            // recurse by getting the addr of the struct
            err := MapToStructv2(field_v.Addr().Interface(), st)

            if err != nil {
                return err
            }

        // we have a slice and a slice of maps, expect struct slice
        } else if field_t.Type.Kind() == reflect.Slice && slice_of_maps_ok{
            sl_t := field_t.Type
            sle_t := field_t.Type.Elem()

            // found slice of map, but did not find slice of struct in defintion
            if sle_t.Kind() != reflect.Struct {
                return fmt.Errorf("field %s is not a slice of struct type", field_t.Name)
            }

            // create a slice with the same paramters of the map
            val := reflect.MakeSlice(sl_t, len(slice_of_maps), cap(slice_of_maps))

            // go over the input, recurse into the struct
            for i := 0; i < len(slice_of_maps); i++ {
                err := MapToStructv2(val.Index(i).Addr().Interface(), slice_of_maps[i])
                if err != nil {
                    return err
                }
            }

            field_v.Set(val)
        // we have a normal value
        } else {
            val := reflect.ValueOf(v)
            if field_t.Type != val.Type() {
                if !val.Type().ConvertibleTo(field_t.Type) {
                    return InvalidTypeError
                }

                val = val.Convert(field_t.Type)
            }

            field_v.Set(val)
        }

    }

    return nil
}


// Convert a struct to a map by json.Marshaling then json.Unmarshaling it
func StructToMap(s interface{}) (map[string]interface{}, error) {
    marsh, err := json.Marshal(s)

    if err != nil {
        return nil, err
    }

    var unmarsh map[string]interface{}

    err = json.Unmarshal(marsh, &unmarsh)

    return unmarsh, err
}

// panics if StructToMap returns an error
func MustStructToMap(s interface{}) map[string]interface{} {
    m, err := StructToMap(s)

    if err != nil {
        panic(err)
    }

    return m
}
