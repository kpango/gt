package gt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func IFtoa(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}
	switch v := value.(type) {
	case fmt.Stringer:
		return v.String(), nil
	case error:
		return v.Error(), nil
	case string:
		return v, nil
	case bool:
		if v {
			return "true", nil
		}
		return "false", nil
	case []byte:
		return string(v), nil
	case []rune:
		return string(v), nil
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'E', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'E', -1, 64), nil
	case complex64:
		return strconv.FormatComplex(complex128(v), 'E', -1, 64), nil
	case complex128:
		return strconv.FormatComplex(v, 'E', -1, 128), nil
	case json.Marshaler:
		b, err := v.MarshalJSON()
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return fmt.Sprint(value), nil
	}
	v := reflect.ValueOf(value)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return IFtoa(v.Interface())
}
