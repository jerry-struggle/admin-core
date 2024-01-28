package runtime

import (
	"github.com/jerry-struggle/admin-core/storage"
	"github.com/jerry-struggle/admin-core/storage/es"
)

// NewSms 创建对应上下文短信对象
func NewEs(es storage.AdapterEs) storage.AdapterEs {
	return &Es{
		es: es,
	}
}

type Es struct {
	es storage.AdapterEs
}

func (e *Es) AddRecord(index string, id int, title, remark, textData, tags string) (string, error) {
	return e.es.AddRecord(index, id, title, remark, textData, tags)
}

func (e *Es) GetRecord(index string, id int) (*es.Knowledge, error) {
	return e.es.GetRecord(index, id)
}

func (e *Es) UpdateRecord(index string, id int, title, remark, textData, tags string) error {
	return e.es.UpdateRecord(index, id, title, remark, textData, tags)
}

func (e *Es) DeleteRecord(index string, id int) error {
	return e.es.DeleteRecord(index, id)
}

func (e *Es) PageRecord(index string, size int, page int, keyword string) (int64, []int, error) {
	return e.es.PageRecord(index, size, page, keyword)
}

func (e *Es) RecordList(index string, keyword string) (int64, []int, error) {
	return e.es.RecordList(index, keyword)
}
