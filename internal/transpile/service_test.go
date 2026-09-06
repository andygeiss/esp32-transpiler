package transpile_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/andygeiss/esp32-transpiler/internal/transpile"
	"github.com/andygeiss/esp32-transpiler/internal/transpile/handlers"
)

// Trim removes every whitespace character, so a comparison ignores layout.
func Trim(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

// Validate transpiles source and compares the sketch against expected,
// ignoring whitespace.
func Validate(source, expected string, t *testing.T) {
	t.Helper()
	var in, out bytes.Buffer
	in.WriteString(source)
	if err := transpile.NewService("", &in, &out).Start(t.Context()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if got, want := Trim(out.String()), Trim(expected); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_Empty_Package(t *testing.T) {
	t.Parallel()
	source := `package test`
	expected := `void setup(){}
	void loop() {}	`
	Validate(source, expected, t)
}

func Test_Function_Declaration(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {}
	func bar() {}
	`
	expected := `void foo(){}
	void bar() {}	`
	Validate(source, expected, t)
}

func Test_Function_Declaration_With_Args(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo(x int) {}
	func bar(y int) {}
	`
	expected := `void foo(int x){}
	void bar(int y) {}	`
	Validate(source, expected, t)
}

func Test_Const_String_Declaration(t *testing.T) {
	t.Parallel()
	source := `package test
	const foo string = "bar"
	`
	expected := `
	const char* foo = "bar";
	`
	Validate(source, expected, t)
}

func Test_Var_String_Declaration(t *testing.T) {
	t.Parallel()
	source := `package test
	var client wifi.Client
	`
	expected := `
	WiFiClient client;
	`
	Validate(source, expected, t)
}

func Test_Function_With_Const_String_Declaration(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		const foo string = "bar"
	}
	`
	expected := `
	void foo() {
		const char* foo = "bar";
	}
	`
	Validate(source, expected, t)
}
func Test_Function_With_Var_String_Declaration(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		var foo string = "bar"
	}
	`
	expected := `
	void foo() {
		const char* foo = "bar";
	}
	`
	Validate(source, expected, t)
}
func Test_Function_With_Function_Call(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		bar()
	}
	`
	expected := `
	void foo() {
		bar();
	}
	`
	Validate(source, expected, t)
}
func Test_Function_With_Function_Call_With_Args(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		bar(1,2,3)
	}
	`
	expected := `
	void foo() {
		bar(1,2,3);
	}
	`
	Validate(source, expected, t)
}
func Test_Function_With_Function_Call_With_String(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		bar("foo")
	}
	`
	expected := `
	void foo() {
		bar("foo");
	}
	`
	Validate(source, expected, t)
}

func Test_Function_With_Package_Function_Call(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		foo.Bar(1,"2")
	}
	`
	expected := `
	void foo() {
		foo.Bar(1,"2");
	}
	`
	Validate(source, expected, t)
}
func Test_Function_With_Assignments(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		x = 1
		y = 2
		z = x + y
	}
	`
	expected := `
	void foo() {
		x = 1;
		y = 2;
		z = x + y;
	}
	`
	Validate(source, expected, t)
}

func Test_Function_With_Package_Selector_Assignments(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		x = bar()
		y = pkg.Bar()
		z = x + y
	}
	`
	expected := `
	void foo() {
		x = bar();
		y = pkg.Bar();
		z = x + y;
	}
	`
	Validate(source, expected, t)
}

func Test_Function_Ident_Mapping(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		serial.Begin()
	}
	`
	expected := `
	void foo() {
		Serial.begin();
	}
	`
	Validate(source, expected, t)
}
func Test_Function_With_Ident_Param(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		foo.Bar(1,"2",digital.Low)
	}
	`
	expected := `
	void foo() {
		foo.Bar(1,"2",LOW);
	}
	`
	Validate(source, expected, t)
}

