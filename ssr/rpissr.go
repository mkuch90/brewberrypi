package ssr

type RpiSSR struct {
	pin int
}

func (s *RpiSSR) On() error {
	return nil
}

func (s *RpiSSR) Off() error {
	return nil
}
