package acquire

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/alexflint/go-restructure"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// Connect performs an interactive connection to the given IMAP server
func connect(server, username, password string) (*client.Client, error) {
	log.Printf("Connecting to %v...", server)
	c, err := client.DialTLS(server, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to %v.", server)

	if err := c.Login(username, password); err != nil {
		err2 := c.Logout()
		if err2 != nil {
			log.Print(err2)
		}
		return nil, err
	}
	log.Printf("Logged in as user %v on %v.", username, server)
	return c, nil
}

func getDmarcMessageSubjects(c *client.Client) (dmarcEmails []uniqueDmarcReportEmailSubject) {
	// Get all messages
	seqset, err := imap.NewSeqSet("1:*")
	if err != nil {
		log.Fatal(err)
	}
	messageChan := make(chan *imap.Message)
	go func() {
		if err := c.Fetch(seqset, []string{"ENVELOPE", "UID"}, messageChan); err != nil {
			log.Fatal(err)
		}
	}()
	for msg := range messageChan {
		var subject dmarcReportEmailSubject
		ok, err := restructure.Find(&subject, msg.Envelope.Subject)
		if err != nil {
			log.Fatal(err)
		}
		if ok {
			email := uniqueDmarcReportEmailSubject{subject, msg.Uid}
			log.Printf("Report for %s from %s with ID %s", email.Domain, email.Submitter, email.ReportID)
			dmarcEmails = append(dmarcEmails, email)
		} else {
			log.Printf("Subject: \"%s\" is not a dmarc report", msg.Envelope.Subject)
		}
	}
	return
}

func extractData(c *client.Client, msg *imap.Message) (string, []byte, error) {
	if msg == nil || msg.BodyStructure == nil {
		return "", nil, fmt.Errorf("nil/bad message: %v", msg)
	}
	if strings.ToLower(msg.BodyStructure.MimeType) == "multipart" {
		for i, part := range msg.BodyStructure.Parts {
			mimeType := strings.ToLower(part.MimeType)
			if mimeType == "application" {
				filename, data, err := getAttachment(c, msg.SeqNum, fmt.Sprintf("[%v]", i+1), part)
				if err != nil {
					log.Println(err)
					continue
				}
				return filename, data, nil
			}
		}
		return "", nil, fmt.Errorf("No application part found in message %v", msg)
	}

	if strings.ToLower(msg.BodyStructure.MimeType) == "application" {
		filename, data, err := getAttachment(c, msg.SeqNum, "[1]", msg.BodyStructure)
		if err != nil {
			return "", nil, err
		}
		return filename, data, nil
	}
	return "", nil, fmt.Errorf("No attachement found in message %v", msg)
}

// GetAttachment returns the specified attachment given the client
func getAttachment(c *client.Client, id uint32, part string, info *imap.BodyStructure) (string, []byte, error) {
	seqset := imap.SeqSet{}
	seqset.AddNum(id)
	messageChan := make(chan *imap.Message, 1)

	reqString := fmt.Sprintf("BODY.PEEK%v", part)
	err := c.Fetch(&seqset, []string{reqString}, messageChan)
	if err != nil {
		return "", nil, err
	}
	msg := <-messageChan
	if msg == nil {
		return "", nil, fmt.Errorf("No message with id %d", id)
	}
	filename, ok := info.Params["name"]
	if !ok {
		filename, ok = info.DispositionParams["filename"]
		if !ok {
			return "", nil, fmt.Errorf("No filename found in message  %v", msg)
		}
	}
	for section, body := range msg.Body {
		if section.String() == fmt.Sprintf("BODY%v", part) {
			bodyReader := io.Reader(body)
			if info.Encoding == "base64" {
				bodyReader = base64.NewDecoder(base64.StdEncoding, bodyReader)
			}
			data, err := ioutil.ReadAll(bodyReader)
			if err != nil {
				return "", nil, err
			}
			return filename, data, nil
		}
	}
	return "", nil, fmt.Errorf("No attachment found in msg %v", msg)
}

// getDmarcMessageBodyStructures returns a list of the byte value of attachments with the MIME type of application
func getDmarcMessageAttachments(c *client.Client, uniqueSubjects []uniqueDmarcReportEmailSubject, outputChan chan *dmarcReportEmail) {
	seqset := imap.SeqSet{}
	for _, dmarcMessage := range uniqueSubjects {
		seqset.AddNum(dmarcMessage.UID)
	}
	subjectMap := make(map[uint32]uniqueDmarcReportEmailSubject)
	for _, subject := range uniqueSubjects {
		subjectMap[subject.UID] = subject
	}

	messageChan := make(chan *imap.Message)
	go func() {
		if err := c.UidFetch(&seqset, []string{"ENVELOPE", "BODYSTRUCTURE"}, messageChan); err != nil {
			log.Fatal(err)
		}
	}()
	messages := []*imap.Message{}

	for msg := range messageChan {
		messages = append(messages, msg)
	}
	for _, msg := range messages {
		filename, data, err := extractData(c, msg)
		if err != nil {
			log.Fatal(err)
		}
		if data != nil && filename != "" {
			outputChan <- &dmarcReportEmail{subjectMap[msg.Uid], data, filename}
		}
	}
	close(outputChan)
	return
}
