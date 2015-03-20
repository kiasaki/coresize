package coresize

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path"

	"github.com/disintegration/gift"
)

const filechunk = 8192 // 8k
var supportedExts = []string{"png", "jpg"}

type ImageFile struct {
	Path string
	Hash string
}

func (i ImageFile) Name() string {
	return path.Base(i.Path)
}

func (i ImageFile) NameWithHash() string {
	return fmt.Sprintf("%s-%s", i.Hash, i.Name())
}

func (i ImageFile) Ext() string {
	return path.Ext(i.Path)[1:]
}

func (i ImageFile) FileType() string {
	// TODO Less naive implementation
	return "image/" + i.Ext()
}

// Returns true is image ext is supported
func (i ImageFile) Supported() bool {
	imageExt := i.Ext()
	for _, ext := range supportedExts {
		if ext == imageExt {
			return true
		}
	}
	return false
}

func (i ImageFile) Render(w io.Writer, x, y int, align string) error {
	// Open file
	file, err := os.Open(i.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return err
	}

	size := img.Bounds().Size()

	// Default to current size if x and y are 0
	if x == 0 && y == 0 {
		x = size.X
		y = size.Y
	}
	// Calculate same aspect ratio y
	if y == 0 {
		g := gift.New(gift.Resize(x, 0, gift.LanczosResampling))
		y = g.Bounds(img.Bounds()).Size().Y
	}
	// Calculate same aspect ratio x
	if x == 0 {
		g := gift.New(gift.Resize(0, y, gift.LanczosResampling))
		x = g.Bounds(img.Bounds()).Size().X
	}

	// Resize image keeping aspect ration, filling space we have
	g := gift.New()
	if size.X > size.Y {
		// image is wider than it's high, resize width after
		g.Add(gift.Resize(0, y, gift.LanczosResampling))
		g.Add(gift.Resize(x, 0, gift.LanczosResampling))
	} else {
		g.Add(gift.Resize(x, 0, gift.LanczosResampling))
		g.Add(gift.Resize(0, y, gift.LanczosResampling))
	}

	// Compute resized image keeping aspect ratio
	dst := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	// Calculate starting point in destination image
	var spx, spy int
	if align[0] == 'c' {
		spy = (y - dst.Bounds().Size().Y) / 2
	} else if align[0] == 'b' {
		spy = y - dst.Bounds().Size().Y
	}
	if align[1] == 'c' {
		spx = (x - dst.Bounds().Size().X) / 2
	} else if align[1] == 'r' {
		spx = x - dst.Bounds().Size().X
	}

	startingPoint := image.Pt(spx, spy)

	// Compute resized image in right canvas
	finalDst := image.NewRGBA(image.Rect(0, 0, x, y))
	sr := dst.Bounds()
	finalDstRect := sr.Add(startingPoint)
	draw.Draw(finalDst, finalDstRect, dst, sr.Min, draw.Src)

	if format == "jpg" {
		err = jpeg.Encode(w, finalDst, &jpeg.Options{
			Quality: 90,
		})
	} else if format == "png" {
		err = png.Encode(w, finalDst)
	} else {
		return errors.New("Unrecognized format: " + format)
	}

	return err
}

func (i *ImageFile) ComputeHash() error {
	file, err := os.Open(i.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// calculate the file size
	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf))
	}

	i.Hash = fmt.Sprintf("%x", hash.Sum(nil))[:8]
	return nil
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
