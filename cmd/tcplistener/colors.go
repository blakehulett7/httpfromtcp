package main

import "fmt"

const ResetColor = "\033[0m"

const Blue = "\033[34m"
const Cyan = "\033[36m"
const Gray = "\033[90m"
const Green = "\033[32m"
const Magenta = "\033[35m"
const Red = "\033[31m"
const White = "\033[37m"
const Yellow = "\033[33m"

const BrightBlue = "\033[94m"
const BrightCyan = "\033[96m"
const BrightGreen = "\033[92m"
const BrightMagenta = "\033[95m"
const BrightRed = "\033[91m"
const BrightWhite = "\033[97m"
const BrightYellow = "\033[93m"

func ColorPrint(color, message string) {
	fmt.Println(color + message + ResetColor)
}
