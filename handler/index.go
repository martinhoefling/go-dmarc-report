package handler

import (
	"net/http"

	"github.com/martinhoefling/go-dmarc-report/report"
)

func Index(w http.ResponseWriter, r *http.Request) {
	feedbacks := report.RequestAllReports()
	renderTemplate(w, "index", feedbacks)
}
