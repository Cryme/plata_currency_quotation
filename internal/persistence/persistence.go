package persistence

var Instance Interface

type CommonPersistenceOperations interface {
	OnStart() error
}

type Interface interface {
	CommonPersistenceOperations
	QuotationRequestPersistentOperations
}
