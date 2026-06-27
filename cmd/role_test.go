package cmd

import (
    "bytes"
    "strings"
    "testing"
    "os"
    "io"

    "github.com/shoobyban/sshman/backend"
    "github.com/stretchr/testify/assert"
)

// resetRoleAssignFlags clears the flags on roleAssignCmd to avoid test interference
func resetRoleAssignFlags() {
    _ = roleAssignCmd.Flags().Set("user", "")
    _ = roleAssignCmd.Flags().Set("host", "")
    _ = roleAssignCmd.Flags().Set("role", "")
}

// captureStdout captures anything written to stdout while fn runs and returns it.
func captureStdout(fn func()) string {
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    outC := make(chan string)
    go func() {
        var buf bytes.Buffer
        io.Copy(&buf, r)
        outC <- buf.String()
    }()
    fn()
    w.Close()
    os.Stdout = old
    return <-outC
}

func TestRoleListCmdOutputsRoles(t *testing.T) {
    // Ensure default roles are initialized
    _ = backend.DefaultConfig()

    out := captureStdout(func() {
        roleListCmd.Run(roleListCmd, []string{})
    })
    if !strings.Contains(out, "Role: admin") {
        t.Fatalf("expected output to contain admin role, got: %s", out)
    }
}

func TestRoleAssignInvalidRolePrintsError(t *testing.T) {
    resetRoleAssignFlags()
    out := captureStdout(func() {
        // set role flag to nonexistent
        _ = roleAssignCmd.Flags().Set("role", "no-such-role")
        // ensure no user/host to hit the missing role branch
        _ = roleAssignCmd.Flags().Set("user", "")
        _ = roleAssignCmd.Flags().Set("host", "")
        roleAssignCmd.Run(roleAssignCmd, []string{})
    })
    assert.Contains(t, out, "Role no-such-role does not exist")
    resetRoleAssignFlags()
}

func TestRoleAssignHostRejected(t *testing.T) {
    resetRoleAssignFlags()
    out := captureStdout(func() {
        // use existing role
        _ = roleAssignCmd.Flags().Set("role", "admin")
        // host supplied without user should trigger host error branch
        _ = roleAssignCmd.Flags().Set("host", "example-host")
        _ = roleAssignCmd.Flags().Set("user", "")
        roleAssignCmd.Run(roleAssignCmd, []string{})
    })
    assert.Contains(t, out, "Hosts cannot have roles")
    resetRoleAssignFlags()
}

func TestRoleAssignNonexistentUser(t *testing.T) {
    resetRoleAssignFlags()
    out := captureStdout(func() {
        // set role to existing role
        _ = roleAssignCmd.Flags().Set("role", "admin")
        // set user to a non-existent user
        _ = roleAssignCmd.Flags().Set("user", "noone@example.com")
        _ = roleAssignCmd.Flags().Set("host", "")
        roleAssignCmd.Run(roleAssignCmd, []string{})
    })
    assert.Contains(t, out, "User noone@example.com not found")
    resetRoleAssignFlags()
}
