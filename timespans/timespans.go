/*
Rating system designed to be used in VoIP Carriers World
Copyright (C) 2012  Radu Ioan Fericean

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
package timespans

import (
	"fmt"
	"time"
	//"log"
)

/*
A unit in which a call will be split that has a specific price related interval attached to it.
*/
type TimeSpan struct {
	TimeStart, TimeEnd time.Time
	ActivationPeriod   *ActivationPeriod
	Interval           *Interval
	MinuteInfo         *MinuteInfo
}

type MinuteInfo struct {
	DestinationId string
	Quantity      float64
	Price         float64
}

/*
Returns the duration of the timespan
*/
func (ts *TimeSpan) GetDuration() time.Duration {
	return ts.TimeEnd.Sub(ts.TimeStart)
}

/*
Returns the cost of the timespan according to the relevant cost interval.
*/
func (ts *TimeSpan) GetCost(cd *CallDescriptor) (cost float64) {
	if ts.MinuteInfo != nil {
		return ts.GetDuration().Seconds() * ts.MinuteInfo.Price
	}
	if ts.Interval == nil {
		return 0
	}
	if ts.Interval.BillingUnit > 0 {
		cost = (ts.GetDuration().Seconds() / ts.Interval.BillingUnit) * ts.Interval.Price
	} else {
		cost = ts.GetDuration().Seconds() * ts.Interval.Price
	}
	if userBudget, err := cd.getUserBudget(); err == nil && userBudget != nil {
		userBudget.mux.RLock()
		if percentageDiscount, err := userBudget.getVolumeDiscount(cd.storageGetter); err == nil && percentageDiscount > 0 {
			cost *= (100 - percentageDiscount) / 100
		}
		userBudget.mux.RUnlock()
	}
	return
}

/*
Returns true if the given time is inside timespan range.
*/
func (ts *TimeSpan) Contains(t time.Time) bool {
	return t.After(ts.TimeStart) && t.Before(ts.TimeEnd)
}

/*
will set the interval as spans's interval if new ponder is greater then span's interval ponder
or if the ponders are equal and new price is lower then spans's interval price
*/
func (ts *TimeSpan) SetInterval(i *Interval) {
	if ts.Interval == nil || ts.Interval.Ponder < i.Ponder {
		ts.Interval = i
	}
	if ts.Interval.Ponder == i.Ponder && i.Price < ts.Interval.Price {
		ts.Interval = i
	}
}

/*
Splits the given timespan according to how it relates to the interval.
It will modify the endtime of the received timespan and it will return
a new timespan starting from the end of the received one.
The interval will attach itself to the timespan that overlaps the interval.
*/
func (ts *TimeSpan) SplitByInterval(i *Interval) (nts *TimeSpan) {
	// if the span is not in interval return nil
	if !(i.Contains(ts.TimeStart) || i.Contains(ts.TimeEnd)) {
		return
	}
	// if the span is enclosed in the interval try to set as new interval and return nil
	if i.Contains(ts.TimeStart) && i.Contains(ts.TimeEnd) {
		ts.SetInterval(i)
		return
	}
	// if only the start time is in the interval split the interval
	if i.Contains(ts.TimeStart) {
		splitTime := i.getRightMargin(ts.TimeStart)
		ts.SetInterval(i)
		if splitTime == ts.TimeStart {
			return
		}
		nts = &TimeSpan{TimeStart: splitTime, TimeEnd: ts.TimeEnd}
		ts.TimeEnd = splitTime

		return
	}
	// if only the end time is in the interval split the interval
	if i.Contains(ts.TimeEnd) {
		splitTime := i.getLeftMargin(ts.TimeEnd)
		if splitTime == ts.TimeEnd {
			return
		}
		nts = &TimeSpan{TimeStart: splitTime, TimeEnd: ts.TimeEnd}
		ts.TimeEnd = splitTime

		nts.SetInterval(i)
		return
	}
	return
}

/*
Splits the given timespan on activation period's activation time.
*/
func (ts *TimeSpan) SplitByActivationPeriod(ap *ActivationPeriod) (newTs *TimeSpan) {
	if !ts.Contains(ap.ActivationTime) {
		return nil
	}
	newTs = &TimeSpan{TimeStart: ap.ActivationTime, TimeEnd: ts.TimeEnd, ActivationPeriod: ap}
	ts.TimeEnd = ap.ActivationTime
	return
}

/*
Splits the given timespan on activation period's activation time.
*/
func (ts *TimeSpan) SplitByMinuteBucket(mb *MinuteBucket) (newTs *TimeSpan) {
	s := ts.GetDuration().Seconds()
	ts.MinuteInfo = &MinuteInfo{mb.DestinationId, s, mb.Price}
	if s <= mb.Seconds {
		mb.Seconds -= s
		return nil
	}
	secDuration, _ := time.ParseDuration(fmt.Sprintf("%vs", mb.Seconds))

	newTimeEnd := ts.TimeStart.Add(secDuration)
	newTs = &TimeSpan{TimeStart: newTimeEnd, TimeEnd: ts.TimeEnd}
	ts.TimeEnd = newTimeEnd
	ts.MinuteInfo.Quantity = mb.Seconds
	mb.Seconds = 0

	return
}
