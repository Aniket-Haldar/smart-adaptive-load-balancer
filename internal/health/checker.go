package health

import(
	"log"
	"net"
	"time"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

type Checker struct{
	Backends []*core.Backend
	Interval time.Duration
	Timeout time.Duration
}

//tcp level hea;thcheck (for l4 and l7)
func(c *Checker)Start(){
	ticker:=time.NewTicker(c.Interval)

	go func(){
		for range ticker.C{
			for _,backend:=range c.Backends{
				go c.checkBackend(backend)
			}
		}
	}()
}

func (c *Checker) checkBackend(b *core.Backend){
	start:=time.Now()

	conn,err:=net.DialTimeout("tcp",b.Address,c.Timeout)
	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if err!=nil{
		b.Alive=false
		log.Println("backend down:",b.Address)
		return
	}

	_=conn.Close()
	b.Alive=true
	b.Latency=time.Since(start)
}
