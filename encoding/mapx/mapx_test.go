package mapx

import (
	"encoding/json"
	"fmt"
	"testing"
)

const JSON_STR = `{
	"a": {
		"b": [
			[1, 2, 3],
			[4, 5, 6]
		]
	},
	"c": {
		"d": "hello",
		"f": [1, 2, 3],
		"g": 123,
		"h": {
			"i": true,
			"j": "world",
			"k": [1, 2, 3]
		}
	}
}`

func TestKVPairsEncode(t *testing.T) {
	var conf map[string]interface{}
	if err := json.Unmarshal([]byte(JSON_STR), &conf); err != nil {
		panic(err)
	}
	fmt.Println(conf)

	rows := Map2KVPairs(conf)
	for _, r := range rows {
		fmt.Println(r)
	}

	m, err := KVPairs2Map(rows)
	if err != nil {
		panic(err)
	}
	fmt.Println(m)
}

func BenchmarkKVPairsEncode(b *testing.B) {
	var conf map[string]interface{}
	if err := json.Unmarshal([]byte(JSON_STR), &conf); err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		Map2KVPairs(conf)
	}
}

func BenchmarkKVPairsDecode(b *testing.B) {
	var conf map[string]interface{}
	if err := json.Unmarshal([]byte(JSON_STR), &conf); err != nil {
		panic(err)
	}
	rows := Map2KVPairs(conf)

	for i := 0; i < b.N; i++ {
		KVPairs2Map(rows)
	}
}
