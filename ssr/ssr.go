package ssr

type SolidStateRelay interface {
	On() error
	Off() error
}
