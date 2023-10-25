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

func (c *Config) DistributeMaxPoints() {
	var unspecifiedPointsChecks []*Check
	totalCheckPoints := 0
	for _, check := range c.Checks {
		totalCheckPoints += check.Points
		if check.NeedsPoints {
			unspecifiedPointsChecks = append(unspecifiedPointsChecks, check)
		}
	}

	pointsRemaining := c.Round.MaxPoints - totalCheckPoints
	pointsPerCheck := pointsRemaining / len(unspecifiedPointsChecks)

	if pointsPerCheck < 1 {
		Fatal(STAGE_DISTRIBUTION,
			fmt.Sprintf(
				"cannot distribute points to unspecified-point vulns without overflowing maximum image points (%d). %s %d",
				c.Round.MaxPoints,
				"please adjust the configuration file: increase 'maxPoints' under '[round]' to at least",
				totalCheckPoints+len(unspecifiedPointsChecks),
			),
		)
	}

	for _, check := range unspecifiedPointsChecks {
		check.Points = pointsPerCheck
		check.NeedsPoints = false
	}
}
