package acquire

import "log"

func processNewAttachements(filterSubjects func([]uniqueDmarcReportEmailSubject) ([]uniqueDmarcReportEmailSubject, error), processReport func(*dmarcReportEmail) error, server, user, password, mailbox string) error {
	connection, err := connect(server, user, password)
	if err != nil {
		return err
	}

	defer func() {
		err2 := connection.Logout()
		if err2 != nil {
			log.Print(err2)
		}
	}()

	mbox, err := connection.Select(mailbox, true)
	if err != nil {
		return err
	}

	log.Printf("Listing all messages in %v", mailbox)
	uniqueSubjects := getDmarcMessageSubjects(connection)
	filteredSubjects, err := filterSubjects(uniqueSubjects)
	if err != nil {
		return err
	}
	log.Printf("There are %v messages in %v, of which %d have subjects matching to a report.", mbox.Messages, mbox.Name, len(uniqueSubjects))
	if len(filteredSubjects) == 0 {
		log.Print("No new reports found")
		return nil
	}
	log.Printf("%d reports are new and will be downloaded", len(filteredSubjects))

	outputChan := make(chan *dmarcReportEmail)
	go getDmarcMessageAttachments(connection, filteredSubjects, outputChan)

	for msg := range outputChan {
		err = processReport(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
