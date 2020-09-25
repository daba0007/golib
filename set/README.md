# SET

set包用于处理集合，目前支持int类型和string类型的集合处理

## SET

```Go
stringList := []int{1, 2, 3}
// 创建一个新的Int集合
// 如果是string集合，则使用 NewStringSet
set := NewIntSet(stringList)
// 增加元素
set.add(4)
// 移除元素
set.remove(4)
// 是否有元素
set.has(3)
// 元素数量
set.count()
// 并集
set.union(NewIntSet(4))
// 减去集合
set.minus(NewIntSet(1, 2))
// 交集
set.intersect(NewIntSet(1, 2, 4))
// 右差集
set.complement(NewIntSet(3, 4))
// 清空集合
set.clear()
// 是否空集
set.empty()
```
