package main

import (
	"bufio"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"cirello.io/errors"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	log.SetPrefix("uqmani: ")
	log.SetFlags(0)
	fd, err := os.Open("kohrah/kohrah.ani")
	checkErr(err)

	var baseImg image.Image
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		l := parse(scanner.Text())
		if baseImg == nil {
			fd, err := os.Open(filepath.Join("kohrah", l.fn))
			checkErr(errors.E(err))
			img, _, err := image.Decode(fd)
			checkErr(err)
			baseImg = img
			continue
		}
		spew.Dump(l)
		rgba := image.NewRGBA(baseImg.Bounds())
		fd, err := os.Open(filepath.Join("kohrah", l.fn))
		checkErr(errors.E(err))
		img2, _, err := image.Decode(fd)
		checkErr(err)

		img2Bounds := img2.Bounds().Size()
		leftTop := image.Point{-l.x, -l.y}
		rightBottom := image.Point{img2Bounds.X + (-l.x), img2Bounds.Y + (-l.y)}
		r2 := image.Rectangle{leftTop, rightBottom}
		draw.Draw(rgba, baseImg.Bounds(), baseImg, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Over)

		out, err := os.Create(l.fn)
		checkErr(err)
		png.Encode(out, rgba)
	}
	checkErr(scanner.Err())
}

type aniLine struct {
	fn string
	x  int
	y  int
}

func parse(s string) aniLine {
	parts := strings.Split(s, " ")
	return aniLine{
		fn: parts[0],
		x:  mustAtoi(parts[3]),
		y:  mustAtoi(parts[4]),
	}
}

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	checkErr(err)
	return v
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
