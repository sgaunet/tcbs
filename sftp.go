package tcbs

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const SSHDDockerImage = "sgaunet/alpine-sshd:latest"

// SFTPServer
type SFTPServer struct {
	sftpServerC        testcontainers.Container
	endpoint           string
	sshdUser           string
	sshdPassword       string
	sshdPrivateKey     string
	sshdFolderWritable string
}

// NewTestSFTPServer creates a new SFTP server for testing with default values
func NewTestSFTPServer() (*SFTPServer, error) {
	var err error
	newTestServer := &SFTPServer{
		sshdUser:           "sshuser", // default user for the sshd container, not modifiable
		sshdPassword:       "sshpassword",
		sshdPrivateKey:     "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCobaPuG8hWdO7t1Z0s7sPa/4DEVU4T1TvdrPuBi0PtYHGIGCUWeEDs6YCHNFjD6IMikWU2Tpi85AOVO2d8pA5XWwD+4D3E/Y1Hivs/orh7fyaBpfga4B91OXJuY+95WAiGl3gW15uEjDJ9u5KVKdY/jbXiLVkpsQZf0EgRrFkJXFZbdhulUVhAUEPKVUdp3ujCDYGrQVivChF2Qq6+c1aV66k0veIwCEdFWxGC92HFwkKgbyBhWHN6dHF2pGlL+6i4e3vUOeBiFQH4eN9CIsRuinK7W1Pzlhxcaw91LufbpPkbKbE8gJyn4zCA8AK/BV9cUw/WvxA6usOYhRN/rQS/h8jCzlZmkV/r7DS7PquHuSwD7+8TaOWX0NwT9hvrKAliK+YWW79DzJ+K0Ln6ZaDun/jhrmliP2DUYIBav1QNaDJ/7u/Z6sV+TUY50ZvX5enW71iA2VEdQR65399wSVlbKraNUeVHo24yHzXsKz9hbNJUZd8P3bSjugIp4oK/uS8=",
		sshdFolderWritable: "tmp",
	}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        SSHDDockerImage,
		ExposedPorts: []string{"22/tcp"},
		WaitingFor:   wait.ForLog("Server listening on 0.0.0.0"),
		Env: map[string]string{
			"AUTHORIZED_KEYS":       newTestServer.sshdPrivateKey,
			"SSHUSER_PASSWORD":      newTestServer.sshdPassword,
			"DATA_FOLDERS":          newTestServer.sshdFolderWritable,
			"DELETEOLDERFILESINMIN": "30",
		},
	}
	newTestServer.sftpServerC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	newTestServer.endpoint, err = newTestServer.sftpServerC.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}
	return newTestServer, nil
}

// Terminate stops the SFTP server
func (t *SFTPServer) Terminate() error {
	return t.sftpServerC.Terminate(context.Background())
}

// GetPrivateKey returns the private key of the SFTP server
func (t *SFTPServer) GetPrivateKey() string {
	return t.sshdPrivateKey
}

// GetEndpoint returns the endpoint of the SFTP server
func (t *SFTPServer) GetEndpoint() string {
	return t.endpoint
}

// GetUsername returns the username of the SFTP server
func (t *SFTPServer) GetUsername() string {
	return t.sshdUser
}

// GetPassword returns the password of the SFTP server
func (t *SFTPServer) GetPassword() string {
	return t.sshdPassword
}
