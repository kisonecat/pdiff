package main

import (
	"errors"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"rsc.io/pdf"
)

func comparePageText(aPage pdf.Page, bPage pdf.Page) error {
	aContent := aPage.Content()
	bContent := bPage.Content()

	aText := ""
	bText := ""

	for _, t := range aContent.Text {
		aText = aText + t.S
	}
	for _, t := range bContent.Text {
		bText = bText + t.S
	}

	if aText != bText {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(aText, bText, false)
		fmt.Println(dmp.DiffPrettyText(diffs))
		return errors.New("Difference in text.")
	}

	return nil
}

func Compare(a string, b string) error {
	aReader, err := pdf.Open(a)
	if err != nil {
		return err
	}

	bReader, err := pdf.Open(b)
	if err != nil {
		return err
	}

	if aReader.NumPage() != bReader.NumPage() {
		return errors.New("Number of pages differs.")
	}

	for pageNumber := 1; pageNumber <= aReader.NumPage(); pageNumber++ {
		aPage := aReader.Page(pageNumber)
		bPage := bReader.Page(pageNumber)
		err = comparePageText(aPage, bPage)
		if err != nil {
			log.Error("Difference on page " + fmt.Sprintf("%d", pageNumber))
			return err
		}

		err = VisualComparison(a, b, pageNumber)
		if err != nil {
			log.Error("Visual difference on page " + fmt.Sprintf("%d", pageNumber))
			return err
		}
	}

	return nil
}
