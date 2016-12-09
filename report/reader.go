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

func parseTime(timestamp string) time.Time {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	utils.CheckError(err)
	return time.Unix(i, 0)
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	utils.CheckError(d.DecodeElement(&v, &start))
	parse := parseTime(strings.TrimSpace(v))
	*c = customTime{parse}
	return nil
}

func (c *customTime) UnmarshalXMLAttr(attr xml.Attr) error {
	parse := parseTime(strings.TrimSpace(attr.Value))
	*c = customTime{parse}
	return nil
}

func parseInt(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 0)
	utils.CheckError(err)
	return i
}

func (c *customInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	utils.CheckError(d.DecodeElement(&v, &start))
	parse := parseInt(strings.TrimSpace(v))
	*c = customInt{parse}
	return nil
}

func (c *customInt) UnmarshalXMLAttr(attr xml.Attr) error {
	parse := parseInt(strings.TrimSpace(attr.Value))
	*c = customInt{parse}
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
