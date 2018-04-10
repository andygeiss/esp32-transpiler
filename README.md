# ESP32 Transpiler

[![Go Report Card](https://goreportcard.com/badge/github.com/andygeiss/esp32-transpiler)](https://goreportcard.com/report/github.com/andygeiss/esp32-transpiler)
[![Build Status](https://travis-ci.org/andygeiss/esp32-transpiler.svg?branch=master)](https://travis-ci.org/andygeiss/esp32-transpiler)

## Purpose

The [Arduino IDE](https://www.arduino.cc/en/Main/Software) is easy to use.
But I faced problems like maintainability and testability at more complicated IoT projects.
I needed to compile and flash the ESP32 before testing my code functionality by doing it 100% manually.

This solution transpiles Golang into Arduino code, which can be compiled to an image by using the ESP32 toolchain.
Now I am able to use a fully automated testing approach instead of doing it 100% manually.

**Important**: The Transpiler only supports a small subset of the [Golang Language Specification](https://golang.org/ref/spec).

## Installation

    go get -u github.com/andygeiss/esp32-transpiler

## Usage

    Usage of esp32-transpiler:
      -source string
            Golang source file
      -target string
            Arduino sketch file
