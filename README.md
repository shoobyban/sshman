# SSH Manager - manage authorized_key file on remote servers

This is a simple tool that I came up after having to on-boarding and off-boarding developers on a
very colourful palette of environments from AWS to 3rd party hosting providers.

As every one of my creations this tool is solving _my_ problem. It does not warranty your problem will be solved,
but in that highly unlikely event please let me know, fixes and pull requests, issues are all very welcome without
again the promise that I'll do anything, I'm normally really busy, apologies.

## Installation

```sh
$ go get github.com/shoobyban/sshman
```

## How does it work?

First of all, from where you will run this tool, you need to be able to access to the server, on a port, 
with a working ssh key (that you don't want to share with anybody else).
First, think about your groups (if you need this feature), limiting users into group of servers, like `live-servers`, `staging-servers`, `production` etc.
This is optional, and any time you can re-register the user with new groups (as long as you have their public key file, note to myself I have that info in the system, small todo).
You register the server into the registry with an alias (and the groups where the server belongs), if you have user ssh `.pub` keys (this is optional) register users with their key file and email address (optionally with the user's groups).
After having a few servers defined (and optionally users) you can run auto discovery.

Configuration will be saved into `~/.ssh/.ssmman`, if you need to move tool to any other server, copy this and the binary and you are set up. Configuration will not have any secure information.

## Usage

### Registering Servers
First, you need servers, that you can already access, with `~/.ssh/authorized_keys` files on the server. Password auth doesn't count.

To register a server, the syntax is 

```sshman register server {alias} {server_address:port} {user} {~/.ssh/working_keyfile.pub} [group1 group2 ...]```

Where groups are optional, it helps when you have several user roles or you want to limit users to certain servers.

Registering a server for example:

```sh
$ sshman register server google my.google.com:22 myuser ~/.ssh/google.pub deploy hosting google
```

`google` will be my alias, it will access `my.google.com` on port 22, with `myuser` user using `~/.ssh/google.pub` from the current user.

### Registering Users

This is optional if you already have all the users on the servers and you just want to be able to move them around or delete them, auto discovery will auto-register the users for you, but adding new users will require this step.

Syntax is 

```sshman register user {email} {sshkey.pub} [group1 group2 ...]```

For example:

```sh
$ sshman register user email@test.com ~/.ssh/user1.pub production-team staging-servers
```

### Auto Discovery users on registered servers

To run auto discovery users on registered servers, or to refresh the configuration if any 3rd party has changed `~/.ssh/authorized_keys` files, run:

```sh
$ sshman update
```

### Adding user to server

After registering user with email, key file and groups, uploading the user to the servers that the user can access:

```sh
$ sshman add email@test.com
```

This command will add user's key to all `~/.ssh/authorized_keys` files on the servers that groups allow. 

**If there is no group information for the user, you will give access to all servers.**

### Deleting user from servers

Any existing user can be deleted from all `~/.ssh/authorized_keys` files from all servers by running 

```sh
$ sshman add email@test.com
```

This will remove the entries from the servers but keep user information in configuration for further modification.

### Listing who's on what server

```sh
$ sshman list auth
```

This will display server alias -> email list mapping, easy to grep or add to reports.

### Listing what user and server is in what group

Easier to explain this with an example scenario:

```sh
$ sshman list groups
production-team servers: [client1.live live2 server3 client1.uat]
production-team users: [email1@test.com email2@company.com]
dev-team servers: [staging.test.com client1.staging]
dev-team users: [junior1@test.com email1@test.com email2@company.com]
```

Notice that group alias is in every line with "servers" and "users" for using `grep` on the list.

### Listing registered servers

Lists server aliases, what server/port, server is in what groups.

```sh 
$ sshman list servers
client1.staging        	staging.client1.com:22              [production-team dev-team]
client1.uat        	    uat.client1.com:22               	[production-team dev-team]
client1.live        	www.client1.com:22               	[production-team]
```

### Listing registered users with groups

```sh
$ sshman list users
```

Will return a mapping of email to groups.

