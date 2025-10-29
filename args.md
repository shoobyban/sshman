# Command Arguments Reorganization Plan

This document outlines a comprehensive plan to reorganize the command arguments structure and behavior of `sshman` to align with industry standards and best practices. The goal is to create a more intuitive, consistent, and user-friendly command-line interface.

## 1. Analysis of the Current Structure

The current command structure, built with `cobra`, has several inconsistencies and areas for improvement:

- **Inconsistent Command Naming:** Commands are a mix of verbs (`add`, `remove`, `rename`) and nouns (`groups`, `web`), leading to a confusing user experience.
- **Inconsistent Argument Handling:** Some commands use positional arguments (`add user`), while others use flags (`assign`). This lack of uniformity makes the CLI difficult to learn and use.
- **Lack of Standardization:** Argument and flag names are not consistent across commands (e.g., `user` vs. `email`).
- **Redundancy:** Commands like `add`, `addHost`, and `addUser` create unnecessary nesting. The parent `add` command is redundant. The same applies to `remove` and `rename`.
- **Confusing `groups` Command:** The `groups` command is used to modify both user and host groups, which is not intuitive.
- **Inconsistent Naming:** The `del` command in `removeUser.go` is inconsistent with the `remove` command.
- **Misleading `update` Command:** The `update` command fetches data, so a name like `fetch` or `sync` would be more descriptive.
- **No Global Flags:** There are no global flags for common options like `--config` or `--verbose`.
 - **Global Flags Present:** The code defines persistent/global flags: `--config`/`-c` and `--verbose`/`-v` (see `cmd/root.go`).
- **Inconsistent Help Messages:** The help messages could be more detailed and provide better examples.

## 2. Proposed Command Structure

To address these issues, we propose a new command structure based on the following principles:

- **Resource-Oriented:** Commands should be organized around resources (e.g., `user`, `host`, `group`).
- **Action-Oriented:** Actions on resources should be subcommands (e.g., `user add`, `host remove`).
- **Consistent Naming:** Use consistent and predictable names for commands, arguments, and flags.
- **Standard Flags:** Use standard flags for common options.

Here is the proposed new command structure:

```
sshman
├── user
│   ├── add <email> <sshkey.pub> [flags]
│   ├── remove <email>
│   ├── rename <old_email> <new_email>
│   ├── list
│   └── groups <email> [groups...]
├── host
│   ├── add <alias> <host:port> <user> <keyfile> [flags]
│   ├── remove <alias>
│   ├── rename <old_alias> <new_alias>
│   ├── list
│   └── groups <alias> [groups...]
├── group
│   └── list
├── role
│   ├── assign --user <email> --role <role>
│   └── list
├── sync
├── tree (show hierarchical resource tree)
├── web
└── version
```

### 2.1. Detailed Command Breakdown

#### `sshman user`

- **`sshman user add <email> <sshkey.pub> [flags]`**: Add a new user.
  - `--group <group>`: (Repeatable) Add the user to one or more groups.
- **`sshman user remove <email>`**: Remove a user.
- **`sshman user rename <old_email> <new_email>`**: Rename a user.
- **`sshman user list`**: List all users.
- **`sshman user groups <email> [groups...]`**: Set the groups for a user. If no groups are provided, the user is removed from all groups.

#### `sshman host`

- **`sshman host add <alias> <host:port> <user> <keyfile> [flags]`**: Add a new host.
  - `--group <group>`: (Repeatable) Add the host to one or more groups.
- **`sshman host remove <alias>`**: Remove a host.
- **`sshman host rename <old_alias> <new_alias>`**: Rename a host.
- **`sshman host list`**: List all hosts.
- **`sshman host groups <alias> [groups...]`**: Set the groups for a host. If no groups are provided, the host is removed from all groups.

#### `sshman group`

- **`sshman group list`**: List all groups.

#### `sshman role`

- **`sshman role assign --user <email> --role <role>`**: Assign a role to a user.
- **`sshman role list`**: List all roles.

Note: Roles are only supported for users — hosts do not have roles.

#### `sshman sync`

- **`sshman sync`**: Fetches all users from `authorized_keys` on all hosts and updates the local configuration. This replaces the old `update` command.

#### `sshman web`

- **`sshman web [flags]`**: Starts the web UI.
  - `--bind <ip>`: Bind to a specific IP address.
  - `--port <port>`: Port for the web UI.
  - `--portfile <file>`: Port filename for dynamic address.

#### `sshman version`

- **`sshman version`**: Prints the version of `sshman`.

## 2.2. Global Flags

- `--config <file>` (or `-c`): Path to the configuration file.
- `--verbose` (or `-v`): Enable verbose output.

## Reality in the codebase (short summary)

The repository largely implements the proposed, resource-oriented CLI, with a few noteworthy differences and duplications that are present in the codebase and should be documented here:

- The code exposes a single `role` command tree implemented in `cmd/role.go` (providing `assign` and `list`). Duplicate `roles` registration has been consolidated in the codebase.
- `sshman sync` is implemented and is the recommended command to fetch users from hosts. There is also an older `update` command (registered as `update` in `cmd/update.go`) which prints a deprecation notice and calls the same update flow; the code intentionally keeps it for backward compatibility.
- Several legacy top-level commands still exist in the tree (for example: `add`, `remove`, `rename`, `read/update`, `groups`, `list`, and `del`) and `cmd/root.go` explicitly marks many of these as Deprecated (it sets their `Deprecated` field when present).
- The `web` command exists and exposes short flags in the implementation: `--bind`/`-b`, `--port`/`-p`, and `--portfile`/`-f`. The code accepts `dynamic` for the `--port` value, which will allocate an ephemeral port and write it to the provided `--portfile`.
- The `role assign` implementation accepts `--user` and `--host` flags in addition to `--role`. The backend treats roles as user-scoped: assigning a role to a host will print an error (hosts are expected to be controlled via groups instead).

These notes reflect the current state of the code and can be used as the basis for either updating the implementation to match the original proposal or updating the proposal to match the implementation.

## 3. Migration Plan

To migrate from the old command structure to the new one, we will follow these steps:

1. **Implement the New Command Structure:** Create new `cobra` commands for the proposed structure.
2. **Deprecate Old Commands:** Mark the old commands as deprecated and print a warning message pointing to the new command.
3. **Update Documentation:** Update all documentation, including the `README.md` and help messages, to reflect the new command structure.
4. **Remove Old Commands:** After a few releases, remove the old commands completely.

This phased approach will ensure a smooth transition for existing users while allowing us to move to a more modern and maintainable command structure.
