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

package engine

import (
	"github.com/cgrates/cgrates/utils"
	"reflect"
	"testing"
	"time"
)

var (
	destinations = `
#Tag,Prefix
GERMANY,49
GERMANY_O2,41
GERMANY_PREMIUM,43
ALL,49
ALL,41
ALL,43
NAT,0256
NAT,0257
NAT,0723
RET,0723
RET,0724
`
	timings = `
WORKDAYS_00,*any,*any,*any,1;2;3;4;5,00:00:00
WORKDAYS_18,*any,*any,*any,1;2;3;4;5,18:00:00
WEEKENDS,*any,*any,*any,6;7,00:00:00
ONE_TIME_RUN,2012,,,,*asap
`
	rates = `
R1,0,0.2,60s,1s,0,*middle,2,10
R2,0,0.1,60s,1s,0,*middle,2,10
R3,0,0.05,60s,1s,0,*middle,2,10
R4,1,1,1s,1s,0,*up,2,10
R5,0,0.5,1s,1s,0,*down,2,10
`
	destinationRates = `
RT_STANDARD,GERMANY,R1
RT_STANDARD,GERMANY_O2,R2
RT_STANDARD,GERMANY_PREMIUM,R2
RT_DEFAULT,ALL,R2
RT_STD_WEEKEND,GERMANY,R2
RT_STD_WEEKEND,GERMANY_O2,R3
P1,NAT,R4
P2,NAT,R5
`
	destinationRateTimings = `
STANDARD,RT_STANDARD,WORKDAYS_00,10
STANDARD,RT_STD_WEEKEND,WORKDAYS_18,10
STANDARD,RT_STD_WEEKEND,WEEKENDS,10
PREMIUM,RT_STD_WEEKEND,WEEKENDS,10
DEFAULT,RT_DEFAULT,WORKDAYS_00,10
EVENING,P1,WORKDAYS_00,10
EVENING,P2,WORKDAYS_18,10
EVENING,P2,WEEKENDS,10
`
	ratingProfiles = `
CUSTOMER_1,0,*out,rif:from:tm,2012-01-01T00:00:00Z,PREMIUM,danb
CUSTOMER_1,0,*out,rif:from:tm,2012-02-28T00:00:00Z,STANDARD,danb
CUSTOMER_2,0,*out,danb:87.139.12.167,2012-01-01T00:00:00Z,STANDARD,danb
CUSTOMER_1,0,*out,danb,2012-01-01T00:00:00Z,PREMIUM,
vdf,0,*out,rif,2012-01-01T00:00:00Z,EVENING,
vdf,0,*out,rif,2012-02-28T00:00:00Z,EVENING,
vdf,0,*out,minu,2012-01-01T00:00:00Z,EVENING,
vdf,0,*out,*any,2012-02-28T00:00:00Z,EVENING,
vdf,0,*out,one,2012-02-28T00:00:00Z,STANDARD,
vdf,0,*out,inf,2012-02-28T00:00:00Z,STANDARD,inf
vdf,0,*out,fall,2012-02-28T00:00:00Z,PREMIUM,rif
`
	actions = `
MINI,*topup,*minutes,*out,100,*unlimited,NAT,*absolute,0,10,10
`
	actionTimings = `
MORE_MINUTES,MINI,ONE_TIME_RUN,10
`
	actionTriggers = `
STANDARD_TRIGGER,*minutes,*out,*min_counter,10,GERMANY_O2,SOME_1,10
STANDARD_TRIGGER,*minutes,*out,*max_balance,200,GERMANY,SOME_2,10
`
	accountActions = `
vdf,minitsboy,*out,MORE_MINUTES,STANDARD_TRIGGER
`
)

var csvr *CSVReader

func init() {
	csvr = NewStringCSVReader(storageGetter, ',', destinations, timings, rates, destinationRates, destinationRateTimings, ratingProfiles, actions, actionTimings, actionTriggers, accountActions)
	csvr.LoadDestinations()
	csvr.LoadTimings()
	csvr.LoadRates()
	csvr.LoadDestinationRates()
	csvr.LoadDestinationRateTimings()
	csvr.LoadRatingProfiles()
	csvr.LoadActions()
	csvr.LoadActionTimings()
	csvr.LoadActionTriggers()
	csvr.LoadAccountActions()
	csvr.WriteToDatabase(false, false)
}