func Test_Function_With_Function_Param(t *testing.T) {
	t.Parallel()
	source := `package test
	func foo() {
		serial.Println(wifi.LocalIP())
	}
	`
	expected := `
	void foo() {
		Serial.println(WiFi.localIP());
	}
	`
	Validate(source, expected, t)
}

func Test_Package_Import(t *testing.T) {
	t.Parallel()
	source := `package test
	import "github.com/andygeiss/esp32-controller"
	import "github.com/andygeiss/esp32-controller/serial"
	import "github.com/andygeiss/esp32-controller/timer"
	import wifi "github.com/andygeiss/esp32-controller/wifi"
	`
	expected := `
	#include <WiFi.h>
	`
	Validate(source, expected, t)
}

func Test_Package_Import_But_Ignore_Controller(t *testing.T) {
	t.Parallel()
	source := `package test
	import controller "github.com/andygeiss/esp32-controller"
	import "github.com/andygeiss/esp32-controller/serial"
	import "github.com/andygeiss/esp32-controller/timer"
	import wifi "github.com/andygeiss/esp32-controller/wifi"
	`
	expected := `
	#include <WiFi.h>
	`
	Validate(source, expected, t)
}

func Test_IfStmt_With_Condition_BasicLit_And_BasicLit(t *testing.T) {
	t.Parallel()
	source := `package test
	func Setup() error {}
	func Loop() error {
		if 1 == 1 {
			serial.Println("1")
		}
	}
`
	expected := `
	void setup() {}
	void loop() {
		if (1 == 1) {
			Serial.println("1");
		}
	}
`
	Validate(source, expected, t)
}

func Test_IfStmt_With_Condition_Ident_And_BasicLit(t *testing.T) {
	t.Parallel()
	source := `package test
	func Setup() error {}
	func Loop() error {
		if x == 1 {
			serial.Println("1")
		}
	}
`
	expected := `
	void setup() {}
	void loop() {
		if (x == 1) {
			Serial.println("1");
		}
	}
`
	Validate(source, expected, t)
}

func Test_IfStmt_With_Condition_CallExpr_And_BasicLit(t *testing.T) {
	t.Parallel()
	source := `package test
	func Setup() error {}
	func Loop() error {
		if x() == 1 {
			serial.Println("1")
		}
	}
`
	expected := `
	void setup() {}
	void loop() {
		if (x() == 1) {
			Serial.println("1");
		}
	}
`
	Validate(source, expected, t)
}

func Test_IfStmt_With_Condition_Const_And_BasicLit(t *testing.T) {
	t.Parallel()
	source := `package test
	const maxX = 1
	func Setup() error {}
	func Loop() error {
		if x == maxX {
			serial.Println("1")
		}
	}
`
	expected := `
	const int maxX = 1;
	void setup() {}
	void loop() {
		if (x == maxX) {
			Serial.println("1");
		}
	}
`
	Validate(source, expected, t)
}

func Test_IfStmt_With_Else(t *testing.T) {
	t.Parallel()
	source := `package test
	const maxX = 1
	func Setup() error {}
	func Loop() error {
		if x == maxX {
			serial.Println("1")
		} else {
			serial.Println("2")
		}
	}
`
	expected := `
	const int maxX = 1;
	void setup() {}
	void loop() {
		if (x == maxX) {
			Serial.println("1");
		} else {
			Serial.println("2");
		}
	}
`
	Validate(source, expected, t)
}

func Test_SwitchStmt_With_Ident_And_BasicLit(t *testing.T) {
	t.Parallel()
	source := `package test
	func Setup() error {}
	func Loop() error {
		switch x {
		case 1:
			serial.Println("1")
		}
	}
`
	expected := `
	void setup() {}
	void loop() {
		switch (x) {
		case 1:
			Serial.println("1");
			break;
		}
	}
`
	Validate(source, expected, t)
}

