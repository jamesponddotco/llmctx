package main

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/errxit-go"
	"git.sr.ht/~jamesponddotco/llmctx/internal/app"
)

func main() {
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		errxit.Exit(err)
	}
}
