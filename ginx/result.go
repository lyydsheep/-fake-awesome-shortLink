package ginx

type Result[T any] struct {
	Msg  string
	Code uint8
	Data T
}
