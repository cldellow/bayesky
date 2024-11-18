package source

type Source interface {
	Next() (string, error)
}
