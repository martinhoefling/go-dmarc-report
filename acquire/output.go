package acquire

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"log"
)

func (subject uniqueDmarcReportEmailSubject) targetDir(base string) string {
	return filepath.Join(base, subject.Domain, subject.Submitter)
}

func (subject uniqueDmarcReportEmailSubject) targetPath(base string) string {
	targetFile := fmt.Sprintf("%s.xml", subject.ReportID)
	return filepath.Join(subject.targetDir(base), targetFile)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func ensureDirExists(path string) error {
	exists, err := pathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		return os.MkdirAll(path, 0700)
	}
	return nil
}

func filterDownloadedSubjects(subjects []uniqueDmarcReportEmailSubject, baseDir string) (filteredSubjects []uniqueDmarcReportEmailSubject, err error) {
	var exists bool
	for _, subject := range subjects {
		exists, err = pathExists(subject.targetPath(baseDir))
		if err != nil {
			return nil, err
		}
		if !exists {
			filteredSubjects = append(filteredSubjects, subject)
		}
	}
	return
}

func writeReport(msg *dmarcReportEmail, baseDir string) error {
	err := ensureDirExists(msg.targetDir(baseDir))
	if err != nil {
		return err
	}
	var xmlBytes []byte
	xmlBytes, err = unpackReport(msg.filename, msg.data)
	if err != nil {
		return err
	}
	log.Printf("Writing new report %s", msg.targetPath(baseDir))
	return ioutil.WriteFile(msg.targetPath(baseDir), xmlBytes, 0600)
}
