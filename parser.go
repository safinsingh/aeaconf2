package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Condition interface {
	Score() bool
}

func DebugCondition(cond Condition) string {
	return DebugCondition1(cond, 1)
}

func indentWrap(str string, indentStr string) string {
	return indentStr + str + indentStr + "}"
}

func DebugCondition1(cond Condition, indent int) string {
	indentStr := strings.Repeat("  ", indent)

	switch c := cond.(type) {
	case *OrExpr:
		return indentWrap(fmt.Sprintf("OR {\n%s\n%s\n",
			DebugCondition1(c.Lhs, indent+1),
			DebugCondition1(c.Rhs, indent+1)), indentStr)
	case *AndExpr:
		return indentWrap(fmt.Sprintf("AND {\n%s\n%s\n",
			DebugCondition1(c.Lhs, indent+1),
			DebugCondition1(c.Rhs, indent+1)), indentStr)
	case *NotFunc:
		return indentWrap(fmt.Sprintf("NOT {\n%s\n",
			DebugCondition1(c.Func, indent+1)), indentStr)
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

		return fmt.Sprintf("%s%s(%s)", indentStr, ty.Name(), strings.Join(parts, ", "))
	}
}

type Check struct {
	Message string
	Points  int

	// the root condition can have a hint
	Condition
}

func (c *Check) Debug() string {
	ret := fmt.Sprintf("%s : %d\n", c.Message, c.Points)
	return ret + DebugCondition(c.Condition)
}

type Parser struct {
	Lexer          *Lexer
	Lookahead      *Token
	LookaheadValid bool
	// currently-parsing check message; used for debugging
	CurrentCheck string
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{Lexer: lexer, Lookahead: nil, LookaheadValid: false, CurrentCheck: "N/A"}
}

func (p *Parser) Peek() *Token {
	if !p.LookaheadValid {
		p.Lookahead = p.Lexer.NextToken()
		p.LookaheadValid = true
	}
	return p.Lookahead
}

func (p *Parser) Consume() *Token {
	token := p.Peek()
	p.LookaheadValid = false
	return token
}

func (p *Parser) SkipTokenType(tokenType TokenType) {
	for {
		token := p.Peek()
		if token.Type != tokenType {
			break
		}
		p.Consume()
	}
}

func (p *Parser) SkipWhitespace() {
	for {
		token := p.Peek()
		if token.Type != TokenNewline && token.Type != TokenIndent {
			break
		}
		p.Consume()
	}
}

func (p *Parser) SkipUntilIndentedBlock() bool {
	for {
		token := p.Peek()
		if token.Type == TokenNewline {
			p.Consume()
			token = p.Peek()
			if token.Type != TokenNewline && token.Type != TokenIndent {
				return false
			}
		} else if token.Type == TokenIndent {
			p.Consume()
			if p.Peek().Type == TokenNewline {
				p.Consume()
			} else {
				// reached non-newline token after indent
				return true
			}
		} else {
			return true
		}
	}
}

// skip newline and indent if a check is hanging onto the next line
// (e.g. the boolean operator terminates the line)
// for example:
//
//	PathExists "/abc" ||
//	FileContains "abc" "abc"
func (p *Parser) SkipIndentIfHanging() {
	if !p.SkipUntilIndentedBlock() {
		p.Errorf("unterminated hanging boolean operator for check: '%s'", p.CurrentCheck)
	}
}

func (p *Parser) Errorf(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	line, column := GetSourceVisualLocation(p.Lexer.Source, p.Lexer.Pos)
	Fatal(STAGE_PARSER, fmt.Sprintf("(line %d, column %d) %s", line, column, message))
}

func (p *Parser) ExpectTokenType(tokenType TokenType, msg string) *Token {
	nextToken := p.Consume()
	if nextToken.Type != tokenType {
		p.Errorf("expected token of type '%s', got '%s' (type '%s'): %s",
			tokenType.Str(), nextToken.Value(), nextToken.Type.Str(), msg)
		return nil // unreachable
	}
	return nextToken
}

func (p *Parser) MaybeParseHint() (string, bool) {
	if p.Peek().Type == TokenLBracket {
		// assume LBracket ([) has been peeked
		p.Consume()
		hintString := p.ExpectTokenType(
			TokenString,
			"expected string inside after opening left-brace ([) denoting the beginning of a hint",
		)
		p.ExpectTokenType(
			TokenRBracket,
			fmt.Sprintf("expecting closing right-brace (]) for hint: %s", hintString.Value()),
		)
		return hintString.Value().(string), true
	}
	return "", false
}

