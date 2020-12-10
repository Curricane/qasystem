package filter

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	err := Init("../data/filter.dat.txt")
	if err != nil {
		t.Errorf("faied to Init data,err is %#v\n", err)
	}

	data := "昨日天气不错但你妈炸了，你的口头禅是草你妈，真是狗日的"
	ret, isr := Replace(data, "***")
	fmt.Printf("isr: %#v, retL%v\n", isr, ret)
}
