package cmd

import (
    "bytes"
    "strings"
    "testing"

    "github.com/shoobyban/sshman/backend"
    "github.com/stretchr/testify/assert"
)

func TestRoleListCmdOutputsRoles(t *testing.T) {
    // Ensure default roles are initialized
    _ = backend.DefaultConfig()

    var buf bytes.Buffer
    roleListCmd.SetOut(&buf)
    roleListCmd.Run(roleListCmd, []string{})
    out := buf.String()
    if !strings.Contains(out, "Role: admin") {
        t.Fatalf("expected output to contain admin role, got: %s", out)
    }
}

func TestRoleAssignInvalidRolePrintsError(t *testing.T) {
    var buf bytes.Buffer
    roleAssignCmd.SetOut(&buf)
    // set role flag to nonexistent
    _ = roleAssignCmd.Flags().Set("role", "no-such-role")
    // ensure no user/host to hit the missing role branch
    _ = roleAssignCmd.Flags().Set("user", "")
    _ = roleAssignCmd.Flags().Set("host", "")
    roleAssignCmd.Run(roleAssignCmd, []string{})
    out := buf.String()
    assert.Contains(t, out, "Role no-such-role does not exist")
}

func TestRoleAssignHostRejected(t *testing.T) {
    var buf bytes.Buffer
    roleAssignCmd.SetOut(&buf)
    // use existing role
    _ = roleAssignCmd.Flags().Set("role", "admin")
    // host supplied without user should trigger host error branch
    _ = roleAssignCmd.Flags().Set("host", "example-host")
    _ = roleAssignCmd.Flags().Set("user", "")
    roleAssignCmd.Run(roleAssignCmd, []string{})
    out := buf.String()
    assert.Contains(t, out, "Error: Hosts cannot have roles")
}
