package main

import (
	"os"
	"time"

	"github.com/marcetin/nauci-go-sa-testovima/math/v9/clockface"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
