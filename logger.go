package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/fatih/color"
)

type CompilerStage int

const (
	STAGE_INI CompilerStage = iota
	STAGE_LEXER
	STAGE_PARSER
	STAGE_DISTRIBUTION
)

func Fatal(stage CompilerStage, message string) {
	var stageStr string
	switch stage {
	case STAGE_INI:
		stageStr = "ini parser"
	case STAGE_LEXER:
		stageStr = "lexer"
	case STAGE_PARSER:
		stageStr = "parser"
	case STAGE_DISTRIBUTION:
		stageStr = "point distribution"
	}

	log.Fatalf("[%s] FATAL: %s", stageStr, message)
}

func DebugCondition(cond Condition) string {
	return DebugCondition1(cond, 0)
}

var colorArray []color.Attribute = []color.Attribute{color.FgGreen, color.FgYellow, color.FgMagenta, color.FgCyan}

func formatHint(hint string) string {
	if hint == "" {
		return ""
	}
	muted := color.New(color.FgHiBlack).SprintFunc()
	return muted(" HINT [ " + hint + " ]")
}

func DebugCondition1(cond Condition, indentLevel int) string {
	indent := strings.Repeat("  ", indentLevel)
	col := color.New(colorArray[indentLevel%len(colorArray)]).SprintFunc()

	switch c := cond.(type) {
	case *OrExpr:
		return fmt.Sprintf(
			"%s \n%s \n%s \n%s%s",
			col(indent+"OR {"),
			DebugCondition1(c.Lhs, indentLevel+1)+col(","),
			DebugCondition1(c.Rhs, indentLevel+1),
			col(indent+"}"),
			formatHint(c.Hint),
		)
	case *AndExpr:
		return fmt.Sprintf(
			"%s \n%s \n%s \n%s%s",
			col(indent+"AND {"),
			DebugCondition1(c.Lhs, indentLevel+1)+col(","),
			DebugCondition1(c.Rhs, indentLevel+1),
			col(indent+"}"),
			formatHint(c.Hint),
		)
	case *NotFunc:
		return fmt.Sprintf(
			"%s \n%s \n%s%s",
			col(indent+"NOT {"),
			DebugCondition1(c.Func, indentLevel+1),
			col(indent+"}"),
			formatHint(c.Hint),
		)
	default:
		val := reflect.ValueOf(cond)
		ty := val.Type()

		if val.Kind() == reflect.Ptr && !val.IsNil() {
			val = val.Elem()
			ty = val.Type()
		}

		var parts []string
		for i := 1; i < val.NumField(); i++ {
			field := ty.Field(i)
			value := val.Field(i)
			parts = append(parts, fmt.Sprintf("%s=\"%v\"", field.Name, value))
		}

		// BaseCondition is always field 0 (guaranteed)
		hint := val.Field(0).FieldByName("Hint").String()
		return fmt.Sprintf("%s%s(%s)%s", indent, ty.Name(), strings.Join(parts, ", "), formatHint(hint))
	}
}
