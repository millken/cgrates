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

package cdrs

import (
	"encoding/json"
	"errors"
	"github.com/cgrates/cgrates/rater"
	"github.com/cgrates/cgrates/utils"
	"strconv"
	"time"
	"fmt"
	"crypto/sha1"
)

const (
	// Freswitch event property names
	CDR_MAP      = "variables"
	DIRECTION    = "direction"
	ORIG_ID      = "sip_call_id" //- originator_id - match cdrs
	SUBJECT      = "cgr_subject"
	ACCOUNT      = "cgr_account"
	DESTINATION  = "cgr_destination"
	REQTYPE      = "cgr_reqtype" //prepaid or postpaid
	TOR          = "cgr_tor"
	UUID         = "uuid" // -Unique ID for this call leg
	CSTMID       = "cgr_cstmid"
	CALL_DEST_NR = "dialed_extension"
	PARK_TIME    = "start_epoch"
	START_TIME   = "answer_epoch"
	END_TIME     = "end_epoch"
	DURATION     = "billsec"
	USERNAME     = "user_name"
	FS_IP        = "sip_local_network_addr"
	RATE         = "rate"
)

type FSCdr map[string]string

func (fsCdr FSCdr) New(body []byte) (rater.CDR, error) {
	fsCdr = make(map[string]string)
	var tmp map[string]interface{}
	var err error
	if err = json.Unmarshal(body, &tmp); err == nil {
		if variables, ok := tmp[CDR_MAP]; ok {
			if variables, ok := variables.(map[string]interface{}); ok {
				for k, v := range variables {
					fsCdr[k] = v.(string)
				}
			}
			return fsCdr, nil
		}
	}
	return nil, err
}

func (fsCdr FSCdr) GetCgrId() string {
	hasher := sha1.New()
	hasher.Write([]byte( fsCdr[FS_IP] ))
	hasher.Write([]byte( fsCdr[UUID] ))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
func (fsCdr FSCdr) GetAccId() string {
	return fsCdr[UUID]
}
func (fsCdr FSCdr) GetCdrHost() string {
	return fsCdr[FS_IP]
}
func (fsCdr FSCdr) GetDirection() string {
	//TODO: implement direction, not related to FS_DIRECTION but traffic towards or from subject/account
	return "OUT"
}
func (fsCdr FSCdr) GetOrigId() string {
	return fsCdr[ORIG_ID]
}
func (fsCdr FSCdr) GetSubject() string {
	return utils.FirstNonEmpty(fsCdr[SUBJECT], fsCdr[USERNAME])
}
func (fsCdr FSCdr) GetAccount() string {
	return utils.FirstNonEmpty(fsCdr[ACCOUNT], fsCdr[USERNAME])
}

// Charging destination number
func (fsCdr FSCdr) GetDestination() string {
	return utils.FirstNonEmpty(fsCdr[DESTINATION], fsCdr[CALL_DEST_NR])
}

// Original dialed destination number, useful in case of unpark
func (fsCdr FSCdr) GetCallDestNr() string {
	return fsCdr[CALL_DEST_NR]
}
func (fsCdr FSCdr) GetTOR() string {
	return utils.FirstNonEmpty(fsCdr[TOR], cfg.DefaultTOR)
}
func (fsCdr FSCdr) GetUUID() string {
	return fsCdr[UUID]
}
func (fsCdr FSCdr) GetTenant() string {
	return utils.FirstNonEmpty(fsCdr[CSTMID], cfg.DefaultTenant)
}
func (fsCdr FSCdr) GetReqType() string {
	return utils.FirstNonEmpty(fsCdr[REQTYPE], cfg.DefaultReqType)
}
func (fsCdr FSCdr) GetExtraParameters() string {
	return "" // ToDo: Add and extract from config
}
func (fsCdr FSCdr) GetFallbackSubj() string {
	return cfg.DefaultSubject
}
func (fsCdr FSCdr) GetStartTime() (t time.Time, err error) {
	st, err := strconv.ParseInt(fsCdr[START_TIME], 0, 64)
	t = time.Unix(0, st*1000)
	return
}

func (fsCdr FSCdr) GetEndTime() (t time.Time, err error) {
	st, err := strconv.ParseInt(fsCdr[END_TIME], 0, 64)
	t = time.Unix(0, st*1000)
	return
}
// Extracts duration as considered by the telecom switch
func (fsCdr FSCdr) GetDuration() int64 {
	dur, _ := strconv.ParseInt(fsCdr[DURATION], 0, 64)
	return dur
}


func (fsCdr FSCdr) GetRate() (float64, error) {
	rate, ok := fsCdr[RATE]
	if !ok {
		return -1, errors.New("Not found")
	}
	return strconv.ParseFloat(rate, 64)
}

func (fsCdr FSCdr) SetRate(rate float64) {
	fsCdr[RATE] = strconv.FormatFloat(rate, 'f', -1, 64)
}
