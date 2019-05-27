package barcode

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/majuansari/gimg/helpers"
	"image"
	"image/color"
	"image/png"
)

//Generate barcode
func GenerateBarcode(code string, width float64, height float64, dpr int) (string, error) {

	if len(code) <= 0 {
		return "", errors.New("empty code received")
	}

	//encodeBarcodeText(code)
	//fmt.Println("Generating Codabar for : ",encodeBarcodeText code)
	var bcode barcode.Barcode
	var barcodeWidth float64
	var barcodeHeight float64
	var aspectRatio float64

	aspectRatio = 4.385

	labelHeight := int(height*.3) * dpr
	fmt.Println(labelHeight, height*.25)

	bcode, err := code39.Encode(code, false, false)
	if err != nil {
		fmt.Printf("String %s cannot be encoded", code)
	}

	barcodeWidth, barcodeHeight = helpers.ValidateOutputWidthAndHeight(width, aspectRatio, height)

	intBarcodeWidth := int(barcodeWidth) * dpr
	intBarcodeHeight := (int(barcodeHeight) * dpr) - labelHeight

	bcode, err = barcode.Scale(bcode, intBarcodeWidth, intBarcodeHeight)

	if err != nil {
		fmt.Println("Code scaling error!", err)
	}

	// create the output file
	//file, _ := os.Create("qrcode.png")
	//defer file.Close()
	// encode the barcode as png
	//png.Encode(file, bcode)

	// create a new blank image with white background
	blankImage := imaging.New(intBarcodeWidth, intBarcodeHeight+(labelHeight), color.NRGBA{255, 255, 255, 255})
	newImgWithBarcode := imaging.Paste(blankImage, bcode, image.Pt(0, 0))

	gc := gg.NewContextForImage(newImgWithBarcode)

	//gc.SetRGB(1, 1, 1)
	//gc.Clear()
	gc.SetRGB(0, 0, 0)

	//@todo check the typecasting
	fontPoints := float64(labelHeight)

	if err := gc.LoadFontFace("assets/fonts/SofiaPro-Bold.ttf", fontPoints); err != nil {
		panic(err)
	}
	gc.DrawStringAnchored(code, float64(intBarcodeWidth)/2, float64(intBarcodeHeight)+(20*float64(dpr)), 0.5, 0.5)

	buffer := new(bytes.Buffer)
	png.Encode(buffer, gc.Image())

	if err != nil {
		fmt.Println(err)
	}
	// everything ok

	return string(buffer.Bytes()), nil

}

//func arrayToString(a []int) string {
//	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", "", -1), "[]")
//}
//
//func encodeBarcodeText(code string) string {
//
//	barcodeTextAsciiArray := []int{}
//	extraNumber := "865688"
//	for _, v := range code {
//		barcodeTextAsciiArray = append(barcodeTextAsciiArray, int(v))
//	}
//
//	barcodeTextAsciiString := arrayToString(barcodeTextAsciiArray)
//	fmt.Println(extraNumber + barcodeTextAsciiString)
//	n := new(big.Int)
//	barcodeTextAsciiInt, _ := n.SetString(extraNumber+barcodeTextAsciiString, 10)
//	fmt.Println(barcodeTextAsciiInt)
//
//	hd := hashids.NewData()
//	hd.Salt = "Y3^9Yl#PnpAJ"
//	hd.MinLength = 6
//	h, _ := hashids.NewWithData(hd)
//	e, _ := h.Encode([]int{11, 22})
//
//	fmt.Println(e)
//	d, _ := h.DecodeWithError(e)
//	return string(d[0])
//}
