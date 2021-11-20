package internal

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/golang/freetype"
	"golang.org/x/image/font/gofont/goregular"
)

func CreateImageWithOverlayTexts(baseImgPath string, textPlaces ...TextPlacement) (string, error) {
	data, err := os.ReadFile(baseImgPath)
	if err != nil {
		return "", err
	}

	baseImg, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	drawImg := image.NewRGBA(image.Rect(0, 0, baseImg.Bounds().Max.X, baseImg.Bounds().Max.Y))
	draw.Draw(drawImg, drawImg.Bounds(), baseImg, image.ZP, draw.Src)

	for i := range textPlaces {
		err = DrawTextOnImage(drawImg, textPlaces[i])
		if err != nil {
			return "", err
		}
	}

	tmpfile, err := os.CreateTemp("", patternizeFilename(filepath.Base(baseImgPath)))
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	err = png.Encode(tmpfile, drawImg)
	if err != nil {
		return tmpfile.Name(), err
	}

	return tmpfile.Name(), tmpfile.Sync()
}

func DrawTextOnImage(img draw.Image, textPlace TextPlacement) error {
	utf8Font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		return err
	}

	// Draw text
	dpi := float64(72)
	fontsize := float64(textPlace.FontSize)
	//spacing := float64(1.5)

	ctx := new(freetype.Context)
	ctx = freetype.NewContext()
	ctx.SetDPI(dpi)
	ctx.SetFont(utf8Font)
	ctx.SetFontSize(fontsize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)

	// get the size of the text
	pt := freetype.Pt(0, 0)                              // int(ctx.PointToFixed(fontsize)>>6))
	ctx.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 0})) // transparent
	extent, err := ctx.DrawString(textPlace.Text, pt)
	if err != nil {
		return err
	}

	//
	pt = freetype.Pt(int(textPlace.X), int(textPlace.Y))
	if textPlace.AnchorX == 0 {
		pt.X -= extent.X / 2
	} else if textPlace.AnchorX > 0 {
		pt.X -= extent.X
	}
	if textPlace.AnchorY == 0 {
		pt.Y -= extent.Y / 2
	} else if textPlace.AnchorY > 0 {
		pt.Y -= extent.Y
	}

	ctx.SetSrc(image.NewUniform(textPlace.Color))
	_, err = ctx.DrawString(textPlace.Text, pt)
	if err != nil {
		return err
	}
	return nil
}
