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
	"fmt"

	"github.com/cgrates/cgrates/utils"

	"testing"
	"time"
)

var (
	//referenceDate = time.Date(2014, 1, 1, 0, 0, 0, 0, time.Local)
	//referenceDate = time.Date(2013, 7, 10, 10, 30, 0, 0, time.Local)
	//referenceDate = time.Date(2013, 12, 31, 23, 59, 59, 0, time.Local)
	referenceDate = time.Now()
	now           = referenceDate
)

func TestActionTimingNothing(t *testing.T) {
	at := &ActionTiming{}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingOnlyHour(t *testing.T) {
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{StartTime: "10:01:00"}}}
	st := at.GetNextStartTime(referenceDate)

	y, m, d := now.Date()
	expected := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingOnlyWeekdays(t *testing.T) {
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{WeekDays: []time.Weekday{time.Monday}}}}
	st := at.GetNextStartTime(referenceDate)

	y, m, d := now.Date()
	h, min, s := now.Clock()
	e := time.Date(y, m, d, h, min, s, 0, time.Local)
	day := e.Day()
	for _, i := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
		e = time.Date(e.Year(), e.Month(), day, e.Hour(), e.Minute(), e.Second(), e.Nanosecond(), e.Location()).AddDate(0, 0, i)
		if e.Weekday() == time.Monday && (e.Equal(now) || e.After(now)) {
			break
		}
	}
	if !st.Equal(e) {
		t.Errorf("Expected %v was %v", e, st)
	}
}

func TestActionTimingHourWeekdays(t *testing.T) {
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{WeekDays: []time.Weekday{time.Monday}, StartTime: "10:01:00"}}}
	st := at.GetNextStartTime(referenceDate)

	y, m, d := now.Date()
	e := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	day := e.Day()
	for _, i := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
		e = time.Date(e.Year(), e.Month(), day, e.Hour(), e.Minute(), e.Second(), e.Nanosecond(), e.Location()).AddDate(0, 0, i)
		if e.Weekday() == time.Monday && (e.Equal(now) || e.After(now)) {
			break
		}
	}
	if !st.Equal(e) {
		t.Errorf("Expected %v was %v", e, st)
	}
}

