package main

import (
	"fmt"
	"reflect"
)

var funcRegistry = make(map[string]reflect.Type)

func init() {
	funcRegistry["PathExists"] = reflect.TypeOf(PathExists{})
	funcRegistry["FileContains"] = reflect.TypeOf(FileContains{})
	funcRegistry["ServiceUp"] = reflect.TypeOf(ServiceUp{})

	CheckFunctionRegistry(funcRegistry)
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
