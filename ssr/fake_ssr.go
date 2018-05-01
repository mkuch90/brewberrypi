package ssr

type FakeSSR struct {
	pin int
}

func (s *FakeSSR) On() error {
	return nil
}

func (s *FakeSSR) Off() error {
	return nil
}
