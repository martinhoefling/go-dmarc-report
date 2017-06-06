package acquire

type dmarcReportEmailSubject struct {
	_         struct{} `^Report [dD]omain:\s*`
	Domain    string   `[^\s]+`
	_         struct{} `\s*Submitter:\s*`
	Submitter string   `[^\s]+`
	_         struct{} `\s*Report-ID:\s*<?`
	ReportID  string   `[^\s^>]+`
	_         struct{} `>?\s*`
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
