package httpobj

// BaseTrackData rpc/kafka 数据格式
type BaseTrackData struct {
	Id string `json:"id"`
	TrackNumber string `json:"n"`
	Lang string `json:"lang"`
	UserId string `json:"userid"`
}


// TrackData 内部传输数据结构
type TrackData []*BaseTrackData

