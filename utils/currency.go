/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2013 ITsysCOM

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/
package utils

import "fmt"

const (
	Subunits = 1000           // cents
	Units    = 100 * Subunits // dollars
)

// The base unit is in 1/1000 of a subunit.
type Currency struct {
	value  int64
	symbol string
}

func (c *Currency) Units() float64 {
	return float64(c.value) / Units
}

func (c *Currency) Subunits() float64 {
	return float64(c.value) / Subunits
}

// ToCurrency creates a Currency from a float (e.g. 2.56 -> 2*Units + 56*Cents)
func ToCurrency(d float64, symbol string) *Currency {
	return &Currency{value: int64(d * Units), symbol: symbol}
}

func (c *Currency) String() string {
	return fmt.Sprintf("%.2f%s", c.Units(), c.symbol)
}

func (c *Currency) Add(other *Currency) {
	c.value += other.value
}

func (c *Currency) AddFloat64(v float64) {
	c.value += ToCurrency(v, "").value
}
