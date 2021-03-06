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
	//"fmt"

	"time"

	"github.com/cgrates/cgrates/utils"
)

/*
A unit in which a call will be split that has a specific price related interval attached to it.
*/
type TimeSpan struct {
	TimeStart, TimeEnd            time.Time
	Cost                          float64
	ratingInfo                    *RatingInfo
	RateInterval                  *RateInterval
	CallDuration                  time.Duration // the call duration so far till TimeEnd
	Increments                    Increments
	MatchedSubject, MatchedPrefix string
}

type Increment struct {
	Duration            time.Duration
	Cost                float64
	BalanceUuids        []string // need more than one for minutes with cost
	BalanceRateInterval *RateInterval
	MinuteInfo          *MinuteInfo
	paid                bool
}

// Holds the minute information related to a specified timespan
type MinuteInfo struct {
	DestinationId string
	Quantity      float64
	//Price         float64
}

type TimeSpans []*TimeSpan

func (timespans *TimeSpans) RemoveOverlapedFromIndex(index int) {
	tss := *timespans
	ts := tss[index]
	endOverlapIndex := index
	for i := index + 1; i < len(tss); i++ {
		if tss[i].TimeEnd.Before(ts.TimeEnd) || tss[i].TimeEnd.Equal(ts.TimeEnd) {
			endOverlapIndex = i
		} else if tss[i].TimeStart.Before(ts.TimeEnd) {
			tss[i].TimeStart = ts.TimeEnd
			break
		}
	}
	if endOverlapIndex > index {
		newSliceEnd := len(tss) - (endOverlapIndex - index)
		// delete overlapped
		copy(tss[index+1:], tss[endOverlapIndex+1:])
		for i := newSliceEnd; i < len(tss); i++ {
			tss[i] = nil
		}
		*timespans = tss[:newSliceEnd]
		return
	}
	*timespans = tss
}

// The paidTs will replace the timespans that are exactly under them from the reciver list
func (timespans *TimeSpans) OverlapWithTimeSpans(paidTs TimeSpans, newTs *TimeSpan, index int) bool {
	tss := *timespans
	// calculate overlaped timespans
	var paidDuration time.Duration
	for _, pts := range paidTs {
		paidDuration += pts.GetDuration()
	}
	if paidDuration > 0 {
		// we must add the rest of the current ts to the remaingTs
		var remainingTs []*TimeSpan
		overlapStartIndex := index
		if newTs != nil {
			remainingTs = append(remainingTs, newTs)
			overlapStartIndex += 1
		}
		for tsi := overlapStartIndex; tsi < len(tss); tsi++ {
			remainingTs = append(remainingTs, tss[tsi])
		}
		overlapEndIndex := 0
		for i, rts := range remainingTs {
			if paidDuration >= rts.GetDuration() {
				paidDuration -= rts.GetDuration()
			} else {
				if paidDuration > 0 {
					// this ts was not fully paid
					fragment := rts.SplitByDuration(paidDuration)
					paidTs = append(paidTs, fragment)
				}
				// find the end position in tss
				overlapEndIndex = overlapStartIndex + i
				break
			}
			// find the end position in tss
			overlapEndIndex = overlapStartIndex + i
		}
		// delete from index to current
		if overlapEndIndex == len(tss)-1 {
			tss = tss[:overlapStartIndex]
		} else {
			if overlapEndIndex+1 < len(tss) {
				tss = append(tss[:overlapStartIndex], tss[overlapEndIndex+1:]...)
			}
		}
		// append the timespans to outer tss
		for i, pts := range paidTs {
			tss = append(tss, nil)
			copy(tss[overlapStartIndex+i+1:], tss[overlapStartIndex+i:])
			tss[overlapStartIndex+i] = pts
		}
		*timespans = tss
		return true
	}
	*timespans = tss
	return false
}

