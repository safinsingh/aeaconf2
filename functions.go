package main

import (
	"fmt"
	"reflect"
)

var FuncRegistry = make(map[string]reflect.Type)

func init() {
	FuncRegistry["PathExists"] = reflect.TypeOf(PathExists{})
	FuncRegistry["FileContains"] = reflect.TypeOf(FileContains{})
	FuncRegistry["ServiceUp"] = reflect.TypeOf(ServiceUp{})

	CheckFunctionRegistry(FuncRegistry)
}

type PathExists struct {
	BaseCondition
	Path string
}

func (p *PathExists) Score() bool {
	return true
}

func (p *PathExists) DefaultString() string {
	return fmt.Sprintf("Path '%s' exists", p.Path)
}

type FileContains struct {
	BaseCondition
	File  string
	Value string
}

func (f *FileContains) Score() bool {
	return true
}

func (f *FileContains) DefaultString() string {
	return fmt.Sprintf("File '%s' contains '%s'", f.File, f.Value)
}

type ServiceUp struct {
	BaseCondition
	Service string
}

func (s *ServiceUp) Score() bool {
	return true
}

func (s *ServiceUp) DefaultString() string {
	return fmt.Sprintf("Service '%s' is running", s.Service)
}

// add more...
