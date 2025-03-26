package utils

type StringErr string

func (e StringErr) Error() string {
	return string(e)
}
