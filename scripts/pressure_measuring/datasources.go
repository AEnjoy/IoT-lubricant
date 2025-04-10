package main

import (
	"math/rand"
	"time"
)

type Data struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`

	Pressure    float32 `json:"pressure"`
	Temperature float32 `json:"temperature"`
	Humidity    float32 `json:"humidity"`

	DisplacementX float32 `json:"displacementX"`
	DisplacementY float32 `json:"displacementY"`
	DisplacementZ float32 `json:"displacementZ"`

	LastOpenTime time.Time `json:"lastOpenTime"`
}

const (
	positionDelta    = 0.00001
	measurementDelta = 0.1
	pressureMin      = 920.0
	pressureMax      = 980.0
	temperatureMin   = -20.0
	temperatureMax   = 1.0
	displacementMax  = 1.5
	humidityMax      = 100.0
)

type dataGenerator struct {
	r                  *rand.Rand
	currentLongitude   float64
	currentLatitude    float64
	currentPressure    float32
	currentTemperature float32
	currentHumidity    float32
	currentDisX        float32
	currentDisY        float32
	currentDisZ        float32
}

func NewDataGenerator() *dataGenerator {
	return &dataGenerator{
		r:                  rand.New(rand.NewSource(time.Now().UnixNano())),
		currentLongitude:   -180 + rand.Float64()*360,
		currentLatitude:    -90 + rand.Float64()*180,
		currentPressure:    pressureMin + rand.Float32()*(pressureMax-pressureMin),
		currentTemperature: temperatureMin + rand.Float32()*(temperatureMax-temperatureMin),
		currentHumidity:    rand.Float32() * humidityMax,
		currentDisX:        rand.Float32() * displacementMax,
		currentDisY:        rand.Float32() * displacementMax,
		currentDisZ:        rand.Float32() * displacementMax,
	}
}

var startTime, dataTime time.Time

func (g *dataGenerator) generate() *Data {
	g.currentLongitude += (g.r.Float64()*2 - 1) * positionDelta
	g.currentLongitude = clampFloat64(g.currentLongitude, -180, 180)

	g.currentLatitude += (g.r.Float64()*2 - 1) * positionDelta
	g.currentLatitude = clampFloat64(g.currentLatitude, -90, 90)

	g.currentPressure += (g.r.Float32()*2 - 1) * measurementDelta
	g.currentPressure = clampFloat32(g.currentPressure, pressureMin, pressureMax)

	g.currentTemperature += (g.r.Float32()*2 - 1) * measurementDelta
	g.currentTemperature = clampFloat32(g.currentTemperature, temperatureMin, temperatureMax)

	g.currentHumidity += (g.r.Float32()*2 - 1) * measurementDelta
	g.currentHumidity = clampFloat32(g.currentHumidity, 0, humidityMax)

	g.currentDisX += (g.r.Float32()*2 - 1) * measurementDelta
	g.currentDisX = clampFloat32(g.currentDisX, 0, displacementMax)
	g.currentDisY += (g.r.Float32()*2 - 1) * measurementDelta
	g.currentDisY = clampFloat32(g.currentDisY, 0, displacementMax)
	g.currentDisZ += (g.r.Float32()*2 - 1) * measurementDelta
	g.currentDisZ = clampFloat32(g.currentDisZ, 0, displacementMax)

	return &Data{
		Longitude:     g.currentLongitude,
		Latitude:      g.currentLatitude,
		Pressure:      g.currentPressure,
		Temperature:   g.currentTemperature,
		Humidity:      g.currentHumidity,
		DisplacementX: g.currentDisX,
		DisplacementY: g.currentDisY,
		DisplacementZ: g.currentDisZ,
		LastOpenTime:  dataTime,
	}
}

func clampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampFloat32(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

var dataCh chan *Data

func StartConcurrentGeneration(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			<-startSig

			generator := NewDataGenerator()
			for {
				dataCh <- generator.generate()
			}
		}()
	}
}
