package semaphore

import "testing"

/*
*
sync.Once用于单次初始化的场景，
singleFlight每次调用都重新执行，用于合并并发请求的场景，可以解决缓存击穿的问题或者是一些幂等性的并发查询问题
缓存击穿：大量请求同时查询一个key，但是这个key正好过期失效了，就会导致大量请求打到数据库上，
*/
func Test1(t *testing.T) {

}
