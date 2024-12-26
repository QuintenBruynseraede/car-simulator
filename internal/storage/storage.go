package storage

type StorageBackend interface {
	ReadString(key string) string
	ReadFloat64(key string) float64
	ReadInt(key string) int
	Write(key string, value any)
	Dump() map[string]string
}
