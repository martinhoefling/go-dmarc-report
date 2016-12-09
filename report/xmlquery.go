package report

type Query struct {
	Feedback Feedback `xml:"feedback"`
}

type Feedback struct {
	ReportMetadata  ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Records         []Record        `xml:"record"`
}

type ReportMetadata struct {
	OrganizationName string    `xml:"org_name"`
	Email            string    `xml:"email"`
	ReportID         string    `xml:"report_id"`
	DateRange        DateRange `xml:"date_range"`
}

type PolicyPublished struct {
	Domain  string    `xml:"domain"`
	DKIM    string    `xml:"adkim"`
	SPF     string    `xml:"aspf"`
	Policy  string    `xml:"p"`
	Percent customInt `xml:"pct"`
}

type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

type Row struct {
	SourceIP        string          `xml:"source_ip"`
	Count           customInt       `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

type AuthResults struct {
	DKIM DKIM `xml:"dkim"`
	SPF  SPF  `xml:"spf"`
}

type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

type DKIM struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

type SPF struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

type DateRange struct {
	Begin customTime `xml:"begin"`
	End   customTime `xml:"end"`
}
