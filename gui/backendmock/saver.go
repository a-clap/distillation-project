package backendmock

type Saver struct {
	values map[string][]byte
}

func NewSaver() *Saver {
	return &Saver{values: map[string][]byte{}}
}

func (l *Saver) Save(key string, data []byte) error {
	l.values[key] = data
	return nil
}

func (l *Saver) Load(key string) (data []byte) {
	d, ok := l.values[key]
	if !ok {
		return nil
	}
	return d
}
