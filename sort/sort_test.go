package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestMergeSortSerial(t *testing.T) {

	s1 := []int{99, 1, 22, 78, 65, 35, 32, 11, 56, 88, 23}
	//s1 := []int{1, 4, 3}
	sequentialMergeSort(s1)
	fmt.Println(s1)
}

// 串行归并排序
func sequentialMergeSort(s []int) {
	if len(s) <= 1 {
		return
	}
	middle := len(s) / 2
	sequentialMergeSort(s[:middle]) // 前半部分
	sequentialMergeSort(s[middle:]) // 后半部分
	merge(s, middle)                // 合并这两部分
}

// 合并两部分整数切片
// 双指针
func merge(s []int, middle int) {
	tmp := make([]int, 0, len(s))
	front := 0
	back := middle
	if front == back {
		return
	}

	for front < middle || back < len(s) {

		if front == middle {
			tmp = append(tmp, s[back:]...)
			break
		}

		if back == len(s) {
			tmp = append(tmp, s[front:middle]...)
			break
		}

		if s[front] < s[back] {
			tmp = append(tmp, s[front])
			front++
		} else {
			tmp = append(tmp, s[back])
			back++
		}

	}

	copy(s, tmp)
}

const max = 2048

// 并行归并排序
func parallelMergesortV1(s []int) {
	if len(s) <= 1 {
		return
	}

	if len(s) <= max {
		sequentialMergeSort(s)
	} else {
		middle := len(s) / 2
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			parallelMergesortV1(s[:middle])
		}()

		go func() {
			defer wg.Done()
			parallelMergesortV1(s[middle:])
		}()

		wg.Wait()
		merge(s, middle)
	}
}
