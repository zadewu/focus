package domain

type FocusRepository interface {
	Init() error
	Create(name string) error
	Switch(name string) error
	Archive(name string) error
	List() ([]Focus, error)
	Current() (string, error)
	AddNote(msg string) error
	GetNotes(name string) ([]Note, error)
	Exists(name string) bool
	RemoteGet(name string) (string, error)
	RemoteSet(name, url string) error
	PushAll(remote string) error
	FetchAll(remote string) error
	CheckoutRemoteBranches(remote string) error
}

type ConfigStore interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

type WorkspaceStore interface {
	Path(name string) string
	Ensure(name string) (string, error)
	ListFiles(name string) ([]File, error)
}

type Exporter interface {
	Export(focus Focus, notes []Note, files []File) error
}
