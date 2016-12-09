package report

import "time"

var FeedbackChannel = make(chan map[string][]Feedback)
var RequestChannel = make(chan Request)

type Request struct {
	Domain    string
	StartDate time.Time
	EndDate   time.Time
}

func Repository(reportPath string) {
	reports := ReadReports(reportPath)
    for request := range RequestChannel {
        if request.Domain == "" {
            FeedbackChannel <- reports
        } else {
            var domainReports = make(map[string][]Feedback)
            domainReports[request.Domain] = reports[request.Domain]
            FeedbackChannel <- domainReports
        }
	}
}

func RequestAllReports() Request {
	return Request{}
}

func RequestDomainReports(domain string) Request {
	return Request{Domain: domain}
}
