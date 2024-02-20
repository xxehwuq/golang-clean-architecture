package random

import gonanoid "github.com/matoous/go-nanoid/v2"

func ID() (string, error) {
	return gonanoid.New()
}
