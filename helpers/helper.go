package helpers

// Validate out width and height and make sure it matches aspect ratio
func ValidateOutputWidthAndHeight(outputWidth float64, aspectRatio float64, outputHeight float64) (float64, float64) {
	if outputWidth == 0 {
		outputWidth = aspectRatio * outputHeight
	}
	if outputHeight == 0 {
		outputHeight = outputWidth / aspectRatio
	}
	return outputWidth, outputHeight
}
