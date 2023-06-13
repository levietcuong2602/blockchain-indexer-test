package mumbai

func (p *Platform) GetVersion() (string, error) {
	return p.client.GetVersion()
}
