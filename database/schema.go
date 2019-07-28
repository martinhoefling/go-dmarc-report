package database

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"time"
)

func EnsureDBSchema(db *pg.DB) {
	opts := orm.CreateTableOptions{IfNotExists: true}

	if err := db.CreateTable((*Feedback)(nil), &opts); err != nil {
		panic(err)
	}

	if err := db.CreateTable((*Record)(nil), &opts); err != nil {
		panic(err)
	}

	if err := db.CreateTable((*DKIMAuthResults)(nil), &opts); err != nil {
		panic(err)
	}

	if err := db.CreateTable((*SPFAuthResults)(nil), &opts); err != nil {
		panic(err)
	}
}

type Feedback struct {
	Version                 string
	OrganizationName        string
	Email                   string
	ReportID                string `sql:",pk"`
	DateBegin               time.Time
	DateEnd                 time.Time
	ContactInfo             string
	Domain                  string
	DKIM                    string
	SPF                     string
	Policy                  string
	SubdomainPolicy         string
	Percent                 int64
	FailureReportingOptions string
}

type Record struct {
	ReportID            string
	SourceIP            string
	Count               int64
	Disposition         string
	DKIMPolicyEvaluated string
	SPFPolicyEvaluated  string
	HeaderFrom          string
	EnvelopeFrom        string
	EnvelopeTo          string
}

type DKIMAuthResults struct {
	ReportID    string
	SourceIP    string
	Domain      string
	Result      string
	HumanResult string
	Selector    string
}

type SPFAuthResults struct {
	ReportID string
	SourceIP string
	Domain   string
	Result   string
	Scope    string
}
