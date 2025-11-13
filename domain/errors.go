package domain

// DomainError represents domain-specific errors with HTTP status mapping
type DomainError interface {
	error
	HTTPStatus() int
}
