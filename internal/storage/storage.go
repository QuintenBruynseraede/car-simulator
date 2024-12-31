package storage

type StorageBackend interface {
	ReadString(key string) (string, error)
	ReadFloat64(key string) (float64, error)
	ReadInt(key string) (int, error)
	Write(key string, value any)
	Dump() map[string]string
}
