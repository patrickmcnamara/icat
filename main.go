package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
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
		// read file and detect it's type
		f, err := os.Open(filename)
		if err != nil {
			errh(err)
		}
		buf := make([]byte, 512)
		f.Read(buf)
		ct := http.DetectContentType(buf)
		f.Seek(0, io.SeekStart)

		// decode the image
		var img image.Image
		switch ct {
		case "image/png":
			img, _ = png.Decode(f)
		case "image/jpeg":
			img, _ = jpeg.Decode(f)
		case "image/gif":
			img, _ = gif.Decode(f)
		default:
			return
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
