package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/user/car-simulator/internal"
	"github.com/user/car-simulator/internal/dst"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_ "go.uber.org/automaxprocs"
)

func main() {
	app := internal.NewApplication(true)

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)

		if errors.Is(err, &dst.DSTAssertionError{}) {
			fmt.Printf("DST Assertion failed:")
			fmt.Printf("%v", err)

			// Dump the state
			state := app.Storage.Dump()
			for key, value := range state {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
		os.Exit(1)
	}
}
