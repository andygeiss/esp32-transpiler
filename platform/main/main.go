package main

import (
	"flag"
	"fmt"
	"github.com/andygeiss/esp32-transpiler/application/transpile"
	"github.com/andygeiss/esp32-transpiler/infrastructure/ino"
	"github.com/andygeiss/log"
	"os"
)

func main() {
	mapping, source, target := getFlags()
	checkFlagsAreValid(mapping, source, target)
	safeTranspile(mapping, source, target)
}

func checkFlagsAreValid(mapping, source, target string) {
	if mapping == "" || source == "" || target == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func getFlags() (string, string, string) {
	mapping := flag.String("mapping", fmt.Sprintf("%s/esp32/mapping.json", os.Getenv("HOME")), "API Mapping file")
	source := flag.String("source", "", "Golang source file")
	target := flag.String("target", "", "Arduino sketch file")
	flag.Parse()
	return *mapping, *source, *target
}

func printUsage() {
	fmt.Print("This program transpiles Golang source into corresponding Arduino sketches.\n\n")
	fmt.Print("Options:\n")
	flag.PrintDefaults()
	fmt.Print("\n")
	fmt.Print("Example:\n")
	fmt.Printf("\tesp32 -source application/blink/controller.go -target application/blink/controller.ino\n\n")
}

func safeTranspile(mapping, source, target string) {
	// Read the Golang source file.
	in, err := os.Open(source)
	if err != nil {
		log.Fatal("Go source file [%s] could not be opened! %v", source, err)
	}
	defer in.Close()
	// Create the Arduino sketch file.
	os.Remove(target)
	out, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		log.Fatal("Arduino sketch file [%s] could not be opened! %v", target, err)
	}
	// Transpiles the Golang source into Arduino sketch.
	m := ino.NewMapping(mapping)
	if err := m.Read(); err != nil {
		log.Fatal("%v", err)
	}
	worker := ino.NewWorker(in, out, m)
	trans := transpile.NewTranspiler(worker)
	if err := trans.Transpile(); err != nil {
		log.Fatal("%v", err)
	}
}
