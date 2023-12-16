package es

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
)

type EsOption struct {
	Host  string
	Index string
	Type  string
}

func NewEs(esCli *elastic.Client, options *EsOption) (*Elasticsearch, error) {
	var err error
	s := &Elasticsearch{Client: esCli, Options: options}
	s.Client, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(options.Host),
	)
	if err != nil {
		return nil, err
	}
	//连接服务测试
	_, _, err = s.Client.Ping(options.Host).Do(context.Background())
	if err != nil {
		return nil, err
	}
	//获取服务端版本号
	_, err = s.Client.ElasticsearchVersion(options.Host)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type Elasticsearch struct {
	Client  *elastic.Client
	Options *EsOption
}

// Knowledge 数据结构体
type Knowledge struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`    // 标题
	Remark   string `json:"remark"`   // 摘要
	TextData string `json:"textData"` // 内容
	Tags     string `json:"tags"`     // 标签
}

// 添加记录
func (e *Elasticsearch) AddRecord(id int, title, remark, textData, tags string) (string, error) {
	//添加数可以用结构体方式和 json字符串方式
	p := Knowledge{id, title, remark, textData, tags}
	put, err := e.Client.Index().Index(e.Options.Index).Type(e.Options.Type).Id(strconv.Itoa(id)).BodyJson(p).Do(context.Background())
	if err != nil {
		return "", err
	}
	return put.Id, nil
}

// 查询记录
func (e *Elasticsearch) GetRecord(id int) (*Knowledge, error) {
	//通过id查找
	get1, err := e.Client.Get().Index(e.Options.Index).Type(e.Options.Type).Id(strconv.Itoa(id)).Do(context.Background())
	if err != nil {
		panic(err)
	}
	var k Knowledge
	if get1.Found {
		err := json.Unmarshal(get1.Source, &k) //将结果集中的数据映射到结构体中
		if err != nil {
			return nil, err
		}
	}
	return &k, nil
}

// 修改
func (e *Elasticsearch) UpdateRecord(id int, title, remark, textData, tags string) error {
	data := make(map[string]interface{})
	if title != "" {
		data["title"] = title
	}
	if remark != "" {
		data["remark"] = remark
	}
	if textData != "" {
		data["textData"] = textData
	}
	if tags != "" {
		data["tags"] = tags
	}
	_, err := e.Client.Update().Index(e.Options.Index).Type(e.Options.Type).Id(strconv.Itoa(id)).
		Doc(data).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 删除一条
func (e *Elasticsearch) DeleteRecord(id int) error {
	_, err := e.Client.Delete().
		Index(e.Options.Index).Type(e.Options.Type).Id(strconv.Itoa(id)).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 分页
func (e *Elasticsearch) PageRecord(size int, page int, keyword string) ([]int, error) {
	var res *elastic.SearchResult
	var err error
	q := elastic.NewBoolQuery()
	q.Must(elastic.NewMatchQuery("title", keyword))
	q.Must(elastic.NewMatchQuery("remark", keyword))
	q.Must(elastic.NewMatchQuery("tags", keyword))
	q.Must(elastic.NewMatchQuery("textData", keyword))
	res, err = e.Client.Search(e.Options.Index).Query(q).Type(e.Options.Type).Size(size).From((page - 1) * size).Do(context.Background())
	if err != nil {
		panic(err)
	}
	ids := make([]int, 0)
	num := res.Hits.TotalHits.Value //搜索到结果总条数
	if num > 0 {
		for _, hit := range res.Hits.Hits {
			id, _ := strconv.Atoi(hit.Id)
			ids = append(ids, id)
		}
	}
	return ids, nil
}