func TestLoadDestinations(t *testing.T) {
	if len(csvr.destinations) != 6 {
		t.Error("Failed to load destinations: ", csvr.destinations)
	}
	for _, d := range csvr.destinations {
		switch d.Id {
		case "NAT":
			if !reflect.DeepEqual(d.Prefixes, []string{`0256`, `0257`, `0723`}) {
				t.Error("Faild to load destinations", d)
			}
		case "ALL":
			if !reflect.DeepEqual(d.Prefixes, []string{`49`, `41`, `43`}) {
				t.Error("Faild to load destinations", d)
			}
		case "RET":
			if !reflect.DeepEqual(d.Prefixes, []string{`0723`, `0724`}) {
				t.Error("Faild to load destinations", d)
			}
		case "GERMANY":
			if !reflect.DeepEqual(d.Prefixes, []string{`49`}) {
				t.Error("Faild to load destinations", d)
			}
		case "GERMANY_O2":
			if !reflect.DeepEqual(d.Prefixes, []string{`41`}) {
				t.Error("Faild to load destinations", d)
			}
		case "GERMANY_PREMIUM":
			if !reflect.DeepEqual(d.Prefixes, []string{`43`}) {
				t.Error("Faild to load destinations", d)
			}
		default:
			t.Error("Unknown destination tag!")
		}
	}
}

func TestLoadTimimgs(t *testing.T) {
	if len(csvr.timings) != 4 {
		t.Error("Failed to load timings: ", csvr.timings)
	}
	timing := csvr.timings["WORKDAYS_00"]
	if !reflect.DeepEqual(timing, &Timing{
		Id:        "WORKDAYS_00",
		Years:     Years{},
		Months:    Months{},
		MonthDays: MonthDays{},
		WeekDays:  WeekDays{1, 2, 3, 4, 5},
		StartTime: "00:00:00",
	}) {
		t.Error("Error loading timing: ", timing)
	}
	timing = csvr.timings["WORKDAYS_18"]
	if !reflect.DeepEqual(timing, &Timing{
		Id:        "WORKDAYS_18",
		Years:     Years{},
		Months:    Months{},
		MonthDays: MonthDays{},
		WeekDays:  WeekDays{1, 2, 3, 4, 5},
		StartTime: "18:00:00",
	}) {
		t.Error("Error loading timing: ", timing)
	}
	timing = csvr.timings["WEEKENDS"]
	if !reflect.DeepEqual(timing, &Timing{
		Id:        "WEEKENDS",
		Years:     Years{},
		Months:    Months{},
		MonthDays: MonthDays{},
		WeekDays:  WeekDays{time.Saturday, time.Sunday},
		StartTime: "00:00:00",
	}) {
		t.Error("Error loading timing: ", timing)
	}
	timing = csvr.timings["ONE_TIME_RUN"]
	if !reflect.DeepEqual(timing, &Timing{
		Id:        "ONE_TIME_RUN",
		Years:     Years{2012},
		Months:    Months{},
		MonthDays: MonthDays{},
		WeekDays:  WeekDays{},
		StartTime: "*asap",
	}) {
		t.Error("Error loading timing: ", timing)
	}
}

