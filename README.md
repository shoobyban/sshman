# SSH Manager - Manage Access Through Authorized Key Files on Remote Hosts

[![Build Status](https://github.com/shoobyban/sshman/actions/workflows/push.yaml/badge.svg?branch=main)](https://github.com/shoobyban/sshman/actions/workflows/push.yaml)
[![Awesome GO](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoobyban/sshman)](https://goreportcard.com/report/github.com/shoobyban/sshman)

This is a simple tool I created to streamline the onboarding and offboarding of engineers across various environments, from AWS to third-party hosting providers.

As with all my creations, this tool solves _my_ problem. While it may not solve yours, I welcome feedback, fixes, pull requests, and issues.

**Caution**: Plan your group memberships carefully. Keep your management key out of any groups to avoid accidentally removing it from a host, which could lock you out.

## Installation

```sh
$ go get github.com/shoobyban/sshman
```

## How Does It Work?

This tool must be run from a host that can access all other hosts using a working SSH key that is not shared with anyone else. Configuration is saved in `~/.ssh/.sshman`. If you need to move the tool to another host, copy this file and the binary, and you're set up. The configuration does not contain any sensitive information.

There are two main resource entities in sshman: users and hosts. Users are identified by their public SSH key and labeled by their email address for simplicity. However, the email address is not used as an email and can be any identifier, such as `sam-key-1` or `sam-key-2`. This is useful when a user has multiple keys for different purposes, although it is not necessary in most cases.

![Users CRUD](docs/screenshot1.png)

The main concept of sshman is the group, which organizes users into "groups of hosts" or hosts into "groups of users." Examples include `live-hosts`, `staging-hosts`, `production`, or `{client1}`, `{client2}`. Groups act as tags; by tagging a user and a host with the same group name, the user gains access to the host.

To add a host to the sshman configuration, provide an alias, an SSH `.pub` key, and the groups the host belongs to (if already defined). Adding the host triggers an auto-discovery feature that downloads all SSH keys from the host as newly defined users and creates pseudo-groups for recognized users with access to that host.

![Adding a User](docs/screenshot2.png)

### Configuration File

Configuration is saved in `~/.ssh/.sshman`. It is a JSON file containing all hosts, users, and groups. "Configuration" might not be the best term for this file.

## Usage

### Adding Hosts

First, ensure you have hosts that you can already access with `~/.ssh/authorized_keys` files. Password authentication is not yet supported, but there are plans to add initial configuration through username and password.

To add a host, use the following syntax:

```bash
sshman add host --alias {alias} --address {host_address:port} --user {user} --key {~/.ssh/working_keyfile.pub} --groups group1,group2
```

Groups are optional and can be added later.

For example:

```bash
sshman add host --alias google --address my.google.com:22 --user myuser --key ~/.ssh/google --groups deploy,hosting,google
```

In this example, `google` is the alias. sshman accesses `my.google.com` on port 22 using the `myuser` user and the private SSH key located at `~/.ssh/google`. The host belongs to the `deploy`, `hosting`, and `google` groups. sshman saves these values in its configuration, accesses the host with the provided credentials, checks for the `~/.ssh/authorized_keys` file, downloads the users, and cross-references them with the current user list, adding new groups as necessary.

### Adding Users

Adding users is optional if all users are already on the hosts and you only need to manage them. Auto-discovery will automatically add users for you, but defining new users requires this step.

Syntax:

```bash
sshman add user --email {email} --key {sshkey.pub} --groups group1,group2
```

Groups are optional and can be added later.

For example:

```bash
sshman add user --email email@test.com --key ~/.ssh/user1.pub --groups production-team,staging-hosts
```

In this example, `email@test.com` is the label. It does not have to be an email address but is easier to identify and has secondary administrative value. The public key in `~/.ssh/user1.pub` is read into the configuration and can be discarded afterward if not used elsewhere. The user belongs to the `production-team` and `staging-hosts` groups. If there are hosts in these groups, the user's public SSH key is added to the `~/.ssh/authorized_keys` files of all relevant hosts.

### Auto-Discovery of Users on Added Hosts

To run auto-discovery of users on added hosts or refresh the configuration if any third party has changed `~/.ssh/authorized_keys` files, run:

```bash
sshman update
```

### Listing Who's on What Host

```bash
sshman list auth
```

This command displays a mapping of host aliases to email lists, making it easy to grep or add to reports.

### Listing What User and Host Belong to Which Group

For example:

```bash
sshman list groups
production-team hosts: [client1.live live2 host3 client1.uat]
production-team users: [email1@test.com email2@company.com]
dev-team hosts: [staging.test.com client1.staging]
dev-team users: [junior1@test.com email1@test.com email2@company.com]
```

Each group alias is listed with its associated hosts and users, making it easy to filter using `grep`.

### Listing Added Hosts

Lists host aliases, their addresses, and the groups they belong to.

```bash
sshman list hosts
client1.staging         staging.client1.com:22              [production-team dev-team]
client1.uat             uat.client1.com:22                  [production-team dev-team]
client1.live            www.client1.com:22                  [production-team]
```

### Listing Added Users with Groups

```bash
sshman list users
```

This command returns a mapping of email addresses to groups.

### Renaming Users and Hosts

Rename a user (modify email) or host (modify alias):

```bash
sshman rename user --email oldemail@host.com --new-email newemail@host.com

sshman rename host --alias oldalias --new-alias newalias
```

### Modifying User and Host Grouping

Modify a user's groups or remove groups to allow global access:

```bash
sshman groups user --email email@host.com --add group1 --remove group2
```

Modify a host's groups or remove it from all groups:

```bash
sshman groups host --alias hostalias --add group1 --remove group2
```

**Note:** Removing a host from a group removes all users who are on the host only because of that group. If the host is in another group, users in both groups are not removed.

### Roles Management

#### Assign a Role
Assign a role to a user:

```bash
sshman roles assign --user <user_email> --role <role_name>
```

- `--user`: Specify the email of the user to assign the role to.
- `--role`: Specify the role to assign.

**Note:** Hosts cannot have roles. Use groups for managing host permissions.

#### List Roles
List all roles and their associated permissions:

```bash
sshman roles list
```

This command displays all roles and the permissions associated with each role.

### Consistent Argument Structure

All commands now use flags for specifying arguments, ensuring clarity and consistency. For example:
- Use `--user` to specify a user email.
- Use `--host` to specify a host alias.
- Use `--role` to specify a role name.

### Example Usage

#### Adding a User

```bash
sshman add user --email user@example.com --key ~/.ssh/id_rsa.pub --groups group1,group2
```

#### Adding a Host

```bash
sshman add host --alias myhost --address myhost.com:22 --user myuser --key ~/.ssh/host_key.pub --groups group1,group2
```

#### Modifying Groups for a User

```bash
sshman groups user --email user@example.com --add group1 --remove group2
```

#### Modifying Groups for a Host

```bash
sshman groups host --alias myhost --add group1 --remove group2
```

### Things To Fix Before Release

- [x] Fix adding users
- [x] Bug: Adding host on frontend does not add keyfile entry into storage, edit afterwards does
- [x] Bug: Renaming host (alias) created a new entry, did not delete old
- [x] Group editing
  - [x] Add group should add users and groups
  - [x] Update group should remove / add resources
  - [x] Delete group should remove resources
- [x] Test all CRUD (users, hosts, groups) together
- [x] Re-read config with file watcher in web mode
- [x] Screenshot with test data (not with sensitive data)
- [x] Reuse stored ssh key for modifying user
- [x] Adding host to download information without the need of running update
- [x] Complete CRUD for missing use cases
- [x] Web interface
- [ ] Full test coverage
- [x] All user related functions should have unit tests
- [ ] All user role related functions should have unit tests
- [x] All host related functions should have unit tests
- [ ] All host group related functions should have unit tests
- [ ] All configuration handling functions should have unit tests
- [ ] All other core functionality should have unit tests
- [ ] Edge case: deleting user should delete the user from all hosts (unless canceled from changeset)
- [ ] Misfeature: Changing keyfile on host does not upload new key with old and delete old
- [ ] Misfeature: Adding host does not check if host config is working
- [ ] Misfeature: Adding host with groups does not upload initial users from group
- [ ] Misfeature: Modifying user groups does not upload / delete hosts

### TODO For Next Release

- [ ] Web authentication
- [ ] Delete host with editing ssh keys
- [ ] Auto-group host specific users (when user is on several hosts, create a group for them, auto-merge groups when possible)
- [ ] CLI to use API (not sure)
- [ ] Web Interface Authentication (where to store creds?)
- [ ] Updated At timestamps
- [ ] Audit log
  - [ ] audit log logging all changes from changeset (sync op) on apply
- [ ] Implement user "role" group of groups for RBAC level of abstraction (developers role = uat-servers+staging-servers group)
- [ ] Testing connection after creating authorized_keys entry

### (Possible) Future Plans

- [ ] Changeset based operation (see [Future plans details](docs/Plans.md))
- [ ] Web Aria tags (at least tagging buttons better and connecting labels)
- [ ] More backend (currently `.ssh/.sshman` JSON configuration file)
- [ ] Adding host key to server using password auth
- [ ] Text UI based on Web frontend
- [ ] State handling (see [Future plans details](docs/Plans.md))
- [ ] Edit multiple items (see [Future plans details](docs/Plans.md))

## Credits

Most of the credit goes to the pain of being a CTO for 17+ years in small and mid-sized companies, where SSH key management is not solved.

The project would have been much harder without the work of [Steve Francia](https://github.com/spf13) and all the cobra and viper contributors, the web UI relies on [Chi](https://github.com/go-chi/chi) and [Vue](https://github.com/vuejs/).

Web UI embedding wouldn't be working without [Gregor Best](https://github.com/farhaven), who nerd-sniped me into helping with a tricky bug on Gophers Slack.

I love the Go community.
