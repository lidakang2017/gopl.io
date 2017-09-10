// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 61.
//!+

// Mandelbrot emits a PNG image of the Mandelbrot fractal.
package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"time"
	"log"
	//"sync"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, +2, +2
	width, height          = 1024, 1024
)

var img *image.RGBA

func main() {

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	/*
	 start := time.Now()
	 for py := 0; py < height; py++ {
		 y := float64(py)/height*(ymax-ymin) + ymin
		 for px := 0; px < width; px++ {
			 x := float64(px)/width*(xmax-xmin) + xmin
			 z := complex(x, y)
			 // Image point (px, py) represents complex value z.
			 img.Set(px, py, mandelbrot(z))
		 }
	 }
	 png.Encode(os.Stdout, img) // NOTE: ignoring errors

	 log.Printf("in sequential it tooked %f", time.Since(start).Seconds())
 */

	start2 := time.Now()
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	for py := range GenPy() {
		go func(py int) {
			y := float64(py)/height*(ymax-ymin) + ymin
			for px := 0; px < width; px++ {

				x := float64(px)/width*(xmax-xmin) + xmin

				go func(x, y float64) {
					z := complex(x, y)
					img.Set(px, py, Mandelbrot(z)) // Image point (px, py) represents complex value z.
					ch1 <- struct{}{}
				}(x, y)
				select {
				case <-ch1:
					continue
				}
			}
			ch2 <- struct{}{}
		}(py)
		select {
		case <-ch2:
			continue
		}
	}
	close(ch1)
	close(ch2)
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
	log.Printf("in parallel it tooked %f", time.Since(start2).Seconds())

}

func GenPy() <-chan int {
	ch := make(chan int)
	go func() {
		for py := 0; py < height; py++ {
			ch <- py
		}
		close(ch)
	}()

	return ch
}

func Mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {

		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}
