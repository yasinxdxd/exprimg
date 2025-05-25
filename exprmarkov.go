package main

import (
	"math"
	"math/rand"
)

type weighted_expr struct {
	kind EXPR_KIND
	prob float32
}

var terminal_kinds = []weighted_expr{
	{EXPR_KIND_VAR_X, 0.15},
	{EXPR_KIND_VAR_Y, 0.15},
	{EXPR_KIND_VAR_T, 0.5},
	{EXPR_KIND_SPECIAL_CONSTANT, 0.1},
	{EXPR_KIND_NUMBER, 0.1},
}

var single_kinds = []weighted_expr{
	{EXPR_KIND_SIN, 0.5},
	{EXPR_KIND_COS, 0.5},
}

var binop_kinds = []weighted_expr{
	{EXPR_KIND_ADD, 0.2},
	{EXPR_KIND_SUB, 0.1},
	{EXPR_KIND_MULT, 0.2},
	{EXPR_KIND_DIV, 0.1},
	{EXPR_KIND_MOD, 0.2},
	{EXPR_KIND_GRATER, 0.05},
	{EXPR_KIND_GRATEREQ, 0.05},
	{EXPR_KIND_LESS, 0.05},
	{EXPR_KIND_LESSEQ, 0.05},
}

var ternary_kinds = []weighted_expr{
	{EXPR_KIND_IF_THEN_ELSE, 1.0}, //0.5
	// {EXPR_KIND_VEC3, 0.5},
}

// Expression category
type exprCategory int

const (
	CATEGORY_TERMINAL exprCategory = iota
	CATEGORY_SINGLE
	CATEGORY_BINOP
	CATEGORY_TERNARY
)

type category_weight struct {
	category exprCategory
	weight   float32
}

var category_weights = []category_weight{
	{CATEGORY_TERMINAL, 0.25},
	{CATEGORY_SINGLE, 0.25},
	{CATEGORY_BINOP, 0.5},
	{CATEGORY_TERNARY, 0.0},
}

func pick_weighted_expr(kinds []weighted_expr) EXPR_KIND {
	r := rand.Float32()
	acc := float32(0)
	for _, e := range kinds {
		acc += e.prob
		if r < acc {
			return e.kind
		}
	}
	return kinds[len(kinds)-1].kind
}

func pick_weighted_category(weights []category_weight) exprCategory {
	total := float32(0)
	for _, w := range weights {
		total += w.weight
	}
	r := rand.Float32() * total
	acc := float32(0)
	for _, w := range weights {
		acc += w.weight
		if r < acc {
			return w.category
		}
	}
	return weights[len(weights)-1].category
}

func generate_expr(depth int) *expr {
	if depth <= 0 {
		kind := pick_weighted_expr(terminal_kinds)
		switch kind {
		case EXPR_KIND_VAR_X:
			return varx()
		case EXPR_KIND_VAR_Y:
			return vary()
		case EXPR_KIND_VAR_T:
			return vart()
		case EXPR_KIND_NUMBER:
			return number((rand.Float32() * 2) - 1)
		case EXPR_KIND_SPECIAL_CONSTANT:
			constants := [...]float32{math.Pi, math.SqrtPi, math.E, math.Log10E, math.Phi, 0.5, 1, 2, 10}
			value := constants[rand.Intn(len(constants))]
			return special_constant(value)
		default:
			panic("unknown terminal kind")
		}
	}

	category := pick_weighted_category(category_weights)

	switch category {
	case CATEGORY_TERMINAL:
		return generate_expr(0)
	case CATEGORY_SINGLE:
		kind := pick_weighted_expr(single_kinds)
		arg := generate_expr(depth - 1)
		switch kind {
		case EXPR_KIND_SIN:
			return sin(arg)
		case EXPR_KIND_COS:
			return cos(arg)
		default:
			panic("unknown single kind")
		}
	case CATEGORY_BINOP:
		kind := pick_weighted_expr(binop_kinds)
		arg1 := generate_expr(depth - 1)
		arg2 := generate_expr(depth - 1)
		switch kind {
		case EXPR_KIND_ADD:
			return add(arg1, arg2)
		case EXPR_KIND_SUB:
			return sub(arg1, arg2)
		case EXPR_KIND_MULT:
			return mult(arg1, arg2)
		case EXPR_KIND_DIV:
			return div(arg1, arg2)
		case EXPR_KIND_MOD:
			return less(arg1, arg2)
		case EXPR_KIND_GRATER:
			return grater(arg1, arg2)
		case EXPR_KIND_LESS:
			return less(arg1, arg2)
		case EXPR_KIND_GRATEREQ:
			return gratereq(arg1, arg2)
		case EXPR_KIND_LESSEQ:
			return lesseq(arg1, arg2)
		default:
			panic("unknown binop kind")
		}
	case CATEGORY_TERNARY:
		kind := pick_weighted_expr(ternary_kinds)
		arg1 := generate_expr(depth - 1)
		arg2 := generate_expr(depth - 1)
		arg3 := generate_expr(depth - 1)
		switch kind {
		case EXPR_KIND_IF_THEN_ELSE:
			return if_then_else(arg1, arg2, arg3)
		case EXPR_KIND_VEC3:
			return vec3(arg1, arg2, arg3)
		default:
			panic("unknown ternary kind")
		}
	default:
		panic("invalid category")
	}
}

func generate_expr_root(depth int) *expr {
	return vec3(
		generate_expr(depth-1),
		generate_expr(depth-1),
		generate_expr(depth-1),
	)
}
