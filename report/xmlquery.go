package report

type Query struct {
	Feedback Feedback `xml:"feedback"`
}

type Feedback struct {
	Version         string          `xml:"version"`
	Metadata        Metadata        `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Records         []Record        `xml:"record"` // Min 1
}

type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

type AuthResults struct {
	DKIM []DKIM `xml:"dkim"` // Min 0
	SPF  []SPF  `xml:"spf"`  // Min 1
}

type DKIM struct {
	Domain      string `xml:"domain"`       // Min 1
	Result      string `xml:"result"`       // Min 1
	HumanResult string `xml:"human_result"` // Min 1
	Selector    string `xml:"selector"`     // Min 0
}

type SPF struct {
	Domain string `xml:"domain"` // Min 1
	Result string `xml:"result"` // Min 1
	Scope  string `xml:"scope"`  // Min 1
}

type Metadata struct {
	OrganizationName string    `xml:"org_name"`
	Email            string    `xml:"email"`
	ReportID         string    `xml:"report_id"`
	DateRange        DateRange `xml:"date_range"`
	ContactInfo      string    `xml:"extra_contact_info"` // Max 1
}

type PolicyPublished struct {
	Domain                  string    `xml:"domain"`
	DKIM                    string    `xml:"adkim"`
	SPF                     string    `xml:"aspf"`
	Policy                  string    `xml:"p"`
	SubdomainPolicy         string    `xml:"sp"`
	Percent                 customInt `xml:"pct"`
	FailureReportingOptions string    `xml:"fo"`
}

type Row struct {
	SourceIP        string          `xml:"source_ip"`
	Count           customInt       `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"` // Min 1
}

type Identifiers struct {
	HeaderFrom   string `xml:"header_from"`   // Min 1
	EnvelopeFrom string `xml:"envelope_from"` // Min 1 ??
	EnvelopeTo   string `xml:"envelope_to"`   // Min 0
}

type PolicyEvaluated struct {
	Disposition string `xml:"disposition"`
	DKIM        string `xml:"dkim"`
	SPF         string `xml:"spf"`
}

type DateRange struct {
	Begin customTime `xml:"begin"`
	End   customTime `xml:"end"`
}
