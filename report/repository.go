package report

import "time"

var requestChannel = make(chan request)

type request struct {
	Domain    string
	StartDate time.Time
	EndDate   time.Time
	Result    chan map[string][]Feedback
}

func Repository(reportPath string) {
	reports := ReadReports(reportPath)
	for req := range requestChannel {
		if req.Domain == "" {
			req.Result <- reports
		} else {
			var domainReports = make(map[string][]Feedback)
			domainReports[req.Domain] = reports[req.Domain]
			req.Result <- domainReports
		}
	}
}

func RequestAllReports() map[string][]Feedback {
	req := request{Result: make(chan map[string][]Feedback, 1)}
	requestChannel <- req
	return <-req.Result
}

func RequestDomainReports(domain string) map[string][]Feedback {
	req := request{Domain: domain, Result: make(chan map[string][]Feedback, 1)}
	requestChannel <- req
	return <-req.Result
}
