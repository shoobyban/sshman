# Internal logic of sshman
This is a top-level document describing the entities and algorithms behind operations.

Current logic has a few constraints:
- every user (email) can only have one key
- users are directly added to groups, no roles are available (like in RBAC)

## Entities

- [x] User: email+public key
- [x]  Host: hostname(or IP)+port+management key
- [x] Group: list of hosts and users groupped together

## Operations

### Adding host
- [x] in: host name, port, label, private key
- [x] reading user list
- [x] creating new users
- [x] adding initial users as protected
- [ ] if private key != main key, upload main key as protected
- [ ] prompting:
    - [ ] for protected users (min 1)
    - [ ] list of new users (no groups)

### Adding user (email)
- [ ] check if email already exists, if yes, prompt to modify user info
- [ ] check if key already exists, if yes, prompt to modify user info
- [ ] create user entry
- [ ] if user existed with different key update group hosts with new key

### Assigning host to group
- [ ] Adding users from group to host

### Assigning user to group
- [ ] adding user to all hosts in group

### Removing host (not changing authorized_keys)
- [ ] removing host from all groups

### Removing user
- [ ] removing user from all hosts, skip if user is protected on that host

### Removing group (not changing authorized_keys)
- [ ] prompt if removing users from host, if true
    - [ ] skip if user is protected on that host
    - [ ] remove group's users from host
- [ ] Removing group label from hosts
- [ ] Removing group label from users
