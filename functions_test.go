package aeaconf2_test

import (
	"fmt"
	"reflect"

	"github.com/safinsingh/aeaconf2"
)

func getFunctionRegistry() map[string]reflect.Type {
	funcRegistry := make(map[string]reflect.Type)

	funcRegistry["PathExists"] = reflect.TypeOf(PathExists{})
	funcRegistry["FileContains"] = reflect.TypeOf(FileContains{})
	funcRegistry["ServiceUp"] = reflect.TypeOf(ServiceUp{})

	aeaconf2.CheckFunctionRegistry(funcRegistry)
	return funcRegistry
}

type PathExists struct {
	aeaconf2.BaseCondition
	Path string
}

func (p *PathExists) Score() (bool, error) {
	return true, nil
}

func (p *PathExists) DefaultString() string {
	return fmt.Sprintf("Path '%s' exists", p.Path)
}

type FileContains struct {
	aeaconf2.BaseCondition
	File  string
	Value string
}

func (f *FileContains) Score() (bool, error) {
	return true, nil
}

func (f *FileContains) DefaultString() string {
	return fmt.Sprintf("File '%s' contains '%s'", f.File, f.Value)
}

type ServiceUp struct {
	aeaconf2.BaseCondition
	Service string
}

func (s *ServiceUp) Score() (bool, error) {
	return true, nil
}

func (s *ServiceUp) DefaultString() string {
	return fmt.Sprintf("Service '%s' is running", s.Service)
}

// add more...
