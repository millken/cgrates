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
package main

import (
	"flag"
	"github.com/cgrates/cgrates/timespans"
	"log"
)

var (
	separator                   = flag.String("separator", ",", "Default field separator")
	redisserver                 = flag.String("redisserver", "tcp:127.0.0.1:6379", "redis server address (tcp:127.0.0.1:6379)")
	redisdb                     = flag.Int("rdb", 10, "redis database number (10)")
	redispass                   = flag.String("pass", "", "redis database password")
	flush                       = flag.Bool("flush", false, "Flush the database before importing")
	monthsFn                    = flag.String("month", "Months.csv", "Months file")
	monthdaysFn                 = flag.String("monthdays", "MonthDays.csv", "Month days file")
	weekdaysFn                  = flag.String("weekdays", "WeekDays.csv", "Week days file")
	destinationsFn              = flag.String("destinations", "Destinations.csv", "Destinations file")
	ratesFn                     = flag.String("rates", "Rates.csv", "Rates file")
	timingsFn                   = flag.String("timings", "Timings.csv", "Timings file")
	ratestimingsFn              = flag.String("ratestimings", "RatesTimings.csv", "Rates timings file")
	ratingprofilesFn            = flag.String("ratingprofiles", "RatingProfiles.csv", "Rating profiles file")
	primaryBalanceActionsFn     = flag.String("primaryBalanceActions", "PrimaryBalanceActions.csv", "Primary balance actions file")
	destinationBalanceActionsFn = flag.String("destinationBalanceActions", "DestinationBalanceActions.csv", "Destination balance actions file")
	sep                         rune
)

func writeToDatabase() {
	storage, err := timespans.NewRedisStorage(*redisserver, *redisdb)
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err)
	}
	if *flush {
		storage.Flush()
	}
	for _, d := range destinations {
		storage.SetDestination(d)
	}
	for k, cds := range ratingProfiles {
		log.Print(k)
		for _, cd := range cds {
			storage.SetActivationPeriodsOrFallback(cd.GetKey(), cd.ActivationPeriods, cd.FallbackKey)
			log.Print(cd.GetKey())
		}
	}
}

func main() {
	flag.Parse()
	sep = []rune(*separator)[0]
	loadDestinations()
	loadRates()
	loadTimings()
	loadRatesTimings()
	loadRatingProfiles()
	loadPrimaryBalanceActions()
	loadDestinationBalanceActions()
	writeToDatabase()
}
