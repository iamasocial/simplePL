package main

import "spl/parser"

type Scope struct {
	vars      map[string]string
	functions map[string]*parser.Node
	parent    *Scope
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		vars:      make(map[string]string),
		functions: make(map[string]*parser.Node),
		parent:    parent,
	}
}