func (incr *Increment) Clone() *Increment {
	nIncr := &Increment{
		Duration:            incr.Duration,
		Cost:                incr.Cost,
		BalanceRateInterval: incr.BalanceRateInterval,
		MinuteInfo:          incr.MinuteInfo,
	}
	nIncr.BalanceUuids = make([]string, len(incr.BalanceUuids))
	copy(nIncr.BalanceUuids, incr.BalanceUuids)
	return nIncr
}

func (incr *Increment) SetMinuteBalance(bUuid string) {
	incr.BalanceUuids[0] = bUuid
}

func (incr *Increment) GetMinuteBalance() string {
	return incr.BalanceUuids[0]
}

func (incr *Increment) SetMoneyBalance(bUuid string) {
	incr.BalanceUuids[1] = bUuid
}

func (incr *Increment) GetMoneyBalance() string {
	return incr.BalanceUuids[1]
}

type Increments []*Increment

func (incs Increments) GetTotalCost() float64 {
	cost := 0.0
	for _, increment := range incs {
		cost += increment.Cost
	}
	return cost
}

// Returns the duration of the timespan
func (ts *TimeSpan) GetDuration() time.Duration {
	return ts.TimeEnd.Sub(ts.TimeStart)
}

// Returns true if the given time is inside timespan range.
func (ts *TimeSpan) Contains(t time.Time) bool {
	return t.After(ts.TimeStart) && t.Before(ts.TimeEnd)
}

/*
Will set the interval as spans's interval if new Weight is lower then span's interval Weight
or if the Weights are equal and new price is lower then spans's interval price
*/
func (ts *TimeSpan) SetRateInterval(i *RateInterval) {
	if ts.RateInterval == nil || ts.RateInterval.Weight < i.Weight {
		ts.RateInterval = i
		return
	}
	iPrice, _, _ := i.GetRateParameters(ts.GetGroupStart())
	tsPrice, _, _ := ts.RateInterval.GetRateParameters(ts.GetGroupStart())
	if ts.RateInterval.Weight == i.Weight && iPrice < tsPrice {
		ts.RateInterval = i
	}
}

// Returns the cost of the timespan according to the relevant cost interval.
// It also sets the Cost field of this timespan (used for refund on session
// manager debit loop where the cost cannot be recalculated)
func (ts *TimeSpan) getCost() float64 {
	if len(ts.Increments) == 0 {
		if ts.RateInterval == nil {
			return 0
		}
		cost := ts.RateInterval.GetCost(ts.GetDuration(), ts.GetGroupStart())
		ts.Cost = utils.Round(cost, ts.RateInterval.Rating.RoundingDecimals, ts.RateInterval.Rating.RoundingMethod)
		return ts.Cost
	} else {
		return ts.Increments[0].Cost * float64(len(ts.Increments))
	}
	return 0
}

func (ts *TimeSpan) createIncrementsSlice() {
	if ts.RateInterval == nil {
		return
	}
	ts.Increments = make([]*Increment, 0)
	// create rated units series
	_, rateIncrement, _ := ts.RateInterval.GetRateParameters(ts.GetGroupStart())
	// we will use the cost calculated cost and devide by nb of increments
	// because ts cost is rounded
	//incrementCost := rate / rateUnit.Seconds() * rateIncrement.Seconds()
	nbIncrements := int(ts.GetDuration() / rateIncrement)
	incrementCost := ts.getCost() / float64(nbIncrements)
	incrementCost = utils.Round(incrementCost, roundingDecimals, utils.ROUNDING_MIDDLE) // just get rid of the extra decimals
	for s := 0; s < nbIncrements; s++ {
		inc := &Increment{
			Duration:     rateIncrement,
			Cost:         incrementCost,
			BalanceUuids: make([]string, 2),
		}
		ts.Increments = append(ts.Increments, inc)
	}
	// put the rounded cost back in timespan
	ts.Cost = incrementCost * float64(nbIncrements)
}

// returns whether the timespan has all increments marked as paid and if not
// it also returns the first unpaied increment
func (ts *TimeSpan) IsPaid() (bool, int) {
	if len(ts.Increments) == 0 {
		return false, 0
	}
	for incrementIndex, increment := range ts.Increments {
		if !increment.paid {
			return false, incrementIndex
		}
	}
	return true, len(ts.Increments)
}

