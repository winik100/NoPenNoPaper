package core

type User struct {
	ID             int
	Name           string
	HashedPassword string
	Role           string
	Materials      Materials //filenames for uploads associated with that user
}

type Materials struct {
	Title    []string
	FileName []string
}
