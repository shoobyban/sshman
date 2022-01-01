package backend

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPConn struct {
	host      SFTPMockHost
	client    *sftp.Client
	mock      bool
	alias     string
	expected  string
	testHosts map[string]SFTPMockHost
	testError bool
}

type SFTPMockHost struct {
	Host string
	User string
	File string
}

type SFTP interface {
	Connect(keyfile, host, user string) error
	Write(data string) error
	Read() ([]byte, error)
	Close()
}

func (s *SFTPConn) Connect(keyfile, host, user string) error {
	if strings.HasPrefix(keyfile, "~/") {
		home, _ := os.UserHomeDir()
		keyfile = filepath.Join(home, "/", keyfile[2:])
	}
	if serv, ok := s.testHosts[host]; ok || s.mock {
		s.host = serv
		s.alias = host
		if s.testHosts == nil {
			s.testHosts = map[string]SFTPMockHost{}
		}
		return nil
	}
	key, err := os.ReadFile(keyfile)
	if err != nil {
		return fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User:    user,
		Auth:    []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout: 3 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	connection, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return err
	}
	s.client, err = sftp.NewClient(connection)
	if err != nil {
		return err
	}
	return nil
}

func (s *SFTPConn) GetHosts() map[string]SFTPMockHost {
	return s.testHosts
}

func (s *SFTPConn) SetError(willError bool) {
	s.testError = willError
}

func (s *SFTPConn) Write(data string) error {
	if data == "" || s.testError {
		return fmt.Errorf("empty data, not writing it")
	}
	if s.mock {
		if (s.expected != "" && data != s.expected) || s.testError {
			return fmt.Errorf("data is not as expected: '%s' instead of '%s'", data, s.expected)
		}
		s.host.File = data
		s.testHosts[s.alias] = s.host
		return nil
	}
	f, err := s.client.OpenFile(".ssh/authorized_keys", os.O_RDWR|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write([]byte(data)); err != nil {
		return err
	}
	return nil
}

func (s *SFTPConn) Read() ([]byte, error) {
	if s.mock {
		if s.testError {
			return nil, fmt.Errorf("test error reading file")
		}
		return []byte(s.host.File), nil
	}
	f, err := s.client.Open(".ssh/authorized_keys")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func (s *SFTPConn) Close() {
	if !s.mock {
		s.client.Close()
	}
}
