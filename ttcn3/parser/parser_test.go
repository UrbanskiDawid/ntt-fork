package parser

import (
	"github.com/nokia/ntt/ttcn3/token"
	"testing"
)

type Expect int

const (
	PASS Expect = iota
	FAIL
)

type Test struct {
	expect Expect
	input  string
}

func TestModules(t *testing.T) {
	modules := []Test{
		{PASS, `module m {}`},
		{PASS, `module m language "str1", "str2" {}`},
	}

	testParse(t, modules, func(p *parser) { p.parseModule() })
}

func TestWithStmts(t *testing.T) {
	withStmts := []Test{
		{PASS, `encode    "str";`},
		{PASS, `variant   "str";`},
		{PASS, `display   "str";`},
		{PASS, `extension "str";`},
		{PASS, `optional  "str";`},
		{PASS, `stepsize  "str";`},
		{PASS, `encode override        "str";`},
		{PASS, `encode @local          "str";`},
		{PASS, `encode @local          "str"."ruleA";`},
		{PASS, `encode ([-])           "str";`},
		{PASS, `encode (a[-])          "str";`},
		{PASS, `encode (group all)     "str";`},
		{PASS, `encode (type all)      "str";`},
		{PASS, `encode (template all)  "str";`},
		{PASS, `encode (const all)     "str";`},
		{PASS, `encode (altstep all)   "str";`},
		{PASS, `encode (testcase all)  "str";`},
		{PASS, `encode (function all)  "str";`},
		{PASS, `encode (signature all) "str";`},
		{PASS, `encode (modulepar all) "str";`},
		{PASS, `encode (type all except {a,b}) "str";`},
	}

	testParse(t, withStmts, func(p *parser) { p.parseWithStmt() })
}

func TestExprs(t *testing.T) {
	exprs := []Test{
		{PASS, `-`},
		{PASS, `a[-]`},
		{PASS, `-1 * x`},
		{PASS, `-x * y`},
		{PASS, `{x := (1+2)*3, y:=a.f()}`},
		{PASS, `{(1+2)*3, a.f()}`},
		{PASS, `{-,-}`},
		{PASS, `(1,*,?,-,2)`},
		{PASS, `t length(5..23)`},
		{PASS, `t length(5..23) ifpresent`},
		{PASS, `t ifpresent`},
		{PASS, `system:p`},
		{FAIL, `modifies t:=23`},
		{FAIL, `complement(all from t)`},
		{PASS, `b := any from c.running -> @index value i`},
		{PASS, `p := decmatch M: {f1:= 10, f2 := '1001'B}`},
		{PASS, `p := decmatch ("UTF-8") M: {f1:= 10, f2 := '1001'B}`},
		{FAIL, `p := @decoded payload`},
		{PASS, `regexp @nocase(x,charstring:"?+(text)?+",0)`},
	}

	testParse(t, exprs, func(p *parser) { p.parseExprList() })
}

func TestFuncDecls(t *testing.T) {
	funcDecls := []Test{
		{PASS, `testcase f() {}`},
		{PASS, `testcase f() runs on A[-] {}`},
		{PASS, `testcase f() runs on C system C {}`},
		{PASS, `function f() {}`},
		{PASS, `function f() return int {}`},
		{PASS, `function f() return template int {}`},
		{PASS, `function f() return template(value) int {}`},
		{PASS, `function f() return value int {}`},
		{PASS, `function f @deterministic () {}`},
		{PASS, `function f() runs on A[-] {}`},
		{PASS, `function f() mtc C {}`},
		{PASS, `function f() runs on C mtc C system C {}`},
		{PASS, `altstep as() { var roi[-] a[4][4]; [] receive; [else] {}}`},
		{FAIL, `external function f();`},
		{PASS, `signature f();`},
		{PASS, `signature f() exception (integer);`},
		{PASS, `signature f() return int;`},
		{PASS, `signature f() return int exception (integer, a.b[0]);`},
		{PASS, `signature f() noblock;`},
		{PASS, `signature f() noblock exception (integer, a.b[0]);`},
	}

	testParse(t, funcDecls, func(p *parser) { p.parseFuncDecl() })
}