func TestLoadRates(t *testing.T) {
	if len(csvr.rates) != 5 {
		t.Error("Failed to load rates: ", csvr.rates)
	}
	rate := csvr.rates["R1"]
	if !reflect.DeepEqual(rate, &LoadRate{
		Tag:                "R1",
		ConnectFee:         0,
		Price:              0.2,
		RateUnit:           time.Minute,
		RateIncrement:      time.Second,
		GroupIntervalStart: 0,
		RoundingMethod:     utils.ROUNDING_MIDDLE,
		RoundingDecimals:   2,
		Weight:             10,
	}) {
		t.Error("Error loading rate: ", csvr.rates)
	}
	rate = csvr.rates["R2"]
	if !reflect.DeepEqual(rate, &LoadRate{
		Tag:                "R2",
		ConnectFee:         0,
		Price:              0.1,
		RateUnit:           time.Minute,
		RateIncrement:      time.Second,
		GroupIntervalStart: 0,
		RoundingMethod:     utils.ROUNDING_MIDDLE,
		RoundingDecimals:   2,
		Weight:             10,
	}) {
		t.Error("Error loading rate: ", csvr.rates)
	}
	rate = csvr.rates["R3"]
	if !reflect.DeepEqual(rate, &LoadRate{
		Tag:                "R3",
		ConnectFee:         0,
		Price:              0.05,
		RateUnit:           time.Minute,
		RateIncrement:      time.Second,
		GroupIntervalStart: 0,
		RoundingMethod:     utils.ROUNDING_MIDDLE,
		RoundingDecimals:   2,
		Weight:             10,
	}) {
		t.Error("Error loading rate: ", csvr.rates)
	}
	rate = csvr.rates["R4"]
	if !reflect.DeepEqual(rate, &LoadRate{
		Tag:                "R4",
		ConnectFee:         1,
		Price:              1.0,
		RateUnit:           time.Second,
		RateIncrement:      time.Second,
		GroupIntervalStart: 0,
		RoundingMethod:     utils.ROUNDING_UP,
		RoundingDecimals:   2,
		Weight:             10,
	}) {
		t.Error("Error loading rate: ", csvr.rates)
	}
	rate = csvr.rates["R5"]
	if !reflect.DeepEqual(rate, &LoadRate{
		Tag:                "R5",
		ConnectFee:         0,
		Price:              0.5,
		RateUnit:           time.Second,
		RateIncrement:      time.Second,
		GroupIntervalStart: 0,
		RoundingMethod:     utils.ROUNDING_DOWN,
		RoundingDecimals:   2,
		Weight:             10,
	}) {
		t.Error("Error loading rate: ", csvr.rates)
	}

}

func TestLoadDestinationRates(t *testing.T) {
	if len(csvr.destinationRates) != 5 {
		t.Error("Failed to load rates: ", csvr.rates)
	}
	drs := csvr.destinationRates["RT_STANDARD"]
	if !reflect.DeepEqual(drs, []*DestinationRate{
		&DestinationRate{
			Tag:             "RT_STANDARD",
			DestinationsTag: "GERMANY",
			Rate:            csvr.rates["R1"],
		},
		&DestinationRate{
			Tag:             "RT_STANDARD",
			DestinationsTag: "GERMANY_O2",
			Rate:            csvr.rates["R2"],
		},
		&DestinationRate{
			Tag:             "RT_STANDARD",
			DestinationsTag: "GERMANY_PREMIUM",
			Rate:            csvr.rates["R2"],
		},
	}) {
		t.Error("Error loading destination rate: ", drs)
	}
	drs = csvr.destinationRates["RT_DEFAULT"]
	if !reflect.DeepEqual(drs, []*DestinationRate{
		&DestinationRate{
			Tag:             "RT_DEFAULT",
			DestinationsTag: "ALL",
			Rate:            csvr.rates["R2"],
		},
	}) {
		t.Error("Error loading destination rate: ", drs)
	}
	drs = csvr.destinationRates["RT_STD_WEEKEND"]
	if !reflect.DeepEqual(drs, []*DestinationRate{
		&DestinationRate{
			Tag:             "RT_STD_WEEKEND",
			DestinationsTag: "GERMANY",
			Rate:            csvr.rates["R2"],
		},
		&DestinationRate{
			Tag:             "RT_STD_WEEKEND",
			DestinationsTag: "GERMANY_O2",
			Rate:            csvr.rates["R3"],
		},
	}) {
		t.Error("Error loading destination rate: ", drs)
	}
	drs = csvr.destinationRates["P1"]
	if !reflect.DeepEqual(drs, []*DestinationRate{
		&DestinationRate{
			Tag:             "P1",
			DestinationsTag: "NAT",
			Rate:            csvr.rates["R4"],
		},
	}) {
		t.Error("Error loading destination rate: ", drs)
	}
	drs = csvr.destinationRates["P2"]
	if !reflect.DeepEqual(drs, []*DestinationRate{
		&DestinationRate{
			Tag:             "P2",
			DestinationsTag: "NAT",
			Rate:            csvr.rates["R5"],
		},
	}) {
		t.Error("Error loading destination rate: ", drs)
	}
}

