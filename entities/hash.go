package entities

type Hash map[string]string

func (hash *Hash) Copy() Hash {
	result := make(Hash, len(*hash))

	for key, value := range *hash {
		result[key] = value
	}

	return result
}
