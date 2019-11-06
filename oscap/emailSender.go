package oscap

import (
	"github.com/scorredoira/email"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
)

// EmailConf defines the details for smtp
type EmailConf struct {
	Smarthost string `yaml:"smtp_smarthost,omitempty"`
	From      string `yaml:"from,omitempty"`
	To        string `yaml:"to,omitempty"`
	Password  string `yaml:"password,omitempty"`
}

// SendFileViaEmail handles send report via email as attachment
func (emailconf *EmailConf) SendFileViaEmail(filePath string) error {
	if emailconf.Smarthost == "" || emailconf.From == "" || emailconf.To == "" {
		return nil
	}

	if errFileExists := fileExists(filePath); errFileExists != nil {
		return errFileExists
	}

	if errSendingEmail := emailconf.sendEmail(filePath); errSendingEmail != nil {
		return errSendingEmail
	}

	return nil

}

func (emailconf *EmailConf) sendEmail(filePath string) error {

	// Get hostname for email header. No need to exit if that fails
	hostname, errGetHostname := os.Hostname()
	if errGetHostname != nil {
		log.Printf("WARN: Could not get Hostname")
	}

	// compose the message
	m := email.NewMessage("Vulnerability report for "+hostname, "See Attachment for the vulnerability report of the server "+hostname)
	m.From = mail.Address{Name: "From", Address: emailconf.From}
	m.To = []string{emailconf.To}

	// add attachments
	if errAttach := m.Attach(filePath); errAttach != nil {
		log.Printf("Error: Could not attach file " + filePath + " to e-mail")
		return errAttach
	}

	// add headers
	// m.AddHeader("X-CUSTOMER-id", "xxxxx")
	log.Printf(strings.Split(emailconf.Smarthost, ":")[0])
	// send it
	auth := smtp.PlainAuth("", emailconf.From, emailconf.Password, strings.Split(emailconf.Smarthost, ":")[0])
	if errSendEmail := email.Send(emailconf.Smarthost, auth, m); errSendEmail != nil {
		log.Printf("Error: Could not send email via smarthost " + emailconf.Smarthost + " from user " + emailconf.From + " to user " + emailconf.To)
		return errSendEmail
	}

	return nil
}
