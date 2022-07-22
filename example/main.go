package main

import (
	"github.com/cchen-byte/trackeSharkes/async"
)

func main() {
	//httpBinLogic := trackerLogic.NewLogic()
	//httpBinLogic := yunexpress.NewYunExpressLogic()
	async.RunAsync()
}
