package httpobj

type Item struct {
	ItemStatus *ItemStatus
	Item       *TrackItem
}

type ItemStatus struct {
	RequestId string
	IsError   bool
}

type KafkaItem struct {
	TrackNumber string    `json:"n"`
	Express     string    `json:"e"`
	UserId      string    `json:"u"`
	TrackItem   TrackItem `json:"d"`
}

type TrackInfo struct {
	Date              string `json:"Date"`              // 时间        		| 可选
	StatusDescription string `json:"StatusDescription"` // 轨迹内容			| 可选
	Details           string `json:"Details"`           // 轨迹内容地点     	| 可选
	CheckpointStatus  string `json:"checkpoint_status"` // 状态节点 状态时结果
	SubStatus         string `json:"substatus"`         // 子状态 | 归属状态节点
	SubStatusNum      int    `json:"substatus_num"`     // 子状态 | 归属状态节点
	Keywords          string `json:"keywords"`          // 状态关键词
	Type              int    `json:"Type"`              //
}

type TrackItem struct {
	CountryName        string       `json:"countryname"`        // 发件国 | 可选
	DestinationCountry string       `json:"destinationcountry"` // 目的国 | 可选
	TrackInfo          []*TrackInfo `json:"trackinfo"`          // 轨迹数据

	// 状态结果,主状态; 效果一致, 优先选择 StatusNum
	StatusDataNum int `json:"stausDataNum"` // 状态结果，主状态
	StatusNum     int `json:"statusNum"`    // *
	MStatusNum    int `json:"mstatusNum"`   //

	StatusInfo            string `json:"statusInfo"`            // 最新轨迹数据拼接(StatusDescription,Details,Date)
	LastUpdateTime        string `json:"lastUpdateTime"`        // 最后更新时间
	FirstUpdateTime       string `json:"firstUpdateTime"`       // 第一次更新时间
	ItemTimeLength        string `json:"itemTimeLength"`        // 运输时长
	StayTimeLength        int    `json:"stayTimeLength"`        // 最新轨迹未更新时间 | 停留时间
	ArrivalFromAbroadTime string `json:"ArrivalfromAbroad"`     // 刚到目的国的时间
	DestinationArrived    string `json:"DestinationArrived"`    // 目的国到达待取时间（-）
	CustomsClearance      string `json:"CustomsClearance"`      // 到达海关时间
	ItemReceived          string `json:"item_received"`         // 收件时间|上网时间
	SubStatus             string `json:"substatus"`             // 总状态的子状态
	SubStatusTime         int    `json:"substatusTime"`         // 总状态的子状态时间
	Function              string `json:"function"`              // 记录线路查询结果
	Lang                  string `json:"lang"`                  // 语言
	NeedReTrack           bool   `json:"need_retrack"`          // 是否重查
	ScheduledDeliveryDate string `json:"ScheduledDeliveryDate"` // 预计到达时间

	// ==============
	LogFunc       string `json:"log_func"`    // _
	Weblink       string `json:"weblink"`     // 官网地址
	ShortCode     string `json:"campanyCode"` // 简码
	Support       int    `json:"support"`     // 是否开放查询
	ParcelType    string `json:"parcelType"`  // 包裹类型
	TrackType     string `json:"trackType"`   //
	TrackLang     string `json:"tracklang"`
	ProxyIsExceed bool   `json:"proxy_is_exceed"`
	City          string `json:"city"`
	State         string `json:"state"`
	CountryCode   string `json:"country_code"` // 国家简码
	ZipCode       string `json:"zipcode"`      // 邮编
}

func NewTrackItem() *TrackItem {
	return &TrackItem{}
}
