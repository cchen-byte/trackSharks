package main

import (
	"fmt"
	"github.com/cchen-byte/trackeSharkes/pkg/traceId"
	"time"
)

func main() {
	t := traceId.TraceId{}
	for i:=1000; i<9000; i++ {
		go func(i int) {
			time.Sleep(time.Second)
			fmt.Println(i, t.GetTraceId())
		}(i)
	}
	time.Sleep(time.Second*21)

}
