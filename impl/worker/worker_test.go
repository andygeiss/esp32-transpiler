package worker_test

import (
	"bytes"
	. "github.com/andygeiss/assert"
	"github.com/andygeiss/esp32-transpiler/impl/worker"
	"strings"
	"testing"
)

// Trim removes all the whitespaces and returns a new string.
func Trim(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}

// Validate the content of a given source with an expected outcome by using a string compare.
// The Worker will be started and used to transform the source into an Arduino sketch format.
func Validate(source, expected string, t *testing.T) {
	var in, out bytes.Buffer
	in.WriteString(source)
	wrk := worker.NewWorker(&in, &out, worker.NewMapping())
	Assert(t, wrk.Start(), IsNil())
	code := out.String()
	tcode, texpected := Trim(code), Trim(expected)
	Assert(t, tcode, IsEqual(texpected))
}

func Test_Empty_Package(t *testing.T) {
	source := `package test`
	expected := `void loop(){}
	void setup() {}	`
	Validate(source, expected, t)
}

func Test_Function_Declaration(t *testing.T) {
	source := `package test
	func foo() {}
	func bar() {}
	`
	expected := `void foo(){}
	void bar() {}	`
	Validate(source, expected, t)
}
func Test_Function_Declaration_With_Args(t *testing.T) {
	source := `package test
	func foo(x int) {}
	func bar(y int) {}
	`
	expected := `void foo(int x){}
	void bar(int y) {}	`
	Validate(source, expected, t)
}
func Test_Const_String_Declaration(t *testing.T) {
	source := `package test
	const foo string = "bar"
	`
	expected := `
	const char* foo = "bar";
	`
	Validate(source, expected, t)
}
func Test_Function_With_Const_String_Declaration(t *testing.T) {
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
	source := `package test
	func foo() {
		var foo string = "bar"
	}
	`
	expected := `
	void foo() {
		char* foo = "bar";
	}
	`
	Validate(source, expected, t)
}
func Test_Function_With_Function_Call(t *testing.T) {
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
	source := `package test
	import "github.com/andygeiss/esp32-mqtt/api/controller"
	import "github.com/andygeiss/esp32-mqtt/api/controller/serial"
	import "github.com/andygeiss/esp32/api/controller/timer"
	import wifi "github.com/andygeiss/esp32/api/controller/wifi"
	`
	expected := `
	#include <WiFi.h>
	`
	Validate(source, expected, t)
}

func Test_Package_Import_But_Ignore_Controller(t *testing.T) {
	source := `package test
	import controller "github.com/andygeiss/esp32-controller"
	import "github.com/andygeiss/esp32-mqtt/api/controller/serial"
	import "github.com/andygeiss/esp32/api/controller/timer"
	import wifi "github.com/andygeiss/esp32/api/controller/wifi"
	`
	expected := `
	#include <WiFi.h>
	`
	Validate(source, expected, t)
}

func Test_BlockStmt_With_Condition_BasicLit_And_BasicLit(t *testing.T) {
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

func Test_BlockStmt_With_Condition_Ident_And_BasicLit(t *testing.T) {
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

func Test_BlockStmt_With_Condition_CallExpr_And_BasicLit(t *testing.T) {
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
