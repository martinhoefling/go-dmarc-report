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
		fileBytes, err := ioutil.ReadFile(filePath)
		utils.CheckError(err)

		r, _ := regexp.Compile("(?is:<feedback.*</feedback>)")
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
