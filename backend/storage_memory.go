package backend

type MemoryStorage struct {
	data *Data
}

func (m *MemoryStorage) Load(data *Data) error {
	m.data = data
	return nil
}

func (m *MemoryStorage) Save(data *Data) error {
	m.data = data
	return nil
}

func (m *MemoryStorage) Watch(notify func()) {
}
