package tests

import (
	"crawler/pkg/random"
	"testing"
)

func TestUserAgent(t *testing.T) {
	t.Run("测试UserAgent随机性", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Log(random.GetRandomUserAgent())
		}
	})

}
