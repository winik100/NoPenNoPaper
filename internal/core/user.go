package core

type User struct {
	ID             int
	Name           string
	HashedPassword string
	Role           string
	Materials      Materials //filenames for uploads associated with that user
}

const RoleAnon = "anonymous"
const RolePlayer = "player"
const RoleGM = "gm"

func (u User) IsGM() bool {
	return u.Role == RoleGM
}

type Materials struct {
	Title    []string
	FileName []string
}
