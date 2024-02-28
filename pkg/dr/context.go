package dr

const KeyCredential = int8(0)
const KeyFilesystem = int8(1)

func NewCredential(username, password string) Credential {
	return Credential{
		username: username,
		password: password,
	}
}

// Credential the data-access credential stored in data transfer context
type Credential struct {
	username string
	password string
}

func (c Credential) Username() string {
	return c.username
}

func (c Credential) Password() string {
	return c.password
}
