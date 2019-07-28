package acquire

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/martinhoefling/go-dmarc-report/database"
	"github.com/martinhoefling/go-dmarc-report/report"
	"log"
)

func filterDownloadedSubjectsByDB(subjects []uniqueDmarcReportEmailSubject, db *pg.DB) (filteredSubjects []uniqueDmarcReportEmailSubject) {
	for _, subject := range subjects {
		querystr := fmt.Sprintf("%%%s%%", subject.ReportID)
		_, err := db.QueryOne(&database.Feedback{}, `SELECT * FROM feedbacks WHERE report_id LIKE ?`, querystr)
		if err != nil {
			log.Printf("New report: %s", subject.ReportID)
			filteredSubjects = append(filteredSubjects, subject)
		}
	}
	return
}

func storeReport(msg *dmarcReportEmail, db *pg.DB) error {
	feedback, err := report.ReadFeedbackXML(msg.xml)
	if err != nil {
		return err
	}
	_, err = db.QueryOne(&database.Feedback{}, `SELECT * FROM feedbacks WHERE report_id = ?`, feedback.Metadata.ReportID)
	if err == nil {
		log.Printf("Report %s already in DB", feedback.Metadata.ReportID)
		return nil
	}

	log.Printf("Report %s will be ingested to DB", feedback.Metadata.ReportID)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	// Rollback tx on error.
	defer tx.Rollback()

	for _, record := range feedback.Records {
		for _, dkim := range record.AuthResults.DKIM {
			if err := db.Insert(&database.DKIMAuthResults{
				ReportID:    feedback.Metadata.ReportID,
				SourceIP:    record.Row.SourceIP,
				Domain:      dkim.Domain,
				Result:      dkim.Result,
				HumanResult: dkim.HumanResult,
				Selector:    dkim.Selector,
			}); err != nil {
				log.Fatal(err)
			}
		}
		for _, spf := range record.AuthResults.SPF {
			if err := db.Insert(&database.SPFAuthResults{
				ReportID: feedback.Metadata.ReportID,
				SourceIP: record.Row.SourceIP,
				Domain:   spf.Domain,
				Result:   spf.Result,
				Scope:    spf.Scope,
			}); err != nil {
				log.Fatal(err)
			}
		}
		if err := db.Insert(&database.Record{
			ReportID:            feedback.Metadata.ReportID,
			SourceIP:            record.Row.SourceIP,
			Count:               record.Row.Count.Int64(),
			Disposition:         record.Row.PolicyEvaluated.Disposition,
			DKIMPolicyEvaluated: record.Row.PolicyEvaluated.DKIM,
			SPFPolicyEvaluated:  record.Row.PolicyEvaluated.SPF,
			HeaderFrom:          record.Identifiers.HeaderFrom,
			EnvelopeFrom:        record.Identifiers.EnvelopeFrom,
			EnvelopeTo:          record.Identifiers.EnvelopeTo,
		}); err != nil {
			log.Fatal(err)
		}
	}
	if err := db.Insert(&database.Feedback{
		Version:                 feedback.Version,
		OrganizationName:        feedback.Metadata.OrganizationName,
		Email:                   feedback.Metadata.Email,
		ReportID:                feedback.Metadata.ReportID,
		DateBegin:               feedback.Metadata.DateRange.Begin.Time,
		DateEnd:                 feedback.Metadata.DateRange.End.Time,
		ContactInfo:             feedback.Metadata.ContactInfo,
		Domain:                  feedback.PolicyPublished.Domain,
		DKIM:                    feedback.PolicyPublished.DKIM,
		SPF:                     feedback.PolicyPublished.SPF,
		Policy:                  feedback.PolicyPublished.Policy,
		SubdomainPolicy:         feedback.PolicyPublished.SubdomainPolicy,
		Percent:                 feedback.PolicyPublished.Percent.Int64(),
		FailureReportingOptions: feedback.PolicyPublished.FailureReportingOptions,
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("Report %s will be committed to DB", feedback.Metadata.ReportID)
	return tx.Commit()
}

func IngestMissingAttachments(server, login, password, mailbox, dburl, dbpass string) error {
	db := database.OpenDBConnection(dburl, dbpass)
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
		}
	}()
	database.EnsureDBSchema(db)

	filter := func(subjects []uniqueDmarcReportEmailSubject) (filteredSubjects []uniqueDmarcReportEmailSubject, err error) {
		return filterDownloadedSubjectsByDB(subjects, db), nil
	}

	process := func(report *dmarcReportEmail) error {
		return storeReport(report, db)
	}

	return processNewAttachements(filter, process, server, login, password, mailbox)
}
