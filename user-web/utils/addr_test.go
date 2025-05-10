package utils

import (
	"testing"
)

func TestGetFreePortFunction(t *testing.T) {
	port, err := GetFreePort()
	if err != nil {
		t.Errorf("获取可用端口失败: %v", err)
	}
	if port <= 0 {
		t.Errorf("获取到的端口号无效: %d", port)
	}
	t.Logf("获取到可用端口: %d", port)
}
