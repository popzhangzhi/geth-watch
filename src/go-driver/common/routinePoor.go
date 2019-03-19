package common

import "fmt"

//任务类型
type Task struct {
	f func(int)
}

func NewTask(f func(int)) *Task {
	T := &Task{
		f: f,
	}
	return T

}

func (T *Task) Execute(workId int) {
	T.f(workId)
}

type Pool struct {
	//接受到的队列
	ReveiceChannel chan *Task
	TaskMaxNumber  int
	//当前正在执行的队列
	CurrentChannel chan *Task
}

func NewRouinePoor(taskMaxNumber int) *Pool {

	P := &Pool{
		ReveiceChannel: make(chan *Task, taskMaxNumber*4),
		TaskMaxNumber:  taskMaxNumber,
		CurrentChannel: make(chan *Task, taskMaxNumber*2),
	}
	return P
}

//协程池创建一个worker并且开始工作，阻塞等待CurrentChannel的值
func (P *Pool) worker(workId int) {
	//worker不断的从CurrentChannel内部任务队列中拿任务
	for task := range P.CurrentChannel {
		//如果拿到任务,则执行task任务
		task.Execute(workId)

	}
}

//协程开始运行入口
func (P *Pool) Run() {
	//根据限制声明出协程数
	for i := 1; i <= P.TaskMaxNumber; i++ {
		go P.worker(i)
	}
	//把外部读到的task传入内部
	for work := range P.ReveiceChannel {
		P.CurrentChannel <- work
	}
	//依次关闭协程
	fmt.Println("关闭pool")
	close(P.CurrentChannel)
	close(P.ReveiceChannel)

}
