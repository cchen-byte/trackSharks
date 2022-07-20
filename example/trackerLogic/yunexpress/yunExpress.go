package yunexpress

import "github.com/cchen-byte/trackeSharkes/httpobj"

type YunExpressLogic struct {

}

func (logic *YunExpressLogic) ConstructFirstRequest (trackData httpobj.TrackData) *httpobj.Request {
	return api5ConstructFirstRequest(trackData)
	//return api2ConstructFirstRequest(trackData)
}

func NewYunExpressLogic() *YunExpressLogic {
	return &YunExpressLogic{}
}