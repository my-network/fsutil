package syncer

import (
	"container/heap"
)

type taskHeapInt []*task
type taskHeap taskHeapInt

func (tHeap *taskHeapInt) Less(i, j int) bool {
	a := (*tHeap)[i]
	b := (*tHeap)[j]
	aTS := a.ExpirationDeadline()
	bTS := b.ExpirationDeadline()
	return aTS.Before(bTS)
}

func (tHeap *taskHeapInt) Len() int {
	return len(*tHeap)
}

func (tHeap *taskHeapInt) Swap(i, j int) {
	(*tHeap)[i], (*tHeap)[j] = (*tHeap)[j], (*tHeap)[i]
	(*tHeap)[i].HeapIdx, (*tHeap)[j].HeapIdx = (*tHeap)[j].HeapIdx, (*tHeap)[i].HeapIdx
}

func (tHeap *taskHeapInt) Push(tI interface{}) {
	t := tI.(*task)
	t.HeapIdx = &[]int{tHeap.Len()}[0]
	*tHeap = append(*tHeap, t)
}

func (tHeap *taskHeapInt) Pop() interface{} {
	result := (*tHeap)[tHeap.Len()-1]
	*tHeap = (*tHeap)[:tHeap.Len()-1]
	result.HeapIdx = nil
	return result
}

func (tHeap *taskHeap) int() *taskHeapInt {
	return (*taskHeapInt)(tHeap)
}

func (tHeap *taskHeap) Len() int {
	return tHeap.int().Len()
}

func (tHeap *taskHeap) Push(t *task) {
	heap.Push(tHeap.int(), t)
}

func (tHeap *taskHeap) Pop() *task {
	return heap.Pop(tHeap.int()).(*task)
}

func (tHeap *taskHeap) Fix(t *task) {
	heap.Fix(tHeap.int(), *t.HeapIdx)
}

func (tHeap *taskHeap) Remove(t *task) {
	heap.Remove(tHeap.int(), *t.HeapIdx)
}