func Test_SwitchStmt_With_Break(t *testing.T) {
	t.Parallel()
	source := `package test
	func Setup() error {}
	func Loop() error {
		switch x {
		case 1:
			serial.Println("1")
			break
		case 2:
			serial.Println("1")
		}
	}
`
	expected := `
	void setup() {}
	void loop() {
		switch (x) {
		case 1:
			Serial.println("1");
			break;
		case 2:
			Serial.println("1");
			break;
		}
	}
`
	Validate(source, expected, t)
}

func Test_ForLoop_WithoutInit_And_Post_Transpiles_To_While(t *testing.T) {
	t.Parallel()
	source := `package test
	import wifi "github.com/andygeiss/esp32-controller/wifi"
	func Setup() error {
		serial.Begin(serial.BaudRate115200)
		wifi.BeginEncrypted("SSID", "PASS")
		for wifi.Status() != wifi.StatusConnected {
			serial.Println("Connecting ...")
		}
		serial.Println("Connected!")
		return nil
	}
	func Loop() error {}
`
	expected := `
	#include <WiFi.h>
	void setup() {
		Serial.begin(115200);
		WiFi.begin("SSID","PASS");
		while(WiFi.status()!=WL_CONNECTED){
			Serial.println("Connecting...");
		}
		Serial.println("Connected!");
	}
	void loop() {}
`
	Validate(source, expected, t)
}

func Test_WiFiWebClient(t *testing.T) {
	t.Parallel()
	source := `package test
	import wifi "github.com/andygeiss/esp32-controller/wifi"
	var client wifi.Client
	func Setup() error {}
	func Loop() error {
		serial.Print("Connecting to ")
		serial.Println(host)
		serial.Print(" ...")
		if (client.Connect(host, 443) == true) {
			serial.Println(" Connected!")
		} else {
			serial.Println(" Failed!")
		}
	}
`
	expected := `#include <WiFi.h>
	WiFiClient client;
	voidsetup(){}
	voidloop(){
		Serial.print("Connecting to");
		Serial.println(host);
		Serial.print(" ...");
		if((client.connect(host, 443) == true)){
			Serial.println(" Connected!");
		} else {
			Serial.println(" Failed!");
		}
	}`
	Validate(source, expected, t)
}

func Test_Nil_Reader_Is_An_Error(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	err := transpile.NewService("", nil, &out).Start(t.Context())
	if !errors.Is(err, transpile.ErrNilReader) {
		t.Errorf("got %v, want %v", err, transpile.ErrNilReader)
	}
}

func Test_Nil_Writer_Is_An_Error(t *testing.T) {
	t.Parallel()
	var in bytes.Buffer
	err := transpile.NewService("", &in, nil).Start(t.Context())
	if !errors.Is(err, transpile.ErrNilWriter) {
		t.Errorf("got %v, want %v", err, transpile.ErrNilWriter)
	}
}

func Test_Broken_Source_Writes_Nothing(t *testing.T) {
	t.Parallel()
	var in, out bytes.Buffer
	in.WriteString("package test\nfunc foo( {}")
	if err := transpile.NewService("", &in, &out).Start(t.Context()); err == nil {
		t.Fatal("got no error, want a parse error")
	}
	if out.Len() != 0 {
		t.Errorf("got %q written, want nothing", out.String())
	}
}

func Test_Cancelled_Context_Writes_Nothing(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	var in, out bytes.Buffer
	in.WriteString("package test\nfunc foo() {}")
	err := transpile.NewService("", &in, &out).Start(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("got %v, want %v", err, context.Canceled)
	}
	if out.Len() != 0 {
		t.Errorf("got %q written, want nothing", out.String())
	}
}

// Refuse checks that source is turned down as outside the subset, and that the
// message says which construct and where.
func Refuse(source, want string, t *testing.T) {
	t.Helper()
	var in, out bytes.Buffer
	in.WriteString(source)
	err := transpile.NewService("controller.go", &in, &out).Start(t.Context())
	var unsupported *handlers.UnsupportedError
	if !errors.As(err, &unsupported) {
		t.Fatalf("got %v, want an *handlers.UnsupportedError", err)
	}
	if !strings.Contains(unsupported.Error(), want) {
		t.Errorf("got %q, want it to mention %q", unsupported.Error(), want)
	}
	if !strings.HasPrefix(unsupported.Error(), "controller.go:") {
		t.Errorf("got %q, want it to start with the source and the line", unsupported.Error())
	}
	if out.Len() != 0 {
		t.Errorf("got %q written, want nothing", out.String())
	}
}

