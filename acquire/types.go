package acquire

type dmarcReportEmailSubject struct {
	_         struct{} `regexp:"^Report [dD]omain:\\s*"`
	Domain    string   `regexp:"[^\\s]+"`
	_         struct{} `regexp:"\\s*Submitter:\\s*"`
	Submitter string   `regexp:"[^\\s]+"`
	_         struct{} `regexp:"\\s*Report-ID:\\s*<?"`
	ReportID  string   `regexp:"[^\\s^>]+"`
	_         struct{} `regexp:">?\\s*"`
}

type uniqueDmarcReportEmailSubject struct {
	dmarcReportEmailSubject
	UID uint32
}

type dmarcReportEmail struct {
	uniqueDmarcReportEmailSubject
	data     []byte
	filename string
}
