package send

type Error struct {
	Msg string `json:"msg"`
}

func NewErr(msg string) Error {
	return Error{Msg: msg}
}