/*
Splits the given timespan according to how it relates to the interval.
It will modify the endtime of the received timespan and it will return
a new timespan starting from the end of the received one.
The interval will attach itself to the timespan that overlaps the interval.
*/
func (ts *TimeSpan) SplitByRateInterval(i *RateInterval) (nts *TimeSpan) {
	//Logger.Debug("here: ", ts, " +++ ", i)
	// if the span is not in interval return nil
	if !(i.Contains(ts.TimeStart, false) || i.Contains(ts.TimeEnd, true)) {
		//Logger.Debug("Not in interval")
		return
	}
	//Logger.Debug(fmt.Sprintf("TS: %+v", ts))
	// split by GroupStart
	if i.Rating != nil {
		i.Rating.Rates.Sort()
		for _, rate := range i.Rating.Rates {
			// Logger.Debug(fmt.Sprintf("Rate: %+v", rate))
			if ts.GetGroupStart() < rate.GroupIntervalStart && ts.GetGroupEnd() > rate.GroupIntervalStart {
				// Logger.Debug(fmt.Sprintf("Splitting"))
				ts.SetRateInterval(i)
				splitTime := ts.TimeStart.Add(rate.GroupIntervalStart - ts.GetGroupStart())
				nts = &TimeSpan{
					TimeStart: splitTime,
					TimeEnd:   ts.TimeEnd,
				}
				nts.copyRatingInfo(ts)
				ts.TimeEnd = splitTime
				nts.SetRateInterval(i)
				nts.CallDuration = ts.CallDuration
				ts.SetNewCallDuration(nts)
				// Logger.Debug(fmt.Sprintf("Group splitting: %+v %+v", ts, nts))
				return
			}
		}
	}
	// if the span is enclosed in the interval try to set as new interval and return nil
	if i.Contains(ts.TimeStart, false) && i.Contains(ts.TimeEnd, true) {
		//Logger.Debug("All in interval")
		ts.SetRateInterval(i)
		return
	}
	// if only the start time is in the interval split the interval to the right
	if i.Contains(ts.TimeStart, false) {
		//Logger.Debug("Start in interval")
		splitTime := i.getRightMargin(ts.TimeStart)

		ts.SetRateInterval(i)
		if splitTime == ts.TimeStart || splitTime.Equal(ts.TimeEnd) {
			return
		}
		nts = &TimeSpan{
			TimeStart: splitTime,
			TimeEnd:   ts.TimeEnd,
		}
		nts.copyRatingInfo(ts)
		ts.TimeEnd = splitTime
		nts.CallDuration = ts.CallDuration
		ts.SetNewCallDuration(nts)
		// Logger.Debug(fmt.Sprintf("right: %+v %+v", ts, nts))
		return
	}
	// if only the end time is in the interval split the interval to the left
	if i.Contains(ts.TimeEnd, true) {
		//Logger.Debug("End in interval")
		//tmpTime := time.Date(ts.TimeStart.)
		splitTime := i.getLeftMargin(ts.TimeEnd)
		splitTime = utils.CopyHour(splitTime, ts.TimeStart)
		if splitTime.Equal(ts.TimeEnd) {
			return
		}
		nts = &TimeSpan{
			TimeStart: splitTime,
			TimeEnd:   ts.TimeEnd,
		}
		nts.copyRatingInfo(ts)
		ts.TimeEnd = splitTime

		nts.SetRateInterval(i)
		nts.CallDuration = ts.CallDuration
		ts.SetNewCallDuration(nts)
		// Logger.Debug(fmt.Sprintf("left: %+v %+v", ts, nts))
		return
	}
	return
}

