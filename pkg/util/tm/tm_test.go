package tm

import (
	"fmt"
	"testing"
	"time"
)

func Test_TruncMinuteTs(t *testing.T) {
	sec := time.Now().Unix()
	min := sec / 60 * 60
	fmt.Printf("sec:%d => min:%d\n", sec, min)
}
