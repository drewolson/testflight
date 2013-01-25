package ws

type TooShortError struct {
}

type TimeoutError struct {
}

func (e TimeoutError) Error() string {
	return "No Message Received in 30 seconds"
}

func (e TooShortError) Error() string {
	return "Unread Messages too Short"
}
