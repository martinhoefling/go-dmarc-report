package report

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var CompiledFeedbackRegexp *regexp.Regexp

type customTime struct {
	time.Time
}

type customInt struct {
	int64
}

func (c customInt) String() string {
	return fmt.Sprintf("%v", c.int64)
}

func (c customInt) Int64() int64 {
	return c.int64
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

func readReport(path string, reports map[string][]Feedback) error {
	fmt.Printf("Loading %s\n", path)
	var file io.Reader
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	feedback, err := ReadFeedbackXML(fileBytes)
	if err != nil {
		return err
	}
	domain := feedback.PolicyPublished.Domain
	domainReports, ok := reports[domain]
	if !ok {
		domainReports = make([]Feedback, 0)
		reports[domain] = append(domainReports, *feedback)
	}
	reports[domain] = append(domainReports, *feedback)
	return nil
}

func ReadFeedbackXML(xmlBytes []byte) (*Feedback, error) {
	if CompiledFeedbackRegexp == nil {
		CompiledFeedbackRegexp = regexp.MustCompile("(?is:<feedback.*</feedback>)")
	}

	validBytes := CompiledFeedbackRegexp.Find(xmlBytes)
	feedback := &Feedback{}
	if err := xml.Unmarshal(validBytes, feedback); err != nil {
		return nil, err
	}
	return feedback, nil
}

func getVisitFunc(reports map[string][]Feedback) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".xml") {
			return readReport(path, reports)
		}
		return nil
	}
}

func ReadReports(reportPath string) (map[string][]Feedback, error) {
	reports := make(map[string][]Feedback)
	if CompiledFeedbackRegexp == nil {
		CompiledFeedbackRegexp = regexp.MustCompile("(?is:<feedback.*</feedback>)")
	}

	fmt.Printf("Loading Reports from %s\n", reportPath)

	err := filepath.Walk(reportPath, getVisitFunc(reports))
	if err != nil {
		return nil, err
	}
	return reports, nil
}
