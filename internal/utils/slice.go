package utils

import (
	"reflect"

	"github.com/Lapakin/edu-planner/internal/adapter/json"
)

func Copy(src, dest any) {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Pointer || destValue.IsNil() {
		panic("destination must be a non-nil pointer")
	}
	jsonData, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(jsonData, dest); err != nil {
		panic(err)
	}
}
