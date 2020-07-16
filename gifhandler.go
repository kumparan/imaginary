package main

import (
	"errors"
	"fmt"

	"github.com/discordapp/lilliput"
	"github.com/kumparan/bimg"
)

// handleGIFResizing handle gif and animated gif resizing using lilliput library
func handleGIFResizing(buf []byte, opts bimg.Options) (out Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch value := r.(type) {
			case error:
				err = value
			case string:
				err = errors.New(value)
			default:
				err = errors.New("lilliput internal error")
			}
			out = Image{}
		}
	}()

	// create new decoder
	decoder, err := lilliput.NewDecoder(buf)

	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		return Image{}, NewError(fmt.Sprintf("error decoding image, %s\n", err), InternalError)
	}
	defer decoder.Close()

	header, err := decoder.Header()
	// this error is much more comprehensive and reflects
	// format errors
	if err != nil {
		return Image{}, NewError(fmt.Sprintf("error reading image header, %s\n", err), InternalError)
	}

	ops := lilliput.NewImageOps(10000)
	defer ops.Close()

	// create a buffer to store the output image, 20MB in this case
	outputImg := make([]byte, 20*1024*1024)

	resizeMethod := lilliput.ImageOpsNoResize
	if opts.Crop {
		resizeMethod = lilliput.ImageOpsFit
	}

	if opts.Width == header.Width() && opts.Height == header.Height() {
		resizeMethod = lilliput.ImageOpsNoResize
	}

	// lilliput image opts
	lilliputImageOpts := &lilliput.ImageOptions{
		FileType:             ".gif",
		Width:                opts.Width,
		Height:               opts.Height,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
	}

	// resize and transcode image
	outputImg, err = ops.Transform(decoder, lilliputImageOpts, outputImg)
	if err != nil {
		return Image{}, NewError(fmt.Sprintf("error transforming image, %s\n", err), InternalError)
	}

	mime := GetImageMimeType(bimg.DetermineImageType(outputImg))
	return Image{Body: outputImg, Mime: mime}, nil
}
