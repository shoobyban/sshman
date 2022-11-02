# Tests

## Configuration / state tests
- [x] Opening config: if there is a configuration file, read it
- [x] Creating config: if no configuration file, create an empty one
- [x] Saving config: write the current configuration into a file

## Host tests
- [x] Adding host: Create a host entry with given arguments, connect to host and download and create users
- [x] Removing host: Delete host entry from configuration, not touching server
- [ ] Modifying host: If group has been changed, add and remove matching users

## User tests
- [ ] Adding user:
  - [x] Adding user to configuration
  - [ ] Uploading user to all matching groups
- [ ] Deleting user: Removing user from configuration, deleting entries from hosts
- [ ] Modifying user: If group has been changed, sync with hosts

## Group tests
- [ ] Adding group: Create group entry with given arguments, add group to given hosts and users
- [ ] Removing group: Delete group entry from configuration, remove it from hosts and users, not touching servers
- [ ] Modifying group:
    - [ ] If users have changed, add and remove matching users, sync with hosts
    - [ ] If hosts have changed, add and remove matching hosts, sync with hosts

## Core functionality tests
- [ ] Download users: Download users from host and add them to configuration, if user is not in group add them as host specific user
- [ ] Sync with host: Download users from host and upload new users to host
