package mapx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type KVPair struct {
	Key   string
	Value string
}

type KVPairs []KVPair

func (r KVPairs) Len() int           { return len(r) }
func (r KVPairs) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r KVPairs) Less(i, j int) bool { return r[i].Key < r[j].Key }

func setArrayElem(arr *[]interface{}, idx int, elem interface{}) {
	if !(idx < len(*arr)) {
		nArr := make([]interface{}, idx+1-len(*arr))
		*arr = append(*arr, nArr...)
	}
	(*arr)[idx] = elem
}

func parseArrayKey(key string) (k string, idx []int) {
	idx = make([]int, 0)
	var i int
	for i = 0; i < len(key) && key[i] != '['; i++ {
	}
	if i == len(key) {
		return key, nil
	}

	k = key[0:i]
	for i < len(key) {
		if key[i] != '[' {
			return key, nil
		}
		left := i

		for ; i < len(key) && key[i] != ']'; i++ {
		}
		if i == len(key) {
			return key, nil
		}

		if d, err := strconv.Atoi(key[left+1 : i]); err != nil {
			return key, nil
		} else {
			idx = append(idx, d)
			i += 1
		}
	}
	return
}

func map2KVPairs(rv reflect.Value, keys []string, rows *KVPairs) {
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Interface:
		if rv.IsNil() {
			*rows = append(*rows, KVPair{strings.Join(keys, "."), "<nil>"})
		} else {
			map2KVPairs(rv.Elem(), keys, rows)
		}
	case reflect.Map:
		for _, mapKey := range rv.MapKeys() {
			mrv := rv.MapIndex(mapKey)
			map2KVPairs(mrv, append(keys, mapKey.String()), rows)
		}
	case reflect.Array, reflect.Slice:
		pkeys := keys[0 : len(keys)-1]
		key := keys[len(keys)-1]
		for i := 0; i < rv.Len(); i++ {
			erv := rv.Index(i)
			ikey := fmt.Sprintf("%s[%d]", key, i)
			map2KVPairs(erv, append(pkeys, ikey), rows)
		}
	case reflect.Float64:
		val := fmt.Sprintf("%v", rv.Float())
		*rows = append(*rows, KVPair{strings.Join(keys, "."), val})
	case reflect.Bool:
		var val string
		if rv.Bool() {
			val = "true"
		} else {
			val = "false"
		}
		*rows = append(*rows, KVPair{strings.Join(keys, "."), val})
	default:
		*rows = append(*rows, KVPair{strings.Join(keys, "."), rv.String()})
	}
}

func Map2KVPairs(m interface{}) KVPairs {
	val := reflect.ValueOf(m)
	if val.Kind() != reflect.Map {
		return nil
	}

	rows := KVPairs(make([]KVPair, 0))
	keys := make([]string, 0)
	map2KVPairs(val, keys, &rows)
	return rows
}

func putElemToArray(idx []int, elem interface{}, arr *[]interface{}) (
	err error) {
	if len(idx) == 1 {
		if idx[0] < len(*arr) && (*arr)[idx[0]] != nil {
			return fmt.Errorf("duplicate value: %v", elem)
		}
		setArrayElem(arr, idx[0], elem)
		return nil
	}

	var narr []interface{}
	if idx[0] < len(*arr) {
		val := (*arr)[idx[0]]
		if val == nil {
			narr = make([]interface{}, 0)
		} else {
			switch val.(type) {
			case []interface{}:
				narr = val.([]interface{})
			default:
				return fmt.Errorf("error value: %v", val)
			}
		}
	} else {
		narr = make([]interface{}, 0)
	}

	if err = putElemToArray(idx[1:], elem, &narr); err != nil {
		return
	}
	setArrayElem(arr, idx[0], narr)
	return
}

func putElemToMapArr(key string, idx []int, elem interface{},
	m map[string]interface{}) (
	err error) {
	var narr []interface{}
	if arr, exist := m[key]; exist {
		switch arr.(type) {
		case []interface{}:
			narr = arr.([]interface{})
		default:
			return fmt.Errorf("error, value: %v", arr)
		}
	} else {
		narr = make([]interface{}, 0)
	}
	if err = putElemToArray(idx, elem, &narr); err != nil {
		return
	}
	m[key] = narr
	return
}

func kvPairs2Map(keys []string, val interface{}, m map[string]interface{}) (
	err error) {
	if len(keys) == 1 {
		if key, idx := parseArrayKey(keys[0]); len(idx) == 0 {
			if _, exist := m[key]; exist {
				return fmt.Errorf("duplicate key: %s", key)
			}
			m[key] = val
		} else if err = putElemToMapArr(key, idx, val, m); err != nil {
			return
		}
		return nil
	}

	if key, idx := parseArrayKey(keys[0]); len(idx) == 0 {
		if mv, exist := m[key]; exist {
			switch mv.(type) {
			case map[string]interface{}:
				return kvPairs2Map(keys[1:], val,
					mv.(map[string]interface{}))
			default:
				return fmt.Errorf("error, value: %v", mv)
			}
		} else {
			nmv := map[string]interface{}{}
			m[key] = nmv
			return kvPairs2Map(keys[1:], val, nmv)
		}
	} else {
		nmv := make(map[string]interface{})
		if err = putElemToMapArr(key, idx, nmv, m); err != nil {
			return
		}
		return kvPairs2Map(keys[1:], val, nmv)
	}
}

func KVPairs2Map(rows KVPairs) (m map[string]interface{}, err error) {
	m = make(map[string]interface{})
	for _, row := range rows {
		var val interface{}
		keys := strings.Split(row.Key, ".")

		if row.Value == "<nil>" {
			val = nil
		} else {
			val = row.Value
		}
		if err = kvPairs2Map(keys, val, m); err != nil {
			return
		}
	}
	return
}
