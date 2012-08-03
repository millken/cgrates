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
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CSVReader struct {
	readerFunc        func(string, rune) (*csv.Reader, *os.File, error)
	actions           map[string][]*Action
	actionsTimings    map[string][]*ActionTiming
	actionsTriggers   map[string][]*ActionTrigger
	accountActions    []*UserBalance
	destinations      []*Destination
	rates             map[string][]*Rate
	timings           map[string][]*Timing
	activationPeriods map[string]*ActivationPeriod
	ratingProfiles    map[string]CallDescriptors
}

func NewFileCSVReader() *CSVReader {
	c := new(CSVReader)
	c.actions = make(map[string][]*Action)
	c.actionsTimings = make(map[string][]*ActionTiming)
	c.actionsTriggers = make(map[string][]*ActionTrigger)
	c.rates = make(map[string][]*Rate)
	c.timings = make(map[string][]*Timing)
	c.activationPeriods = make(map[string]*ActivationPeriod)
	c.ratingProfiles = make(map[string]CallDescriptors)
	c.readerFunc = openFileCSVReader
	return c
}

func NewStringCSVReader() *CSVReader {
	c := NewFileCSVReader()
	c.readerFunc = openStringCSVReader
	return c
}

func openFileCSVReader(fn string, comma rune) (csvReader *csv.Reader, fp *os.File, err error) {
	fp, err = os.Open(fn)
	if err != nil {
		return
	}
	csvReader = csv.NewReader(fp)
	csvReader.Comma = comma
	csvReader.TrailingComma = true
	return
}

func openStringCSVReader(data string, comma rune) (csvReader *csv.Reader, fp *os.File, err error) {
	csvReader = csv.NewReader(strings.NewReader(data))
	csvReader.Comma = comma
	csvReader.TrailingComma = true
	return
}

func (csvr *CSVReader) WriteToDatabase(storage StorageGetter, flush, verbose bool) {
	if flush {
		storage.Flush()
	}
	if verbose {
		log.Print("Destinations")
	}
	for _, d := range csvr.destinations {
		storage.SetDestination(d)
		if verbose {
			log.Print(d.Id, " : ", d.Prefixes)
		}
	}
	if verbose {
		log.Print("Rating profiles")
	}
	for _, cds := range csvr.ratingProfiles {
		for _, cd := range cds {
			storage.SetActivationPeriodsOrFallback(cd.GetKey(), cd.ActivationPeriods, cd.FallbackKey)
			if verbose {
				log.Print(cd.GetKey())
			}
		}
	}
	if verbose {
		log.Print("Action timings")
	}
	for k, ats := range csvr.actionsTimings {
		storage.SetActionTimings(ACTION_TIMING_PREFIX+":"+k, ats)
		if verbose {
			log.Println(k)
		}
	}
	if verbose {
		log.Print("Actions")
	}
	for k, as := range csvr.actions {
		storage.SetActions(k, as)
		if verbose {
			log.Println(k)
		}
	}
	if verbose {
		log.Print("Account actions")
	}
	for _, ub := range csvr.accountActions {
		storage.SetUserBalance(ub)
		if verbose {
			log.Println(ub.Id)
		}
	}
}

func (csvr *CSVReader) LoadDestinations(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {

		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		var dest *Destination
		for _, d := range csvr.destinations {
			if d.Id == tag {
				dest = d
				break
			}
		}
		if dest == nil {
			dest = &Destination{Id: tag}
			csvr.destinations = append(csvr.destinations, dest)
		}
		dest.Prefixes = append(dest.Prefixes, record[1])
	}
}

func (csvr *CSVReader) LoadRates(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		r, err := NewRate(record[1], record[2], record[3], record[4], record[5])
		if err != nil {
			continue
		}
		csvr.rates[tag] = append(csvr.rates[tag], r)
	}
}

func (csvr *CSVReader) LoadTimings(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}

		t := NewTiming(record[1:]...)
		csvr.timings[tag] = append(csvr.timings[tag], t)
	}
}

func (csvr *CSVReader) LoadRateTimings(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}

		ts, exists := csvr.timings[record[2]]
		if !exists {
			log.Printf("Could not get timing for tag %v", record[2])
			continue
		}
		for _, t := range ts {
			rt := NewRateTiming(record[1], t, record[3])
			rs, exists := csvr.rates[record[1]]
			if !exists {
				log.Printf("Could not rate for tag %v", record[2])
				continue
			}
			for _, r := range rs {
				_, exists := csvr.activationPeriods[tag]
				if !exists {
					csvr.activationPeriods[tag] = &ActivationPeriod{}
				}
				csvr.activationPeriods[tag].AddIntervalIfNotPresent(rt.GetInterval(r))
			}
		}
	}
}

