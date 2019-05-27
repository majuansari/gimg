package image

import (
	"errors"
	"github.com/majuansari/gimg/clilliput"
	"github.com/majuansari/gimg/helpers"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": map[int]int{clilliput.JpegQuality: 85},
	".png":  map[int]int{clilliput.PngCompression: 7},
	".webp": map[int]int{clilliput.WebpQuality: 85},
}

func ResizeImage(outputWidth, outputHeight float64, imageUrl string) (string, error) {
	var outputFilename string
	var stretch bool

	imageBuffer, loadError := loadImageResponseInBuffer(imageUrl)

	if loadError == nil {
		outputImg, err := resizeImageBuffer(imageBuffer, outputFilename, outputWidth, outputHeight, stretch)
		return string(outputImg), err
	}

	return "", loadError
}

//Resize the image loaded in image buffer
func resizeImageBuffer(imageBuffer []byte, outputFilename string, outputWidth float64, outputHeight float64, stretch bool) ([]byte, error) {
	decoder, err := clilliput.NewOpenCVDecoder(imageBuffer)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		return []byte{}, err
	}
	defer decoder.Close()

	header, err := decoder.Header()
	ops := clilliput.NewImageOps(8192)
	defer ops.Close()
	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)
	// use user supplied filename to guess output type if provided
	// otherwise don"t transcode (use existing type)
	outputType := "." + strings.ToLower(decoder.Description())
	if outputFilename != "" {
		outputType = filepath.Ext(outputFilename)
	}

	aspectRatio := header.Width() / header.Height()

	outputWidth, outputHeight = helpers.ValidateOutputWidthAndHeight(outputWidth, aspectRatio, outputHeight)

	resizeMethod := clilliput.ImageOpsFit
	if stretch {
		resizeMethod = clilliput.ImageOpsResize
	}
	opts := &clilliput.ImageOptions{
		FileType:             outputType,
		Width:                outputWidth,
		Height:               outputHeight,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: false,
		EncodeOptions:        EncodeOptions[outputType],
	}
	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	return outputImg, err
}

//Load the image from image url to buffer
func loadImageResponseInBuffer(imageUrl string) ([]byte, error) {
	if imageUrl == "" {
		return []byte{}, errors.New("no input filename provided, quitting")
	}
	var client http.Client
	resp, err := client.Get(imageUrl)
	if err != nil {
		return []byte{}, err

	}
	defer resp.Body.Close()
	// decoder wants []byte, so read the whole file into a buffer
	inputBuf, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return []byte{}, err2
	}
	return inputBuf, err
}

func GetMimeType(fileName string) (string, error) {
	types := map[string]string{
		".png":  "image/png",
		".jpe":  "image/jpeg",
		".jpeg": "image/jpeg",
		".jpg":  "image/jpeg",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".ico":  "image/vnd.microsoft.icon",
		".tiff": "image/tiff",
		".tif":  "image/tiff",
		".svg":  "image/svg+xml",
		".svgz": "image/svg+xml",
	}
	extension := filepath.Ext(fileName)

	if len(types[extension]) >= 0 {
		return types[extension], nil
	}

	return "image/jpeg", nil

}
