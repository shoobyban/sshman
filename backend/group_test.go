package backend

import "testing"

func TestGroupHasHost(t *testing.T) {
	g := Group{
		Name: "group1",
		Hosts: []*Host{
			{Alias: "host1"},
			{Alias: "host2"},
		},
	}
	if !g.HasHost("host1") {
		t.Errorf("Expected group to contain host1")
	}
	if g.HasHost("host3") {
		t.Errorf("Expected group to not contain host3")
	}
}
