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

package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/config"
)

var (
	cgrConfig, _ = config.NewDefaultCGRConfig()
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
	runs       = flag.Int("runs", 10000, "stress cycle number")
	parallel   = flag.Int("parallel", 0, "run n requests in parallel")
	ratingdb_type = flag.String("ratingdb_type", cgrConfig.RatingDBType, "The type of the RatingDb database <redis>")
	ratingdb_host = flag.String("ratingdb_host", cgrConfig.RatingDBHost, "The RatingDb host to connect to.")
	ratingdb_port = flag.String("ratingdb_port", cgrConfig.RatingDBPort, "The RatingDb port to bind to.")
	ratingdb_name = flag.String("ratingdb_name", cgrConfig.RatingDBName, "The name/number of the RatingDb to connect to.")
	ratingdb_user = flag.String("ratingdb_user", cgrConfig.RatingDBUser, "The RatingDb user to sign in as.")
	ratingdb_pass = flag.String("ratingdb_passwd", cgrConfig.RatingDBPass, "The RatingDb user's password.")
	accountdb_type = flag.String("accountdb_type", cgrConfig.AccountDBType, "The type of the AccountingDb database <redis>")
	accountdb_host = flag.String("accountdb_host", cgrConfig.AccountDBHost, "The AccountingDb host to connect to.")
	accountdb_port = flag.String("accountdb_port", cgrConfig.AccountDBPort, "The AccountingDb port to bind to.")
	accountdb_name = flag.String("accountdb_name", cgrConfig.AccountDBName, "The name/number of the AccountingDb to connect to.")
	accountdb_user = flag.String("accountdb_user", cgrConfig.AccountDBUser, "The AccountingDb user to sign in as.")
	accountdb_pass = flag.String("accountdb_passwd", cgrConfig.AccountDBPass, "The AccountingDb user's password.")
	dbdata_encoding = flag.String("dbdata_encoding", cgrConfig.DBDataEncoding, "The encoding used to store object data in strings.")
	raterAddress  = flag.String("rater_address", "", "Rater address for remote tests. Empty for internal rater.")
	rpcEncoding   = flag.String("rpc_encoding", cgrConfig.RPCEncoding, "Rpc encoding to use when talking to remote rater <json|gob>")
	tor = flag.String("tor", "call", "The type of record to use in queries.")
	tenant = flag.String("tenant", "call", "The type of record to use in queries.")
	subject = flag.String("subject", "1001", "The rating subject to use in queries.")
	destination = flag.String("destination", "+4986517174963", "The destination to use in queries.")
	
	nilDuration = time.Duration(0)
)

func durInternalRater( cd *engine.CallDescriptor) (time.Duration, error) {
	ratingDb, err := engine.ConfigureRatingStorage(*ratingdb_type, *ratingdb_host, *ratingdb_port, *ratingdb_name, *ratingdb_user, *ratingdb_pass, *dbdata_encoding)
	if err != nil {
		return nilDuration, fmt.Errorf("Could not connect to rating database: %s", err.Error())
	}
	defer ratingDb.Close()
	engine.SetRatingStorage(ratingDb)
	accountDb, err := engine.ConfigureAccountingStorage(*accountdb_type, *accountdb_host, *accountdb_port, *accountdb_name, *accountdb_user, *accountdb_pass, *dbdata_encoding)
	if err != nil {
		return nilDuration, fmt.Errorf("Could not connect to accounting database: %s", err.Error())
	}
	defer accountDb.Close()
	engine.SetAccountingStorage(accountDb)
	if err := ratingDb.CacheRating(nil, nil, nil); err != nil {
		return nilDuration, fmt.Errorf("Cache rating error: %s", err.Error())
	}
	log.Printf("Runnning %d cycles...", *runs)
	var result *engine.CallCost
	j := 0
	start := time.Now()
	for i := 0; i < *runs; i++ {
		result, err = cd.GetCost()
		if *memprofile != "" {
			runtime.MemProfileRate = 1
			runtime.GC()
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.WriteHeapProfile(f)
			f.Close()
			break
		}
		j = i
	}
	log.Print(result, j, err)
	memstats := new(runtime.MemStats)
	runtime.ReadMemStats(memstats)
	log.Printf("memstats before GC: Kbytes = %d footprint = %d",
		memstats.HeapAlloc/1024, memstats.Sys/1024)
	return time.Since(start), nil
}


func durRemoteRater( cd *engine.CallDescriptor) (time.Duration, error) {
	result := engine.CallCost{}
	var client *rpc.Client
	var err error
	if *rpcEncoding=="json" {
		client, err = jsonrpc.Dial("tcp", *raterAddress)
	} else {
		client, err = rpc.Dial("tcp", *raterAddress)
	}
	if err != nil {
		return nilDuration, fmt.Errorf("Could not connect to engine: ", err.Error())
	}
	defer client.Close()
	start := time.Now()
	if *parallel > 0 {
		// var divCall *rpc.Call
		var sem = make(chan int, *parallel)
		var finish = make(chan int)
		for i := 0; i < *runs; i++ {
			go func() {
				sem <- 1
				client.Call("Responder.GetCost", cd, &result)
				<-sem
				finish <- 1
				// divCall = client.Go("Responder.GetCost", cd, &result, nil)
			}()
		}
		for i := 0; i < *runs; i++ {
			<-finish
		}
		// <-divCall.Done
	} else {
		for j := 0; j < *runs; j++ {
			client.Call("Responder.GetCost", cd, &result)
		}
	}
	log.Println(result)
	return time.Since(start), nil
}
	
	


func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	cd := &engine.CallDescriptor{
		TimeStart:    time.Date(2013, time.December, 13, 22, 30, 0, 0, time.UTC),
		TimeEnd:      time.Date(2013, time.December, 13, 22, 31, 0, 0, time.UTC),
		CallDuration: 60 * time.Second,
		Direction:    "*out",
		TOR:          *tor,
		Tenant:       *tenant,
		Subject:      *subject,
		Destination:  *destination,
	}
	var duration time.Duration
	var err error
	if len(*raterAddress) == 0 {
		duration, err = durInternalRater(cd)
	} else {
		duration, err = durRemoteRater(cd)
	}
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("Elapsed: %d resulted: %f req/s.", duration, float64(*runs)/duration.Seconds())
	}
}
