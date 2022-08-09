package id

// CyclingIdxFunc
// @title       可循环的ID、和递增函数
// @description 可循环的ID、和递增函数，不支持并发。
// @auth        小卒    2022/08/03 10:57
// @return      idx   *uint64  "IP端口的信息"
// @return      after func()   "递增函数"
func CyclingIdxFunc() (*uint64, *uint64, func()) {
	// 一次循环使用的ID的最大值
	const maxIdxPerCycle = ^uint64(0) - 99999999
	// 表示第几次循环
	cycle := uint64(0)
	// 当前ID值
	idx := uint64(0)
	// id函数递增
	after := func() {
		if idx >= maxIdxPerCycle {
			idx = 0
			cycle++
		}
		idx++
	}
	return &cycle, &idx, after
}
