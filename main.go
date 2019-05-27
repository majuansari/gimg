package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/majuansari/gimg/barcode"
	"github.com/majuansari/gimg/config"
	"github.com/majuansari/gimg/crypto"
	shukranImage "github.com/majuansari/gimg/image"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	Log *log.Logger
)

//Http Handler
func handleRequest(w http.ResponseWriter, r *http.Request) {

	var response, mimeKey string
	var err error

	pathSegments := strings.Split(r.URL.Path, "/")

	query := r.URL.Query()
	presetKey := query.Get("p")
	dpr := query.Get("dpr")
	intDpr, _ := strconv.ParseInt(dpr, 0, 64)

	if pathSegments[1] == "card" {
		if intDpr <= 1 {
			intDpr = 2
		}
		width, height := getWidthAndHeight(int(intDpr), config.BarcodePresets[presetKey])

		code := crypto.Decrypt(pathSegments[2], os.Getenv("CRYPTO_PASS"), os.Getenv("CRYPTO_SALT"))
		response, err = barcode.GenerateBarcode(code, width, height, int(intDpr))
		mimeKey = "image/png"
	} else {
		if len(config.ImgPresets[presetKey]) > 0 && intDpr <= 4 {
			width, height := getWidthAndHeight(int(intDpr), config.ImgPresets[presetKey])

			imageUrl := os.Getenv("CDN_BASE_URL") + r.URL.Path

			mimeKey, _ = shukranImage.GetMimeType(imageUrl)
			response, err = shukranImage.ResizeImage(width, height, imageUrl)
		} else {
			err = errors.New("preset or dpr issue")
		}

	}

	if err == nil {
		w.Header().Set("Cache-Control", "public, s-maxage=31536000, max-age=31536000")
	} else {
		Log.Println(err)
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	}
	w.Header().Set("Content-Type", mimeKey)

	w.Write([]byte(response))

}

func newLog(logpath string) *log.Logger {
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	return log.New(file, "", log.LstdFlags|log.Lshortfile)
}

//Get width and height from preset
func getWidthAndHeight(intDpr int, preset map[string]int) (float64, float64) {
	if intDpr <= 0 {
		intDpr = 1
	}
	width := float64(preset["w"] * int(intDpr))
	height := float64(preset["h"] * int(intDpr))
	return width, height
}

func main() {
	err := godotenv.Load()

	Log = newLog("logs/info.log")

	if err != nil {
		Log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/favicon.ico", handleFavicon)

	fmt.Println("Starting Server")
	if err := http.ListenAndServe(":3434", nil); err != nil {
		panic(err)
	}
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {

}
