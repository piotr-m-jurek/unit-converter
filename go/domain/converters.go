package domain

import (
	"fmt"
	"strconv"
)

type Converter interface {
	Compute(value string, unitFrom, unitTo string) (float64, error)
}

type LengthConverter struct{}
type TemperatureConverter struct{}
type WeightConverter struct{}

func (c LengthConverter) Compute(length, unitFrom, unitTo string) (float64, error) {
	val, err := strconv.ParseFloat(length, 64)
	if err != nil {
		return 0, err
	}

	meters, err := c.getMeters(val, unitFrom)
	if err != nil {
		return 0, err
	}

	var output float64

	switch unitTo {
	case "cm":
		output = meters * 100
	case "m":
		output = meters
	case "km":
		output = meters / 1000
	}

	return output, nil
}

func (c LengthConverter) getMeters(length float64, unitFrom string) (float64, error) {
	switch unitFrom {
	case "cm":
		return length / 100, nil
	case "m":
		return length, nil
	case "km":
		return length * 1000, nil
	default:
		return 0, fmt.Errorf("invalid unit: %s", unitFrom)
	}
}

func (c WeightConverter) Compute(weight, unitFrom, unitTo string) (float64, error) {
	val, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		return 0, err
	}
	grams, err := c.getGrams(val, unitFrom)
	if err != nil {
		return 0, err
	}

	var output float64
	switch unitTo {
	case "g":
		output = grams
	case "kg":
		output = grams / 1000
	case "lb":
		output = grams * 2.20462
	}
	return output, nil
}

func (c TemperatureConverter) Compute(temperature, unitFrom, unitTo string) (float64, error) {
	temp, err := strconv.ParseFloat(temperature, 64)
	if err != nil {
		return 0, err
	}
	kelvins, err := c.getKelvins(temp, unitFrom)
	if err != nil {
		return 0, err
	}

	var output float64
	switch unitTo {
	case "c":
		output = kelvins - 273.15
	case "f":
		output = kelvins*9/5 - 459.67
	}

	return output, nil
}

func (c TemperatureConverter) getKelvins(temperature float64, unitFrom string) (float64, error) {
	switch unitFrom {
	case "c":
		return temperature + 273.15, nil
	case "f":
		return (temperature + 459.67) * 5 / 9, nil
	default:
		return 0, fmt.Errorf("invalid unit: %s", unitFrom)
	}
}

func (c WeightConverter) getGrams(weight float64, unitFrom string) (float64, error) {
	switch unitFrom {
	case "g":
		return weight, nil
	case "kg":
		return weight * 1000, nil
	case "lb":
		return weight * 453.592, nil
	default:
		return 0, fmt.Errorf("invalid unit: %s", unitFrom)
	}
}
