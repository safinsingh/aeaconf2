package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type TokenType int

const (
	TokenIndent TokenType = iota
	TokenNewline

	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket

	TokenColon
	TokenSemicolon
	TokenUnderscore

	TokenAnd
	TokenOr

	TokenString
	TokenNumber
	TokenIdent

	TokenEOF
)

func (ty TokenType) Str() string {
	switch ty {
	case TokenNewline:
		return "TokenNewline"
	case TokenIndent:
		return "TokenIndent"
	case TokenLParen:
		return "TokenLParen"
	case TokenRParen:
		return "TokenRParen"
	case TokenLBracket:
		return "TokenLBracket"
	case TokenRBracket:
		return "TokenRBracket"
	case TokenColon:
		return "TokenColon"
	case TokenSemicolon:
		return "TokenSemicolon"
	case TokenUnderscore:
		return "TokenUnderscore"
	case TokenAnd:
		return "TokenAnd"
	case TokenOr:
		return "TokenOr"
	case TokenIdent:
		return "TokenIdent"
	case TokenString:
		return "TokenString"
	case TokenNumber:
		return "TokenNumber"
	case TokenEOF:
		return "TokenEOF"
	default:
		panic("unknown token type")
	}
}

type Token struct {
	Type   TokenType
	Lexeme []byte
}

func NewToken(tokenType TokenType, lexeme []byte) *Token {
	return &Token{Type: tokenType, Lexeme: lexeme}
}

func (t *Token) Debug() string {
	if t.Type == TokenNewline {
		return "Token{type: TokenNewline, lexeme: '\\n'}"
	}
	return fmt.Sprintf("Token{type: %s, lexeme: '%s'}", t.Type.Str(), string(t.Lexeme))
}

func (t *Token) Value() any {
	if t.Type == TokenString {
		return string(t.Lexeme[1 : len(t.Lexeme)-1])
	} else if t.Type == TokenNumber {
		num, err := strconv.Atoi(string(t.Lexeme))
		if err != nil {
			// should never happen
			panic("invalid number")
		}
		return num
	} else {
		return string(t.Lexeme)
	}
}

type Lexer struct {
	Source     []byte
	Pos        int
	LineOffset int
}

func NewLexer(source []byte, lineOffset int) *Lexer {
	return &Lexer{Source: source, Pos: 0, LineOffset: lineOffset}
}

func (l *Lexer) Errorf(format string, a ...any) {
	message := fmt.Sprintf(format, a...)
	line, column := l.GetSourceVisualLocation()
	Fatal(STAGE_LEXER, fmt.Sprintf("(line %d, column %d) %s", line, column, message))
}

func (l *Lexer) ExpectCharacter(ch byte) {
	if l.Pos < len(l.Source) {
		if l.Source[l.Pos] == ch {
			l.Pos++
		} else {
			l.Errorf("unexpected character: '%c', expected character '%c'", l.Source[l.Pos], ch)
		}
	} else {
		l.Errorf("unexpected end of file, expected character '%c'", ch)
	}
}

func (l *Lexer) AdvanceToken(tokenType TokenType, lexeme byte) *Token {
	l.Pos++
	return NewToken(tokenType, []byte{lexeme})
}

func (l *Lexer) AdvanceToken2(tokenType TokenType, currentLexeme byte, expect byte) *Token {
	l.Pos++
	l.ExpectCharacter(expect)
	return NewToken(tokenType, []byte{currentLexeme, expect})
}

func (l *Lexer) LexString(ch byte) *Token {
	initialPos := l.Pos
	l.Pos++
	for l.Pos < len(l.Source) {
		if l.Source[l.Pos] == ch {
			break
		}
		if l.Source[l.Pos] == '\\' {
			l.Pos++
			if l.Pos >= len(l.Source) {
				l.Errorf("string not terminated, expected '%c'", ch)
			}
			// Only escape for double quotes
			if l.Source[l.Pos] == '"' {
				l.Pos++
			}
		} else {
			l.Pos++
		}
	}
	if l.Pos >= len(l.Source) {
		l.Errorf("string not terminated, expected '%c'", ch)
	}
	l.Pos++
	return NewToken(TokenString, l.Source[initialPos:l.Pos])
}

func (l *Lexer) LexIdent() *Token {
	initialPos := l.Pos
	l.Pos++
	for l.Pos < len(l.Source) && !unicode.IsSpace(rune(l.Source[l.Pos])) {
		l.Pos++
	}
	return NewToken(TokenIdent, l.Source[initialPos:l.Pos])
}

func (l *Lexer) LexNumber() *Token {
	initialPos := l.Pos
	l.Pos++
	for l.Pos < len(l.Source) && unicode.IsNumber(rune(l.Source[l.Pos])) {
		l.Pos++
	}
	return NewToken(TokenNumber, l.Source[initialPos:l.Pos])
}

func (l *Lexer) NextToken() *Token {
	if l.Pos >= len(l.Source) {
		return NewToken(TokenEOF, []byte{})
	}

	if l.Pos == 0 || l.Source[l.Pos-1] == '\n' {
		initialPos := l.Pos
		for l.Pos < len(l.Source) && (l.Source[l.Pos] == '\t' || l.Source[l.Pos] == ' ') {
			l.Pos++
		}
		if l.Pos-initialPos > 0 {
			return NewToken(TokenIndent, l.Source[initialPos:l.Pos])
		}
	}

	ch := l.Source[l.Pos]

	switch ch {
	case '\n':
		return l.AdvanceToken(TokenNewline, ch)
	case '(':
		return l.AdvanceToken(TokenLParen, ch)
	case ')':
		return l.AdvanceToken(TokenRParen, ch)
	case '[':
		return l.AdvanceToken(TokenLBracket, ch)
	case ']':
		return l.AdvanceToken(TokenRBracket, ch)
	case ':':
		return l.AdvanceToken(TokenColon, ch)
	case ';':
		return l.AdvanceToken(TokenSemicolon, ch)
	case '_':
		return l.AdvanceToken(TokenUnderscore, ch)
	case '&':
		return l.AdvanceToken2(TokenAnd, ch, '&')
	case '|':
		return l.AdvanceToken2(TokenOr, ch, '|')
	case '"', '\'':
		return l.LexString(ch)
	}

	if unicode.IsLetter(rune(ch)) {
		return l.LexIdent()
	}

	if ch == '-' || unicode.IsNumber(rune(ch)) {
		return l.LexNumber()
	}

	// skip comments entirely
	if ch == '/' {
		if l.Pos+1 < len(l.Source) && l.Source[l.Pos+1] == '/' {
			l.Pos += 2 // skip second '/'
			for l.Pos < len(l.Source) && l.Source[l.Pos] != '\n' {
				l.Pos++
			}
			return l.NextToken()
		}
	}

	// Any "space" character besides a newline (which is already handled)
	if unicode.IsSpace(rune(ch)) {
		l.Pos++
		return l.NextToken()
	}

	l.Errorf("unhandled character : '%c'", ch)
	return nil
}
