package db

type KeyValueStoreDB interface {
	Get(k string) (string, error)
	Set(k string, v string) error
}
