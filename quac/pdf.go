package quac

import (
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/rigelrozanski/thranch/quac/idea"
)

func ExportToPDF() error {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.SetFont("Courier", "", 7)
	writeHeight := float64(3)

	ideas := idea.GetAllIdeasNonConsuming()
	ideasText := ideas.WithText()
	ideasImage := ideas.WithImage()
	for _, idea := range ideasText {
		pdf.Write(writeHeight, fmt.Sprintf("%v\n%s\n_______________________________\n",
			idea.Filename, idea.GetContent()))
	}

	pdf.Write(writeHeight, fmt.Sprintf("________________IMAGES_______________\n"))

	var opt gofpdf.ImageOptions
	for _, idea := range ideasImage {
		pdf.Write(writeHeight, fmt.Sprintf("%v\n", idea.Filename))
		pdf.ImageOptions(idea.Path(), -10, 1, 0, 0, true, opt, 0, "")
		pdf.Write(writeHeight, fmt.Sprintf("_______________________________\n"))
	}

	err := pdf.OutputFileAndClose(os.ExpandEnv("$HOME/Desktop/quack_export.pdf"))
	if err != nil {
		return err
	}
	return nil
}
