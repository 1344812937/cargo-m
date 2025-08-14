package tasks

import (
	"cargo-m/internal/service"
	"github.com/robfig/cron/v3"
)

type CronTask struct {
	cron *cron.Cron
}

// NewCronTask 注册当前定时任务
func NewCronTask(mavenService *service.MavenService) *CronTask {
	println("定时任务注册")
	res := &CronTask{cron: cron.New()}
	_, err := res.cron.AddFunc("*/5 * * * *", func() {
		println("扫描本地目录任务开始")
		mavenService.GetLocalMavenRepo()
		println("扫描本地目录任务结束")
	})
	if err != nil {
		panic(err)
	}
	println("定时任务注册成功!")
	return res
}

func (t *CronTask) Start() {
	t.cron.Start()
	// defer t.cron.Stop()
}

func (t *CronTask) Stop() {
	t.cron.Stop()
}
