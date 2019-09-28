package week_1

import (
	"math/rand"
	"testing"
	"time"
)

func TestBubbleSort(t *testing.T) {
	t.Log("冒泡排序用例")
	length := 10
	data := make([]int, length)
	rand.Seed(time.Now().Unix()) //设置随机数种子，可以保证每次随机都是随机的，默认值是 1，所以每次生成的随机数会是一样的。
	for i := 0; i < length; i++ {
		data[i] = rand.Intn(100)
	}
	t.Logf("初始化待排序数据：%v", data)

	/**
	冒泡排序
	双层循环遍历，每一次的内循环会求出一个靠左/右的数字
	两层循环加一次常数级别的操作，故复杂度: O(n) * O(n) * O(1) = O(n^2)
	*/
	sortTimes := 0
	lengthK := length - 1
	for i := 0; i < lengthK; i++ {
		doSwap := false
		for j := 0; j < lengthK; j++ {
			sortTimes++
			if data[j] > data[j+1] { //顺序
				tmp := data[j+1]
				data[j+1] = data[j]
				data[j] = tmp
				doSwap = true
			}
		}
		if !doSwap { //有一次发现空循环了，后续的循环就不需要了，已经完成排序了
			break
		}
	}
	t.Logf("循环次数：%d", sortTimes)
	t.Logf("排序后的数据输出：%v", data)
}
