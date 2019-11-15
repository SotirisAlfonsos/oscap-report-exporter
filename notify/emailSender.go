package notify

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/scorredoira/email"
	"net/mail"
	"net/smtp"
	"os"
	"oscap-report-exporter/common"
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
func (emailconf *EmailConf) SendFileViaEmail(filePath string, logger log.Logger) error {
	if emailconf.Smarthost == "" || emailconf.From == "" || emailconf.To == "" {
		level.Debug(logger).Log("msg", "skiping email send. No smarthost, from or to defined")
		return nil
	}

	if errFileExists := common.FileExists(filePath); errFileExists != nil {
		return errFileExists
	}

	if errSendingEmail := emailconf.sendEmail(filePath, logger); errSendingEmail != nil {
		return errSendingEmail
	}

	return nil

}

func (emailconf *EmailConf) sendEmail(filePath string, logger log.Logger) error {

	// Get hostname for email header. No need to exit if that fails
	hostname, errGetHostname := os.Hostname()
	if errGetHostname != nil {
		level.Warn(logger).Log("msg", "could not get Hostname")
	}

	// compose the message
	m := email.NewMessage("Vulnerability report for "+hostname, "See Attachment for the vulnerability report of the server "+hostname)
	m.From = mail.Address{Name: "Oscap exporter", Address: emailconf.From}
	m.To = []string{emailconf.To}

	// add attachments
	if err := m.Attach(filePath); err != nil {
		return errors.Wrap(err, "could not attach file "+filePath+" to e-mail")
	}

	// send it
	if err := email.Send(emailconf.Smarthost, emailconf.configureAuth(), m); err != nil {
		return errors.Wrap(err, "could not send email via smarthost "+emailconf.Smarthost+" from user "+emailconf.From+" to user "+emailconf.To)
	}

	return nil
}

func (emailconf *EmailConf) configureAuth() smtp.Auth {
	if emailconf.Password != "" {
		return smtp.PlainAuth("", emailconf.From, emailconf.Password, strings.Split(emailconf.Smarthost, ":")[0])
	}
	return nil
}
