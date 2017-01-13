package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martinhoefling/go-dmarc-report/report"
)

type SingleReportPage struct {
	Domain string
	Report report.Feedback
}

func SingleReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := vars["domain"]
	reportID := vars["id"]
	feedback := report.RequestDomainReport(domain, reportID)
	p := SingleReportPage{Domain: domain, Report: feedback}
	renderTemplate(w, "singlereport", p)
}
