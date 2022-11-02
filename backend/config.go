package backend

type Config interface {
	AddHost(host *Host, withUsers bool) error
	AddUser(newuser *User, host string) error
	AddUserByEmail(email string) bool
	AddUserToHosts(newuser *User)
	DeleteGroup(label string) bool
	DeleteHost(alias string) bool
	DeleteUser(email string) bool
	DeleteUserByID(id string) bool
	FromGroup(host *Host, email string) bool
	GetGroups() map[string]LabelGroup
	GetHost(alias string) *Host
	GetHosts(group string) []*Host
	GetUser(lsum string) *User
	GetUserByEmail(email string) (string, *User)
	GetUserByKey(lsum string) *User
	GetUsers(group string) []*User
	Hosts() map[string]*Host
	Log() *ILog
	PrepareHost(args ...string) (*Host, error)
	PrepareUser(email, filename string, groups ...string) (*User, error)
	Regenerate(aliases ...string)
	RemoveUserFromHosts(deluser *User) error
	SetHost(alias string, host *Host)
	StopUpdate()
	Update(aliases ...string)
	UpdateGroup(groupLabel string, users, hosts []string)
	UpdateHost(host *Host) error
	UpdateUser(newuser *User) error
	UserExists(lsum string) bool
	Users() map[string]*User
	WatchFile(notify func())
	Write()
}
