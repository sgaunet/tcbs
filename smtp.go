package tcbs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const SMTPDDockerImage = "axllent/mailpit:v1.14.0"

// SMTPServer
type SMTPServer struct {
	smtpServerC testcontainers.Container
	endpoint    string
}

// NewTestSMTPServer creates a new SMTP server for testing with default values
func NewTestSMTPServer() (*SMTPServer, error) {
	var err error
	newTestServer := &SMTPServer{}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        SMTPDDockerImage,
		ExposedPorts: []string{"1025/tcp"},
		WaitingFor:   wait.ForLog("[http] accessible via http://localhost:8025/"),
		Env:          map[string]string{},
	}
	newTestServer.smtpServerC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	now := time.Now()
	for time.Since(now) < 10*time.Second && newTestServer.endpoint == "" {
		newTestServer.endpoint, err = newTestServer.smtpServerC.Endpoint(ctx, "")
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(500 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, err
	}
	return newTestServer, nil
}

// Terminate stops the SMTP server
func (t *SMTPServer) Terminate() error {
	return t.smtpServerC.Terminate(context.Background())
}

// GetEndpoint returns the endpoint of the SFTP server
func (t *SMTPServer) GetSMTPEndpoint() string {
	return t.endpoint
}

// GetPort returns the port of the SFTP server
func (t *SMTPServer) GetSMTPPortStr() string {
	// endpoint is formatted as "host:port"
	// host := strings.Split(t.endpoint, ":")[0]
	port := strings.Split(t.endpoint, ":")[1]
	return port
}

func (t *SMTPServer) GetSMTPPortInt() int {
	// endpoint is formatted as "host:port"
	// host := strings.Split(t.endpoint, ":")[0]
	p := t.GetSMTPPortStr()
	port, err := strconv.Atoi(p)
	if err != nil {
		return 0
	}
	return port
}

func (t *SMTPServer) GetSMTPHost() string {
	// endpoint is formatted as "host:port"
	return strings.Split(t.endpoint, ":")[0]
}

func (t *SMTPServer) GetSMTPUser() string {
	return ""
}

func (t *SMTPServer) GetSMTPPassword() string {
	return ""
}
