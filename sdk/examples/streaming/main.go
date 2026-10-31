package main

import (
	"fmt"

	"github.com/icosmos-space/llmatrix/sdk/streaming"
)

func main() {
	// Streaming wiring lands with Engine integration.
	fmt.Println("streaming example scaffold:", streaming.EventTypeDelta)
}
