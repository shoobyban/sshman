package backend

import (
	"fmt"
)

type SFTPMockServer struct {
	Host string
	User string
	File string
}

type SFTPMock struct {
	server      SFTPMockServer
	alias       string
	expected    string
	testServers map[string]SFTPMockServer
}

func (s *SFTPMock) Connect(keyfileContent, host, user string) error {
	if serv, ok := s.testServers[host]; ok {
		s.server = serv
		s.alias = host
	} else {
		return fmt.Errorf("no such server")
	}
	return nil
}

func (s *SFTPMock) GetServers() map[string]SFTPMockServer {
	return s.testServers
}

func (s *SFTPMock) Write(data string) error {
	if data == "" {
		return fmt.Errorf("empty data, not writing it")
	}
	if s.expected != "" && data != s.expected {
		return fmt.Errorf("data is not as expected: '%s' instead of '%s'", data, s.expected)
	}
	s.server.File = data
	s.testServers[s.alias] = s.server
	return nil
}

func (s *SFTPMock) Read() ([]byte, error) {
	return []byte(s.server.File), nil
}

func (s *SFTPMock) Close() {

}
