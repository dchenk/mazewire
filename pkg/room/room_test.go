package room

type TestDataStore struct{}

func (TestDataStore) Store([]byte) (int64, error) {
	return 0, nil
}

func (TestDataStore) GetMulti([]int64) ([]Datum, error) {
	return nil, nil
}

func (TestDataStore) Get(int64) (*Datum, error) {
	return nil, nil
}
