package acquire

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"strings"
)

func unpackReport(filename string, data []byte) ([]byte, error) {
	dataBuffer := bytes.NewBuffer(data)
	if strings.HasSuffix(filename, ".gz") {
		file, err := gzip.NewReader(dataBuffer)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(file)
	}
	if strings.HasSuffix(filename, ".zip") {
		bt, err := ioutil.ReadAll(dataBuffer)
		if err != nil {
			return nil, err
		}
		r, err := zip.NewReader(bytes.NewReader(bt), int64(len(bt)))
		if err != nil {
			return nil, err
		}
		if len(r.File) != 1 {
			return nil, fmt.Errorf("Not exactly one file in zip %s", filename)
		}
		f := r.File[0]
		if !strings.HasSuffix(f.Name, ".xml") {
			return nil, fmt.Errorf("Not an xml file in zip %s", filename)
		}
		file, err3 := f.Open()
		if err3 != nil {
			return nil, err3
		}
		return ioutil.ReadAll(file)
	}
	return nil, fmt.Errorf("Unknown extension of file %s", filename)
}