func TestLoadDestinationRateTimings(t *testing.T) {
	if len(csvr.activationPeriods) != 4 {
		t.Error("Failed to load rate timings: ", csvr.activationPeriods)
	}
	rplan := csvr.activationPeriods["STANDARD"]
	expected := &RatingPlan{
		ActivationTime: time.Time{},
		RateIntervals: RateIntervalList{
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{1, 2, 3, 4, 5},
				StartTime:  "00:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 0,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.2,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.1,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
				},
				RoundingMethod:   utils.ROUNDING_MIDDLE,
				RoundingDecimals: 2,
			},
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{1, 2, 3, 4, 5},
				StartTime:  "18:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 0,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.1,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.05,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
				},
				RoundingMethod:   utils.ROUNDING_MIDDLE,
				RoundingDecimals: 2,
			},
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{6, 0},
				StartTime:  "00:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 0,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.1,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.05,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
				},
				RoundingMethod:   utils.ROUNDING_MIDDLE,
				RoundingDecimals: 2,
			},
		},
	}
	if !reflect.DeepEqual(rplan, expected) {
		t.Errorf("Error loading rating plan: %#v", csvr.activationPeriods["STANDARD"])
	}
	rplan = csvr.activationPeriods["PREMIUM"]
	expected = &RatingPlan{
		ActivationTime: time.Time{},
		RateIntervals: RateIntervalList{
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{6, 0},
				StartTime:  "00:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 0,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.1,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.05,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
				},
				RoundingMethod:   utils.ROUNDING_MIDDLE,
				RoundingDecimals: 2,
			},
		},
	}
	if !reflect.DeepEqual(rplan, expected) {
		t.Errorf("Error loading rating plan: %#v", csvr.activationPeriods["PREMIUM"])
	}
	rplan = csvr.activationPeriods["DEFAULT"]
	expected = &RatingPlan{
		ActivationTime: time.Time{},
		RateIntervals: RateIntervalList{
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{1, 2, 3, 4, 5},
				StartTime:  "00:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 0,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.1,
						RateIncrement:      time.Second,
						RateUnit:           time.Minute,
					},
				},
				RoundingMethod:   utils.ROUNDING_MIDDLE,
				RoundingDecimals: 2,
			},
		},
	}
	if !reflect.DeepEqual(rplan, expected) {
		t.Errorf("Error loading rating plan: %#v", csvr.activationPeriods["DEFAULT"])
	}
	rplan = csvr.activationPeriods["EVENING"]
	expected = &RatingPlan{
		ActivationTime: time.Time{},
		RateIntervals: RateIntervalList{
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{1, 2, 3, 4, 5},
				StartTime:  "00:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 1,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              1,
						RateIncrement:      time.Second,
						RateUnit:           time.Second,
					},
				},
				RoundingMethod:   utils.ROUNDING_UP,
				RoundingDecimals: 2,
			},
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{1, 2, 3, 4, 5},
				StartTime:  "18:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 0,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.5,
						RateIncrement:      time.Second,
						RateUnit:           time.Second,
					},
				},
				RoundingMethod:   utils.ROUNDING_DOWN,
				RoundingDecimals: 2,
			},
			&RateInterval{
				Years:      Years{},
				Months:     Months{},
				MonthDays:  MonthDays{},
				WeekDays:   WeekDays{6, 0},
				StartTime:  "00:00:00",
				EndTime:    "",
				Weight:     10,
				ConnectFee: 0,
				Rates: RateGroups{
					&Rate{
						GroupIntervalStart: 0,
						Value:              0.5,
						RateIncrement:      time.Second,
						RateUnit:           time.Second,
					},
				},
				RoundingMethod:   utils.ROUNDING_DOWN,
				RoundingDecimals: 2,
			},
		},
	}
	if !reflect.DeepEqual(rplan, expected) {
		t.Errorf("Error loading rating plan: %#v", csvr.activationPeriods["EVENING"])
	}
}

func TestLoadRatingProfiles(t *testing.T) {
	if len(csvr.ratingProfiles) != 9 {
		t.Error("Failed to load rating profiles: ", len(csvr.ratingProfiles), csvr.ratingProfiles)
	}
	rp := csvr.ratingProfiles["*out:CUSTOMER_1:0:rif:from:tm"]
	expected := &RatingProfile{}
	if reflect.DeepEqual(rp, expected) {
		t.Error("Error loading rating profile: ", csvr.ratingProfiles["*out:CUSTOMER_1:0:rif:from:tm"])
	}
}

