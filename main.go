package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	unicommon "github.com/unidoc/unipdf/v3/common"
	pdf "github.com/unidoc/unipdf/v3/model"
)

var (
	sourceDirectory *string
	output          *string
)

func init() {
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelInfo))

	sourceDirectory = flag.String("source", "", "Source directory containing PDF files to merge")
	output = flag.String("outfile", "", "Outfile path and name")
}

func main() {

	flag.Parse()

	files, err := ioutil.ReadDir(*sourceDirectory)
	if err != nil {
		panic(err)
	}

	var sources = []string{}
	for _, f := range files {
		sources = append(sources, filepath.Join(*sourceDirectory, f.Name()))
	}

	err = mergePDF(sources, *output)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Complete, see output file: %s\n", *output)
}

func mergePDF(inputPaths []string, outputPath string) error {

	pdfWriter := pdf.NewPdfWriter()

	for _, inputPath := range inputPaths {
		f, err := os.Open(inputPath)
		if err != nil {
			return err
		}

		defer f.Close()

		pdfReader, err := pdf.NewPdfReader(f)
		if err != nil {
			return err
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
