package ethereum

import "github.com/unanoc/blockchain-indexer/pkg/node/version"

func (p *Platform) GetVersion() (string, error) {
	v, err := p.client.GetVersion()
	if err != nil {
		return "", err
	}

	return version.NewParser().Parse(v), nil
}
