package aeaconf2

import (
	"fmt"
	"reflect"
)

type Condition interface {
	Score() bool
	// for autogenerated check messaged
	DefaultString() string
}

type BaseCondition struct {
	Hint string
}

type AndExpr struct {
	BaseCondition
	Lhs Condition
	Rhs Condition
}

func (a *AndExpr) Score() bool {
	return a.Lhs.Score() && a.Rhs.Score()
}

func (a *AndExpr) DefaultString() string {
	return fmt.Sprintf("(%s AND %s)", a.Lhs.DefaultString(), a.Rhs.DefaultString())
}

type OrExpr struct {
	BaseCondition
	Lhs Condition
	Rhs Condition
}

func (o *OrExpr) Score() bool {
	return o.Lhs.Score() || o.Rhs.Score()
}

func (o *OrExpr) DefaultString() string {
	return fmt.Sprintf("(%s OR %s)", o.Lhs.DefaultString(), o.Rhs.DefaultString())
}

type NotFunc struct {
	BaseCondition
	// Func is always a function call
	Func Condition
}

func (n *NotFunc) Score() bool {
	return !n.Func.Score()
}

func (n *NotFunc) DefaultString() string {
	return fmt.Sprintf("NOT (%s)", n.Func.DefaultString())
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
