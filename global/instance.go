package global

import "sync"

//  确保是单例的
var (
	myonce     sync.Once
	myInstance *MyInstance // 注意必须得用指针这里，才能避免拷贝保证所有请求操作的是同一个实例
)

type MyInstance struct {
	SucceedCount      int
	FailedCount       int
	CurrentTraceIdVal int
	Lock              sync.Mutex
}

func GetMyInstance() *MyInstance {
	myonce.Do(func() {
		myInstance = &MyInstance{}
	})
	return myInstance
}
