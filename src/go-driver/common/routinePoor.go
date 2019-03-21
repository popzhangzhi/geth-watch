package common

import (
	"strconv"
	"sync"
	"time"
)

//任务类型
type Task struct {
	//任务执行函数
	f func(map[string]string)
	//执行时外部自定义注入的参数
	params map[string]string
}

func NewTask(f func(map[string]string), params map[string]string) *Task {
	if params["taskId"] == "" {
		panic("必须有taskId（任务唯一标识），本批次中唯一，无需全局唯一")
	}
	T := &Task{
		f:      f,
		params: params,
	}
	return T

}

func (T *Task) Execute(arg map[string]string) {
	//fmt.Println(arg)
	T.f(arg)

}

//运行任务后结果
type PoolRel struct {
	rel  map[string]int
	Lock sync.Mutex
}

type Pool struct {
	//接受到的队列
	ReveiceChannel chan *Task
	TaskMaxNumber  int
	//当前正在执行的队列
	CurrentChannel chan *Task
	PoolRel        PoolRel
}

//创建协程池，阻塞，非缓存的形式
func NewRouinePoor(taskMaxNumber int) *Pool {

	P := &Pool{
		ReveiceChannel: make(chan *Task),
		TaskMaxNumber:  taskMaxNumber,
		CurrentChannel: make(chan *Task),
		PoolRel: PoolRel{
			rel: make(map[string]int),
		},
	}
	return P
}

//协程池创建一个worker并且开始工作，阻塞等待CurrentChannel的值
func (P *Pool) worker(workId int) {
	//worker不断的从CurrentChannel内部任务队列中拿任务
	for task := range P.CurrentChannel {
		//如果拿到任务,创建自定义参数，执行task任务
		arg := make(map[string]string)
		arg["workId"] = strconv.Itoa(workId)
		for k, v := range task.params {
			arg[k] = v
		}

		task.Execute(arg)

		P.workDone(arg["taskId"])
	}
}

//记录对应taskId的任务运行结束，不记录成功或者失败，协程已成功结束
func (P *Pool) workDone(taskId string) {
	P.PoolRel.Lock.Lock()
	defer P.PoolRel.Lock.Unlock()
	P.PoolRel.rel[taskId] = 1
}

//协程开始运行入口
func (P *Pool) Run() {
	//根据限制声明出协程数
	for i := 1; i <= P.TaskMaxNumber; i++ {
		go P.worker(i)
	}
	//把外部读到的task传入内部,主进程上阻塞，但在边界值情况下，task有可能没有时间在协程中运行就被程序就退出了
	var count int = 0
	for work := range P.ReveiceChannel {

		count++
		P.CurrentChannel <- work
	}
	//利用共享所，保证所有任务支持完成后才会退出
	for len(P.PoolRel.rel) != count {
		time.Sleep(10 * time.Second)
	}

	close(P.CurrentChannel)

}

//关闭channel
func (P *Pool) Close() {

	close(P.ReveiceChannel)
}
