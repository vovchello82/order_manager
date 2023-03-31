package publisher

type Publisher interface {
	EmitObject(payload string) error
}
