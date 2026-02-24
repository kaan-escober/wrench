package main

import (
	"fmt"
	"os"

	"github.com/kaan-escober/wrench/internal/ui"
)

func main() {
	if err := ui.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
