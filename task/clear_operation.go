package task

import (
	"github.com/bearllflee/go_shop/global"
	"github.com/bearllflee/go_shop/model"
	"github.com/bearllflee/go_shop/service"
)

var operationRecordService = service.OperationRecordServiceApp

func ClearOperationRecord() error {
	var ids []uint
	var records []model.OperationRecord
	// 按创建时间升序排序
	global.DB.Model(&model.OperationRecord{}).Order("created_at asc").Limit(10).Find(&records)
	for _, record := range records {
		ids = append(ids, uint(record.ID))
	}
	return operationRecordService.DeleteOperationRecordByIds(ids)
}
