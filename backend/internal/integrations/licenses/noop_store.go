package licenses

type noopStore struct{}

func (n *noopStore) Ping() error { return nil }