func TestModuleDefs(t *testing.T) {
	moduleDefs := []Test{
		{PASS, `import from m all;`},
		{PASS, `import from m language "str1", "str2" all;`},
		{PASS, `import from m all except {}`},
		{PASS, `import from m all except {
                        template  all;
                        const     all;
                        altstep   all;
                        testcase  all;
                        function  all;
                        signature all;
                        modulepar all;
                        import    all;
                        type      all }`},
		{PASS, `import from m all except { group all }`},
		{PASS, `import from m all except { group x,y }`},
		{PASS, `import from m {
                        template  all;
                        const     all;
                        altstep   all;
                        testcase  all;
                        function  all;
                        signature all;
                        modulepar all;
                        import    all;
                        type      all }`},
		{PASS, `import from m {
                        group x except { group all }, y }`},

		{FAIL, `friend module m;`},
		{FAIL, `public modulepar integer x;`},
		{FAIL, `private function fn() {}`},
	}
	testParse(t, moduleDefs, func(p *parser) { p.parseModuleDef() })
}

func TestValueDecls(t *testing.T) {
	valueDecls := []Test{
		{PASS, `const integer x;`},
		{PASS, `const int x := 1;`},
		{PASS, `const int x := 1, yi := 2;`},
		{PASS, `const int x[len] := 1, y := 2;`},
		{PASS, `const a[-] x := 1;`},
		{PASS, `const a[1] x[2][3] := x[4];`},
		{PASS, `var int x, y := 2, z;`},
		{PASS, `var template          int x;`},
		{PASS, `var template(omit)    int x;`},
		{PASS, `var template(value)   int x;`},
		{PASS, `var template(present) int x;`},
		{PASS, `var omit    int x;`},
		{PASS, `var value   int x;`},
		{PASS, `var present int x;`},
		{PASS, `var value @lazy int x;`},
		{PASS, `var value @lazy int x, y := ?;`},
		{PASS, `template int x := ?;`},
		{PASS, `template int x modifies y := ?;`},
		{PASS, `template int x(int i) := i;`},
		{PASS, `template @lazy int x := ?;`},
		{PASS, `template @lazy int x(int i) := i;`},
		{PASS, `template @lazy int  x(int i) modifies y := ?;`},
		{PASS, `template @lazy a[-] x(int i) modifies y := ?;`},
		{PASS, `template(omit)    int x := ?;`},
		{PASS, `template(value)   int x := ?;`},
		{PASS, `template(present) int x := ?;`},
		{PASS, `timer x, y := 1.0, y;`},
		{PASS, `port P x[len], y := 1, z := 2 ;`},
		{PASS, `modulepar RoI[-] x, y:=23, z;`},
	}

	testParse(t, valueDecls, func(p *parser) { p.parseDecl() })
}

func TestFormalPars(t *testing.T) {
	formalPars := []Test{
		{PASS, `()`},
		{PASS, `(int y)`},
		{PASS, `(int x, int y)`},
		{PASS, `(in int x, out int y, inout int z)`},
		{PASS, `(in template(value) @fuzzy timer x := 1, out timer y)`},
		{PASS, `(out timer y, in template(value) @fuzzy timer x := 1)`},
		{PASS, `(out timer y := -, in value @fuzzy timer x := 1)`},
		{PASS, `(out timer y := -, in value timer x := (1,2,3))`},
	}
	testParse(t, formalPars, func(p *parser) { p.parseParameters() })
}

