# Buffer

有两个版本 v1 和 v2



## v1

利用了 `runtime.SetFinalizer`，当变量在被系统回收的前一刻，我们清空了 buffer 的数据，并将归还给 sync.Pool

```go
runtime.SetFinalizer(buffer, func(buffer *Buffer) {
	buffer.Reset()
	pool.Put(buffer.inner)
	buffer.inner = nil
})
```

优点：编码简单

缺点：回收慢，导致创建很多 buffer



## v2

- 本库提供了 SimplePool 和 StatisticalPool
  - SimplePool 简单的 Pool
  - StatisticalPool 支持统计的 Pool
  - 用户自定义 Pool 需要实现 Pool 接口

- 自定义对象，实现 `sync.RefCountable` 和 `sync.Releasable`

  ```go
  type Buffer struct {
  	inner *bytes.Buffer
  	sync.RefCountable
  	sync.Releasable
  }
  ```

  > sync.RefCountable 不需要用户自己实现，相应的，用户需要告诉 Pool 如何构建自定义对象，例如：
  >
  > ```go
  > newFunc := func(factory sync.RefCountableFactory) sync.RefCountable {
  > 	buffer := &Buffer{inner: &bytes.Buffer{}}
  > 	buffer.inner.Grow(defaultBufferSize)
  > 	buffer.RefCountable = factory(buffer)
  > 	return buffer
  > }
  > // pool_ = sync.NewSimplePool(newFunc)
  > pool_ = sync.NewStatisticalPool(newFunc)
  > ```
  >
  > 内部提供 factory sync.RefCountableFactory 用于简化编程，buffer.RefCountable = factory(buffer) 完成了初始化

- 使用

  ```go
  // 创建 Buffer
  buffer := v2.NewBuffer()
  // 写入数据
  _, _ = buffer.Write([]byte("hello"))
  
  // 增加引用，传递给 myFunc
  buffer.Ref()
  // myFunc 不使用 buffer 之后，应该显式调用 buffer.DeRef() 释放
  // 当释放之后，就不应该再使用它
  myFunc(buffer)
  
  // 这里也需要释放
  buffer.DeRef()
  ```

优点：回收快，性能高，特别适用于分发只读的共享数据

缺点：代码复杂，需要在合适的时机调用 DeRef() 释放，任何地方修改了 buffer 都会影响到其他使用者，即 buffer 是共享的

> **说明：**
>
> 因为 v2 适用于分发只读的共享数据，所以用户应该考虑如何将数据设计成只读的，避免在共享过程中其他使用者破坏数据。

