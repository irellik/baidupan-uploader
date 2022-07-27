package flags

import "flag"

var (
	Concurrent int
	Rage       int
)

func InitFlag() {
	flag.IntVar(&Concurrent, "c", 10, "并发数")
	flag.IntVar(&Rage, "r", 256, "分片大小")
}