// A Go case runs one clause; a C++ case runs every clause after it too, so the
// break the Go source never needed has to be written in.
func Test_SwitchStmt_Cases_Do_Not_Fall_Through(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		switch x {
		case 1:
			serial.Println("one")
		case 2:
			serial.Println("two")
		default:
			serial.Println("other")
		}
		return nil
	}
`
	expected := `
	void loop() {
		switch (x) {
		case 1: Serial.println("one"); break;
		case 2: Serial.println("two"); break;
		default: Serial.println("other"); break;
		}
	}
`
	Validate(source, expected, t)
}

// Go's fallthrough is what C++ does on its own: the break is left out.
func Test_SwitchStmt_Fallthrough_Leaves_The_Break_Out(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		switch x {
		case 1:
			serial.Println("one")
			fallthrough
		case 2:
			serial.Println("two")
		}
		return nil
	}
`
	expected := `
	void loop() {
		switch (x) {
		case 1: Serial.println("one");
		case 2: Serial.println("two"); break;
		}
	}
`
	Validate(source, expected, t)
}

// C++ has no "case 1, 2:"; it wants a label each.
func Test_SwitchStmt_Case_With_Several_Values(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		switch x {
		case 1, 2:
			serial.Println("low")
		}
		return nil
	}
`
	expected := `
	void loop() {
		switch (x) {
		case 1: case 2: Serial.println("low"); break;
		}
	}
`
	Validate(source, expected, t)
}

func Test_ForLoop_Without_A_Condition_Loops_Forever(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		for {
			serial.Println("x")
		}
	}
`
	expected := `
	void loop() {
		while (true) { Serial.println("x"); }
	}
`
	Validate(source, expected, t)
}

func Test_ForLoop_With_A_Plain_Condition_Transpiles_To_While(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		for ok {
			serial.Println("x")
		}
	}
`
	expected := `
	void loop() {
		while (ok) { Serial.println("x"); }
	}
`
	Validate(source, expected, t)
}

func Test_ForLoop_Keeps_Its_Init_And_Post(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		for i := 0; i < 10; i++ {
			serial.Println("x")
		}
	}
`
	expected := `
	void loop() {
		for (auto i = 0;i<10;i++) { Serial.println("x"); }
	}
`
	Validate(source, expected, t)
}

// break and continue leave a loop at different places, so one must not become
// the other.
func Test_Continue_Is_Not_A_Break(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		for i < 3 {
			continue
		}
		for i < 3 {
			break
		}
	}
`
	expected := `
	void loop() {
		while (i<3) { continue; }
		while (i<3) { break; }
	}
`
	Validate(source, expected, t)
}

func Test_IfStmt_With_Else_If(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		if x == 1 {
			serial.Println("1")
		} else if x == 2 {
			serial.Println("2")
		} else {
			serial.Println("3")
		}
	}
`
	expected := `
	void loop() {
		if (x == 1) { Serial.println("1"); }
		else if (x == 2) { Serial.println("2"); }
		else { Serial.println("3"); }
	}
`
	Validate(source, expected, t)
}

// Dropping the parentheses would change what the sketch works out first.
func Test_Parentheses_Are_Kept(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		z = (x + y) * 2
	}
`
	expected := `
	void loop() { z = (x+y)*2; }
`
	Validate(source, expected, t)
}

func Test_Unary_Operators(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		if !ok {
			serial.Println("-")
		}
		x = -y
		z = ^y
		w = y &^ 3
	}
`
	expected := `
	void loop() {
		if (!ok) { Serial.println("-"); }
		x = -y;
		z = ~y;
		w = y&~3;
	}
`
	Validate(source, expected, t)
}

