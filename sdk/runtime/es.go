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

func (e *Es) AddRecord(id int, title, remark, textData, tags string) (string, error) {
	return e.es.AddRecord(id, title, remark, textData, tags)
}

func (e *Es) GetRecord(id int) (*es.Knowledge, error) {
	return e.es.GetRecord(id)
}

func (e *Es) UpdateRecord(id int, title, remark, textData, tags string) error {
	return e.es.UpdateRecord(id, title, remark, textData, tags)
}

func (e *Es) DeleteRecord(id int) error {
	return e.es.DeleteRecord(id)
}

func (e *Es) PageRecord(size int, page int, keyword string) (int64, []int, error) {
	return e.es.PageRecord(size, page, keyword)
}
