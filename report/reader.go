package report

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"compress/gzip"

	"io"
	"os"

	"archive/zip"

	"bytes"

	"github.com/martinhoefling/go-dmarc-report/utils"
)

type customTime struct {
	time.Time
}

type customInt struct {
	int64
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return err
	}
	*c = customTime{time.Unix(i, 0)}
	return nil
}

func (c *customTime) UnmarshalXMLAttr(attr xml.Attr) error {
	i, err := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
	if err != nil {
		return err
	}
	*c = customTime{time.Unix(i, 0)}
	return nil
}

func (c *customInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return err
	}
	*c = customInt{i}
	return nil
}

func (c *customInt) UnmarshalXMLAttr(attr xml.Attr) error {
	i, err := strconv.ParseInt(strings.TrimSpace(attr.Value), 10, 64)
	if err != nil {
		return err
	}
	*c = customInt{i}
	return nil
}

func ReadReports(reportPath string) map[string][]Feedback {
	files, err := ioutil.ReadDir(reportPath)
	utils.CheckError(err)
	reports := make(map[string][]Feedback)
	for _, f := range files {
		filePath, err := filepath.Abs(path.Join(reportPath, f.Name()))
		utils.CheckError(err)
		fmt.Printf("Loading %s\n", filePath)
		var file io.Reader
		file, err = os.Open(filePath)
		utils.CheckError(err)
		if strings.HasSuffix(filePath, ".gz") {
			file, err = gzip.NewReader(file)
			utils.CheckError(err)
		}
		if strings.HasSuffix(filePath, ".zip") {
			bt, err2 := ioutil.ReadAll(file)
			utils.CheckError(err2)
			r, err2 := zip.NewReader(bytes.NewReader(bt), int64(len(bt)))
			utils.CheckError(err2)
			if len(r.File) != 1 {
				fmt.Printf("Not exactly one file in zip %s", filePath)
				continue
			}
			f := r.File[0]
			if !strings.HasSuffix(f.Name, ".xml") {
				fmt.Printf("Not an xml file in zip %s", filePath)
				continue
			}
			file, err2 = f.Open()
			utils.CheckError(err2)
		}
		r, _ := regexp.Compile("(?is:<feedback.*</feedback>)")
		fileBytes, err := ioutil.ReadAll(file)
		utils.CheckError(err)
		validBytes := r.Find(fileBytes)
		var q Query
		utils.CheckError(xml.Unmarshal(validBytes, &q.Feedback))
		domain := q.Feedback.PolicyPublished.Domain
		domainReports, ok := reports[domain]
		if !ok {
			domainReports = make([]Feedback, 0)
			reports[domain] = append(domainReports, q.Feedback)
		}
		reports[domain] = append(domainReports, q.Feedback)
	}
	return reports
}
