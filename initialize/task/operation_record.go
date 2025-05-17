package task

import (
	"log"
	"math"
	"time"

	"github.com/bearllflee/go_shop/global"
	"github.com/bearllflee/go_shop/task"
)

func ClearOperationRecord(cronString string) {
	global.Cron.AddFunc(cronString, func() {
		// task.ClearOperationRecord()
		ExecuteWithRetry(task.ClearOperationRecord, 3)
	})
}

func ExecuteWithRetry(job func() error, maxRetries int) {
	for i := range maxRetries {
		err := job()
		if err == nil {
			return
		}

		log.Printf("第 %d 次尝试失败，:%v", i+1, err)
		time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
	}
	log.Printf("%d次尝试均失败", maxRetries)
}
