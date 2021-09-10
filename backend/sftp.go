package backend

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPConn struct {
	server      SFTPMockServer
	client      *sftp.Client
	mock        bool
	alias       string
	expected    string
	testServers map[string]SFTPMockServer
}

type SFTPMockServer struct {
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
	if serv, ok := s.testServers[host]; ok || s.mock {
		s.server = serv
		s.alias = host
		if s.testServers == nil {
			s.testServers = map[string]SFTPMockServer{}
		}
		return nil
	}
	key, err := ioutil.ReadFile(keyfile)
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

func (s *SFTPConn) GetServers() map[string]SFTPMockServer {
	return s.testServers
}

func (s *SFTPConn) Write(data string) error {
	if data == "" {
		return fmt.Errorf("empty data, not writing it")
	}
	if s.mock {
		if s.expected != "" && data != s.expected {
			return fmt.Errorf("data is not as expected: '%s' instead of '%s'", data, s.expected)
		}
		s.server.File = data
		s.testServers[s.alias] = s.server
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
		return []byte(s.server.File), nil
	}
	f, err := s.client.Open(".ssh/authorized_keys")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func (s *SFTPConn) Close() {
	if !s.mock {
		s.client.Close()
	}
}
