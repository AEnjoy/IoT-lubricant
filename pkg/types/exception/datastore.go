package exception

type Operation interface {
	Do(*Exception) error
}
