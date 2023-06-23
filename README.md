# to_entangle
由go语言实现数据的双向绑定


`go get github.com/InvisibleFuture/to_entangle`


### 双表映射

```go
package main

import (
    toentangle "github.com/InvisibleFuture/to_entangle"
)

// 创建一个 Entangle
entangle := NewToEntangle("data/test")

// 添加一对数据, 使其双向绑定
entangle.Add("a", "b")
entangle.Add("a", "c")

// 获取 a 的全部绑定数据
arr, _ := entangle.Get("a")
fmt.Println(arr)

// 获取 b 的全部绑定数据
arr, _ = entangle.Get("b")
fmt.Println(arr)

// 获取 c 的全部绑定数据
arr, _ = entangle.Get("c")
fmt.Println(arr)

// 移除所有绑定
entangle.Remove("a", "b")
entangle.Remove("a", "c")

// 获取 a 的全部绑定数据
arr, _ = entangle.Get("a")
fmt.Println(arr)

// 获取 s 的全部绑定数据
arr, _ = entangle.Get("s")
fmt.Println(arr)

```


### 单表映射

```go
package main

import (
    toentangle "github.com/InvisibleFuture/to_entangle"
)

// 创建一个 Entangle
entangle := NewEntangle("data/test")

// 添加一对数据, 使其双向绑定
entangle.Add("a", "b")
entangle.Add("a", "c")

// 获取 a 的全部绑定数据
arr, _ := entangle.Get("a")
fmt.Println(arr)

// 获取 b 的全部绑定数据
arr, _ = entangle.Get("b")
fmt.Println(arr)

// 获取 c 的全部绑定数据
arr, _ = entangle.Get("c")
fmt.Println(arr)

// 移除所有绑定
entangle.Remove("a", "b")
entangle.Remove("a", "c")

// 获取 a 的全部绑定数据
arr, _ = entangle.Get("a")
fmt.Println(arr)

// 获取 s 的全部绑定数据
arr, _ = entangle.Get("s")
fmt.Println(arr)

```
