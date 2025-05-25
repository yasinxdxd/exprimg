package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
)

type EXPR_KIND int

const (
	EXPR_KIND_INVALID EXPR_KIND = iota

	TERMINAL_BEGIN
	EXPR_KIND_NUMBER           // terminal
	EXPR_KIND_SPECIAL_CONSTANT // terminal
	EXPR_KIND_VAR_X            // terminal
	EXPR_KIND_VAR_Y            // terminal
	TERMINAL_END

	SINGLE_BEGIN
	// single
	EXPR_KIND_SIN
	EXPR_KIND_COS
	EXPR_KIND_SQRT
	SINGLE_END

	BINOP_BEGIN
	// binop
	EXPR_KIND_VEC2
	EXPR_KIND_ADD
	EXPR_KIND_SUB
	EXPR_KIND_MULT
	EXPR_KIND_DIV
	EXPR_KIND_MOD
	EXPR_KIND_POW
	EXPR_KIND_EQ
	EXPR_KIND_NEQ
	EXPR_KIND_GRATER
	EXPR_KIND_LESS
	EXPR_KIND_GRATEREQ
	EXPR_KIND_LESSEQ
	EXPR_KIND_DOT
	EXPR_KIND_CROSS
	BINOP_END

	TERNARY_BEGIN
	// ternary
	EXPR_KIND_VEC3
	EXPR_KIND_CLAMP
	EXPR_KIND_MIX
	EXPR_KIND_IF_THEN_ELSE
	TERNARY_END
)

func expr_kind_to_str(kind EXPR_KIND) string {
	switch kind {
	case EXPR_KIND_NUMBER:
		return "EXPR_KIND_NUMBER"
	case EXPR_KIND_SPECIAL_CONSTANT:
		return "EXPR_KIND_SPECIAL_CONSTANT"
	case EXPR_KIND_VAR_X:
		return "EXPR_KIND_VAR_X"
	case EXPR_KIND_VAR_Y:
		return "EXPR_KIND_VAR_Y"
	case EXPR_KIND_SIN:
		return "EXPR_KIND_SIN"
	case EXPR_KIND_COS:
		return "EXPR_KIND_COS"
	case EXPR_KIND_SQRT:
		return "EXPR_KIND_SQRT"
	case EXPR_KIND_VEC2:
		return "EXPR_KIND_VEC2"
	case EXPR_KIND_ADD:
		return "EXPR_KIND_ADD"
	case EXPR_KIND_SUB:
		return "EXPR_KIND_SUB"
	case EXPR_KIND_MULT:
		return "EXPR_KIND_MULT"
	case EXPR_KIND_DIV:
		return "EXPR_KIND_DIV"
	case EXPR_KIND_MOD:
		return "EXPR_KIND_MOD"
	case EXPR_KIND_POW:
		return "EXPR_KIND_POW"
	case EXPR_KIND_EQ:
		return "EXPR_KIND_EQ"
	case EXPR_KIND_NEQ:
		return "EXPR_KIND_NEQ"
	case EXPR_KIND_GRATER:
		return "EXPR_KIND_GRATER"
	case EXPR_KIND_LESS:
		return "EXPR_KIND_LESS"
	case EXPR_KIND_GRATEREQ:
		return "EXPR_KIND_GRATEREQ"
	case EXPR_KIND_LESSEQ:
		return "EXPR_KIND_LESSEQ"
	case EXPR_KIND_DOT:
		return "EXPR_KIND_DOT"
	case EXPR_KIND_CROSS:
		return "EXPR_KIND_CROSS"
	case EXPR_KIND_VEC3:
		return "EXPR_KIND_VEC3"
	case EXPR_KIND_CLAMP:
		return "EXPR_KIND_CLAMP"
	case EXPR_KIND_MIX:
		return "EXPR_KIND_MIX"
	case EXPR_KIND_IF_THEN_ELSE:
		return "EXPR_KIND_IF_THEN_ELSE"
	default:
		_ = fmt.Errorf("runtime error: %d is not a valid kind for expr_kind_to_str()", kind)
		return ""
	}
}

type expr struct {
	kind EXPR_KIND
	prob float32
	gen  func(code string) string

	single_expr struct {
		arg1 *expr
	}
	binop_expr struct {
		arg1 *expr
		arg2 *expr
	}
	ternary_expr struct {
		arg1 *expr
		arg2 *expr
		arg3 *expr
	}
	terminal_expr struct {
		value float32
	}
}

func is_terminal(e expr) bool {
	if e.kind > TERMINAL_BEGIN && e.kind < TERMINAL_END {
		return true
	}
	return false
}

