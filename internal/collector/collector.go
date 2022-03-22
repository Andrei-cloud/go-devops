package collector

import (
	"math/rand"
	"runtime"
	"time"
)

type Collector interface {
	Collect()
	GetGauges() map[string]float64
	GetCounter() map[string]int64
}

type collector struct {
	counter int64
	gauges  map[string]float64
}

var _ Collector = &collector{}

func NewCollector() *collector {
	c := &collector{}
	c.gauges = make(map[string]float64)
	return c
}

func (c *collector) Collect() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	c.gauges["Alloc"] = float64(m.Alloc)
	c.gauges["BuckHashSys"] = float64(m.BuckHashSys)
	c.gauges["Frees"] = float64(m.Frees)
	c.gauges["GCCPUFraction"] = float64(m.GCCPUFraction)
	c.gauges["GCSys"] = float64(m.GCSys)
	c.gauges["HeapAlloc"] = float64(m.HeapAlloc)
	c.gauges["HeapIdle"] = float64(m.HeapIdle)
	c.gauges["HeapInuse"] = float64(m.HeapInuse)
	c.gauges["HeapObjects"] = float64(m.HeapObjects)
	c.gauges["HeapReleased"] = float64(m.HeapReleased)
	c.gauges["HeapSys"] = float64(m.HeapSys)
	c.gauges["LastGC"] = float64(m.LastGC)
	c.gauges["Lookups"] = float64(m.Lookups)
	c.gauges["MCacheInuse"] = float64(m.MCacheInuse)
	c.gauges["MCacheSys"] = float64(m.MCacheSys)
	c.gauges["MSpanInuse"] = float64(m.MSpanInuse)
	c.gauges["MSpanSys"] = float64(m.MSpanSys)
	c.gauges["Mallocs"] = float64(m.Mallocs)
	c.gauges["NextGC"] = float64(m.NextGC)
	c.gauges["NumForcedGC"] = float64(m.NumForcedGC)
	c.gauges["NumGC"] = float64(m.NumGC)
	c.gauges["OtherSys"] = float64(m.OtherSys)
	c.gauges["PauseTotalNs"] = float64(m.PauseTotalNs)
	c.gauges["StackInuse"] = float64(m.StackInuse)
	c.gauges["StackSys"] = float64(m.StackSys)
	c.gauges["Sys"] = float64(m.Sys)
	c.gauges["TotalAlloc"] = float64(m.TotalAlloc)
	c.gauges["RandomValue"] = randomValue()

	c.counter++
}

func (c *collector) GetGauges() map[string]float64 {
	return c.gauges
}

func (c *collector) GetCounter() map[string]int64 {
	return map[string]int64{
		"PollCount": c.counter,
	}
}

func randomValue() float64 {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	return r.Float64()
}
