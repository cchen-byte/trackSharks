package yunexpress

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

type Logic struct {
	flag string
	BatchSize int
}

func (logic *Logic) ConstructFirstRequest (reqMeta *httpobj.RequestMeta, trackData *httpobj.TrackData) *httpobj.Request {
	//fmt.Println(logic.flag)
	return api5ConstructFirstRequest(reqMeta, trackData)
	//return api2ConstructFirstRequest(trackData)
}

func NewYunExpressLogic(flag string) *Logic {
	return &Logic{
		flag: flag,
		BatchSize: 50,
	}
}