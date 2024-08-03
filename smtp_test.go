package tcbs_test

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/sgaunet/tcbs"
	"github.com/stretchr/testify/assert"
	mail "github.com/xhit/go-simple-mail/v2"
)

func TestNewSmtpContainer(t *testing.T) {
	container, err := tcbs.NewTestSMTPServer()
	if err != nil {
		t.Fatalf("could not create smtp container: %v", err)
	}
	defer container.Terminate()

	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	// defer cancel()
	SMTPClient := mail.NewSMTPClient()
	SMTPClient.Host = container.GetSMTPHost()
	SMTPClient.Port = container.GetSMTPPortInt()
	SMTPClient.Username = container.GetSMTPUser()
	SMTPClient.Password = container.GetSMTPPassword()
	// SMTPClient.Encryption = mail.EncryptionTLS
	SMTPClient.Encryption = mail.EncryptionNone
	SMTPClient.KeepAlive = true
	SMTPClient.ConnectTimeout = 10 * time.Second
	SMTPClient.SendTimeout = 10 * time.Second
	SMTPClient.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	client, err := SMTPClient.Connect()
	assert.Nil(t, err)
	defer client.Close()
}
