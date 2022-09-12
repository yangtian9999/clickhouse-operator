package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
)

func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func ArraySearch(target string, str_array []string) bool {
	for _, str := range str_array {
		if target == str {
			return true
		}
	}
	return false
}

func Md5CheckSum(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:16])
}


type Collection interface {
	Union(oth interface{})interface{}
	Intersect(oth interface{})interface{}
	Difference(oth interface{})interface{}
}

type Map map[string]interface{}

// Union elem both in this and oth, If the same key have different value, use this
func (this  Map)Union(oth interface{})interface{} {
	out := make(Map)
	other := oth.(Map)
	for k, v := range other {
		out[k] = v
	}

	for k, v := range this {
		out[k] = v
	}
	return out
}

// Intersect if this and oth both have key, but have not equal value, discover with this's value
func (this  Map)Intersect(oth interface{})interface{} {
	out := make(Map)
	other := oth.(Map)
	for k, v := range this {
		if _, ok := other[k]; ok {
			out[k] = v
		}
	}

	return out
}

// Difference elemt in this, not in oth
func (this  Map)Difference(oth interface{})interface{} {
	out := make(Map)
	other := oth.(Map)
	for k, v := range this {
		if _, ok := other[k]; !ok {
			out[k] = v
		}
	}
	return out
}
