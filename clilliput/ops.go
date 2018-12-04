package clilliput

// ShukranImgOps is a reusable object that can resize and encode images.
type ShukranImgOps struct {
	*ImageOps
	frames     []*Framebuffer
	frameIndex int
}

func (o *ShukranImgOps) active() *Framebuffer {
	return o.frames[o.frameIndex]
}

func (o *ShukranImgOps) secondary() *Framebuffer {
	return o.frames[1-o.frameIndex]
}

func (o *ShukranImgOps) swap() {
	o.frameIndex = 1 - o.frameIndex
}

func (o *ShukranImgOps) decode(d Decoder) error {
	active := o.active()
	return d.DecodeTo(active)
}

func (o *ShukranImgOps) fit(d Decoder, width, height float64) error {
	active := o.active()
	secondary := o.secondary()
	err := active.FitDirection(width, height, secondary, FitTopLeft)
	if err != nil {
		return err
	}
	o.swap()
	return nil
}

func (o *ShukranImgOps) resize(d Decoder, width, height float64) error {
	active := o.active()
	secondary := o.secondary()
	err := active.ResizeTo(width, height, secondary)
	if err != nil {
		return err
	}
	o.swap()
	return nil
}

func (o *ShukranImgOps) normalizeOrientation(orientation ImageOrientation) {
	active := o.active()
	active.OrientationTransform(orientation)
}

// Transform performs the requested transform operations on the Decoder specified by d.
// The result is written into the output buffer dst. A new slice pointing to dst is returned
// with its length set to the length of the resulting image. Errors may occur if the decoded
// image is too large for ImageOps or if Encoding fails.
//
// It is important that .Decode() not have been called already on d.
func (o *ShukranImgOps) Transform(d Decoder, opt *ImageOptions, dst []byte) ([]byte, error) {
	_, err := d.Header()
	if err != nil {
		return nil, err
	}

	enc, err := newOpenCVEncoder(opt.FileType, d, dst)

	if err != nil {
		return nil, err
	}
	defer enc.Close()

	for {
		emptyFrame := false

		err = o.decode(d)
		if err != nil {
			//Commenting out failure case
			//if err != io.EOF {
			//	return nil, err
			//}
			// io.EOF means we are out of frames, so we should signal to encoder to wrap up
			emptyFrame = true
		}

		//o.normalizeOrientation(h.Orientation())

		if opt.ResizeMethod == ImageOpsFit {
			o.fit(d, opt.Width, opt.Height)
		} else if opt.ResizeMethod == ImageOpsResize {
			o.resize(d, opt.Width, opt.Height)
		}

		var content []byte
		if emptyFrame {
			content, err = o.encodeEmpty(enc, opt.EncodeOptions)
		} else {
			content, err = o.encode(enc, opt.EncodeOptions)
		}

		if err != nil {
			return nil, err
		}

		if content != nil {
			return content, nil
		}

		// content == nil and err == nil -- this is encoder telling us to do another frame

		// for mulitple frames/gifs we need the decoded frame to be active again
		o.swap()
	}
}

// NewTestImageOps creates a new ShukranImgOps object that will operate
// on images up to maxSize on each axis.
func NewImageOps(maxSize int) *ShukranImgOps {
	frames := make([]*Framebuffer, 2)
	//fmt.Println(NewFramebuffer(maxSize, maxSize))
	frames[0] = NewFramebuffer(maxSize, maxSize)
	frames[1] = NewFramebuffer(maxSize, maxSize)
	return &ShukranImgOps{
		frames:     frames,
		frameIndex: 0,
	}
}

// Clear resets all pixel data in ImageOps. This need not be called
// between calls to Transform. You may choose to call this to remove
// image data from memory.
func (o *ShukranImgOps) Clear() {
	o.frames[0].Clear()
	o.frames[1].Clear()
}

// Close releases resources associated with ImageOps
func (o *ShukranImgOps) Close() {
	o.frames[0].Close()
	o.frames[1].Close()
}

func (o *ShukranImgOps) encode(e Encoder, opt map[int]int) ([]byte, error) {
	active := o.active()
	return e.Encode(active, opt)
}

func (o *ShukranImgOps) encodeEmpty(e Encoder, opt map[int]int) ([]byte, error) {
	return e.Encode(nil, opt)
}
