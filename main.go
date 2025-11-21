package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andygeiss/esp32-transpiler/transpile"
)

func main() {
	source, target := getFlags()
	checkFlagsAreValid(source, target)
	safeTranspile(source, target)
}

func checkFlagsAreValid(source, target string) {
	if source == "" || target == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func getFlags() (string, string) {
	source := flag.String("source", "", "Golang source file")
	target := flag.String("target", "", "Arduino sketch file")
	flag.Parse()
	return *source, *target
}

func printUsage() {
	fmt.Print("This program transpiles Golang source into corresponding Arduino sketches.\n\n")
	fmt.Print("Options:\n")
	flag.PrintDefaults()
	fmt.Print("\n")
	fmt.Print("Example:\n")
	fmt.Printf("\tesp32 -source impl/blink/controller.go -target impl/blink/controller.transpile\n\n")
}

func safeTranspile(source, target string) {
	// Read the Golang source file.
	in, err := os.Open(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Go source file [%s] could not be opened! %v", source, err)
		os.Exit(1)
	}
	defer in.Close()

	// Create the Arduino sketch file.
	os.Remove(target)
	out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Arduino sketch file [%s] could not be opened! %v", target, err)
		os.Exit(1)
	}

	// Transpiles the Golang source into Arduino sketch.
	service := transpile.NewService(in, out)
	if err := service.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
