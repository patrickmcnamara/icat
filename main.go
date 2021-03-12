package main

import (
	"fmt"
	"image"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// get terminal size
	tw, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		errh(err)
	}

	// loop through image arguments
	for i, filename := range os.Args[1:] {
		// open image file
		f, err := os.Open(filename)
		if err != nil {
			errh(err)
		}

		// decode image
		img, _, err := image.Decode(f)
		if err != nil {
			errh(err)
		}

		// calculate image sizing
		iw := img.Bounds().Max.X
		ih := img.Bounds().Max.Y
		r := (iw / tw) + 1

		// print filename and info if there is more than one image
		if len(os.Args) > 2 {
			fmt.Printf("%s (%dx%d, %.2f%%)\n", filename, iw/r, ih/r, float64(r)/100)
		}

		// print image
		for y := 0; y < ih-r; y += r {
			for x := 0; x < iw-r; x += r {
				c := img.At(x, y)
				r, g, b, _ := c.RGBA()
				r /= 256
				g /= 256
				b /= 256
				fmt.Printf("\x1b[38;2;%d;%d;%dm█\x1b[0m", r, g, b)
			}
			fmt.Println()
		}

		// extra newline between images
		if i != len(os.Args[1:])-1 {
			fmt.Println()
		}
	}
}

func errh(err error) {
	fmt.Fprintf(os.Stderr, "icat: %s\n", err)
	os.Exit(1)
}
