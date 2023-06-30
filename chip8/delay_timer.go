package chip8

import "time"

type DelayTimer struct {
	timer byte
}

func NewDelayTimer() *DelayTimer {
	dt := new(DelayTimer)
	dt.timer = 0
	return dt
}

func (dt *DelayTimer) Start() {
	go dt.decrementTimer()
}

func (dt *DelayTimer) decrementTimer() {
	for {
		if dt.timer > 0 {
			dt.timer--
			println("Decrement timer, timer is ", dt.timer)
		}
		//time.Sleep(time.Second * 1)
		time.Sleep(time.Microsecond * 16667)
	}
}

func (dt *DelayTimer) setTimer(b byte) {
	dt.timer = b
}
