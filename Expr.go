package main

import (
	"fmt"
)

type object interface{}

type Visitor interface {
	visitBinary(Binary) (object, error)
	visitGrouping(Grouping) (object, error)
	visitLiteral(Literal) (object, error)
	visitUnary(Unary) (object, error)
}

type Expr interface {
	accept(Visitor) (object, error)
	print() string
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (t Binary) accept(v Visitor) (object, error) {
	return v.visitBinary(t)
}

func (t Binary) print() string {
	return fmt.Sprintf("(%s %s %s)", t.operator.lexeme, t.left.print(), t.right.print())
}

type Grouping struct {
	expression Expr
}

func (t Grouping) accept(v Visitor) (object, error) {
	return v.visitGrouping(t)
}

func (t Grouping) print() string {
	return fmt.Sprintf("(grouping %s)", t.expression.print())
}

type Unary struct {
	operator Token
	right    Expr
}

func (t Unary) accept(v Visitor) (object, error) {
	return v.visitUnary(t)
}

func (t Unary) print() string {
	return fmt.Sprintf("(%s %s)", t.operator.lexeme, t.right.print())
}

type Literal struct {
	value object
}

func (t Literal) accept(v Visitor) (object, error) {
	return v.visitLiteral(t)
}

func (t Literal) print() string {
	return fmt.Sprintf("%v", t.value)
}
