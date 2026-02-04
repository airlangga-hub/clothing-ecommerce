package main

import (
	"bytes"
	"fmt"
	"strings"
)

func PrintStdOut(buf *bytes.Buffer) {
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	tableWidth := len(lines[0])
	separator := "+" + strings.Repeat("-", tableWidth-2) + "+"

	fmt.Println(separator)
	fmt.Println(lines[0])
	fmt.Println(separator)

	for _, line := range lines[1:] {
		fmt.Println(line)
	}
	fmt.Println(separator)
}
