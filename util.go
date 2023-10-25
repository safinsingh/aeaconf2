package main

import (
	"fmt"
	"reflect"
)

func GetSourceVisualLocation(source []byte, pos int) (int, int) {
	line := 1
	column := 1
	for i := 0; i < pos; i++ {
		if source[i] == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return line, column
}

func BuildAndTree(conditions []Condition) Condition {
	var result Condition

	for _, cond := range conditions {
		if result == nil {
			result = cond
		} else {
			result = &AndExpr{Lhs: result, Rhs: cond}
		}
	}

	return result
}

func CheckFunctionRegistry(funcs map[string]reflect.Type) {
	for funcName, ty := range funcs {
		if ty.NumField() == 0 {
			panic(fmt.Sprintf("ICE: function '%s' has invalid # of arguments: 0 (must include BaseCondition)", funcName))
		}

		field0 := ty.Field(0)
		if field0.Type != reflect.TypeOf(BaseCondition{}) {
			panic(fmt.Sprintf(
				"ICE: function '%s' has invalid first struct field '%s' (type '%s'): must be BaseCondition",
				funcName,
				field0.Name,
				field0.Type))
		}
	}
}
