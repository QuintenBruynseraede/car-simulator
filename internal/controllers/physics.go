package controllers

import (
	"context"
	"time"

	"github.com/user/car-simulator/internal/events"
	"github.com/user/car-simulator/internal/storage"
	"go.uber.org/zap"
)

const (
	gravity        = 9.81
	roadResistance = 0.015
	airDensity     = 1.225
)

type Car struct {
	tuningConstant        float64 // Increase in rpm per second that the gas pedal is pressed
	maxRPM                float64
	idleRPM               float64
	mass                  float64 // kg
	dragCoef              float64 //
	frontalArea           float64 // m^2
	tireRadius            float64 // m
	finalDrive            float64
	gearRatio             float64
	shiftUpThresholdRPM   float64
	shiftDownThresholdRPM float64
	brakeStrength         float64 // Decelaration in m/s^2 when braking
	// Dynamic
	velocity     float64
	rpm          float64
	throttle     float64
	gear         int
	lastUpdate   time.Time
	brakePressed bool
}

type PhysicsController struct{}

func NewCar() *Car {
	return &Car{
		tuningConstant:        1500,
		maxRPM:                8000,
		idleRPM:               800,
		mass:                  1500,
		dragCoef:              0.35,
		frontalArea:           2.0,
		tireRadius:            0.4,
		finalDrive:            3.5,
		gearRatio:             4.0,
		shiftUpThresholdRPM:   3500,
		shiftDownThresholdRPM: 1500,
		brakeStrength:         20,
		gear:                  1,
		lastUpdate:            time.Now(),
	}
}

func NewPhysicsController() *PhysicsController {
	return &PhysicsController{}
}

func (c *Car) Update(td float64) {
	c.rpm += c.tuningConstant * td * c.throttle
	if c.rpm > c.maxRPM {
		c.rpm = c.maxRPM
	}
	if c.rpm < c.idleRPM {
		c.rpm = c.idleRPM
	}

	speedMS := c.velocity / 3.6 // m/2
	engineForce := (c.rpmToTorque() * c.gearRatio * c.finalDrive) / c.tireRadius

	// Calculate drag force
	dragForce := 0.5 * c.dragCoef * c.frontalArea * airDensity * speedMS * speedMS
	rollingResistance := 0.015 * c.mass * gravity
	netForce := engineForce*c.throttle - dragForce - rollingResistance
	acceleration := netForce / c.mass

	if c.brakePressed {
		acceleration -= c.brakeStrength
	}

	speedMS += acceleration * td
	c.velocity = speedMS * 3.6 // Convert back to km/h
	if c.velocity < 0 {
		c.velocity = 0
	}

	// Shift if needed
	if c.rpm > c.shiftUpThresholdRPM && c.gear < 5 {
		c.gear += 1
		c.rpm -= c.shiftUpThresholdRPM - c.shiftDownThresholdRPM - 100
	} else if c.rpm < c.shiftDownThresholdRPM && c.gear > 1 {
		c.gear -= 1
		c.rpm += 1000
	}
}

func (c *Car) rollingResistance() float64 {
	return roadResistance * c.mass * gravity
}

func (c *Car) dragForce() float64 {
	return 0.5 * c.dragCoef * c.frontalArea * airDensity * c.velocity * c.velocity
}

func (c *Car) rpmToTorque() float64 {
	// Define a simple torque curve
	switch {
	case c.rpm < 1000:
		return 50
	case c.rpm < 3000:
		return 150 + (c.rpm-1000)*0.1
	case c.rpm < 5000:
		return 200
	default:
		return 200 - (c.rpm-5000)*0.1
	}
}

func (c *Car) GearRatio() float64 {
	switch c.gear {
	case 1:
		return 3.8
	case 2:
		return 2.1
	case 3:
		return 1.4
	case 4:
		return 1.0
	case 5:
		return 0.8
	default:
		return 0.7
	}
}

func (p *PhysicsController) Run(_ context.Context, kvs storage.StorageBackend, eventBus *events.EventBus, logger *zap.Logger) error {
	car := NewCar()

	gas_pressed := make(chan events.Event)
	eventBus.Subscribe(events.EventGasPedalPressed, gas_pressed)

	gas_released := make(chan events.Event)
	eventBus.Subscribe(events.EventGasPedalReleased, gas_released)

	brake_pressed := make(chan events.Event)
	eventBus.Subscribe(events.EventBrakePedalPressed, brake_pressed)

	brake_released := make(chan events.Event)
	eventBus.Subscribe(events.EventBrakePedalReleased, brake_released)

	for {
		select {
		case <-gas_pressed:
			car.throttle = 1
		case <-gas_released:
			if kvs.ReadString(KeyEngineOn) == "true" {
				car.throttle = -0.3
			} else {
				car.throttle = -2
			}
		case <-brake_pressed:
			car.brakePressed = true
		case <-brake_released:
			car.brakePressed = false
		default:
			td := time.Since(car.lastUpdate).Seconds()
			car.Update(td)
			car.lastUpdate = time.Now()
			kvs.Write(KeyRPM, car.rpm)
			kvs.Write(KeyVelocity, car.velocity)
			kvs.Write(KeyGear, car.gear)
		}
	}
}
