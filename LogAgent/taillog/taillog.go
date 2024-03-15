package taillog

import (
	"GO/LogAgent/kafka"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hpcloud/tail"
)

// Manager Etcd 更新时LogEntry中的数据会改变，需要动态的启停LogTask,其中Map用来CURD
type Manager struct {
	LogEntries []LogEntry
	LogTasks   []*LogTask
	//key:path+"_"+topic
	TaskMap   map[string]*LogTask
	NewConfCh chan []byte
}

var Mgr Manager

// LogEntry 用来json unmarshal的
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// LogTask 用来开启TAIL 任务获取Log日志并发往Kafka的 同时还需要支持Etcd 动态热更新
type LogTask struct {
	topic   string
	path    string
	tailObj *tail.Tail
	ctx     context.Context //context.Context 是个interface 不是结构体
	cancel  context.CancelFunc
}

func Init() {
	//初始化TaskMap,以及通道
	Mgr.TaskMap = make(map[string]*LogTask, 16)
	//Question:又是一个大Question:我通道没初始化就把值传过来了，导致接受不到值,多写err:=
	//	单通道就是为了避免异步问题的
	Mgr.NewConfCh = make(chan []byte, 1)
	//	开启这一步时,LogEntry已有数据,把LogEntry的数据加到LogTask中
	for _, v := range Mgr.LogEntries {
		t := LogTask{
			path:  v.Path,
			topic: v.Topic,
		}
		Mgr.LogTasks = append(Mgr.LogTasks, &t)
	}
	//	初始化LogTask
	for _, vv := range Mgr.LogTasks {
		vv.init()
	}
	//	这个是配置热更新管理端
	//	Question:这里也需要异步执行,否则会一直卡住
	go Mgr.run()
}

func (t *LogTask) init() {
	//	LogTask有了基本的path和topic之后，可执行init
	//	tail 读取数据
	tailObj, err := newTailObj(t.path)
	if err != nil {
		fmt.Println("Tail file failed,err:", err)
		fmt.Println("filepath:", t.path)
	}
	t.tailObj = tailObj
	//	设置context
	ctx, cancel := context.WithCancel(context.Background())
	t.ctx = ctx
	t.cancel = cancel
	//	更新Manager的map
	Mgr.TaskMap[t.path+"_"+t.topic] = t
	//	运行起来
	//	Question:需要异步执行Run才行
	go t.run()
}

// 内置函数，便于新建tail对象
func newTailObj(filename string) (tailObj *tail.Tail, err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tailObj, err = tail.TailFile(filename, config)
	return
}

// 获取数据
func (t *LogTask) run() {
	for {
		select {
		case line, ok := <-t.tailObj.Lines:
			if !ok {
				fmt.Printf("tail file close reopen,filename:%s\n", t.tailObj.Filename)
				continue
			}
			//fmt.Println("line:" + line.Text)
			//把数据发送给Kafka
			kafka.SendToKafka(line.Text, t.topic)
		case <-t.ctx.Done():
			//	说明要结束任务了
			return
		}
	}
}

// manager run:保持接受Etcd的变化
// 注意异步,这里应该要上锁
func (m *Manager) run() {
	//如果接受到新的数据
	var newConf []LogEntry
	for v := range m.NewConfCh {
		fmt.Println("新的配置到达")
		if v == nil {
			//Question:这里出现一个问题:如果传入的数据是空的,那么就会unmarshal失败
			fmt.Println("新的配置是空的!")
			newConf = nil
		} else {
			err := json.Unmarshal(v, &newConf)
			if err != nil {
				//这里或许continue更好
				panic(err)
			}
		}

		//1.先删除
		for _, old := range m.LogEntries {
			var exist = false
			for _, n := range newConf {
				if old.Path+"_"+old.Topic == n.Path+"_"+n.Topic {
					exist = true
				}
			}
			if !exist {
				m.TaskMap[old.Path+"_"+old.Topic].cancel() //结束这个任务
				//	还得从Map中把这个key-Value删除
				delete(m.TaskMap, old.Path+"_"+old.Topic)
				fmt.Printf("任务:%v 被结束了!\n", old)
			}
		}

		//2.再新增
		for _, n := range newConf {
			_, ok := m.TaskMap[n.Path+"_"+n.Topic]
			//如果有新配置,创建
			if !ok {
				//	在LogTask中增加
				task := &LogTask{
					topic: n.Topic,
					path:  n.Path,
				}
				task.init()
				m.LogTasks = append(m.LogTasks, task)
				//	在TaskMap中增加
				m.TaskMap[task.path+"_"+task.topic] = task
			}
		}

		//更新LogEntry
		m.LogEntries = newConf
	}
}
