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

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(options.Host),
	)
	if err != nil {
		return nil, err
	}
	//连接服务测试
	_, _, err = client.Ping(options.Host).Do(context.Background())
	if err != nil {
		return nil, err
	}
	//获取服务端版本号
	_, err = client.ElasticsearchVersion(options.Host)
	if err != nil {
		return nil, err
	}
	return &Elasticsearch{Client: client, Options: options}, nil
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
func (e *Elasticsearch) AddRecord(index string, id int, title, remark, textData, tags string) (string, error) {
	//添加数可以用结构体方式和 json字符串方式
	p := Knowledge{id, title, remark, textData, tags}
	put, err := e.Client.Index().Index(index).Id(strconv.Itoa(id)).BodyJson(p).Do(context.Background())
	if err != nil {
		return "", err
	}
	return put.Id, nil
}

// 查询记录
func (e *Elasticsearch) GetRecord(index string, id int) (*Knowledge, error) {
	//通过id查找
	get1, err := e.Client.Get().Index(index).Id(strconv.Itoa(id)).Do(context.Background())
	if err != nil {
		return nil, err
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
func (e *Elasticsearch) UpdateRecord(index string, id int, title, remark, textData, tags string) error {
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
	_, err := e.Client.Update().Index(index).Id(strconv.Itoa(id)).
		Doc(data).
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 删除一条
func (e *Elasticsearch) DeleteRecord(index string, id int) error {
	_, err := e.Client.Delete().
		Index(index).Id(strconv.Itoa(id)).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 列表
func (e *Elasticsearch) RecordList(index string, keyword string) ([]int, error) {
	var res *elastic.SearchResult
	var err error
	ids := make([]int, 0)
	q := elastic.NewBoolQuery()
	q.Should(elastic.NewMatchQuery("title", keyword),
		elastic.NewMatchQuery("remark", keyword),
		elastic.NewMatchQuery("tags", keyword),
		elastic.NewMatchQuery("textData", keyword))
	res, err = e.Client.Search().Index(index).Query(q).Size(100).
		Sort("title", true).
		Sort("tags", true).
		Sort("remark", true).
		Sort("textData", true).
		Do(context.Background())
	if err != nil {
		return ids, err
	}
	for _, hit := range res.Hits.Hits {
		id, _ := strconv.Atoi(hit.Id)
		ids = append(ids, id)
	}
	return ids, nil
}

// 分页
func (e *Elasticsearch) PageRecord(index string, size int, page int, keyword string) (int64, []int, error) {
	var res *elastic.SearchResult
	var err error
	ids := make([]int, 0)
	q := elastic.NewBoolQuery()
	q.Should(elastic.NewMatchQuery("title", keyword),
		elastic.NewMatchQuery("remark", keyword),
		elastic.NewMatchQuery("tags", keyword),
		elastic.NewMatchQuery("textData", keyword))
	res, err = e.Client.Search().Index(index).Query(q).Size(size).From((page-1)*size).
		Sort("title", true).
		Sort("tags", true).
		Sort("remark", true).
		Sort("textData", true).Do(context.Background())
	if err != nil {
		return 0, ids, err
	}
	num := res.Hits.TotalHits.Value //搜索到结果总条数
	if num > 0 {
		for _, hit := range res.Hits.Hits {
			id, _ := strconv.Atoi(hit.Id)
			ids = append(ids, id)
		}
	}
	return num, ids, nil
}