func Test_IncDec_Statements(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		x++
		y--
	}
`
	expected := `void loop() { x++; y--; }`
	Validate(source, expected, t)
}

func Test_Short_Variable_Declaration(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		n := random.Num(10)
	}
`
	expected := `void loop() { auto n = random(10); }`
	Validate(source, expected, t)
}

// The old handler kept only literal values, so a call on the right of a var
// went missing without a word.
func Test_Var_Declaration_From_A_Call(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		var n int = random.Num(10)
	}
`
	expected := `void loop() { int n = random(10); }`
	Validate(source, expected, t)
}

// C++ has no int by default, so the type the literal implies is written out.
func Test_Declaration_Without_A_Type_Names_One(t *testing.T) {
	t.Parallel()
	source := `package test
	const a = 1
	const b = 2.5
	const c = "s"
	var d = x + 1
	`
	expected := `
	const int a = 1;
	const float b = 2.5;
	const char* c = "s";
	auto d = x+1;
	`
	Validate(source, expected, t)
}

// The const keyword belongs on every line of a group, not just the first.
func Test_Const_Group_Keeps_Const_On_Each(t *testing.T) {
	t.Parallel()
	source := `package test
	const (
		a int = 1
		b int = 2
	)
	`
	expected := `const int a = 1;const int b = 2;`
	Validate(source, expected, t)
}

func Test_Function_With_A_Return_Type(t *testing.T) {
	t.Parallel()
	source := `package test
	func add(a int, b int) int {
		return a + b
	}
	`
	expected := `int add(int a,int b) { return a+b; }`
	Validate(source, expected, t)
}

// "return nil" is how a controller says nothing went wrong, and setup and loop
// give nothing back.
func Test_Return_Nil_Is_Dropped(t *testing.T) {
	t.Parallel()
	source := `package test
	func Setup() error {
		return nil
	}
	`
	expected := `void setup() {}`
	Validate(source, expected, t)
}

func Test_Chained_Selector(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		serial.Println(wifi.LocalIP().String())
	}
	`
	expected := `void loop() { Serial.println(WiFi.localIP().String()); }`
	Validate(source, expected, t)
}

// A raw string and Go's own number spellings are not C++ spellings.
func Test_Literals_Are_Rewritten_For_CPlusPlus(t *testing.T) {
	t.Parallel()
	source := "package test\nfunc Loop() error {\n\tserial.Println(`raw \"quoted\"`)\n\tdelay(1_000)\n\tpin(0o755)\n}\n"
	expected := `void loop() { Serial.println("raw \"quoted\""); delay(1000); pin(0755); }`
	Validate(source, expected, t)
}

func Test_Blank_And_Dot_Imports_Include_Nothing(t *testing.T) {
	t.Parallel()
	source := `package test
	import _ "github.com/andygeiss/esp32-controller/serial"
	import . "github.com/andygeiss/esp32-controller/timer"
	`
	// Nothing reaches the sketch, so it is the empty one: the runtime still
	// calls setup and loop.
	Validate(source, `void setup() {} void loop() {}`, t)
}

func Test_Unsupported_Constructs_Are_Refused(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, body, want string
	}{
		{"a range loop", "for i := range xs {}", "a range loop"},
		{"a go statement", "go f()", "a go statement"},
		{"a defer statement", "defer f()", "a defer statement"},
		{"a goto", "goto done", "goto"},
		{"a labelled break", "for x { break outer }", "a labelled break"},
		{"a switch without a value", "switch { case x > 1: }", "a switch without a value"},
		{"an if with an init statement", "if n := f(); n > 1 {}", "an if with an init statement"},
		{"a multiple assignment", "x, y = 1, 2", "an assignment to more than one variable"},
		{"a channel receive", "x = <-ch", "the unary operator <-"},
		{"an address", "x = &y", "the unary operator &"},
		{"a composite literal", "x = T{}", "T{}"},
		{"an index", "x = xs[0]", "xs[0]"},
		{"iota", "const a = iota", "iota"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Refuse("package test\nfunc Loop() error {\n"+tt.body+"\n}\n", tt.want, t)
		})
	}
}

