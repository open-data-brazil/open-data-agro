package storage

// NewLocalBronzeStoreForTest exposes the local bronze store for ingest integration tests.
func NewLocalBronzeStoreForTest(root string) BronzeStore {
	return newLocalBronzeStore(root)
}
