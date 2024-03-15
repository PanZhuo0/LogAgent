package taillog

import (
	"GO/LogAgent/kafka"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
)

// Manager Etcd 更新时LogEntry中的数据会改变，需要动态的启停LogTask,其中Map用来CURD
type Manager struct {
	LogEntries []LogEntry
	LogTasks   []*LogTask
	TaskMap    map[string]LogTask
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
		//	tail 读取数据
		tailObj, err := newTailObj(vv.path)
		if err != nil {
			fmt.Println("Tail file failed,err:", err)
			fmt.Println("filepath:", vv.path)
		}
		vv.tailObj = tailObj
		//	设置context
		ctx, cancel := context.WithCancel(context.Background())
		vv.ctx = ctx
		vv.cancel = cancel
		//	运行起来
		//	Question:需要异步执行Run才行
		go vv.run()
	}
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