func TestActionTimingOnlyMonthdays(t *testing.T) {

	y, m, d := now.Date()
	tomorrow := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{MonthDays: utils.MonthDays{1, 25, 2, tomorrow.Day()}}}}
	st := at.GetNextStartTime(referenceDate)
	expected := tomorrow
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingHourMonthdays(t *testing.T) {

	y, m, d := now.Date()
	testTime := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	tomorrow := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	if now.After(testTime) {
		y, m, d = tomorrow.Date()
	}
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{MonthDays: utils.MonthDays{now.Day(), tomorrow.Day()}, StartTime: "10:01:00"}}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingOnlyMonths(t *testing.T) {

	y, m, _ := now.Date()
	nextMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	t.Log("NextMonth: ", nextMonth)
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{Months: utils.Months{time.February, time.May, nextMonth.Month()}}}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingHourMonths(t *testing.T) {

	y, m, d := now.Date()
	testTime := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	nextMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	if now.After(testTime) {
		m = nextMonth.Month()
		y = nextMonth.Year()

	}
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{
		Months:    utils.Months{now.Month(), nextMonth.Month()},
		StartTime: "10:01:00"}}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingHourMonthdaysMonths(t *testing.T) {

	y, m, d := now.Date()
	testTime := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	nextMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	tomorrow := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)

	if now.After(testTime) {
		y, m, d = tomorrow.Date()
	}
	nextDay := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	month := nextDay.Month()
	if nextDay.Before(now) {
		if now.After(testTime) {
			month = nextMonth.Month()
		}
	}
	at := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			Months:    utils.Months{now.Month(), nextMonth.Month()},
			MonthDays: utils.MonthDays{now.Day(), tomorrow.Day()},
			StartTime: "10:01:00",
		},
	}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(y, month, d, 10, 1, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingFirstOfTheMonth(t *testing.T) {

	y, m, _ := now.Date()
	nextMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	at := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			MonthDays: utils.MonthDays{1},
		},
	}}
	st := at.GetNextStartTime(referenceDate)
	expected := nextMonth
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingOnlyYears(t *testing.T) {

	y, m, d := now.Date()
	nextYear := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{Years: utils.Years{now.Year(), nextYear.Year()}}}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingHourYears(t *testing.T) {

	y, m, d := now.Date()
	testTime := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	nextYear := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
	year := now.Year()
	if now.After(testTime) {
		year = nextYear.Year()
	}
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{Years: utils.Years{now.Year(), nextYear.Year()}, StartTime: "10:01:00"}}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(year, m, d, 10, 1, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingHourMonthdaysYear(t *testing.T) {

	y, m, d := now.Date()
	testTime := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	nextYear := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
	tomorrow := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	if now.After(testTime) {
		y, m, d = tomorrow.Date()
	}
	nextDay := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	year := nextDay.Year()
	if nextDay.Before(now) {
		if now.After(testTime) {
			year = nextYear.Year()
		}
	}
	at := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			Years:     utils.Years{now.Year(), nextYear.Year()},
			MonthDays: utils.MonthDays{now.Day(), tomorrow.Day()},
			StartTime: "10:01:00",
		},
	}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(year, m, d, 10, 1, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingHourMonthdaysMonthYear(t *testing.T) {

	y, m, d := now.Date()
	testTime := time.Date(y, m, d, 10, 1, 0, 0, time.Local)
	nextYear := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
	nextMonth := time.Date(y, m, 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 0)
	tomorrow := time.Date(y, m, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	day := now.Day()
	if now.After(testTime) {
		day = tomorrow.Day()
	}
	nextDay := time.Date(y, m, day, 10, 1, 0, 0, time.Local)
	month := now.Month()
	if nextDay.Before(now) {
		if now.After(testTime) {
			month = nextMonth.Month()
		}
	}
	nextDay = time.Date(y, month, day, 10, 1, 0, 0, time.Local)
	year := now.Year()
	if nextDay.Before(now) {
		if now.After(testTime) {
			year = nextYear.Year()
		}
	}
	at := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			Years:     utils.Years{now.Year(), nextYear.Year()},
			Months:    utils.Months{now.Month(), nextMonth.Month()},
			MonthDays: utils.MonthDays{now.Day(), tomorrow.Day()},
			StartTime: "10:01:00",
		},
	}}
	st := at.GetNextStartTime(referenceDate)
	expected := time.Date(year, month, day, 10, 1, 0, 0, time.Local)
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingFirstOfTheYear(t *testing.T) {
	y, _, _ := now.Date()
	nextYear := time.Date(y, 1, 1, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
	at := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			Years:     utils.Years{nextYear.Year()},
			Months:    utils.Months{time.January},
			MonthDays: utils.MonthDays{1},
			StartTime: "00:00:00",
		},
	}}
	st := at.GetNextStartTime(referenceDate)
	expected := nextYear
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingFirstMonthOfTheYear(t *testing.T) {
	y, _, _ := now.Date()
	nextYear := time.Date(y, 1, 1, 0, 0, 0, 0, time.Local).AddDate(1, 0, 0)
	at := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			Months: utils.Months{time.January},
		},
	}}
	st := at.GetNextStartTime(referenceDate)
	expected := nextYear
	if !st.Equal(expected) {
		t.Errorf("Expected %v was %v", expected, st)
	}
}

func TestActionTimingCheckForASAP(t *testing.T) {
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{StartTime: ASAP}}}
	if !at.CheckForASAP() {
		t.Errorf("%v should be asap!", at)
	}
}

func TestActionTimingIsOneTimeRun(t *testing.T) {
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{StartTime: ASAP}}}
	if !at.CheckForASAP() {
		t.Errorf("%v should be asap!", at)
	}
	if !at.IsOneTimeRun() {
		t.Errorf("%v should be one time run!", at)
	}
}

func TestActionTimingOneTimeRun(t *testing.T) {
	at := &ActionTiming{Timing: &RateInterval{Timing: &RITiming{StartTime: ASAP}}}
	at.CheckForASAP()
	nextRun := at.GetNextStartTime(referenceDate)
	if nextRun.IsZero() {
		t.Error("next time failed for asap")
	}
}

