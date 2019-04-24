package util

var (
    Math = &math{}
)

type math struct{}

func (p *math) AbsInt64(x int64) int64 {
    if x >= 0 {
        return x
    }
    return -x
}
