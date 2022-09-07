package types

type ChainType string

const (
	BSC = ChainType("bsc")
)

//nolint:gochecknoglobals
var ChainTypes = map[ChainType]bool{
	BSC: true,
}

func IsChainSupported(chain ChainType) bool {
	_, exists := ChainTypes[chain]

	return exists
}