func is_single(e expr) bool {
	if e.kind > SINGLE_BEGIN && e.kind < SINGLE_END {
		return true
	}
	return false
}

func is_binop(e expr) bool {
	if e.kind > BINOP_BEGIN && e.kind < BINOP_END {
		return true
	}
	return false
}

func is_ternary(e expr) bool {
	if e.kind > TERNARY_BEGIN && e.kind < TERNARY_END {
		return true
	}
	return false
}

func expr_valid(kind EXPR_KIND, valid_kinds ...EXPR_KIND) bool {
	for _, k := range valid_kinds {
		if kind == k {
			return true
		}
	}
	return false
}

func create_expr_terminal(kind EXPR_KIND) expr {
	if !expr_valid(kind, EXPR_KIND_VAR_X, EXPR_KIND_VAR_Y, EXPR_KIND_NUMBER, EXPR_KIND_SPECIAL_CONSTANT) {
		err := fmt.Errorf("%d is not a valid kind for create_expr_terminal()", kind)
		fmt.Println(err.Error())
		return expr{kind: EXPR_KIND_INVALID}
	}

	switch kind {
	case EXPR_KIND_VAR_X:
		break
	case EXPR_KIND_VAR_Y:
		break
	case EXPR_KIND_NUMBER:
		break
	case EXPR_KIND_SPECIAL_CONSTANT:
		break
	}

	e := expr{
		kind:          kind,
		terminal_expr: struct{ value float32 }{value: 0.0},
	}
	return e
}

func create_expr_single(kind EXPR_KIND) expr {
	if !expr_valid(kind, EXPR_KIND_SIN, EXPR_KIND_COS, EXPR_KIND_SQRT) {
		err := fmt.Errorf("%d is not a valid kind for create_expr_single()", kind)
		fmt.Println(err.Error())
		return expr{kind: EXPR_KIND_INVALID}
	}

	switch kind {
	case EXPR_KIND_SIN:
		fallthrough
	case EXPR_KIND_COS:
		fallthrough
	case EXPR_KIND_SQRT:
		{
			e := expr{
				kind: kind,
				single_expr: struct {
					arg1 *expr
				}{
					arg1: new(expr),
				},
			}
			return e
		}
	}
	return expr{kind: EXPR_KIND_INVALID}
}

func create_expr_binop(kind EXPR_KIND) expr {
	if !expr_valid(kind,
		EXPR_KIND_ADD,
		EXPR_KIND_SUB,
		EXPR_KIND_MULT,
		EXPR_KIND_DIV,
		EXPR_KIND_MOD,
		EXPR_KIND_POW,
		EXPR_KIND_EQ,
		EXPR_KIND_NEQ,
		EXPR_KIND_GRATER,
		EXPR_KIND_LESS,
		EXPR_KIND_GRATEREQ,
		EXPR_KIND_LESSEQ,
		EXPR_KIND_DOT,
		EXPR_KIND_CROSS,
	) {
		err := fmt.Errorf("%d is not a valid kind for create_expr_binop()", kind)
		fmt.Println(err.Error())
		return expr{kind: EXPR_KIND_INVALID}
	}

	switch kind {
	case EXPR_KIND_ADD:
		fallthrough
	case EXPR_KIND_SUB:
		fallthrough
	case EXPR_KIND_MULT:
		fallthrough
	case EXPR_KIND_DIV:
		fallthrough
	case EXPR_KIND_EQ:
		fallthrough
	case EXPR_KIND_NEQ:
		fallthrough
	case EXPR_KIND_GRATER:
		fallthrough
	case EXPR_KIND_LESS:
		fallthrough
	case EXPR_KIND_GRATEREQ:
		fallthrough
	case EXPR_KIND_LESSEQ:
		fallthrough
	case EXPR_KIND_DOT:
		fallthrough
	case EXPR_KIND_CROSS:
		{
			e := expr{
				kind: kind,
				binop_expr: struct {
					arg1 *expr
					arg2 *expr
				}{
					arg1: new(expr),
					arg2: new(expr),
				},
			}
			return e
		}
	}
	return expr{kind: EXPR_KIND_INVALID}
}

func create_expr_ternary(kind EXPR_KIND) expr {
	if !expr_valid(kind, EXPR_KIND_VEC3, EXPR_KIND_CLAMP, EXPR_KIND_MIX, EXPR_KIND_IF_THEN_ELSE) {
		err := fmt.Errorf("%d is not a valid kind for create_expr_ternary()", kind)
		fmt.Println(err.Error())
		return expr{kind: EXPR_KIND_INVALID}
	}

	switch kind {
	case EXPR_KIND_VEC3:
		fallthrough
	case EXPR_KIND_CLAMP:
		fallthrough
	case EXPR_KIND_MIX:
		fallthrough
	case EXPR_KIND_IF_THEN_ELSE:
		{
			e := expr{
				kind: kind,
				ternary_expr: struct {
					arg1 *expr
					arg2 *expr
					arg3 *expr
				}{
					arg1: new(expr),
					arg2: new(expr),
					arg3: new(expr),
				},
			}
			return e
		}
	}
	return expr{kind: EXPR_KIND_INVALID}
}

