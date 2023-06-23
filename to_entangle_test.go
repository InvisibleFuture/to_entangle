package toentangle

import (
	"fmt"
	"os"
	"testing"
)

func TestToEntangle_Get(t *testing.T) {
	// 创建一个 ToEntangle
	entangle := New("data")

	// 添加一对数据, 使其双向绑定
	entangle.Add("a", "b")
	entangle.Add("a", "c")

	// 获取 a 的全部绑定数据
	arr, _ := entangle.Get("a")
	fmt.Println(arr)
	if len(arr) != 2 || arr[0] != "b" || arr[1] != "c" {
		t.Errorf("Get(\"a\") = %v; want [\"b\", \"c\"]", arr)
	}

	// 获取 b 的全部绑定数据
	arr, _ = entangle.Get("b")
	fmt.Println(arr)
	if len(arr) != 1 || arr[0] != "a" {
		t.Errorf("Get(\"b\") = %v; want [\"a\"]", arr)
	}

	// 获取 c 的全部绑定数据
	arr, _ = entangle.Get("c")
	fmt.Println(arr)
	if len(arr) != 1 || arr[0] != "a" {
		t.Errorf("Get(\"c\") = %v; want [\"a\"]", arr)
	}

	// 清理 leveldb
	entangle.Close()

	// 删除 data 文件夹
	os.RemoveAll("data")
}
