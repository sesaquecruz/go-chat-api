package chat

type ChatError string

func (e ChatError) Error() string {
	return string(e)
}
