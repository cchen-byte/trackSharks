package main

import (
	"github.com/cchen-byte/trackeSharkes/async"
	"github.com/cchen-byte/trackeSharkes/example/trackerLogic"
)

func main() {
	httpBinLogic := trackerLogic.NewLogic()
	async.RunAsync(httpBinLogic)
}
