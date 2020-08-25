package parallel


type Parallel struct {
	pool chan bool
	done chan bool
}

func (p *Parallel) Init(max int) {
	p.pool = make(chan bool, max)
	p.done = make(chan bool, max)
	go func() {
		for {
			<-p.done
			<-p.pool
		}
	}()
}

func (p *Parallel) Do() {
	p.pool <- true
}
func (p *Parallel) Done() {
	p.done <- true
}