func Test_Unsupported_Declarations_Are_Refused(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, decl, want string
	}{
		{"a method", "func (c *Ctl) Loop() error {}", "a method"},
		{"a type declaration", "type Ctl struct{}", "a type declaration"},
		{"a pointer type", "func f(c *Ctl) {}", "the type *Ctl"},
		{"a slice type", "func f(xs []int) {}", "the type []int"},
		{"two results", "func f() (int, error) {}", "a function returning more than one value"},
		{"an error result", "func f() error {}", "an error result on a function other than Setup or Loop"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Refuse("package test\n"+tt.decl+"\n", tt.want, t)
		})
	}
}

// Test_The_README_Example pins the sketch byte for byte, spaces and all. Every
// other test compares with the whitespace taken out, which would let a missing
// space between a type and a name through.
func Test_The_README_Example(t *testing.T) {
	t.Parallel()
	source := `package controller

import "github.com/andygeiss/esp32-controller/serial"

// Setup runs once when the board powers up.
func Setup() error {
	serial.Begin(serial.BaudRate115200)
	return nil
}

// Loop runs over and over, as fast as it returns.
func Loop() error {
	serial.Println("hello")
	return nil
}
`
	const expected = "void setup() {Serial.begin(115200);}\nvoid loop() {Serial.println(\"hello\");}\n"

	var in, out bytes.Buffer
	in.WriteString(source)
	if err := transpile.NewService("controller.go", &in, &out).Start(t.Context()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if out.String() != expected {
		t.Errorf("got %q, want %q", out.String(), expected)
	}
}

func Test_Character_Literal(t *testing.T) {
	t.Parallel()
	Validate("package test\nfunc Loop() error { serial.Print('x') }\n",
		`void loop() { Serial.print('x'); }`, t)
}

func Test_Unsupported_Expressions_Are_Refused(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, body, want string
	}{
		{"a spread call", "f(xs...)", "a call spreading its last argument with ..."},
		{"a joined string", `serial.Println("a" + name)`, "joining strings with +"},
		{"a function literal", "f := func() {}", "a function literal"},
		{"a rune that is not a byte", "serial.Print('é')", "the character 'é'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Refuse("package test\nfunc Loop() error {\n"+tt.body+"\n}\n", tt.want, t)
		})
	}
}

// The board calls setup and loop itself, so a signature it cannot call is
// turned down here rather than by the linker.
func Test_Setup_And_Loop_Keep_The_Shape_The_Board_Calls(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, decl, want string
	}{
		{"Setup with a parameter", "func Setup(x int) error {}", "Setup with a parameter"},
		{"Loop with a parameter", "func Loop(x int) error {}", "Loop with a parameter"},
		{"Setup returning an int", "func Setup() int {}", "Setup returning anything but error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			Refuse("package test\n"+tt.decl+"\n", tt.want, t)
		})
	}
}

func Test_Setup_And_Loop_Without_A_Result(t *testing.T) {
	t.Parallel()
	Validate("package test\nfunc Setup() {}\nfunc Loop() {}\n", "void setup() {} void loop() {}", t)
}