func TestTypes(t *testing.T) {
	types := []Test{
		// Subtypes
		{PASS, `type integer t`},
		{PASS, `type int t (0..255)`},
		{PASS, `type int t length(2)`},
		{PASS, `type a[0] t (0,1) length(2)`},

		// List Types
		{PASS, `type set of int s`},
		{PASS, `type set length(2) of int s`},
		{PASS, `type set length(2) of int s length(2)`},
		//{PASS, `type set length(2) of int s (0,1,2) length(2)`},
		{PASS, `type set of set of int s`},
		{PASS, `type set length(1) of set length(2) of int() s length(3)`},

		// Struct Types
		{PASS, `type set s {}`},
		{PASS, `type set s {int a optional }`},
		{PASS, `type set s {set length(1) of set length(2) of int() f1[-][-] length(3) optional`},
		{PASS, `type union s {@default set of int f1 optional}`},
		{PASS, `type enumerated a[1][2] {e, e[3], e[-](1)`},

		// Port Types
		{PASS, `type port p message {address a.b[-]}`},
		{PASS, `type port p message {inout float, a.b[-]}`},
		{PASS, `type port p message {map param (out int i:=1)}`},
		{PASS, `type port p message {unmap param (out int i:=1)}`},
		{PASS, `type port p procedure {}`},
		{PASS, `type port p mixed {}`},

		// Component Types
		{PASS, `type component C {}`},
		{PASS, `type component C extends C[-], m.Base {}`},

		// Behaviour Types
		{PASS, `type function fn() runs on self return template int`},
		{PASS, `type altstep  as() runs on self return int`},
		{PASS, `type testcase tc() runs on C system TSI`},
	}
	testParse(t, types, func(p *parser) { p.parseType() })
}

func TestStmts(t *testing.T) {
	stmts := []Test{
		// Structural Statements
		{PASS, `repeat;`},
		{PASS, `break;`},
		{PASS, `continue;`},
		{PASS, `return;`},
		{PASS, `return x() * 1;`},
		{PASS, `label L1;`},
		{PASS, `goto L2;`},
		{PASS, `for (var int i := 0; i<23; i := i+1) {}`},
		{PASS, `for (i:=x; i<23; i:=i+1) {}`},
		{PASS, `while (23) {}`},
		{PASS, `do {} while (23);`},
		{PASS, `if (1) {}`},
		{PASS, `if (1) {} else {}`},
		{PASS, `if (1) {} else if (2) {} else {}`},
		{PASS, `select union (p.x()) { case(1) {} case else {}}`},
		{PASS, `select  (23) {case(1) {} case else {}}`},
		{PASS, `interleave {}`},
		{PASS, `alt {}`},
		{PASS, `alt { [] receive; [23<foo()] p.timeout { var i x:=23; } [else] {}}`},

		// Value Declaration Statements
		{PASS, `var comp C := C.create;`},
		{PASS, `var comp C := C.create("han solo") alive;`},

		// Expr Statements
		{PASS, `send() to 80;`},
		{PASS, `send() to v_dst;`},
		{PASS, `receive from ip.address:?;`},
		{PASS, `receive from ip.address:? -> @index x;`},
		{PASS, `testcase.stop;`},
		{PASS, `stop;`},
		{FAIL, `map (system:p1, c:p);`},
		{FAIL, `map (p1, p2) param ("localhost", 80);`},
		{FAIL, `unmap;`},
		{FAIL, `unmap (true);`},
		{FAIL, `unmap (true) param (-,-);`},
		{PASS, `p.call(foo) to 80;`},
		{PASS, `p[i].call(S:{});`},
		{PASS, `any port.check;`},
		{PASS, `p.check(receive);`},
		{FAIL, `p.check(from x -> timestamp bar);`},
		{FAIL, `p.check(-> @index value i);`},
		{FAIL, `p.check(receive from x -> value ("foo"));`},
		{FAIL, `p.check(getreply(23 value x) from x -> sender(foo));`},
		{PASS, `p.getreply(23);`},
		{FAIL, `p.reply(23 value x);`},
	}
	testParse(t, stmts, func(p *parser) { p.parseStmt() })
}

func testParse(t *testing.T, tests []Test, f func(p *parser)) {
	for _, tt := range tests {
		err := anyParse(tt.input, f)
		if tt.expect == PASS && err != nil {
			t.Errorf("Parse(%#q):\n\t%v\n\n", tt.input, err)
		}
		if tt.expect == FAIL && err == nil {
			t.Errorf("breakage vanished: Parse(%#q)", tt.input)
		}
	}
}

func anyParse(input string, f func(p *parser)) error {
	var p parser
	p.init(token.NewFileSet(), "", []byte(input), 0)
	f(&p)
	p.errors.Sort()
	return p.errors.Err()
}
