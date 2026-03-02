package random

import (
	"time"
	"math/rand"
)

func NewRandomString(size int) string {
    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// символы для генерации
    chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
        "abcdefghijklmnopqrstuvwxyz" +
        "0123456789_")

    b := make([]rune, size)	// буф символов
    for i := range b {
        b[i] = chars[rnd.Intn(len(chars))]
    }

    return string(b)
}