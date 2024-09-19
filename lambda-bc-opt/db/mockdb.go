package db


type MockRedisDB struct {}

// KeyValueStoreDB
func (rdb *MockRedisDB) Get(k string) (string, error) {
	result := "0"
	return result, nil
}
func (rdb *MockRedisDB) Set(k string, v string) error {
	return nil
}

func ConsMockRedisDB() *MockRedisDB {
	return &MockRedisDB{}
}
