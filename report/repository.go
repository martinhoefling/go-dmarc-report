package report

import "time"

var requestChannel = make(chan request)

type request struct {
	Domain    string
	reportID  string
	StartDate time.Time
	EndDate   time.Time
	Result    chan map[string][]Feedback
}

func Repository(reportPath string) {
	dmarcReports := ReadReports(reportPath)
	for req := range requestChannel {
		if req.Domain == "" {
			req.Result <- dmarcReports
		} else {
			var domainReports = make(map[string][]Feedback)
			if req.reportID == "" {
				domainReports[req.Domain] = dmarcReports[req.Domain]
			} else {
				for _, dmarcReport := range dmarcReports[req.Domain] {
					if dmarcReport.Metadata.ReportID == req.reportID {
						domainReports[req.Domain] = []Feedback{dmarcReport}
					}
				}
			}
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

func RequestDomainReport(domain string, id string) Feedback {
	req := request{Domain: domain, reportID: id, Result: make(chan map[string][]Feedback, 1)}
	requestChannel <- req
	result := <-req.Result
	return result[domain][0]
}
