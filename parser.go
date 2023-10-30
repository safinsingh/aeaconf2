package main

import (
	"fmt"
	"reflect"
)

type Parser struct {
	Lexer          *Lexer
	Lookahead      *Token
	LookaheadValid bool
	// currently-parsing check message; used for debugging
	CurrentCheckMessage string

	// map from function names to corresponding reflect type
	FuncRegistry map[string]reflect.Type
}

func NewParser(lexer *Lexer, funcRegistry map[string]reflect.Type) *Parser {
	return &Parser{Lexer: lexer, Lookahead: nil, LookaheadValid: false, FuncRegistry: funcRegistry}
}

func (p *Parser) Errorf(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	line, column := p.Lexer.GetSourceVisualLocation()
	Fatal(STAGE_PARSER, fmt.Sprintf("(line %d, column %d) %s", line, column, message))
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

func (p *Parser) SkipUntilNewlineBlock() {
	for {
		token := p.Peek()
		if token.Type == NewTokenline {
			p.Consume()
			token = p.Peek()
			if token.Type != TokenIndent && token.Type != NewTokenline {
				return
			}
		} else if token.Type == TokenIndent {
			p.Consume()
			token = p.Peek()
			if token.Type != NewTokenline {
				p.Errorf("expected non-indented line to begin new check")
			}
		} else if token.Type == TokenEOF {
			p.Errorf("unexpected EOF: expected non-indented line to begin new check")
		} else {
			return
		}
	}
}

func (p *Parser) SkipUntilIndentedBlock() bool {
	for {
		token := p.Peek()
		if token.Type == TokenIndent {
			p.Consume()
			if p.Peek().Type != NewTokenline {
				return true
			}
		} else if token.Type == NewTokenline {
			p.Consume()
			token = p.Peek()
			if token.Type != TokenIndent && token.Type != NewTokenline {
				return false
			}
		} else if token.Type == TokenEOF {
			p.Errorf("unexpected EOF: expected indented line to begin new condition")
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
func (p *Parser) SkipUntilIndentedBlockIfHanging() {
	if p.Peek().Type == NewTokenline {
		p.SkipUntilIndentedBlock()
	}
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

func (p *Parser) MaybeParseHint() string {
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
		return hintString.Value().(string)
	}
	return ""
}

func (p *Parser) ParseCondition() Condition {
	lhs := p.ParseAnd()
	for p.Peek().Type == TokenOr {
		p.Consume()
		p.SkipUntilIndentedBlockIfHanging()

		rhs := p.ParseAnd()
		lhs = &OrExpr{Lhs: lhs, Rhs: rhs}
	}
	return lhs
}

func (p *Parser) ParseAnd() Condition {
	lhs := p.ParseFactor()
	for p.Peek().Type == TokenAnd {
		p.Consume()
		p.SkipUntilIndentedBlockIfHanging()

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
		p.ExpectTokenType(
			TokenRParen,
			fmt.Sprintf("expected closing right parenthese for condition for check: '%s'", p.CurrentCheckMessage),
		)
	} else if next.Type == TokenIdent {
		cond = p.ParseFunc()
	} else {
		p.Errorf(
			"invalid boolean expression for check '%s': expected function or parenthesized expression, got '%s'",
			p.CurrentCheckMessage,
			next.Debug(),
		)
	}

	hint := p.MaybeParseHint()
	if len(hint) != 0 {
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
	notIdx := len(funcName) - len("Not")
	if funcName[notIdx:] == "Not" {
		notFunc = true
		funcName = funcName[:notIdx]
	}

	funcType, ok := p.FuncRegistry[funcName]
	if !ok {
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
				"expected argument for condition '%s': '%s' (arg %d out of %d total arguments)",
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
	p.SkipUntilNewlineBlock()

	// current check has empty message; avoids any "magic" generation-needing message
	currentCheckMessageEmpty := false

	// parse check name
	if p.Peek().Type == TokenUnderscore {
		p.CurrentCheckMessage = "<anonymous>"
		currentCheckMessageEmpty = true
		p.Consume()
	} else {
		p.CurrentCheckMessage = p.ExpectTokenType(
			TokenString,
			"expected check title as a string or placeholder ('_')",
		).Value().(string)
	}

	p.ExpectTokenType(
		TokenColon,
		fmt.Sprintf("expected a colon following the check message: '%s'", p.CurrentCheckMessage),
	)

	// parse points
	points := 0
	pointsEmpty := false
	// if point number isn't a placeholder
	if p.Peek().Type == TokenUnderscore {
		pointsEmpty = true
		p.Consume()
	} else {
		points = p.ExpectTokenType(
			TokenNumber,
			fmt.Sprintf("expected integer point value or placeholder ('_') to follow colon for check: '%s'",
				p.CurrentCheckMessage),
		).Value().(int)
	}

	// parse hint if it exists
	rootHint := p.MaybeParseHint()

	var finalCond Condition
	// if single-line check
	if p.Peek().Type == TokenSemicolon {
		p.Consume()
		finalCond = p.ParseCondition()
	} else {
		var andedConditions []Condition
		for p.SkipUntilIndentedBlock() {
			cond := p.ParseCondition()
			andedConditions = append(andedConditions, cond)
		}

		// no conditions parsed
		if len(andedConditions) == 0 {
			p.Errorf(
				"unexpected end of file: expected an indented condition block for check '%s'",
				p.CurrentCheckMessage,
			)
		}

		finalCond = BuildAndTree(andedConditions)
	}

	checkString := p.CurrentCheckMessage
	if currentCheckMessageEmpty {
		checkString = finalCond.DefaultString()
	}

	return &Check{Message: checkString, Points: points, PointsEmpty: pointsEmpty, Condition: finalCond, Hint: rootHint}
}

func (p *Parser) Checks() []*Check {
	var checks []*Check
	for p.Peek().Type != TokenEOF {
		checks = append(checks, p.NextCheck())
		p.SkipUntilNewlineBlock()
	}
	return checks
}
