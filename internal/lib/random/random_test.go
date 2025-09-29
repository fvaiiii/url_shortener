package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	cases := []struct {
		name string
		size int
	}{
		{
			name: "s1",
			size: 1,
		},
		{
			name: "s5",
			size: 5,
		},
		{
			name: "s10",
			size: 10,
		},
		{
			name: "s15",
			size: 15,
		},
		{
			name: "s20",
			size: 20,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			str1 := NewRandomString(c.size)
			str2 := NewRandomString(c.size)

			assert.Len(t, str1, c.size)
			assert.Len(t, str2, c.size)

			assert.NotEqual(t, str1, str2)
		})
	}
}
