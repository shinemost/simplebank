package mail

import (
	"testing"

	"github.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {

	//before to test send email with gmail , you must change this variable:EMAIL_SENDER_ADDRESS to correct address in the app.env file.
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailEmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"

	content := `
	<h1>Hello World</h1>
	<p>This is a test message  from <a href="http://github.com/shinemost">shinemost</a></p>
	`
	to := []string{"supertain147@163.com"}

	attatchFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attatchFiles)

	require.NoError(t, err)

}