func (csvr *CSVReader) LoadRatingProfiles(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tenant" {
			// skip header line
			continue
		}
		if len(record) != 7 {
			log.Printf("Malformed rating profile: %v", record)
			continue
		}
		tenant, tor, direction, subject, fallbacksubject := record[0], record[1], record[2], record[3], record[4]
		at, err := time.Parse(time.RFC3339, record[6])
		if err != nil {
			log.Printf("Cannot parse activation time from %v", record[6])
			continue
		}
		for _, d := range csvr.destinations {
			for _, p := range d.Prefixes { //destinations
				// Search for a CallDescriptor with the same key
				var cd *CallDescriptor
				key := fmt.Sprintf("%s:%s:%s:%s:%s", direction, tenant, tor, subject, p)
				for _, c := range csvr.ratingProfiles[p] {
					if c.GetKey() == key {
						cd = c
					}
				}
				if cd == nil {
					cd = &CallDescriptor{
						Direction:   direction,
						Tenant:      tenant,
						TOR:         tor,
						Subject:     subject,
						Destination: p,
					}
					csvr.ratingProfiles[p] = append(csvr.ratingProfiles[p], cd)
				}
				ap, exists := csvr.activationPeriods[record[5]]
				if !exists {
					log.Print("Could not load ratinTiming for tag: ", record[5])
					continue
				}
				newAP := &ActivationPeriod{}
				//copy(newAP.Intervals, ap.Intervals)
				newAP.Intervals = append(newAP.Intervals, ap.Intervals...)
				newAP.ActivationTime = at
				cd.AddActivationPeriodIfNotPresent(newAP)
				if fallbacksubject != "" &&
					csvr.ratingProfiles[p].getKey(fmt.Sprintf("%s:%s:%s:%s:%s", direction, tenant, tor, subject, FallbackDestination)) == nil {
					cd = &CallDescriptor{
						Direction:   direction,
						Tenant:      tenant,
						TOR:         tor,
						Subject:     subject,
						Destination: FallbackDestination,
						FallbackKey: fmt.Sprintf("%s:%s:%s:%s", direction, tenant, tor, fallbacksubject),
					}
					csvr.ratingProfiles[p] = append(csvr.ratingProfiles[p], cd)
				}
			}
		}
	}
}

func (csvr *CSVReader) LoadActions(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		units, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Could not parse action units: %v", err)
			continue
		}
		var a *Action
		if record[2] != MINUTES {
			a = &Action{
				ActionType: record[1],
				BalanceId:  record[2],
				Direction:  record[3],
				Units:      units,
			}
		} else {
			price, percent := 0.0, 0.0
			value, err := strconv.ParseFloat(record[7], 64)
			if err != nil {
				log.Printf("Could not parse action price: %v", err)
				continue
			}
			if record[6] == PERCENT {
				percent = value
			}
			if record[6] == ABSOLUTE {
				price = value
			}
			minutesWeight, err := strconv.ParseFloat(record[8], 64)
			if err != nil {
				log.Printf("Could not parse action minutes weight: %v", err)
				continue
			}
			weight, err := strconv.ParseFloat(record[9], 64)
			if err != nil {
				log.Printf("Could not parse action weight: %v", err)
				continue
			}
			a = &Action{
				ActionType: record[1],
				BalanceId:  record[2],
				Direction:  record[3],
				Weight:     weight,
				MinuteBucket: &MinuteBucket{
					Seconds:       units,
					Weight:        minutesWeight,
					Price:         price,
					Percent:       percent,
					DestinationId: record[5],
				},
			}
		}
		csvr.actions[tag] = append(csvr.actions[tag], a)
	}
}

func (csvr *CSVReader) LoadActionTimings(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}

		ts, exists := csvr.timings[record[2]]
		if !exists {
			log.Printf("Could not load the timing for tag: %v", record[2])
			continue
		}
		weight, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("Could not parse action timing weight: %v", err)
			continue
		}
		for _, t := range ts {
			at := &ActionTiming{
				Id:     GenUUID(),
				Tag:    record[2],
				Weight: weight,
				Timing: &Interval{
					Months:    t.Months,
					MonthDays: t.MonthDays,
					WeekDays:  t.WeekDays,
					StartTime: t.StartTime,
				},
				ActionsId: record[1],
			}
			csvr.actionsTimings[tag] = append(csvr.actionsTimings[tag], at)
		}
	}
}

func (csvr *CSVReader) LoadActionTriggers(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		tag := record[0]
		if tag == "Tag" {
			// skip header line
			continue
		}
		value, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("Could not parse action trigger value: %v", err)
			continue
		}
		weight, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			log.Printf("Could not parse action trigger weight: %v", err)
			continue
		}
		at := &ActionTrigger{
			BalanceId:      record[1],
			Direction:      record[2],
			ThresholdValue: value,
			DestinationId:  record[4],
			ActionsId:      record[5],
			Weight:         weight,
		}
		csvr.actionsTriggers[tag] = append(csvr.actionsTriggers[tag], at)
	}
}

func (csvr *CSVReader) LoadAccountActions(fn string, comma rune) {
	csvReader, fp, err := csvr.readerFunc(fn, comma)
	if err != nil {
		return
	}
	if fp != nil {
		defer fp.Close()
	}
	for record, err := csvReader.Read(); err == nil; record, err = csvReader.Read() {
		if record[0] == "Tenant" {
			continue
		}
		tag := fmt.Sprintf("%s:%s:%s", record[2], record[0], record[1])
		aTriggers, exists := csvr.actionsTriggers[record[4]]
		if !exists {
			log.Printf("Could not get action triggers for tag %v", record[4])
			continue
		}
		aTimingsTag := record[3]
		ub := &UserBalance{
			Type:           UB_TYPE_PREPAID,
			Id:             tag,
			ActionTriggers: aTriggers,
		}
		csvr.accountActions = append(csvr.accountActions, ub)

		aTimings, exists := csvr.actionsTimings[aTimingsTag]
		if !exists {
			log.Printf("Could not get action triggers for tag %v", aTimingsTag)
			// must not continue here
		}
		for _, at := range aTimings {
			at.UserBalanceIds = append(at.UserBalanceIds, tag)
		}
	}
}