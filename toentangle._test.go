package toentangle

import (
	"fmt"
	"os"
	"testing"
)

func TestToEntangle_Get(t *testing.T) {
	// 创建一个 ToEntangle
	toentangle := NewToEntangle("data")

	// 添加一对数据, 使其双向绑定
	toentangle.Add("a", "b")
	toentangle.Add("a", "c")

	// 获取 a 的全部绑定数据
	arr, _ := toentangle.GetA("a")
	fmt.Println(arr)
	if len(arr) != 2 {
		t.Errorf("GetA(\"a\") = %v; want [\"b\", \"c\"]", arr)
	}

	// 获取 b 的全部绑定数据
	arr, _ = toentangle.GetB("b")
	fmt.Println(arr)
	if len(arr) != 1 || arr[0] != "a" {
		t.Errorf("GetB(\"b\") = %v; want [\"a\"]", arr)
	}

	// 获取 c 的全部绑定数据
	arr, _ = toentangle.GetB("c")
	fmt.Println(arr)
	if len(arr) != 1 || arr[0] != "a" {
		t.Errorf("GetB(\"c\") = %v; want [\"a\"]", arr)
	}

	// 移除所有绑定
	toentangle.Remove("a", "b")
	toentangle.Remove("a", "c")

	// 获取 a 的全部绑定数据
	arr, _ = toentangle.GetA("a")
	fmt.Println(arr)
	if len(arr) != 0 {
		t.Errorf("GetA(\"a\") = %v; want []", arr)
	}

	// 获取 s 的全部绑定数据
	arr, _ = toentangle.GetB("s")
	fmt.Println(arr)
	if len(arr) != 0 {
		t.Errorf("GetB(\"s\") = %v; want []", arr)
	}

	// 清理 leveldb
	toentangle.Close()

	// 删除 data 文件夹
	os.RemoveAll("data")
}
