package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTP struct {
	client *sftp.Client
}

func connect(keyfile, host, user string) (*SFTP, error) {
	if strings.HasPrefix(keyfile, "~/") {
		home, _ := os.UserHomeDir()
		keyfile = filepath.Join(home, "/", keyfile[2:])
	}
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout: 3 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//log.Printf("host: %s %v\n", hostname, remote)
			return nil
		},
	}
	connection, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %s", err)
	}

	client, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	return &SFTP{client: client}, nil
}

func (s *SFTP) Write(data string) error {
	if data == "" {
		return fmt.Errorf("empty data, not writing it")
	}
	f, err := s.client.OpenFile(".ssh/authorized_keys", os.O_RDWR|os.O_TRUNC)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(data)); err != nil {
		return err
	}
	f.Close()
	return nil
}

func (s *SFTP) Read() ([]byte, error) {
	f, err := s.client.Open(".ssh/authorized_keys")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func (s *SFTP) Close() {
	s.client.Close()
}
