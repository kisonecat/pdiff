package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
)

func FastCompare(img1, img2 *image.NRGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func VisualComparison(a string, b string, pageNumber int) error {
	aFile, err := ioutil.TempFile(os.TempDir(), "pdiff-a")
	defer os.Remove(aFile.Name())
	if err != nil {
		return err
	}

	bFile, err := ioutil.TempFile(os.TempDir(), "pdiff-b")
	defer os.Remove(bFile.Name())
	if err != nil {
		return err
	}

	cmdName := "mutool"
	aCmdArgs := []string{"draw", "-r", "50", "-c", "rgba", "-o", aFile.Name(), "-F", "png", a, fmt.Sprintf("%d", pageNumber)}
	aCmd := exec.Command(cmdName, aCmdArgs...)
	aCmd.Start()

	bCmdArgs := []string{"draw", "-r", "50", "-c", "rgba", "-o", bFile.Name(), "-F", "png", b, fmt.Sprintf("%d", pageNumber)}
	bCmd := exec.Command(cmdName, bCmdArgs...)
	bCmd.Start()

	// these should be parallelized
	aCmd.Wait()
	bCmd.Wait()

	aImageFile, err := os.Open(aFile.Name())
	defer aImageFile.Close()
	if err != nil {
		return err
	}
	aImage, err := png.Decode(aImageFile)
	if err != nil {
		return err
	}

	var aRgba *image.NRGBA
	var ok bool
	if aRgba, ok = aImage.(*image.NRGBA); !ok {
		return fmt.Errorf("Could not load RGBA format.")
	}

	bImageFile, err := os.Open(bFile.Name())
	defer bImageFile.Close()
	if err != nil {
		return err
	}
	bImage, err := png.Decode(bImageFile)
	if err != nil {
		return err
	}

	var bRgba *image.NRGBA
	if bRgba, ok = bImage.(*image.NRGBA); !ok {
		return fmt.Errorf("Could not load RGBA format.")
	}

	metric, err := FastCompare(aRgba, bRgba)
	if err != nil {
		return err
	}

	if metric > 1000 {
		return fmt.Errorf("On page %d, squared difference of %d\n", pageNumber, metric)
	}

	return nil
}
