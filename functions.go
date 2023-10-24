package main

import (
	"reflect"
)

var funcRegistry = make(map[string]reflect.Type)

func init() {
	funcRegistry["PathExists"] = reflect.TypeOf(PathExists{})
	funcRegistry["FileContains"] = reflect.TypeOf(FileContains{})
	funcRegistry["ServiceUp"] = reflect.TypeOf(ServiceUp{})
}

type PathExists struct {
	BaseCondition
	Path string
}

func (p PathExists) Score() bool {
	return true
}

type FileContains struct {
	BaseCondition
	File  string
	Value string
}

func (f FileContains) Score() bool {
	return true
}

type ServiceUp struct {
	BaseCondition
	Service string
}

func (f ServiceUp) Score() bool {
	return true
}

// add more...