func TestActionTimingLogFunction(t *testing.T) {
	a := &Action{
		ActionType: "*log",
		BalanceId:  "test",
		Balance:    &Balance{Value: 1.1},
	}
	at := &ActionTiming{
		actions: []*Action{a},
	}
	err := at.Execute()
	if err != nil {
		t.Errorf("Could not execute LOG action: %v", err)
	}
}

func TestActionTimingPriotityListSortByWeight(t *testing.T) {
	at1 := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			Years:     utils.Years{2100},
			Months:    utils.Months{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December},
			MonthDays: utils.MonthDays{1},
			StartTime: "00:00:00",
		},
		Weight: 20,
	}}
	at2 := &ActionTiming{Timing: &RateInterval{
		Timing: &RITiming{
			Years:     utils.Years{2100},
			Months:    utils.Months{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December},
			MonthDays: utils.MonthDays{2},
			StartTime: "00:00:00",
		},
		Weight: 10,
	}}
	var atpl ActionTimingPriotityList
	atpl = append(atpl, at2, at1)
	atpl.Sort()
	if atpl[0] != at1 || atpl[1] != at2 {
		t.Error("Timing list not sorted correctly: ", at1, at2, atpl)
	}
}

func TestActionTimingPriotityListWeight(t *testing.T) {
	at1 := &ActionTiming{
		Timing: &RateInterval{
			Timing: &RITiming{
				Months:    utils.Months{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December},
				MonthDays: utils.MonthDays{1},
				StartTime: "00:00:00",
			},
		},
		Weight: 10.0,
	}
	at2 := &ActionTiming{
		Timing: &RateInterval{
			Timing: &RITiming{
				Months:    utils.Months{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December},
				MonthDays: utils.MonthDays{1},
				StartTime: "00:00:00",
			},
		},
		Weight: 20.0,
	}
	var atpl ActionTimingPriotityList
	atpl = append(atpl, at2, at1)
	atpl.Sort()
	if atpl[0] != at1 || atpl[1] != at2 {
		t.Error("Timing list not sorted correctly: ", atpl)
	}
}

func TestActionTimingsRemoveMember(t *testing.T) {
	at1 := &ActionTiming{
		Id:             "some uuid",
		Tag:            "test",
		UserBalanceIds: []string{"one", "two", "three"},
		ActionsId:      "TEST_ACTIONS",
	}
	at2 := &ActionTiming{
		Id:             "some uuid22",
		Tag:            "test2",
		UserBalanceIds: []string{"three", "four"},
		ActionsId:      "TEST_ACTIONS2",
	}
	ats := ActionPlan{at1, at2}
	if outAts := RemActionTiming(ats, "", "four"); len(outAts[1].UserBalanceIds) != 1 {
		t.Error("Expecting fewer balance ids", outAts[1].UserBalanceIds)
	}
	if ats = RemActionTiming(ats, "", "three"); len(ats) != 1 {
		t.Error("Expecting fewer actionTimings", ats)
	}
	if ats = RemActionTiming(ats, "some_uuid22", ""); len(ats) != 1 {
		t.Error("Expecting fewer actionTimings members", ats)
	}
	ats2 := ActionPlan{at1, at2}
	if ats2 = RemActionTiming(ats2, "", ""); len(ats2) != 0 {
		t.Error("Should have no members anymore", ats2)
	}
}

