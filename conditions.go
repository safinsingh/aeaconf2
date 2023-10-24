package main

import (
	"reflect"
)

type BaseCondition struct {
	Hint string
}

// I hate go
func GetConditionHint(cond Condition) string {
	val := reflect.ValueOf(cond)
	baseCond := val.FieldByName("BaseCondition")
	if baseCond.IsValid() {
		hint := baseCond.FieldByName("Hint")
		if hint.IsValid() {
			return hint.String()
		}
	}

	panic("ICE: could not get condition hint")
}

// I still hate go
func SetConditionHint(cond Condition, newHint string) {
	val := reflect.ValueOf(cond)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	baseCond := val.FieldByName("BaseCondition")
	if baseCond.IsValid() && baseCond.CanSet() {
		hint := baseCond.FieldByName("Hint")
		if hint.IsValid() && hint.CanSet() {
			hint.SetString(newHint)
			return
		}
	}
	panic("ICE: could not set condition hint")
}

type AndExpr struct {
	BaseCondition
	Lhs Condition
	Rhs Condition
}

func (a *AndExpr) Score() bool {
	return a.Lhs.Score() && a.Rhs.Score()
}

type OrExpr struct {
	BaseCondition
	Lhs Condition
	Rhs Condition
}

func (a *OrExpr) Score() bool {
	return a.Lhs.Score() || a.Rhs.Score()
}

type NotFunc struct {
	BaseCondition
	Func Condition
}

func (n *NotFunc) Score() bool {
	return !n.Func.Score()
}
