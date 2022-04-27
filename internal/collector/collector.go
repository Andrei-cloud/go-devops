package collector

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type Collector interface {
	Collect()
	CollectExtra()
	GetGauges() map[string]float64
	GetCounter() map[string]int64
}

type collector struct {
	counter int64
	gauges  map[string]float64
	mu      sync.RWMutex
}

var _ Collector = &collector{}

func NewCollector() *collector {
	c := &collector{}
	c.gauges = make(map[string]float64)
	return c
}

func (c *collector) setGauge(key string, value float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.gauges[key] = value
}

func (c *collector) Collect() {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	c.setGauge("Alloc", float64(m.Alloc))
	c.setGauge("BuckHashSys", float64(m.BuckHashSys))
	c.setGauge("Frees", float64(m.Frees))
	c.setGauge("GCCPUFraction", float64(m.GCCPUFraction))
	c.setGauge("GCSys", float64(m.GCSys))
	c.setGauge("HeapAlloc", float64(m.HeapAlloc))
	c.setGauge("HeapIdle", float64(m.HeapIdle))
	c.setGauge("HeapInuse", float64(m.HeapInuse))
	c.setGauge("HeapObjects", float64(m.HeapObjects))
	c.setGauge("HeapReleased", float64(m.HeapReleased))
	c.setGauge("HeapSys", float64(m.HeapSys))
	c.setGauge("LastGC", float64(m.LastGC))
	c.setGauge("Lookups", float64(m.Lookups))
	c.setGauge("MCacheInuse", float64(m.MCacheInuse))
	c.setGauge("MCacheSys", float64(m.MCacheSys))
	c.setGauge("MSpanInuse", float64(m.MSpanInuse))
	c.setGauge("MSpanSys", float64(m.MSpanSys))
	c.setGauge("Mallocs", float64(m.Mallocs))
	c.setGauge("NextGC", float64(m.NextGC))
	c.setGauge("NumForcedGC", float64(m.NumForcedGC))
	c.setGauge("NumGC", float64(m.NumGC))
	c.setGauge("OtherSys", float64(m.OtherSys))
	c.setGauge("PauseTotalNs", float64(m.PauseTotalNs))
	c.setGauge("StackInuse", float64(m.StackInuse))
	c.setGauge("StackSys", float64(m.StackSys))
	c.setGauge("Sys", float64(m.Sys))
	c.setGauge("TotalAlloc", float64(m.TotalAlloc))
	c.setGauge("RandomValue", randomValue())

	c.counter++
}

func (c *collector) CollectExtra() {
	var p []float64
	m, err := mem.VirtualMemory()
	if err == nil {
		c.setGauge("TotalMemory", float64(m.Total))
		c.setGauge("FreeMemory", float64(m.Free))
	}

	if p, err = cpu.Percent(0, false); err != nil {
		return
	}
	for i, v := range p {
		c.setGauge(fmt.Sprintf("CPUutilization%d", i+1), v)
	}
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