func expr_stack(stack []*expr, root *expr) []*expr {
	if is_terminal(*root) {
		stack = append(stack, root)
		return stack
	}
	switch root.kind {
	// single:
	case EXPR_KIND_SIN:
		fallthrough
	case EXPR_KIND_COS:
		fallthrough
	case EXPR_KIND_SQRT:
		stack = append(stack, root)
		stack = expr_stack(stack, root.single_expr.arg1)
		break

	// binop:
	case EXPR_KIND_VEC2:
		fallthrough
	case EXPR_KIND_ADD:
		fallthrough
	case EXPR_KIND_SUB:
		fallthrough
	case EXPR_KIND_MULT:
		fallthrough
	case EXPR_KIND_DIV:
		fallthrough
	case EXPR_KIND_MOD:
		fallthrough
	case EXPR_KIND_POW:
		fallthrough
	case EXPR_KIND_EQ:
		fallthrough
	case EXPR_KIND_NEQ:
		fallthrough
	case EXPR_KIND_GRATER:
		fallthrough
	case EXPR_KIND_LESS:
		fallthrough
	case EXPR_KIND_GRATEREQ:
		fallthrough
	case EXPR_KIND_LESSEQ:
		fallthrough
	case EXPR_KIND_DOT:
		fallthrough
	case EXPR_KIND_CROSS:
		stack = append(stack, root)
		stack = expr_stack(stack, root.binop_expr.arg1)
		stack = expr_stack(stack, root.binop_expr.arg2)
		break

	// ternary:
	case EXPR_KIND_VEC3:
		fallthrough
	case EXPR_KIND_CLAMP:
		fallthrough
	case EXPR_KIND_MIX:
		fallthrough
	case EXPR_KIND_IF_THEN_ELSE:
		stack = append(stack, root)
		stack = expr_stack(stack, root.ternary_expr.arg1)
		stack = expr_stack(stack, root.ternary_expr.arg2)
		stack = expr_stack(stack, root.ternary_expr.arg3)
		break
	default:
		err := fmt.Errorf("runtime error: %d is not a valid kind for evaluate()", root.kind)
		fmt.Println(err.Error())
		return stack
	}
	return stack

}

