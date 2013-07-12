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

// This file deals with tp_* data definition

type TPRate struct {
	TPid      string     // Tariff plan id
	RateId    string     // Rate id
	RateSlots []RateSlot // One or more RateSlots
}

type RateSlot struct {
	ConnectFee     float64 // ConnectFee applied once the call is answered
	Rate           float64 // Rate applied
	RatedUnits     int     //  Number of billing units this rate applies to
	RateIncrements int     // This rate will apply in increments of duration
	Weight         float64 // Rate's priority when dealing with grouped rates
}

type TPDestinationRate struct {
	TPid              string            // Tariff plan id
	DestinationRateId string            // DestinationRate profile id
	DestinationRates  []DestinationRate // Set of destinationid-rateid bindings
}

type DestinationRate struct {
	DestinationId string // The destination identity
	RateId        string // The rate identity
}

type TPDestRateTiming struct {
	TPid             string           // Tariff plan id
	DestRateTimingId string           // DestinationRate profile id
	DestRateTimings  []DestRateTiming // Set of destinationid-rateid bindings
}

type DestRateTiming struct {
	DestRatesId string  // The DestinationRate identity
	TimingId    string  // The timing identity
	Weight      float64 // Binding priority taken into consideration when more DestinationRates are active on a time slot
}

type TPRateProfile struct {
	TPid                   string                  // Tariff plan id
	RateProfileId          string                  // RateProfile id
	Tenant                 string                  // Tenant's Id
	TOR                    string                  // TypeOfRecord
	Direction              string                  // Traffic direction, OUT is the only one supported for now
	Subject                string                  // Rating subject, usually the same as account
	RatesFallbackSubject   string                  // Fallback on this subject if rates not found for destination
	RatingActivations []RatingActivation // Activate rate profiles at specific time
}

type RatingActivation struct {
	ActivationTime   int64 // Time when this profile will become active, defined as unix epoch time
	DestRateTimingId string // Id of DestRateTiming profile
}

type AttrTPRateProfileIds struct {
	TPid string // Tariff plan id
	Tenant string // Tenant's Id
	TOR   string  // TypeOfRecord
	Direction string // Traffic direction
	Subject   string // Rating subject, usually the same as account
}
