package main

import (
	"github.com/cchen-byte/trackeSharkes/async"
	"github.com/cchen-byte/trackeSharkes/example/trackerLogic/yunexpress"
)

func main() {
	//httpBinLogic := trackerLogic.NewLogic()
	httpBinLogic := yunexpress.NewYunExpressLogic()
	async.RunAsync(httpBinLogic)
}