func TestActionTriggerMatchNil(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	var a *Action
	if !at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerMatchAllBlank(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	a := &Action{}
	if !at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerMatchMinuteBucketBlank(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	a := &Action{Direction: OUTBOUND, BalanceId: CREDIT}
	if !at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerMatchMinuteBucketFull(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	a := &Action{ExtraParameters: fmt.Sprintf(`{"ThresholdType":"%v", "ThresholdValue": %v}`, TRIGGER_MAX_BALANCE, 2)}
	if !at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerMatchAllFull(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	a := &Action{Direction: OUTBOUND, BalanceId: CREDIT, ExtraParameters: fmt.Sprintf(`{"ThresholdType":"%v", "ThresholdValue": %v}`, TRIGGER_MAX_BALANCE, 2)}
	if !at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerMatchSomeFalse(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	a := &Action{Direction: INBOUND, BalanceId: CREDIT, ExtraParameters: fmt.Sprintf(`{"ThresholdType":"%v", "ThresholdValue": %v}`, TRIGGER_MAX_BALANCE, 2)}
	if at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerMatcBalanceFalse(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	a := &Action{Direction: OUTBOUND, BalanceId: CREDIT, ExtraParameters: fmt.Sprintf(`{"ThresholdType":"%v", "ThresholdValue": %v}`, TRIGGER_MAX_BALANCE, 3.0)}
	if at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerMatcAllFalse(t *testing.T) {
	at := &ActionTrigger{
		Direction:      OUTBOUND,
		BalanceId:      CREDIT,
		ThresholdType:  TRIGGER_MAX_BALANCE,
		ThresholdValue: 2,
	}
	a := &Action{Direction: INBOUND, BalanceId: MINUTES, ExtraParameters: fmt.Sprintf(`{"ThresholdType":"%v", "ThresholdValue": %v}`, TRIGGER_MAX_COUNTER, 3)}
	if at.Match(a) {
		t.Errorf("Action trigger [%v] does not match action [%v]", at, a)
	}
}

func TestActionTriggerPriotityList(t *testing.T) {
	at1 := &ActionTrigger{Weight: 10}
	at2 := &ActionTrigger{Weight: 20}
	at3 := &ActionTrigger{Weight: 30}
	var atpl ActionTriggerPriotityList
	atpl = append(atpl, at2, at1, at3)
	atpl.Sort()
	if atpl[0] != at1 || atpl[2] != at3 || atpl[1] != at2 {
		t.Error("List not sorted: ", atpl)
	}
}

func TestActionResetTriggres(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 10}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	resetTriggersAction(ub, nil)
	if ub.ActionTriggers[0].Executed == true || ub.ActionTriggers[1].Executed == true {
		t.Error("Reset triggers action failed!")
	}
}

func TestActionResetTriggresExecutesThem(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 10}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	resetTriggersAction(ub, nil)
	if ub.ActionTriggers[0].Executed == true || ub.BalanceMap[CREDIT][0].Value == 12 {
		t.Error("Reset triggers action failed!")
	}
}

func TestActionResetTriggresActionFilter(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 10}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	resetTriggersAction(ub, &Action{BalanceId: SMS})
	if ub.ActionTriggers[0].Executed == false || ub.ActionTriggers[1].Executed == false {
		t.Error("Reset triggers action failed!")
	}
}

func TestActionSetPostpaid(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_PREPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	setPostpaidAction(ub, nil)
	if ub.Type != UB_TYPE_POSTPAID {
		t.Error("Set postpaid action failed!")
	}
}

func TestActionSetPrepaid(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_POSTPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	setPrepaidAction(ub, nil)
	if ub.Type != UB_TYPE_PREPAID {
		t.Error("Set prepaid action failed!")
	}
}

func TestActionResetPrepaid(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_POSTPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: SMS, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: SMS, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	resetPrepaidAction(ub, nil)
	if ub.Type != UB_TYPE_PREPAID ||
		ub.BalanceMap[CREDIT].GetTotalValue() != 0 ||
		len(ub.UnitCounters) != 0 ||
		ub.BalanceMap[MINUTES+OUTBOUND][0].Value != 0 ||
		ub.ActionTriggers[0].Executed == true || ub.ActionTriggers[1].Executed == true {
		t.Log(ub.BalanceMap)
		t.Error("Reset prepaid action failed!")
	}
}

func TestActionResetPostpaid(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_PREPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: SMS, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: SMS, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	resetPostpaidAction(ub, nil)
	if ub.Type != UB_TYPE_POSTPAID ||
		ub.BalanceMap[CREDIT].GetTotalValue() != 0 ||
		len(ub.UnitCounters) != 0 ||
		ub.BalanceMap[MINUTES+OUTBOUND][0].Value != 0 ||
		ub.ActionTriggers[0].Executed == true || ub.ActionTriggers[1].Executed == true {
		t.Error("Reset postpaid action failed!")
	}
}

func TestActionTopupResetCredit(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_PREPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT + OUTBOUND: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Direction: OUTBOUND, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: CREDIT, Direction: OUTBOUND, Balance: &Balance{Value: 10}}
	topupResetAction(ub, a)
	if ub.Type != UB_TYPE_PREPAID ||
		ub.BalanceMap[CREDIT+OUTBOUND].GetTotalValue() != 10 ||
		len(ub.UnitCounters) != 1 ||
		len(ub.BalanceMap[MINUTES+OUTBOUND]) != 2 ||
		ub.ActionTriggers[0].Executed != true || ub.ActionTriggers[1].Executed != true {
		t.Errorf("Topup reset action failed: %+v", ub.BalanceMap[CREDIT+OUTBOUND][0])
	}
}

func TestActionTopupResetMinutes(t *testing.T) {
	ub := &UserBalance{
		Id:   "TEST_UB",
		Type: UB_TYPE_PREPAID,
		BalanceMap: map[string]BalanceChain{
			CREDIT + OUTBOUND:  BalanceChain{&Balance{Value: 100}},
			MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Direction: OUTBOUND, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: MINUTES, Direction: OUTBOUND, Balance: &Balance{Value: 5, Weight: 20, DestinationId: "NAT"}}
	topupResetAction(ub, a)
	if ub.Type != UB_TYPE_PREPAID ||
		ub.BalanceMap[MINUTES+OUTBOUND].GetTotalValue() != 5 ||
		ub.BalanceMap[CREDIT+OUTBOUND].GetTotalValue() != 100 ||
		len(ub.UnitCounters) != 1 ||
		len(ub.BalanceMap[MINUTES+OUTBOUND]) != 1 ||
		ub.ActionTriggers[0].Executed != true || ub.ActionTriggers[1].Executed != true {
		t.Errorf("Topup reset minutes action failed: %+v", ub.BalanceMap[MINUTES+OUTBOUND][0])
	}
}

func TestActionTopupCredit(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_PREPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT + OUTBOUND: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Direction: OUTBOUND, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: CREDIT, Direction: OUTBOUND, Balance: &Balance{Value: 10}}
	topupAction(ub, a)
	if ub.Type != UB_TYPE_PREPAID ||
		ub.BalanceMap[CREDIT+OUTBOUND].GetTotalValue() != 110 ||
		len(ub.UnitCounters) != 1 ||
		len(ub.BalanceMap[MINUTES+OUTBOUND]) != 2 ||
		ub.ActionTriggers[0].Executed != true || ub.ActionTriggers[1].Executed != true {
		t.Error("Topup action failed!", ub.BalanceMap[CREDIT+OUTBOUND].GetTotalValue())
	}
}

func TestActionTopupMinutes(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_PREPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: MINUTES, Direction: OUTBOUND, Balance: &Balance{Value: 5, Weight: 20, DestinationId: "NAT"}}
	topupAction(ub, a)
	if ub.Type != UB_TYPE_PREPAID ||
		ub.BalanceMap[MINUTES+OUTBOUND].GetTotalValue() != 15 ||
		ub.BalanceMap[CREDIT].GetTotalValue() != 100 ||
		len(ub.UnitCounters) != 1 ||
		len(ub.BalanceMap[MINUTES+OUTBOUND]) != 2 ||
		ub.ActionTriggers[0].Executed != true || ub.ActionTriggers[1].Executed != true {
		t.Error("Topup minutes action failed!", ub.BalanceMap[MINUTES+OUTBOUND])
	}
}

func TestActionDebitCredit(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_PREPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT + OUTBOUND: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Direction: OUTBOUND, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: CREDIT, Direction: OUTBOUND, Balance: &Balance{Value: 10}}
	debitAction(ub, a)
	if ub.Type != UB_TYPE_PREPAID ||
		ub.BalanceMap[CREDIT+OUTBOUND].GetTotalValue() != 90 ||
		len(ub.UnitCounters) != 1 ||
		len(ub.BalanceMap[MINUTES+OUTBOUND]) != 2 ||
		ub.ActionTriggers[0].Executed != true || ub.ActionTriggers[1].Executed != true {
		t.Error("Debit action failed!", ub.BalanceMap[CREDIT+OUTBOUND].GetTotalValue())
	}
}

func TestActionDebitMinutes(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_PREPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}, &ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: MINUTES, Direction: OUTBOUND, Balance: &Balance{Value: 5, Weight: 20, DestinationId: "NAT"}}
	debitAction(ub, a)
	if ub.Type != UB_TYPE_PREPAID ||
		ub.BalanceMap[MINUTES+OUTBOUND][0].Value != 5 ||
		ub.BalanceMap[CREDIT].GetTotalValue() != 100 ||
		len(ub.UnitCounters) != 1 ||
		len(ub.BalanceMap[MINUTES+OUTBOUND]) != 2 ||
		ub.ActionTriggers[0].Executed != true || ub.ActionTriggers[1].Executed != true {
		t.Error("Debit minutes action failed!", ub.BalanceMap[MINUTES+OUTBOUND][0])
	}
}

