package pdf2zpl

import (
	"encoding/base64"
	"github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"simonwaldherr.de/go/zplgfa"
)

func Base64ToZpl(base64PDF string) string {
	data, err := base64.StdEncoding.DecodeString(base64PDF)
	if err != nil {
		log.Fatalf("Error decoding base64 PDF: %v", err)
	}

	tmpFile, err := ioutil.TempFile("", "temp_pdf_*.pdf")
	if err != nil {
		log.Fatalf("Error creating temporary PDF file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		log.Fatalf("Error writing to temporary PDF file: %v", err)
	}

	doc, err := fitz.New(tmpFile.Name())
	if err != nil {
		log.Fatalf("Error opening PDF: %v", err)
	}
	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		log.Fatalf("Error extracting image from PDF: %v", err)
	}

	newWidth := uint(500)
	resizedImg := resize.Resize(newWidth, 0, img, resize.Lanczos3)

	imgFile, err := ioutil.TempFile("", "temp_img_*.png")
	if err != nil {
		log.Fatalf("Error creating temporary image file: %v", err)
	}
	defer os.Remove(imgFile.Name())

	err = png.Encode(imgFile, resizedImg)
	if err != nil {
		log.Fatalf("Error encoding image to PNG: %v", err)
	}
	imgFile.Close()

	filename := filepath.Base(imgFile.Name())
	err = os.Rename(imgFile.Name(), filename)
	if err != nil {
		log.Fatalf("Error saving PNG image: %v", err)
	}
	defer os.Remove(filename) // Add this line
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening PNG image file: %v", err)
	}
	defer f.Close()

	img, err = png.Decode(f)
	if err != nil {
		log.Fatalf("Error decoding PNG image: %v", err)
	}

	flat := zplgfa.FlattenImage(img)

	gfimg := zplgfa.ConvertToZPL(flat, zplgfa.CompressedASCII)

	return gfimg
}
