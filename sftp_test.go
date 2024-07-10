package tcbs_test

import (
	"log"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/melbahja/goph"
	"github.com/sgaunet/tcbs"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestNewSftpContainer(t *testing.T) {
	container, err := tcbs.NewTestSFTPServer()
	assert.Nil(t, err)
	defer container.Terminate()
	// init ssh client
	auth := goph.Password(container.GetPassword())
	host := strings.Split(container.GetEndpoint(), ":")[0]
	port := strings.Split(container.GetEndpoint(), ":")[1]
	portUint, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	client, err := goph.NewConn(&goph.Config{
		User:     container.GetUsername(),
		Addr:     host,
		Port:     uint(portUint),
		Auth:     auth,
		Callback: VerifyHost,
	})
	assert.Nil(t, err)
	client.Close()
}

// no fingerprint check
func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	// return goph.AddKnownHost(host, remote, key, "")
	return nil
}
