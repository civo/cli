package utility

import (
	"fmt"
	"reflect"
	"strings"
)

func mapToStringKeys(data reflect.Value) ([]string, error) {
	var keys []string

	i := 0
	iter := data.MapRange()
	for iter.Next() {
		k := iter.Key().Interface().(string)
		keys = append(keys, k)
		i++
	}

	return keys, nil
}

// FindPartialKey finds a partial match within a list of strings
func FindPartialKey(search string, data interface{}) (string, error) {
	keys, err := mapToStringKeys(reflect.ValueOf(data))
	if err != nil {
		return "", err
	}

	var result string

	for _, k := range keys {
		if strings.Contains(k, search) {
			if result == "" {
				result = k
			} else {
				return "", fmt.Errorf("unable to find %s because there were multiple matches", search)
			}
		}
	}

	if result == "" {
		return "", fmt.Errorf("unable to find %s at all in the list", search)
	}

	return result, nil
}
