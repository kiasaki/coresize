package coresize

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path"
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

	image, format, err := image.Decode(file)
	if err != nil {
		return err
	}
	fmt.Println(format)

	if format == "jpg" {
		err = jpeg.Encode(w, image, &jpeg.Options{
			Quality: 90,
		})
	} else if format == "png" {
		err = png.Encode(w, image)
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
