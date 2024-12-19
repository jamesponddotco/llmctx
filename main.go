package main

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/errxit-go"
	"git.sr.ht/~jamesponddotco/llmctx/internal/app"
)

func main() {
	cli := app.New(os.Args[1:])

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		errxit.Exit(err)
	}
}
