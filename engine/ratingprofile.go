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
	"errors"
	"fmt"
)

type RatingProfile struct {
	Id                                                                                             string
	FallbackKey                                                                                    string // FallbackKey is used as complete combination of Tenant:TOR:Direction:Subject
	DestinationMap                                                                                 map[string][]*RatingPlan
	Tag, Tenant, TOR, Direction, Subject, DestRatesTimingTag, RatesFallbackSubject, ActivationTime string // used only for loading
}

// Adds an activation period that applyes to current rating profile if not already present.
func (rp *RatingProfile) AddRatingPlanIfNotPresent(destInfo string, plans ...*RatingPlan) {
	if rp.DestinationMap == nil {
		rp.DestinationMap = make(map[string][]*RatingPlan, 1)
	}
	for _, plan := range plans {
		found := false
		for _, existingPlan := range rp.DestinationMap[destInfo] {
			if plan.Equal(existingPlan) {
				existingPlan.AddRateInterval(plan.RateIntervals...)
				found = true
				break
			}
		}
		if !found {
			rp.DestinationMap[destInfo] = append(rp.DestinationMap[destInfo], plan)
		}
	}
}

func (rp *RatingProfile) GetRatingPlansForPrefix(destPrefix string) (foundPrefix string, aps []*RatingPlan, err error) {
	bestPrecision := 0
	for dId, v := range rp.DestinationMap {
		precision, err := storageGetter.DestinationContainsPrefix(dId, destPrefix)
		if err != nil {
			Logger.Err(fmt.Sprintf("Error checking destination: %v", err))
			continue
		}
		if precision > bestPrecision {
			bestPrecision = precision
			aps = v
		}
	}

	if bestPrecision > 0 {
		return destPrefix[:bestPrecision], aps, nil
	}

	return "", nil, errors.New("not found")
}