func evaluate(reverse_code_stack []*expr, x float32, y float32) color.RGBA {

	value_stack := []float32{}

	for i, e := range reverse_code_stack {
		switch e.kind {
		// terminal
		case EXPR_KIND_NUMBER:
			fallthrough
		case EXPR_KIND_SPECIAL_CONSTANT:
			value := e.terminal_expr.value
			value_stack = append(value_stack, value)
			break
		case EXPR_KIND_VAR_X:
			value_stack = append(value_stack, x)
			break
		case EXPR_KIND_VAR_Y:
			value_stack = append(value_stack, y)
			break

		// single:
		case EXPR_KIND_SIN:
			value_stack[len(value_stack)-1] = float32(math.Sin(float64(value_stack[len(value_stack)-1])))
			break
		case EXPR_KIND_COS:
			value_stack[len(value_stack)-1] = float32(math.Cos(float64(value_stack[len(value_stack)-1])))
			break
		case EXPR_KIND_SQRT:
			value_stack[len(value_stack)-1] = float32(math.Sqrt(float64(value_stack[len(value_stack)-1])))
			break

		// binop:
		case EXPR_KIND_VEC2:
			break
		case EXPR_KIND_ADD:
			value_stack[len(value_stack)-2] = value_stack[len(value_stack)-1] + value_stack[len(value_stack)-2]
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_SUB:
			value_stack[len(value_stack)-2] = value_stack[len(value_stack)-1] - value_stack[len(value_stack)-2]
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_MULT:
			value_stack[len(value_stack)-2] = value_stack[len(value_stack)-1] * value_stack[len(value_stack)-2]
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_DIV:
			value_stack[len(value_stack)-2] = value_stack[len(value_stack)-1] / value_stack[len(value_stack)-2]
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_MOD:
			value_stack[len(value_stack)-2] = float32(math.Mod(float64(value_stack[len(value_stack)-1]), float64(value_stack[len(value_stack)-2])))
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_POW:
			break
		case EXPR_KIND_EQ:
			break
		case EXPR_KIND_NEQ:
			break
		case EXPR_KIND_GRATER:
			if value_stack[len(value_stack)-1] > value_stack[len(value_stack)-2] {
				value_stack[len(value_stack)-2] = 1
			} else {
				value_stack[len(value_stack)-2] = 0
			}
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_LESS:
			if value_stack[len(value_stack)-1] < value_stack[len(value_stack)-2] {
				value_stack[len(value_stack)-2] = 1
			} else {
				value_stack[len(value_stack)-2] = 0
			}
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_GRATEREQ:
			if value_stack[len(value_stack)-1] >= value_stack[len(value_stack)-2] {
				value_stack[len(value_stack)-2] = 1
			} else {
				value_stack[len(value_stack)-2] = 0
			}
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_LESSEQ:
			if value_stack[len(value_stack)-1] <= value_stack[len(value_stack)-2] {
				value_stack[len(value_stack)-2] = 1
			} else {
				value_stack[len(value_stack)-2] = 0
			}
			value_stack = value_stack[1:]
			break
		case EXPR_KIND_DOT:
			break
		case EXPR_KIND_CROSS:
			break

		// ternary:
		case EXPR_KIND_VEC3:
			if i == len(reverse_code_stack)-1 {
				r := uint8(value_stack[len(value_stack)-1] * 255)
				g := uint8(value_stack[len(value_stack)-2] * 255)
				b := uint8(value_stack[len(value_stack)-3] * 255)
				return color.RGBA{r, g, b, 0xff}
			}
			break
		case EXPR_KIND_CLAMP:
			break
		case EXPR_KIND_MIX:
			break
		case EXPR_KIND_IF_THEN_ELSE:
			if value_stack[len(value_stack)-3] != 0 {
				value_stack[len(value_stack)-3] = value_stack[len(value_stack)-2]
			} else {
				value_stack[len(value_stack)-3] = value_stack[len(value_stack)-1]
			}
			value_stack = value_stack[2:]
			break
		default:
			err := fmt.Errorf("runtime error: %d is not a valid kind for evaluate()", e.kind)
			fmt.Println(err.Error())
			break
		}
	}

	return color.RGBA{0xA0, 0x20, 0xF0, 0xff} // my err color
}

func print_ast(e *expr, indent int) {
	fmt.Printf("%*s", indent*2, "")
	id := strings.ToLower(expr_kind_to_str(e.kind))
	if len(id) > 10 {
		id = id[10:]
	}
	fmt.Printf("%s", id)
	if is_terminal(*e) {
		if e.terminal_expr.value != 0 {
			fmt.Printf(":%f", e.terminal_expr.value)
		}
		fmt.Println()
		return
	}
	fmt.Printf(" (\n")
	if is_single(*e) {
		print_ast(e.single_expr.arg1, indent+1)
	} else if is_binop(*e) {
		print_ast(e.binop_expr.arg1, indent+1)
		print_ast(e.binop_expr.arg2, indent+1)
	} else if is_ternary(*e) {
		print_ast(e.ternary_expr.arg1, indent+1)
		print_ast(e.ternary_expr.arg2, indent+1)
		print_ast(e.ternary_expr.arg3, indent+1)
	}

	fmt.Printf("%*s)\n", indent*2, "")
}

func run(root *expr) []*expr {
	stack := []*expr{}
	stack = expr_stack(stack, root)
	// reverse it
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}
	return stack
}

func gen_image(width int, height int, function func() *expr) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	stack := run(function())

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			dx := float32(x) / float32(width)
			dy := float32(y) / float32(height)

			c := evaluate(stack, dx, dy)
			img.Set(x, y, c)
		}
	}

	// encode as png.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

func test_code() *expr {
	e := vec3(
		number(1),
		number(0),
		number(0),
	)
	return e
}

func test_code2() *expr {
	/*
		IMPORTANT THINKING: When the occurring probablities of an expression according to the current expression
		are implemented through markov chains, it is very weird but all the type information like bool, float, vec2, vec3
		will be intrinsically held by these chains. This shows that despite markov chains are 1 layered chains they are so powerful
		of storing that kind of semantic data
	*/
	e := vec3(
		number(0),
		if_then_else(
			grater(varx(), vary()), number(1), number(0),
		),
		number(0),
	)
	return e
}

func test_uv() *expr {
	e := vec3(
		varx(),
		vary(),
		number(0),
	)
	return e
}

func test_random() *expr {
	e := generate_expr_root(8)
	print_ast(e, 0)
	return e
}

func main() {
	width := 400
	height := 400

	gen_image(width, height, test_random)
	// gen_image(width, height, test_code2)

}
