// SPDX-License-Identifier: MIT

// Package sliceutil 提供对数组和切片的相关功能
package sliceutil

import (
	"fmt"
	"reflect"
)

// Reverse 反转数组中的元素
func Reverse(slice interface{}) {
	v := getSliceValue(slice, true)
	l := v.Len()
	swap := reflect.Swapper(v.Interface()) // 采用 v.Interface{}，而不是 slice，slice 可能是指针

	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Delete 删除 slice 中符合 eq 条件的元素
//
// slice 的类型只能是切片或是切片指针，其它任意类型都将 panic，数组也不行；
// eq 对比函数，用于确定指定的元素是否可以删除，返回 true 表示可以删除；
// size 返回新的数组大小，用户可以从原始数组上生成新的数组：
//  slice[:size]
func Delete(slice interface{}, eq func(i int) bool) (size int) {
	v := getSliceValue(slice, true)
	l := v.Len()

	var cnt int
	swap := reflect.Swapper(v.Interface()) // 采用 v.Interface{}，而不是 slice，slice 可能是指针
	last := l - 1
	for i := 0; i <= last; i++ {
		if !eq(i) {
			continue
		}

		for j := i; j < last; j++ {
			swap(j, j+1)
		}
		cnt++
		i--
		last--
	}

	return l - cnt
}

// QuickDelete 删除 slice 中符合 eq 条件的元素
//
// 功能与 Delete 相同，但是性能相对 Delete 会好一些，同时也不再保证元素与原数组相同。
func QuickDelete(slice interface{}, eq func(i int) bool) (size int) {
	v := getSliceValue(slice, true)
	l := v.Len()

	var cnt int
	swap := reflect.Swapper(v.Interface())
	last := l - 1
	for i := 0; i <= last; i++ {
		if !eq(i) {
			continue
		}

		swap(i, last)
		cnt++
		last--
		i--
	}

	return l - cnt
}

// Count 检测数组中指定值的数量
//
// slice 需要检测的数组或是切片，其它类型会 panic；
// eq 对比函数，i 表示数组的下标，需要在函数将该下标表示的值与你需要的值进行比较是否相等；
func Count(slice interface{}, eq func(i int) bool) (count int) {
	v := getSliceValue(slice, false)
	l := v.Len()

	for i := 0; i < l; i++ {
		if eq(i) {
			count++
		}
	}
	return
}

// Dup 检测数组或是切片中是否包含重复的值
//
// slice 需要检测的数组或是切片，其它类型会 panic；
// eq 对比数组中两个值是否相等，相等需要返回 true；
// 返回值表示存在相等值时，第二个值在数组中的下标值；
func Dup(slice interface{}, eq func(i, j int) bool) int {
	v := getSliceValue(slice, false)
	l := v.Len()
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			if eq(i, j) {
				return j
			}
		}
	}

	return -1
}

// Contains container 是否包含了 sub 中的所有元素
//
// container 与 sub 都必须是数组或是切片类型。
// 如果只是需要判断某一个值是否在 container 中，可以使用 Count() 函数。
//
// 数组或是切 片的元素类型未必要相同，只要值是可比较的就行。
func Contains(container, sub interface{}) bool {
	c := getSliceValue(container, false)
	s := getSliceValue(sub, false)

	cl := c.Len()
	sl := s.Len()
	if sl > cl {
		return false
	}

LOOP:
	for i := 0; i < sl; i++ {
		sv := s.Index(i)
		st := sv.Type()

		for j := 0; j < cl; j++ {
			cv := c.Index(j)
			ct := cv.Type()

			if sv.Interface() == cv.Interface() ||
				(st.ConvertibleTo(ct) && cv.Interface() == sv.Convert(ct).Interface()) ||
				(ct.ConvertibleTo(st) && sv.Interface() == cv.Convert(st).Interface()) {
				continue LOOP
			}
		}
		return false
	}
	return true
}

func getSliceValue(slice interface{}, onlySlice bool) reflect.Value {
	v := reflect.ValueOf(slice)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if onlySlice && v.Kind() != reflect.Slice {
		panic(fmt.Sprint("参数 slice 只能是 slice"))
	}

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		panic(fmt.Sprint("参数 slice 只能是 slice 或是 array"))
	}

	return v
}
