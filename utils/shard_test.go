package utils

import (
	"fmt"
	"testing"
)

func TestShard(t *testing.T) {
	cases := []struct {
		inputKey     string
		inputPartNum int
		outPartion   int
	}{
		{"000000000092", 64, 2},
		{"00000100009h", 64, 2},
		{"000i00000092", 64, 2},
		{"00000m000092", 64, 2},
		{"adfafalk", 64, 2},
		{"s13erka98042", 64, 2},
	}

	for _, c := range cases {
		testname := fmt.Sprintf("%s-%d", c.inputKey, c.inputPartNum)
		t.Run(testname, func(t *testing.T) {
			out := Shard(c.inputKey, c.inputPartNum)
			fmt.Println(c.inputKey, c.inputPartNum, out)
			// if out != c.outPartion {
			// 	t.Errorf("c %s: want %d but get %d", testname, out, c.outPartion)
			// }
		})
	}
}
