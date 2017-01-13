package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/howeyc/gopass"
)

func Connect(server, username string) (*client.Client, error) {
	log.Print("IMAP Password: ")
	password, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}

	log.Println("Connecting to server...")
	// Connect to server
	c, err := client.DialTLS(server, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected")

	// Login
	if err := c.Login(username, string(password)); err != nil {
		err2 := c.Logout()
		if err2 != nil {
			log.Print(err2)
		}
		return nil, err
	}
	log.Println("Logged in")
	return c, nil
}

func ListMessages(c *client.Client) (list [][]byte) {
	// Get all messages
	seqset, err := imap.NewSeqSet("1:*")
	if err != nil {
		log.Fatal(err)
	}
	messageChan := make(chan *imap.Message)
	go func() {
		// c.Fetch closes the messages channel when done.
		if err := c.Fetch(seqset, []string{"ENVELOPE", "BODYSTRUCTURE"}, messageChan); err != nil {
			log.Fatal(err)
		}
	}()
	messages := []*imap.Message{}
	for msg := range messageChan {
		messages = append(messages, msg)
	}

	for _, msg := range messages {
		if msg == nil || msg.BodyStructure == nil {
			log.Printf("nil/bad message: %v", msg)
			continue
		}
		if strings.ToLower(msg.BodyStructure.MimeType) == "multipart" {
			found := false
			for i, part := range msg.BodyStructure.Parts {
				mimeType := strings.ToLower(part.MimeType)
				if mimeType == "application" {
					_, data, err := GetAttachment(c, msg.SeqNum, fmt.Sprintf("[%v]", i+1), part)
					if err != nil {
						log.Println(err)
						continue
					}
					list = append(list, data)
					found = true
				}
			}
			if !found {
				fmt.Println("No application part found :/")
				//fmt.Println("No application part found :/ Parts:")
				//for _, part := range msg.BodyStructure.Parts {
				//fmt.Println(part)
				//}
				//fmt.Println("------------------------------------")
			}
		} else if strings.ToLower(msg.BodyStructure.MimeType) == "application" {
			_, data, err := GetAttachment(c, msg.SeqNum, "[1]", msg.BodyStructure)
			if err != nil {
				log.Println(err)
				continue
			}
			list = append(list, data)
		}
	}
	return
}

func GetAttachment(c *client.Client, id uint32, part string, info *imap.BodyStructure) (string, []byte, error) {
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
		return "", nil, errors.New("No message returned")
	}
	filename, ok := info.Params["name"]
	if ok {
		fmt.Println("Filename:")
		fmt.Println(filename)
	} else {
		fmt.Println("No filename :(")
	}
	fmt.Println("---------------------------")
	for section, body := range msg.Body {
		if section.String() == fmt.Sprintf("BODY%v", part) {
			var bodyReader io.Reader
			bodyReader = body
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
	return "", nil, errors.New("No attachment found")
}

func GetAllAttachments(server, user, mailbox string) error {
	c, err := Connect(server, user)
	if err != nil {
		return err
	}

	// Don't forget to logout
	defer func() {
		err2 := c.Logout()
		if err2 != nil {
			log.Print(err2)
		}
	}()

	mbox, err := c.Select(mailbox, true)
	if err != nil {
		return err
	}

	log.Printf("Listing all messages in %v", mailbox)
	messageIds := ListMessages(c)
	log.Printf("There are %v messages in %v, of which %v have relevant attachments", mbox.Messages, mbox.Name, len(messageIds))
	for i, data := range messageIds {
		err = ioutil.WriteFile(fmt.Sprintf("attachment-%v.dat", i), data, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := GetAllAttachments("mail.foo.to:993", "johannes@ebke.org", "INBOX/admin/ebke.org/postmaster")
	//err := GetAllAttachments("mail.foo.to:993", "johannes@ebke.org", "INBOX/admin/symann.org/postmaster")
	if err != nil {
		log.Fatal(err)
	}
}
