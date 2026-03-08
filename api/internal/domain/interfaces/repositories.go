//go:generate moq -out repositories_mock.go -pkg interfaces . Repositories
package interfaces

type Repositories interface {
	PostQuerier
	PostCommander
	Paginator
}
