// Writes an SVG clockface of the current time to Stdout.
package main

import (
	"os"
	"time"

	"github.com/marcetin/nauci-go-sa-testovima/math/vFinal/clockface/svg"
)

func main() {
	t := time.Now()
	svg.Write(os.Stdout, t)
}
