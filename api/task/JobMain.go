package task

import (
	"github.com/beego/beego/v2/adapter/logs"
	"github.com/beego/beego/v2/adapter/toolbox"
	"sync"
	"v1-go-api/graph/models"
)

const (
	CRONTAB_SECOND = "*/1 * * * * *"
	CRONTAB_MINUTE = "0 */1 * * * *"
	_              = iota
	POOL_TYPE_PAIR
	POOL_TYPE_TOKEN
)

var (
	LOCK          sync.RWMutex
	GraphDataList map[string]*models.GraphData
	logger        = logs.GetBeeLogger()
)

func init() {
	GraphDataList = make(map[string]*models.GraphData)
}

func Start() {
	task := toolbox.NewTask("task", CRONTAB_SECOND, Do)
	toolbox.AddTask("Do", task)
	toolbox.StartTask()
}
