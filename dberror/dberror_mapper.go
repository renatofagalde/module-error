package dberror

type DBErrorMapper interface {
	Map(err error) error
}
