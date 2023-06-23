package toentangle

import (
	"sync"

	jsoniter "github.com/json-iterator/go"
	leveldb "github.com/syndtr/goleveldb/leveldb"
)

// 两表之间的双向绑定
// 1. 存入一对数据, 使其双向绑定
// 2. 移除一对数据, 使其双向解绑
// 3. 移除一个数据, 使其全部解绑
// 4. 获取一个数据的全部绑定数据

type ToEntangle struct {
	lock sync.Mutex
	dba  *leveldb.DB
	dbb  *leveldb.DB
}

// GetA 获取一个数据的全部绑定数据
func (t *ToEntangle) GetA(a string) (arr []string, err error) {
	// 获取 a 的数据, 如果不存在, 则直接返回空数组
	data, err := t.dba.Get([]byte(a), nil)
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

// GetB 获取一个数据的全部绑定数据
func (t *ToEntangle) GetB(b string) (arr []string, err error) {
	// 获取 b 的数据, 如果不存在, 则直接返回空数组
	data, err := t.dbb.Get([]byte(b), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return []string{}, nil
		}
		return nil, err
	}
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

	var set_data = func(a string, b string, db *leveldb.DB) (err error) {
		data, err := db.Get([]byte(a), nil)
		if err != nil {
			if err == leveldb.ErrNotFound {
				data, _ = jsoniter.Marshal([]string{b})
				return db.Put([]byte(a), data, nil)
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

		// 编码后存入
		data, _ = jsoniter.Marshal(arr)
		return db.Put([]byte(a), data, nil)
	}

	// 设置 a 的数据
	err = set_data(a, b, t.dba)
	if err != nil {
		return err
	}

	// 设置 b 的数据
	return set_data(b, a, t.dbb)
}

// 单向解绑
func (t *ToEntangle) remove(a string, b string, db *leveldb.DB) (err error) {
	data, err := db.Get([]byte(a), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil
		}
		return err
	}

	// data 为json字符串的数组, 解码后删除 b
	var arr []string
	err = jsoniter.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	// 转换为 map, 删除 b
	m := make(map[string]bool)
	for _, v := range arr {
		m[v] = true
	}
	delete(m, b)

	// 转换回数组
	arr = make([]string, 0, len(m))
	for k := range m {
		arr = append(arr, k)
	}

	// 编码后存入
	data, _ = jsoniter.Marshal(arr)
	return db.Put([]byte(a), data, nil)
}

// Remove 移除一对数据, 使其双向解绑
func (t *ToEntangle) Remove(a string, b string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// 移除 a 的数据
	err = t.remove(a, b, t.dba)
	if err != nil {
		return err
	}

	// 移除 b 的数据
	return t.remove(b, a, t.dbb)
}

// RemoveA 移除一个数据, 使其全部解绑
func (t *ToEntangle) RemoveA(a string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// 获取 a 的数据
	data, err := t.dba.Get([]byte(a), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil
		}
		return err
	}

	// data 为json字符串的数组, 解码后删除全部数据
	var arr []string
	err = jsoniter.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	// 删除全部数据
	for _, v := range arr {
		err = t.remove(a, v, t.dbb)
		if err != nil {
			return err
		}
	}

	// 删除 a 的数据
	return t.dba.Delete([]byte(a), nil)
}

// RemoveB 移除一个数据, 使其全部解绑
func (t *ToEntangle) RemoveB(b string) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// 获取 b 的数据
	data, err := t.dbb.Get([]byte(b), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil
		}
		return err
	}

	// data 为json字符串的数组, 解码后删除全部数据
	var arr []string
	err = jsoniter.Unmarshal(data, &arr)
	if err != nil {
		return err
	}

	// 删除全部数据
	for _, v := range arr {
		err = t.remove(b, v, t.dba)
		if err != nil {
			return err
		}
	}

	// 删除 b 的数据
	return t.dbb.Delete([]byte(b), nil)
}

// New 创建一个双向绑定 ToEntangle
func NewToEntangle(path string) *ToEntangle {
	dba, _ := leveldb.OpenFile(path+"/a", nil)
	dbb, _ := leveldb.OpenFile(path+"/b", nil)
	return &ToEntangle{
		dba: dba,
		dbb: dbb,
	}
}

// Close 关闭 leveldb
func (t *ToEntangle) Close() {
	t.dba.Close()
	t.dbb.Close()
}
