package eval_test

import (
	"testing"

	"github.com/nokia/ntt/internal/loc"
	"github.com/nokia/ntt/internal/ttcn3/ast/eval"
	"github.com/nokia/ntt/internal/ttcn3/parser"
	"github.com/nokia/ntt/runtime"
)

func TestInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0", 0},
		{"-0", 0},
		{"+0", 0},
		{"10", 10},
		{"-10", -10},
		{"+10", 10},
		{"1+2*3", 7},
		{"(1+2)*3", 9},
	}
	for _, tt := range tests {
		val := testEval(t, tt.input)
		if val == nil {
			t.Errorf("Evaluation of %q returned nil", tt.input)
			continue
		}
		testInt(t, val, tt.expected)
	}

}

func TestBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"not true", false},
		{"not not true", true},
		{"not not not true", false},
		{"not false", true},
		{"not not false", false},
		{"not not not false", true},
		{"1<1", false},
		{"1<=1", true},
		{"1<2", true},
		{"1==1", true},
		{"1==2", false},
		{"1!=1", false},
		{"1!=2", true},
		{"2-1 < 2", true},
		{"2+1==1+2", true},
		{"true==false", false},
		{"true!=false", true},
	}
	for _, tt := range tests {
		val := testEval(t, tt.input)
		if val == nil {
			t.Errorf("Evaluation of %q returned nil", tt.input)
			continue
		}
		testBool(t, val, tt.expected)
	}

}
func TestIfStmt(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		val := testEval(t, tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testInt(t, val, int64(expected))
		default:
			if val != nil {
				t.Errorf("object is not nil. got=%T (%+v)", val, val)
			}
		}
	}
}

func testEval(t *testing.T, input string) runtime.Object {
	fset := loc.NewFileSet()
	nodes, err := parser.Parse(fset, "<stdin>", input)
	if err != nil {
		t.Fatalf("testEval: %s", err.Error())
	}
	if len(nodes) != 1 {
		t.Fatalf("testEval: Multiple nodes have been returned, but only one is supported.")
	}
	return eval.Eval(nodes[0], runtime.NewEnv())
}

func testInt(t *testing.T, obj runtime.Object, expected int64) bool {
	i, ok := obj.(runtime.Int)
	if !ok {
		t.Errorf("object is not runtime.Int. got=%T (%+v)", obj, obj)
		return false
	}

	if !i.IsInt64() {
		t.Errorf("object is to big to compare. got=%s", i)
		return false
	}

	if val := i.Int64(); val != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", val, expected)
		return false
	}

	return true
}

func testBool(t *testing.T, obj runtime.Object, expected bool) bool {
	b, ok := obj.(runtime.Bool)
	if !ok {
		t.Errorf("object is not runtime.Bool. got=%T (%+v)", obj, obj)
		return false
	}

	if val := b.Bool(); val != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", val, expected)
		return false
	}

	return true
}
