package memory

import (
	"sync"
	"url-shortener/internal/storage"
)

// Реализация хранилища в памяти
type Storage struct {
	mu sync.RWMutex
	aliasUrl map[string]string
	urlAlias map[string]string
}

// Конструктор для Storage
func New() *Storage {
	return &Storage{
		aliasUrl: make(map[string]string),
		urlAlias: make(map[string]string),
	}
}

// Сохранение URL и его алиаса
func (st *Storage) SaveURL(urlToSave string, alias string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	// URL уже существует
	if _, ok := st.urlAlias[urlToSave]; ok {
		return storage.ErrURLExists
	}

	// Алиас уже занят другим URL
	if _, ok  := st.aliasUrl[alias]; ok {
		return storage.ErrURLExists
	}

	st.aliasUrl[alias] = urlToSave
	st.urlAlias[urlToSave] = alias

	return nil
}

// Получение URL по алиасу
func (st *Storage) GetURL(alias string) (string, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	_, ok := st.aliasUrl[alias]
	if !ok {
		return "", storage.ErrURLNotFound
	}

	return st.aliasUrl[alias], nil
}

func (st *Storage) DeleteURL(alias string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	url, ok := st.aliasUrl[alias]
	if !ok {
		return storage.ErrURLNotFound
	}
	
	delete(st.aliasUrl, alias)
	delete(st.urlAlias, url)
	return nil
}