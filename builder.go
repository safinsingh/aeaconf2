package main

import (
	"bytes"
	"reflect"
)

type AeaconfBuilder struct {
	ChecksRaw    []byte
	FuncRegistry map[string]reflect.Type
	MaxPoints    int
	LineOffset   int
}

func NewAeaconfBuilder() *AeaconfBuilder {
	return &AeaconfBuilder{}
}

func DefaultAeaconfBuilder(checksRaw []byte, funcRegistry map[string]reflect.Type) *AeaconfBuilder {
	return &AeaconfBuilder{ChecksRaw: checksRaw, FuncRegistry: funcRegistry, MaxPoints: 100, LineOffset: 0}
}

func (a *AeaconfBuilder) SetChecksRaw(checksRaw []byte) *AeaconfBuilder {
	a.ChecksRaw = checksRaw
	return a
}

func (a *AeaconfBuilder) SetFuncRegistry(funcRegistry map[string]reflect.Type) *AeaconfBuilder {
	a.FuncRegistry = funcRegistry
	return a
}

func (a *AeaconfBuilder) SetMaxPoints(maxPoints int) *AeaconfBuilder {
	a.MaxPoints = maxPoints
	return a
}

func (a *AeaconfBuilder) SetLineOffset(lineOffset int) *AeaconfBuilder {
	a.LineOffset = lineOffset
	return a
}

func (a *AeaconfBuilder) GetChecks() []*Check {
	l := NewLexer(bytes.TrimSpace(a.ChecksRaw), a.LineOffset)
	p := NewParser(l, a.FuncRegistry)
	checks := p.Checks()
	DistributeMaxPoints(checks, a.MaxPoints)

	return checks
}