func TestActionResetAllCounters(t *testing.T) {
	ub := &UserBalance{
		Id:   "TEST_UB",
		Type: UB_TYPE_POSTPAID,
		BalanceMap: map[string]BalanceChain{
			CREDIT: BalanceChain{&Balance{Value: 100}},
			MINUTES: BalanceChain{
				&Balance{Value: 10, Weight: 20, DestinationId: "NAT"},
				&Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	resetCountersAction(ub, nil)
	if ub.Type != UB_TYPE_POSTPAID ||
		ub.BalanceMap[CREDIT].GetTotalValue() != 100 ||
		len(ub.UnitCounters) != 1 ||
		len(ub.UnitCounters[0].Balances) != 2 ||
		len(ub.BalanceMap[MINUTES]) != 2 ||
		ub.ActionTriggers[0].Executed != true {
		t.Errorf("Reset counters action failed: %+v", ub.UnitCounters)
	}
	if len(ub.UnitCounters) < 1 {
		t.FailNow()
	}
	mb := ub.UnitCounters[0].Balances[0]
	if mb.Weight != 20 || mb.Value != 0 || mb.DestinationId != "NAT" {
		t.Errorf("Balance cloned incorrectly: %v!", mb)
	}
}

func TestActionResetCounterMinutes(t *testing.T) {
	ub := &UserBalance{
		Id:   "TEST_UB",
		Type: UB_TYPE_POSTPAID,
		BalanceMap: map[string]BalanceChain{
			CREDIT:  BalanceChain{&Balance{Value: 100}},
			MINUTES: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: MINUTES}
	resetCounterAction(ub, a)
	if ub.Type != UB_TYPE_POSTPAID ||
		ub.BalanceMap[CREDIT].GetTotalValue() != 100 ||
		len(ub.UnitCounters) != 2 ||
		len(ub.UnitCounters[1].Balances) != 2 ||
		len(ub.BalanceMap[MINUTES]) != 2 ||
		ub.ActionTriggers[0].Executed != true {
		for _, b := range ub.UnitCounters[1].Balances {
			t.Logf("B: %+v", b)
		}
		t.Error("Reset counters action failed!", ub.UnitCounters[1])
	}
	if len(ub.UnitCounters) < 2 || len(ub.UnitCounters[1].Balances) < 1 {
		t.FailNow()
	}
	mb := ub.UnitCounters[1].Balances[0]
	if mb.Weight != 20 || mb.Value != 0 || mb.DestinationId != "NAT" {
		t.Errorf("Minute bucked cloned incorrectly: %v!", mb)
	}
}

func TestActionResetCounterCREDIT(t *testing.T) {
	ub := &UserBalance{
		Id:             "TEST_UB",
		Type:           UB_TYPE_POSTPAID,
		BalanceMap:     map[string]BalanceChain{CREDIT: BalanceChain{&Balance{Value: 100}}, MINUTES + OUTBOUND: BalanceChain{&Balance{Value: 10, Weight: 20, DestinationId: "NAT"}, &Balance{Weight: 10, DestinationId: "RET"}}},
		UnitCounters:   []*UnitsCounter{&UnitsCounter{BalanceId: CREDIT, Direction: OUTBOUND, Balances: BalanceChain{&Balance{Value: 1}}}, &UnitsCounter{BalanceId: SMS, Direction: OUTBOUND, Balances: BalanceChain{&Balance{Value: 1}}}},
		ActionTriggers: ActionTriggerPriotityList{&ActionTrigger{BalanceId: CREDIT, Direction: OUTBOUND, ThresholdValue: 2, ActionsId: "TEST_ACTIONS", Executed: true}},
	}
	a := &Action{BalanceId: CREDIT, Direction: OUTBOUND}
	resetCounterAction(ub, a)
	if ub.Type != UB_TYPE_POSTPAID ||
		ub.BalanceMap[CREDIT].GetTotalValue() != 100 ||
		len(ub.UnitCounters) != 2 ||
		len(ub.BalanceMap[MINUTES+OUTBOUND]) != 2 ||
		ub.ActionTriggers[0].Executed != true {
		t.Error("Reset counters action failed!", ub.UnitCounters)
	}
}

func TestActionTriggerLogging(t *testing.T) {
	at := &ActionTrigger{
		Id:             "some_uuid",
		BalanceId:      CREDIT,
		Direction:      OUTBOUND,
		ThresholdValue: 100.0,
		DestinationId:  "NAT",
		Weight:         10.0,
		ActionsId:      "TEST_ACTIONS",
	}
	as, err := accountingStorage.GetActions(at.ActionsId, false)
	if err != nil {
		t.Error("Error getting actions for the action timing: ", as, err)
	}
	storageLogger.LogActionTrigger("rif", RATER_SOURCE, at, as)
	//expected := "rif*some_uuid;MONETARY;OUT;NAT;TEST_ACTIONS;100;10;false*|TOPUP|MONETARY|OUT|10|0"
	var key string
	atMap, _ := accountingStorage.GetAllActionTimings()
	for k, v := range atMap {
		_ = k
		_ = v
		/*if strings.Contains(k, LOG_ACTION_TRIGGER_PREFIX) && strings.Contains(v, expected) {
			key = k
			break
		}*/
	}
	if key != "" {
		t.Error("Action timing was not logged")
	}
}

func TestActionTimingLogging(t *testing.T) {
	i := &RateInterval{
		Timing: &RITiming{
			Months:    utils.Months{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December},
			MonthDays: utils.MonthDays{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
			WeekDays:  utils.WeekDays{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			StartTime: "18:00:00",
			EndTime:   "00:00:00",
		},
		Weight: 10.0,
		Rating: &RIRate{
			ConnectFee: 0.0,
			Rates:      RateGroups{&Rate{0, 1.0, 1 * time.Second, 60 * time.Second}},
		},
	}
	at := &ActionTiming{
		Id:             "some uuid",
		Tag:            "test",
		UserBalanceIds: []string{"one", "two", "three"},
		Timing:         i,
		Weight:         10.0,
		ActionsId:      "TEST_ACTIONS",
	}
	as, err := accountingStorage.GetActions(at.ActionsId, false)
	if err != nil {
		t.Error("Error getting actions for the action trigger: ", err)
	}
	storageLogger.LogActionTiming(SCHED_SOURCE, at, as)
	//expected := "some uuid|test|one,two,three|;1,2,3,4,5,6,7,8,9,10,11,12;1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31;1,2,3,4,5;18:00:00;00:00:00;10;0;1;60;1|10|TEST_ACTIONS*|TOPUP|MONETARY|OUT|10|0"
	var key string
	atMap, _ := accountingStorage.GetAllActionTimings()
	for k, v := range atMap {
		_ = k
		_ = v
		/*if strings.Contains(k, LOG_ACTION_TIMMING_PREFIX) && strings.Contains(string(v), expected) {
			key = k
		}*/
	}
	if key != "" {
		t.Error("Action trigger was not logged")
	}
}

func TestActionMakeNegative(t *testing.T) {
	a := &Action{Balance: &Balance{Value: 10}}
	genericMakeNegative(a)
	if a.Balance.Value > 0 {
		t.Error("Failed to make negative: ", a)
	}
	genericMakeNegative(a)
	if a.Balance.Value > 0 {
		t.Error("Failed to preserve negative: ", a)
	}
}

func TestMinBalanceTriggerSimple(t *testing.T) {
	//ub := &UserBalance{}
}

/********************************** Benchmarks ********************************/

func BenchmarkUUID(b *testing.B) {
	m := make(map[string]int, 1000)
	for i := 0; i < b.N; i++ {
		uuid := utils.GenUUID()
		if len(uuid) == 0 {
			b.Fatalf("GenUUID error %s", uuid)
		}
		b.StopTimer()
		c := m[uuid]
		if c > 0 {
			b.Fatalf("duplicate uuid[%s] count %d", uuid, c)
		}
		m[uuid] = c + 1
		b.StartTimer()
	}
}
