// SPDX-License-Identifier: GPL-3.0

package main

import (
	"context"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)

	go func() {
		<-interruptChan
		cancel()
	}()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
