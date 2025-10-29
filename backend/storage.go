package backend

type Storage interface {
	Load(data *Data) error
	Save(data *Data) error
	Watch(notify func())
}
