package backend

type Config interface {
	GetUserByEmail(email string) (string, *User)
	Write()
	GetUsers(group string) []*User
	AddUserToHosts(newuser *User)
	AddHost(host *Host, withUsers bool) error
	SetHost(alias string, host *Host)
	Hosts() map[string]*Host
	Users() map[string]*User
	UserExists(lsum string) bool
	GetUserByKey(lsum string) *User
	RemoveUserFromHosts(deluser *User) error
	PrepareHost(args ...string) (*Host, error)
	DeleteUserByID(id string) bool
	DeleteUser(email string) bool
	PrepareUser(email, filename string, groups ...string) (*User, error)
	AddUser(newuser *User, host string) error
	UpdateUser(newuser *User) error
	DeleteHost(alias string) bool
	Update(aliases ...string)
	UpdateHost(host *Host) error
	Regenerate(aliases ...string)
	GetGroups() map[string]LabelGroup
	AddUserByEmail(email string) bool
	GetUser(lsum string) *User
	GetHost(alias string) *Host
	DeleteGroup(label string) bool
	UpdateGroup(groupLabel string, users, hosts []string)
	FromGroup(host *Host, email string) bool
	StopUpdate()
	WatchFile(notify func())
	Log() *ILog
}
