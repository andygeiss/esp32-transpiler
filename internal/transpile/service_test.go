package transpile_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/andygeiss/esp32-transpiler/internal/transpile"
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
	if err := transpile.NewService(&in, &out).Start(t.Context()); err != nil {
		t.Fatalf("Start: %v", err)
	}
	if got, want := Trim(out.String()), Trim(expected); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_Empty_Package(t *testing.T) {
	t.Parallel()
	source := `package test`
	expected := `void loop(){}
	void setup() {}	`
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
		char* foo = "bar";
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
	const maxX = 1;
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
	const maxX = 1;
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
		if(client.connect(host, 443) == true){
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
	err := transpile.NewService(nil, &out).Start(t.Context())
	if !errors.Is(err, transpile.ErrNilReader) {
		t.Errorf("got %v, want %v", err, transpile.ErrNilReader)
	}
}

func Test_Nil_Writer_Is_An_Error(t *testing.T) {
	t.Parallel()
	var in bytes.Buffer
	err := transpile.NewService(&in, nil).Start(t.Context())
	if !errors.Is(err, transpile.ErrNilWriter) {
		t.Errorf("got %v, want %v", err, transpile.ErrNilWriter)
	}
}

func Test_Broken_Source_Writes_Nothing(t *testing.T) {
	t.Parallel()
	var in, out bytes.Buffer
	in.WriteString("package test\nfunc foo( {}")
	if err := transpile.NewService(&in, &out).Start(t.Context()); err == nil {
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
	err := transpile.NewService(&in, &out).Start(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("got %v, want %v", err, context.Canceled)
	}
	if out.Len() != 0 {
		t.Errorf("got %q written, want nothing", out.String())
	}
}
