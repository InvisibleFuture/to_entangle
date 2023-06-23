package toentangle

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
	leveldb "github.com/syndtr/goleveldb/leveldb"
)

// 1. 存入一对数据, 使其双向绑定
// 2. 移除一对数据, 使其双向解绑
// 3. 移除一个数据, 使其全部解绑
// 4. 获取一个数据的全部绑定数据

type ToEntangle struct {
	lock sync.Mutex
	db   *leveldb.DB
}

// Get 获取一个数据的全部绑定数据
func (t *ToEntangle) Get(a string) (arr []string, err error) {
	// 获取 a 的数据, 如果不存在, 则直接返回空数组
	data, err := t.db.Get([]byte(a), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return []string{}, nil
		}
		return nil, err
	}

	// data 为json字符串的数组, 解码后返回
	err = jsoniter.Unmarshal(data, &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// Add 添加一对数据, 使其双向绑定
func (t *ToEntangle) Add(a string, b string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	var set_data = func(a string, b string) (err error) {
		// 获取 a 的数据
		data, err := t.db.Get([]byte(a), nil)
		if err != nil {
			// 如果不存在, 则创建并直接添加 b
			if err == leveldb.ErrNotFound {
				data, _ = jsoniter.Marshal([]string{b})
				return t.db.Put([]byte(a), data, nil)
			}
			return err
		}

		// data 为json字符串的数组, 解码后添加 b
		var arr []string
		err = jsoniter.Unmarshal(data, &arr)
		if err != nil {
			return err
		}

		// 转换为 map, 添加 b, 以去重
		m := make(map[string]bool)
		for _, v := range arr {
			m[v] = true
		}
		m[b] = true

		// 转换回数组
		arr = make([]string, 0, len(m))
		for k := range m {
			arr = append(arr, k)
		}

		// 重新编码后存入
		data, _ = jsoniter.Marshal(arr)
		return t.db.Put([]byte(a), data, nil)
	}

	err = set_data(a, b)
	if err != nil {
		return err
	}
	err = set_data(b, a)
	if err != nil {
		return err
	}
	return err
}

// 单向解绑
func (t *ToEntangle) remove_item(a string, b string) (err error) {
	// 获取 a 的数据, 如果不存在, 则直接返回
	data, err := t.db.Get([]byte(a), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil
		}
		return err
	}

	// data 为json字符串的数组, 解码后移除 b
	var arr []string
	err = jsoniter.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	// 如果只有一个, 则直接移除条目
	if len(arr) == 1 {
		return t.db.Delete([]byte(a), nil)
	}

	// 直接遍历移除
	for i, v := range arr {
		if v == b {
			arr = append(arr[:i], arr[i+1:]...)
			break
		}
	}

	// 重新编码后存入
	data, _ = jsoniter.Marshal(arr)
	return t.db.Put([]byte(a), data, nil)
}

// Remove 移除一对数据, 使其双向解绑
func (t *ToEntangle) Remove(a string, b string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if err = t.remove_item(a, b); err != nil {
		return err
	}

	if err = t.remove_item(b, a); err != nil {
		return err
	}
	return err
}

// RemoveAll 移除一个数据, 使其全部解绑
func (t *ToEntangle) RemoveAll(a string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// 获取 a 的数据
	data, err := t.db.Get([]byte(a), nil)
	if err != nil {
		// 如果不存在, 则直接返回
		if err == leveldb.ErrNotFound {
			return nil
		}
		return err
	}

	// data 为json字符串的数组, 解码后移除 b
	var arr []string
	err = jsoniter.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	// 遍历移除
	for _, v := range arr {
		if err = t.remove_item(a, v); err != nil {
			return err
		}
		if err = t.remove_item(v, a); err != nil {
			return err
		}
	}
	return nil
}

// New 创建一个 ToEntangle
func New(path string) *ToEntangle {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}
	return &ToEntangle{
		db: db,
	}
}

// Close 关闭 leveldb
func (t *ToEntangle) Close() {
	t.db.Close()
}
