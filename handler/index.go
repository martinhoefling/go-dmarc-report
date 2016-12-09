package handler

import (
	"net/http"

	"github.com/martinhoefling/go-dmarc-report/report"
)

func Index(w http.ResponseWriter, r *http.Request) {
	report.RequestChannel <- report.RequestAllReports()
	feedbacks := <-report.FeedbackChannel
	renderTemplate(w, "index", feedbacks)
}