// Split the timespan at the given increment start
func (ts *TimeSpan) SplitByIncrement(index int) *TimeSpan {
	if index <= 0 || index >= len(ts.Increments) {
		return nil
	}
	timeStart := ts.GetTimeStartForIncrement(index)
	newTs := &TimeSpan{
		RateInterval: ts.RateInterval,
		TimeStart:    timeStart,
		TimeEnd:      ts.TimeEnd,
	}
	newTs.copyRatingInfo(ts)
	newTs.CallDuration = ts.CallDuration
	ts.TimeEnd = timeStart
	newTs.Increments = ts.Increments[index:]
	ts.Increments = ts.Increments[:index]
	ts.SetNewCallDuration(newTs)
	return newTs
}

// Split the timespan at the given second
func (ts *TimeSpan) SplitByDuration(duration time.Duration) *TimeSpan {
	if duration <= 0 || duration >= ts.GetDuration() {
		return nil
	}
	timeStart := ts.TimeStart.Add(duration)
	newTs := &TimeSpan{
		RateInterval: ts.RateInterval,
		TimeStart:    timeStart,
		TimeEnd:      ts.TimeEnd,
	}
	newTs.copyRatingInfo(ts)
	newTs.CallDuration = ts.CallDuration
	ts.TimeEnd = timeStart
	// split the increment
	for incrIndex, incr := range ts.Increments {
		if duration-incr.Duration >= 0 {
			duration -= incr.Duration
		} else {

			splitIncrement := ts.Increments[incrIndex].Clone()
			splitIncrement.Duration -= duration
			ts.Increments[incrIndex].Duration = duration
			newTs.Increments = Increments{splitIncrement}
			if incrIndex < len(ts.Increments)-1 {
				newTs.Increments = append(newTs.Increments, ts.Increments[incrIndex+1:]...)
			}
			ts.Increments = ts.Increments[:incrIndex+1]
			break
		}
	}
	ts.SetNewCallDuration(newTs)
	return newTs
}

// Splits the given timespan on activation period's activation time.
func (ts *TimeSpan) SplitByRatingPlan(rp *RatingInfo) (newTs *TimeSpan) {
	if !ts.Contains(rp.ActivationTime) {
		return nil
	}
	newTs = &TimeSpan{
		TimeStart: rp.ActivationTime,
		TimeEnd:   ts.TimeEnd,
	}
	newTs.copyRatingInfo(ts)
	newTs.CallDuration = ts.CallDuration
	ts.TimeEnd = rp.ActivationTime
	ts.SetNewCallDuration(newTs)
	// Logger.Debug(fmt.Sprintf("RP SPLITTING: %+v %+v", ts, newTs))
	return
}

// Returns the starting time of this timespan
func (ts *TimeSpan) GetGroupStart() time.Duration {
	s := ts.CallDuration - ts.GetDuration()
	if s < 0 {
		s = 0
	}
	return s
}

func (ts *TimeSpan) GetGroupEnd() time.Duration {
	return ts.CallDuration
}

// sets the CallDuration attribute to reflect new timespan
func (ts *TimeSpan) SetNewCallDuration(nts *TimeSpan) {
	d := ts.CallDuration - nts.GetDuration()
	if d < 0 {
		d = 0
	}
	ts.CallDuration = d
}

func (nts *TimeSpan) copyRatingInfo(ts *TimeSpan) {
	if ts.ratingInfo == nil {
		return
	}
	nts.ratingInfo = ts.ratingInfo
	nts.MatchedSubject = ts.ratingInfo.MatchedSubject
	nts.MatchedPrefix = ts.ratingInfo.MatchedPrefix
}

// returns a time for the specified second in the time span
func (ts *TimeSpan) GetTimeStartForIncrement(index int) time.Time {
	return ts.TimeStart.Add(time.Duration(int64(index) * ts.Increments[0].Duration.Nanoseconds()))
}

func (ts *TimeSpan) RoundToDuration(duration time.Duration) {
	if duration < ts.GetDuration() {
		duration = utils.RoundTo(duration, ts.GetDuration())
	}
	if duration > ts.GetDuration() {
		initialDuration := ts.GetDuration()
		ts.TimeEnd = ts.TimeStart.Add(duration)
		ts.CallDuration = ts.CallDuration + (duration - initialDuration)
	}
}
