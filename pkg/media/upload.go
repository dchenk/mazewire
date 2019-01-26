package media

import (
	"encoding/hex"

	"github.com/satori/go.uuid"
)

// genID generates a 32-character hex string from a UUID.
func genID() (string, error) {
	u, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(u[:]), nil
}
