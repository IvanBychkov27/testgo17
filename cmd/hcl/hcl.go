package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

var (
	defaultUpdateInterval    uint = 60      // Проверять файлы в папке каждые 60 секунд
	defaultControlAddress         = ":2000" //
	defaultDBTimeout         uint = 30      // таймаут в секундах
	defaultSendBatchInterval uint = 30
	defaultConcurrent        uint = 16
)

type ConnectionData struct {
	Schema   string   `hcl:"schema"`
	Host     []string `hcl:"host"`
	Port     int      `hcl:"port"`
	Database string   `hcl:"database"`
	Table    string   `hcl:"table"`
	Username string   `hcl:"username"`
	Password string   `hcl:"password"`
	Cert     string   `hcl:"cert"`
	Timeout  uint     `hcl:"timeout,optional"`
}

type Task struct {
	Name              string         `hcl:"name,label"`
	SourceFolder      string         `hcl:"sourcePath"`
	ErrorFolder       string         `hcl:"errorPath"`
	UpdateInterval    uint           `hcl:"updateInterval"`
	BatchSize         uint           `hcl:"batchSize,optional"`
	MoveToError       bool           `hcl:"moveToError"`
	Prefix            string         `hcl:"prefix,optional"`
	Suffix            string         `hcl:"suffix,optional"`
	MinAge            uint           `hcl:"minAge,optional"`
	SendBatchInterval uint           `hcl:"sendBatchInterval,optional"`
	Concurrent        uint           `hcl:"concurrent,optional"`
	FilesCountLimit   uint           `hcl:"filesCountLimit,optional"`
	ConnectionData    ConnectionData `hcl:"connection,block"`
}

type Config struct {
	Tasks          []*Task `hcl:"task,block"`
	ControlAddress string  `hcl:"controlAddress,optional"`
	IsDebug        bool    `hcl:"debug,optional"`
}

func main() {
	cfg := &Config{}

	fileName := "cmd/hcl/config.hcl"
	errDecode := hclsimple.DecodeFile(fileName, nil, cfg)
	if errDecode != nil {
		fmt.Println("error", errDecode)
	}

	fmt.Println("len cfg.Task     =", len(cfg.Tasks))
	fmt.Println("cfg.Task[0].Name =", cfg.Tasks[0].Name)
	fmt.Println("cfg.Task[0].ConnectionData.Host =", cfg.Tasks[0].ConnectionData.Host)
	fmt.Println("cfg.Task[0].ConnectionData.Host[0] =", cfg.Tasks[0].ConnectionData.Host[0])
	fmt.Println("cfg.Task[0].ConnectionData.Host[1] =", cfg.Tasks[0].ConnectionData.Host[1])
	fmt.Println("cfg.Task[0].ConnectionData.Host[2] =", cfg.Tasks[0].ConnectionData.Host[2])

	fmt.Println()
	fmt.Println("cfg.Task[1].Name =", cfg.Tasks[1].Name)
	fmt.Println("cfg.Task[1].ConnectionData.Host =", cfg.Tasks[1].ConnectionData.Host)

}

func parseTask(cfg *Config) (*Config, error) {
	if cfg.ControlAddress == "" {
		cfg.ControlAddress = defaultControlAddress
	}

	names := map[string]struct{}{}

	for _, task := range cfg.Tasks {
		_, ok := names[task.Name]
		if ok {
			return nil, fmt.Errorf("duplicated task name: %s", task.Name)
		}
		names[task.Name] = struct{}{}

		if task.SourceFolder == "" {
			return nil, fmt.Errorf("empty source folder for task %s", task.Name)
		}
		if task.MoveToError && task.ErrorFolder == "" {
			return nil, fmt.Errorf("empty error folder for task %s", task.Name)
		}
		if task.UpdateInterval == 0 {
			task.UpdateInterval = defaultUpdateInterval
		}
		if task.SendBatchInterval == 0 {
			task.SendBatchInterval = defaultSendBatchInterval
		}
		if task.Concurrent == 0 {
			task.Concurrent = defaultConcurrent
		}
		if task.ConnectionData.Timeout == 0 {
			task.ConnectionData.Timeout = defaultDBTimeout
		}
	}

	return cfg, nil
}

/*
data :=`
task "Binder" {

  sourcePath = "/opt/stat/d-binder/buffer"
  errorPath = "/opt/stat/d-binder/error"
  updateInterval = 60
  batchSize = 8
  moveToError = true
  prefix = ""
  suffix = ".buffer"
  minAge = 0
  sendBatchInterval = 120
  concurrent = 1
  filesCountLimit = 0

  connection {
    schema = "https"
    host = ["rc1c-32ykinimy6oha17h.mdb.yandexcloud.net","rc1c-33ykinimy6oha17h.mdb.yandexcloud.net","rc1c-34ykinimy6oha17h.mdb.yandexcloud.net"]
    port = 8443
    database = "plat"
    table = "d_table"
    username = "platform"
    password = "JMFi"
    cert = "/YandexInternalRootCA.crt"
    timeout = 600
  }
}
`
*/
