package mailer

import(
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//we want to store template into binary formate into // this variable

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer{
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data any)error{
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil{
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil{
		return nil
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil{
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)

	// declaring new message function for new mail
	// here we declare header of the mail and body of the mail for sending to user
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// if we fail to send message to smtp server then 
	// try again
	for i := 1; i <= 3; i++{
		err = m.dialer.DialAndSend(msg)
		if nil == err{
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	
	return nil
}
