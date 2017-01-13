package handler

import (
	"net/http"

	"time"

	"github.com/gorilla/mux"
	"github.com/martinhoefling/go-dmarc-report/report"
)

const lastReportsCount = 10

type DomainPage struct {
	Domain              string
	Reports             []report.Feedback
	ReportsLastWeek     []report.Feedback
	ReportsLastMonth    []report.Feedback
	ReportsLastSixMonth []report.Feedback
	ReportsLastYear     []report.Feedback
	LastReports         []report.Feedback
}

func filterByDate(p DomainPage) DomainPage {
	p.ReportsLastWeek = []report.Feedback{}
	p.ReportsLastMonth = []report.Feedback{}
	p.ReportsLastSixMonth = []report.Feedback{}
	p.ReportsLastYear = []report.Feedback{}
	for _, rep := range p.Reports {
		if time.Now().Before(rep.Metadata.DateRange.End.AddDate(0, 0, 7)) {
			p.ReportsLastWeek = append(p.ReportsLastWeek, rep)
		}
		if time.Now().Before(rep.Metadata.DateRange.End.AddDate(0, 1, 0)) {
			p.ReportsLastMonth = append(p.ReportsLastMonth, rep)
		}
		if time.Now().Before(rep.Metadata.DateRange.End.AddDate(0, 6, 0)) {
			p.ReportsLastSixMonth = append(p.ReportsLastSixMonth, rep)
		}
		if time.Now().Before(rep.Metadata.DateRange.End.AddDate(1, 0, 0)) {
			p.ReportsLastYear = append(p.ReportsLastYear, rep)
		}
	}
	return p
}

func Domain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	feedbacks := report.RequestDomainReports(domain)
	p := DomainPage{Domain: domain, Reports: feedbacks[domain]}
	p = filterByDate(p)
	var fromReport int
	if len(p.Reports) > lastReportsCount {
		fromReport = len(p.Reports) - 10
	} else {
		fromReport = 0
	}
	p.LastReports = p.Reports[fromReport:]
	renderTemplate(w, "domain", p)
}
