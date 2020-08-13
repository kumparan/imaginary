package main

import (
	"bytes"
	"image"
	"math"
	"time"

	"github.com/kumparan/bimg"

	log "github.com/sirupsen/logrus"

	"github.com/disintegration/imaging"
	"willnorris.com/go/gifresize"
)

var resampleFilter = imaging.Lanczos

func ProcessGIF(imageBuf []byte, opts bimg.Options) (out Image, err error) {
	start := time.Now()
	resultBuffer := new(bytes.Buffer)

	fn := func(img image.Image) image.Image {
		return transformGIFFrame(img, opts)
	}

	err = gifresize.Process(resultBuffer, bytes.NewReader(imageBuf), fn)
	if err != nil {
		return Image{}, err
	}

	// keep this loggerto monitor performance?
	log.Println("gif process duration: ", time.Since(start))

	resultByte := resultBuffer.Bytes()
	mime := GetImageMimeType(bimg.DetermineImageType(resultByte))
	return Image{Body: resultByte, Mime: mime}, nil
}

func transformGIFFrame(m image.Image, opts bimg.Options) image.Image {
	// Parse crop and resize parameters before applying any transforms.
	// This is to ensure that any percentage-based values are based off the
	// size of the original image.
	rect := m.Bounds()
	w, h, resize := resizeGIFParams(m, opts)

	// crop if needed
	if !m.Bounds().Eq(rect) {
		m = imaging.Crop(m, rect)
	}
	// resize if needed
	if resize {
		if !opts.Crop {
			m = imaging.Fit(m, w, h, resampleFilter)
		} else {
			if w == 0 || h == 0 {
				m = imaging.Resize(m, w, h, resampleFilter)
			} else {
				m = imaging.Thumbnail(m, w, h, resampleFilter)
			}
		}
	}

	// rotate
	rotate := float64(opts.Rotate) - math.Floor(float64(opts.Rotate)/360)*360
	switch rotate {
	case 90:
		m = imaging.Rotate90(m)
	case 180:
		m = imaging.Rotate180(m)
	case 270:
		m = imaging.Rotate270(m)
	}

	// flip
	if opts.Flip {
		m = imaging.FlipV(m)
	}
	if opts.Flop {
		m = imaging.FlipH(m)
	}

	return m
}

func resizeGIFParams(m image.Image, opts bimg.Options) (w, h int, resize bool) {
	// convert percentage width and height values to absolute values
	imgW := m.Bounds().Dx()
	imgH := m.Bounds().Dy()
	w = EvaluateFloat(float64(opts.Width), imgW)
	h = EvaluateFloat(float64(opts.Height), imgH)

	// if requested width and height match the original, skip resizing
	if (w == imgW || w == 0) && (h == imgH || h == 0) {
		return 0, 0, false
	}

	return w, h, true
}

func EvaluateFloat(f float64, max int) int {
	if 0 < f && f < 1 {
		return int(float64(max) * f)
	}
	if f < 0 {
		return 0
	}
	return int(f)
}