/*
CUSTOMER_1,0,*out,rif:from:tm,2012-01-01T00:00:00Z,PREMIUM,danb
CUSTOMER_1,0,*out,rif:from:tm,2012-02-28T00:00:00Z,STANDARD,danb
CUSTOMER_2,0,*out,danb:87.139.12.167,2012-01-01T00:00:00Z,STANDARD,danb
CUSTOMER_1,0,*out,danb,2012-01-01T00:00:00Z,PREMIUM,
vdf,0,*out,rif,2012-01-01T00:00:00Z,EVENING,
vdf,0,*out,rif,2012-02-28T00:00:00Z,EVENING,
vdf,0,*out,minu,2012-01-01T00:00:00Z,EVENING,
vdf,0,*out,*any,2012-02-28T00:00:00Z,EVENING,
vdf,0,*out,one,2012-02-28T00:00:00Z,STANDARD,
vdf,0,*out,inf,2012-02-28T00:00:00Z,STANDARD,inf
vdf,0,*out,fall,2012-02-28T00:00:00Z,PREMIUM,rif
*/

func TestLoadActions(t *testing.T) {
	if len(csvr.actions) != 1 {
		t.Error("Failed to load actions: ", csvr.actions)
	}
	a := csvr.actions["MINI"][0]
	expected := &Action{
		Id:               a.Id,
		ActionType:       TOPUP,
		BalanceId:        MINUTES,
		Direction:        OUTBOUND,
		ExpirationString: UNLIMITED,
		Weight:           10,
		Balance: &Balance{
			Id:               a.Balance.Id,
			Value:            100,
			Weight:           10,
			SpecialPriceType: PRICE_ABSOLUTE,
			SpecialPrice:     0,
			DestinationId:    "NAT",
		},
	}
	if !reflect.DeepEqual(a, expected) {
		t.Error("Error loading action: ", a)
	}
}

func TestLoadActionTimings(t *testing.T) {
	if len(csvr.actionsTimings) != 1 {
		t.Error("Failed to load action timings: ", csvr.actionsTimings)
	}
	atm := csvr.actionsTimings["MORE_MINUTES"][0]
	expected := &ActionTiming{
		Id:             atm.Id,
		Tag:            "ONE_TIME_RUN",
		UserBalanceIds: []string{"*out:vdf:minitsboy"},
		Timing: &RateInterval{
			Years:     Years{2012},
			Months:    Months{},
			MonthDays: MonthDays{},
			WeekDays:  WeekDays{},
			StartTime: ASAP,
		},
		Weight:    10,
		ActionsId: "MINI",
	}
	if !reflect.DeepEqual(atm, expected) {
		t.Error("Error loading action timing: ", atm)
	}
}

func TestLoadActionTriggers(t *testing.T) {
	if len(csvr.actionsTriggers) != 1 {
		t.Error("Failed to load action triggers: ", csvr.actionsTriggers)
	}
	atr := csvr.actionsTriggers["STANDARD_TRIGGER"][0]
	expected := &ActionTrigger{
		Id:             atr.Id,
		BalanceId:      MINUTES,
		Direction:      OUTBOUND,
		ThresholdType:  TRIGGER_MIN_COUNTER,
		ThresholdValue: 10,
		DestinationId:  "GERMANY_O2",
		Weight:         10,
		ActionsId:      "SOME_1",
		Executed:       false,
	}
	if !reflect.DeepEqual(atr, expected) {
		t.Error("Error loading action trigger: ", atr)
	}
	atr = csvr.actionsTriggers["STANDARD_TRIGGER"][1]
	expected = &ActionTrigger{
		Id:             atr.Id,
		BalanceId:      MINUTES,
		Direction:      OUTBOUND,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 200,
		DestinationId:  "GERMANY",
		Weight:         10,
		ActionsId:      "SOME_2",
		Executed:       false,
	}
	if !reflect.DeepEqual(atr, expected) {
		t.Error("Error loading action trigger: ", atr)
	}
}

func TestLoadAccountActions(t *testing.T) {
	if len(csvr.accountActions) != 1 {
		t.Error("Failed to load account actions: ", csvr.accountActions)
	}
	aa := csvr.accountActions[0]
	expected := &UserBalance{
		Id:             "*out:vdf:minitsboy",
		Type:           UB_TYPE_PREPAID,
		ActionTriggers: csvr.actionsTriggers["STANDARD_TRIGGER"],
	}
	if !reflect.DeepEqual(aa, expected) {
		t.Error("Error loading account action: ", aa)
	}
}

/*
vdf,minitsboy,*out,MORE_MINUTES,STANDARD_TRIGGER
*/
