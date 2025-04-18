// Code generated by godddx, DO AVOID EDIT.
package uniqueid

// Storer data persistence
type Storer interface {
	UniqueID() UniqueIDStorer
}

// Core business domain
type Core struct {
	store Storer
	m     *IDManager
}

// NewCore create business domain
func NewCore(store Storer, length int) Core {
	return Core{
		store: store,
		m:     NewIDManager(store.UniqueID(), length),
	}
}

// UniqueID 获取全局唯一 ID
func (c Core) UniqueID(prefix string) string {
	return c.m.UniqueID(prefix)
}
