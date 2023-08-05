package process

type endCondition interface {
	end() (bool, int64)
}

var (
	_ endCondition = (*endConditionTime)(nil)
	_ endCondition = (*endConditionTemperature)(nil)
)

type getTemperature func() float64
type getTime func() int64

type endConditionTime struct {
	getTime  getTime
	start    int64
	duration int64
	leftTime int64
}

func newEndConditionTime(duration int64, time getTime) *endConditionTime {
	e := &endConditionTime{
		getTime:  time,
		duration: duration,
		leftTime: duration,
		start:    time(),
	}

	e.reset()
	return e
}

func (e *endConditionTime) end() (bool, int64) {
	left := e.left()
	return left == 0, left
}

func (e *endConditionTime) reset() {
	e.start = e.getTime()
	e.leftTime = e.duration
}

func (e *endConditionTime) left() int64 {
	t := (e.start + e.leftTime) - e.getTime()
	if t < 0 {
		return 0
	}
	return t
}

type endConditionTemperature struct {
	waiting        bool
	endTime        *endConditionTime
	threshold      float64
	getTemperature getTemperature
}

func newEndConditionTemperature(duration int64, time getTime, threshold float64, temperature getTemperature) *endConditionTemperature {
	return &endConditionTemperature{
		waiting:        false,
		endTime:        newEndConditionTime(duration, time),
		threshold:      threshold,
		getTemperature: temperature,
	}

}

func (e *endConditionTemperature) left() int64 {
	if e.waiting {
		return e.endTime.left()
	}
	return e.endTime.duration
}

func (e *endConditionTemperature) end() (bool, int64) {
	overThreshold := e.getTemperature() > e.threshold
	// Did anything change since last call?
	if overThreshold != e.waiting {
		if overThreshold {
			// Okay, we are now over threshold
			// Start waiting for time
			e.waiting = true
			e.endTime.reset()
		} else {
			// Sadly, we are below threshold
			e.waiting = false
		}
	}
	if overThreshold {
		return e.endTime.end()
	}

	return false, e.endTime.duration
}