func (p *Parser) ParseCondition() Condition {
	lhs := p.ParseAnd()
	for p.Peek().Type == TokenOr {
		p.Consume()
		p.SkipIndentIfHanging()

		rhs := p.ParseAnd()
		lhs = &OrExpr{Lhs: lhs, Rhs: rhs}
	}
	return lhs
}

func (p *Parser) ParseAnd() Condition {
	lhs := p.ParseFactor()
	for p.Peek().Type == TokenAnd {
		p.Consume()
		p.SkipIndentIfHanging()

		rhs := p.ParseFactor()
		lhs = &AndExpr{Lhs: lhs, Rhs: rhs}
	}
	return lhs
}

func (p *Parser) ParseFactor() Condition {
	next := p.Peek()

	var cond Condition
	if next.Type == TokenLParen {
		p.Consume()
		cond = p.ParseCondition()
		p.ExpectTokenType(TokenRParen, fmt.Sprintf("expected closing right parenthese for condition for check: '%s'", p.CurrentCheck))
	} else if next.Type == TokenIdent {
		cond = p.ParseFunc()
	} else {
		p.Errorf("unhandled token while parsing boolean expression on check '%s': %s", p.CurrentCheck, next.Debug())
	}

	if hint, ok := p.MaybeParseHint(); ok {
		SetConditionHint(cond, hint)
	}
	return cond
}

func (p *Parser) ParseFunc() Condition {
	funcName := p.Consume().Value().(string)
	if len(funcName) <= len("Not") {
		p.Errorf("invalid check name: %s", funcName)
	}

	notFunc := false
	notIdx := len(funcName) - 3
	if funcName[notIdx:] == "Not" {
		notFunc = true
		funcName = funcName[:notIdx]
	}

	funcType, found := funcRegistry[funcName]
	if !found {
		p.Errorf("invalid function name: %s", funcName)
	}
	numArgs := funcType.NumField()

	ptr := reflect.New(funcType)
	elem := ptr.Elem()

	// Skip "BaseCondition" field
	for i := 1; i < numArgs; i++ {
		field := funcType.Field(i)
		arg := p.ExpectTokenType(
			TokenString,
			fmt.Sprintf(
				"expected argument for condition '%s': '%s' (at %d out of %d total arguments)",
				funcName,
				field.Name,
				i,
				numArgs,
			),
		).Value().(string)

		elem.Field(i).SetString(arg)
	}

	fun := ptr.Interface().(Condition)

	if notFunc {
		return &NotFunc{Func: fun}
	}
	return fun
}

func (p *Parser) NextCheck() *Check {
	p.SkipWhitespace()

	checkString := p.ExpectTokenType(
		TokenString,
		"expected check title as a string",
	).Value().(string)
	p.CurrentCheck = checkString

	p.ExpectTokenType(
		TokenColon,
		fmt.Sprintf("expected a colon following the check message: '%s'", checkString),
	)

	checkPoints := p.ExpectTokenType(
		TokenNumber,
		fmt.Sprintf("expected numeric point value to follow colon for check: '%s'", checkString),
	).Value().(int)

	rootHint, hinted := p.MaybeParseHint()

	var finalCond Condition
	if p.Peek().Type == TokenSemicolon {
		p.Consume()
		finalCond = p.ParseCondition()
	} else {
		var andedConditions []Condition
		for p.SkipUntilIndentedBlock() && p.Peek().Type != TokenEOF {
			cond := p.ParseCondition()
			andedConditions = append(andedConditions, cond)
		}

		// no conditions parsed
		if len(andedConditions) == 0 {
			p.Errorf("unexpected end of file: expected an indented condition block for check '%s'", p.CurrentCheck)
		}

		finalCond = BuildAndTree(andedConditions)
	}

	if hinted {
		SetConditionHint(finalCond, rootHint)
	}

	p.CurrentCheck = "N/A"
	return &Check{Message: checkString, Points: checkPoints, Condition: finalCond}
}

func (p *Parser) Checks() []*Check {
	var checks []*Check
	for p.Peek().Type != TokenEOF {
		checks = append(checks, p.NextCheck())
	}
	return checks
}