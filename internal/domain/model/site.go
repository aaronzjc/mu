package model

type CrawType int

const (
	CrawHtml CrawType = 1 // 网站是HTML
	CrawApi  CrawType = 2 // 网站是JSON接口

	ByType  = 1 // 通过服务器类型
	ByHosts = 2 // 服务器IPs

	Disable = 0 // 禁用
	Enable  = 1 // 启用
)

type Site struct {
	ID         int    `gorm:"primaryKey"`
	Name       string `gorm:"name"`
	Root       string `gorm:"root"`
	Key        string `gorm:"key"`
	Desc       string `gorm:"desc"`
	Type       int8   `gorm:"type"`
	Tags       string `gorm:"tags"`
	Cron       string `gorm:"cron"`
	Enable     int8   `gorm:"enable"`
	NodeOption int8   `gorm:"node_option"`
	NodeType   int8   `gorm:"node_type"`
	NodeHosts  string `gorm:"node_hosts"`
	ReqHeaders string `gorm:"req_headers"`
}

func (s *Site) TableName() string {
	return "site"
}

type NewsItem struct {
	Key       string            `json:"key"`
	Title     string            `json:"title"`
	Desc      string            `json:"desc"`
	Rank      float64           `json:"rank"`
	OriginUrl string            `json:"origin_url"`
	Card      uint8             `json:"card_type"`
	Ext       map[string]string `json:"ext"`
}
type News struct {
	T    string
	List []NewsItem
}
