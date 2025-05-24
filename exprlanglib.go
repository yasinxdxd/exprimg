package main

func varx() *expr {
	res := create_expr_terminal(EXPR_KIND_VAR_X)
	return &res
}

func vary() *expr {
	res := create_expr_terminal(EXPR_KIND_VAR_Y)
	return &res
}

func number(val float32) *expr {
	res := create_expr_terminal(EXPR_KIND_NUMBER)
	res.terminal_expr.value = val
	return &res
}

func special_constant(val float32) *expr {
	res := create_expr_terminal(EXPR_KIND_SPECIAL_CONSTANT)
	res.terminal_expr.value = val
	return &res
}

func sin(x *expr) *expr {
	res := create_expr_single(EXPR_KIND_SIN)
	res.single_expr.arg1 = x
	return &res
}

func cos(x *expr) *expr {
	res := create_expr_single(EXPR_KIND_COS)
	res.single_expr.arg1 = x
	return &res
}

func vec2(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_VEC2)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func add(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_ADD)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func sub(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_SUB)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func mult(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_MULT)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func div(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_DIV)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func eq(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_EQ)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func neq(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_NEQ)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func grater(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_GRATER)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func gratereq(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_GRATEREQ)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func less(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_LESS)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func lesseq(x *expr, y *expr) *expr {
	res := create_expr_binop(EXPR_KIND_LESSEQ)
	res.binop_expr.arg1 = x
	res.binop_expr.arg2 = y
	return &res
}

func if_then_else(x *expr, y *expr, z *expr) *expr {
	res := create_expr_ternary(EXPR_KIND_IF_THEN_ELSE)
	res.ternary_expr.arg1 = x
	res.ternary_expr.arg2 = y
	res.ternary_expr.arg3 = z
	return &res
}

func vec3(x *expr, y *expr, z *expr) *expr {
	res := create_expr_ternary(EXPR_KIND_VEC3)
	res.ternary_expr.arg1 = x
	res.ternary_expr.arg2 = y
	res.ternary_expr.arg3 = z
	return &res
}