// Fuzz_Start_Never_Panics holds the line the panic in the handlers is allowed
// to stop at. They unwind through a recover to report a construct they cannot
// translate, so any source the parser accepts has to come back as a sketch or
// as an error — never as a crash.
func Fuzz_Start_Never_Panics(f *testing.F) {
	seeds := []string{
		"package test",
		"package test\nfunc Loop() error { for {} }",
		"package test\nfunc Loop() error { for i := 0; i < 3; i++ {} }",
		"package test\nfunc Loop() error { if x { } else if y { } else { } }",
		"package test\nfunc Loop() error { switch x { case 1, 2: fallthrough\ndefault: } }",
		"package test\nfunc Loop() error { switch { } }",
		"package test\nfunc Loop() error { go f(); defer g() }",
		"package test\nfunc Loop() error { for i := range xs {} }",
		"package test\nfunc Loop() error { x, y = 1, 2 }",
		"package test\nfunc Loop() error { f(xs...) }",
		"package test\nfunc Loop() error { x = <-ch; y = &z; w = *p }",
		"package test\nfunc Loop() error { x = xs[0]; y = T{}; z = func() {} }",
		"package test\nconst ( a = iota\nb )",
		"package test\nvar a, b = 1, 2",
		"package test\ntype T struct{}",
		"package test\nfunc (c *T) Loop() error {}",
		"package test\nfunc f() (int, error) {}",
		"package test\nimport wifi \"w\"\nimport _ \"s\"",
		"package test\nfunc Loop() error { serial.Println(`raw`); delay(0o7); pin(1_0) }",
		"package test\nlabel: for {}",
		"package test\nfunc Loop() error { select {} }",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, source string) {
		var in, out bytes.Buffer
		in.WriteString(source)
		// Whatever it decides, it must decide it without crashing.
		_ = transpile.NewService("fuzz.go", &in, &out).Start(t.Context())
	})
}

// "return nil" inside a case leaves the whole function, not just the switch,
// so it has to stay a return: a break would run whatever follows the switch.
func Test_SwitchStmt_Case_Ending_In_Return_Nil_Leaves_The_Function(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		switch x {
		case 1:
			serial.Println("one")
			return nil
		case 2:
			serial.Println("two")
		}
		return nil
	}
`
	expected := `
	void loop() {
		switch (x) {
		case 1: Serial.println("one"); return;
		case 2: Serial.println("two"); break;
		}
	}
`
	Validate(source, expected, t)
}

// An early "return nil" stops the function; the code after the switch must not
// run, the way it would not in Go.
func Test_Early_Return_Nil_Skips_What_Follows(t *testing.T) {
	t.Parallel()
	source := `package test
	func Loop() error {
		if done {
			return nil
		}
		serial.Println("working")
		return nil
	}
`
	expected := `
	void loop() {
		if (done) { return; }
		Serial.println("working");
	}
`
	Validate(source, expected, t)
}

// A case that really does return leaves the switch on its own.
func Test_SwitchStmt_Case_Ending_In_A_Value_Return_Needs_No_Break(t *testing.T) {
	t.Parallel()
	source := `package test
	func pick(x int) int {
		switch x {
		case 1:
			return 10
		}
		return 0
	}
`
	expected := `
	int pick(int x) {
		switch (x) {
		case 1: return 10;
		}
		return 0;
	}
`
	Validate(source, expected, t)
}

func Test_AndNot_Assignment(t *testing.T) {
	t.Parallel()
	Validate("package test\nfunc Loop() error { x &^= mask }\n",
		`void loop() { x &= ~(mask); }`, t)
}

// A void setup cannot hand an error back, and C++ would not compile the try.
func Test_Setup_Returning_A_Value_Is_Refused(t *testing.T) {
	t.Parallel()
	Refuse("package test\nfunc Setup() error {\n\tif err != nil { return err }\n\treturn nil\n}\n",
		"Setup returning a value", t)
}

// A bare return leaves setup early, which void allows.
func Test_Setup_With_A_Bare_Return(t *testing.T) {
	t.Parallel()
	Validate("package test\nfunc Setup() error {\n\tif x { return }\n\treturn nil\n}\n",
		`void setup() { if (x) { return; } }`, t)
}

func Test_Discarded_Assignment_Keeps_The_Call(t *testing.T) {
	t.Parallel()
	Validate("package test\nfunc Loop() error { _ = f() }\n", `void loop() { f(); }`, t)
}

func Test_A_Function_Without_A_Body_Is_Refused(t *testing.T) {
	t.Parallel()
	Refuse("package test\nfunc f(x int) int\n", "a function without a body", t)
}
