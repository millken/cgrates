/*
Rating system designed to be used in VoIP Carriems World
Copyright (C) 2013 ITsysCOM

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either vemsion 3 of the License, or
(at your option) any later vemsion.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package rater

import (
	"errors"
	"fmt"
	"github.com/cgrates/cgrates/utils"
	"strings"
	"time"
)

type MapStorage struct {
	dict map[string][]byte
	ms   Marshaler
}

func NewMapStorage() (DataStorage, error) {
	return &MapStorage{dict: make(map[string][]byte), ms: new(MyMarshaler)}, nil
}

func (ms *MapStorage) Close() {}

func (ms *MapStorage) Flush() error {
	ms.dict = make(map[string][]byte)
	return nil
}

func (ms *MapStorage) GetRatingProfile(key string) (rp *RatingProfile, err error) {
	if values, ok := ms.dict[RATING_PROFILE_PREFIX+key]; ok {
		rp = new(RatingProfile)
		err = ms.ms.Unmarshal(values, rp)
	} else {
		return nil, errors.New("not found")
	}
	return
}

func (ms *MapStorage) SetRatingProfile(rp *RatingProfile) (err error) {
	result, err := ms.ms.Marshal(rp)
	ms.dict[RATING_PROFILE_PREFIX+rp.Id] = result
	return
}

func (ms *MapStorage) GetDestination(key string) (dest *Destination, err error) {
	if values, ok := ms.dict[DESTINATION_PREFIX+key]; ok {
		dest = &Destination{Id: key}
		err = ms.ms.Unmarshal(values, dest)
	} else {
		return nil, errors.New("not found")
	}
	return
}
func (ms *MapStorage) SetDestination(dest *Destination) (err error) {
	result, err := ms.ms.Marshal(dest)
	ms.dict[DESTINATION_PREFIX+dest.Id] = result
	return
}

// Extracts destinations from StorDB on specific tariffplan id
func (ms *MapStorage) GetTPDestination(tpid, destTag string) (*Destination, error) {
	return nil, nil
}

func (ms *MapStorage) SetTPDestination(tpid string, dest *Destination) error {
	return errors.New(utils.ERR_NOT_IMPLEMENTED)
}

func (ms *MapStorage) GetActions(key string) (as Actions, err error) {
	if values, ok := ms.dict[ACTION_PREFIX+key]; ok {
		err = ms.ms.Unmarshal(values, &as)
	} else {
		return nil, errors.New("not found")
	}
	return
}

func (ms *MapStorage) SetActions(key string, as Actions) (err error) {
	result, err := ms.ms.Marshal(&as)
	ms.dict[ACTION_PREFIX+key] = result
	return
}

func (ms *MapStorage) GetUserBalance(key string) (ub *UserBalance, err error) {
	if values, ok := ms.dict[USER_BALANCE_PREFIX+key]; ok {
		ub = &UserBalance{Id: key}
		err = ms.ms.Unmarshal(values, ub)
	} else {
		return nil, errors.New("not found")
	}
	return
}

func (ms *MapStorage) SetUserBalance(ub *UserBalance) (err error) {
	result, err := ms.ms.Marshal(ub)
	ms.dict[USER_BALANCE_PREFIX+ub.Id] = result
	return
}

func (ms *MapStorage) GetActionTimings(key string) (ats ActionTimings, err error) {
	if values, ok := ms.dict[ACTION_TIMING_PREFIX+key]; ok {
		err = ms.ms.Unmarshal(values, &ats)
	} else {
		return nil, errors.New("not found")
	}
	return
}

func (ms *MapStorage) SetActionTimings(key string, ats ActionTimings) (err error) {
	if len(ats) == 0 {
		// delete the key
		delete(ms.dict, ACTION_TIMING_PREFIX+key)
		return
	}
	result, err := ms.ms.Marshal(&ats)
	ms.dict[ACTION_TIMING_PREFIX+key] = result
	return
}

func (ms *MapStorage) GetAllActionTimings() (ats map[string]ActionTimings, err error) {
	ats = make(map[string]ActionTimings)
	for key, value := range ms.dict {
		if !strings.Contains(key, ACTION_TIMING_PREFIX) {
			continue
		}
		var tempAts ActionTimings
		err = ms.ms.Unmarshal(value, &tempAts)
		ats[key[len(ACTION_TIMING_PREFIX):]] = tempAts
	}

	return
}

func (ms *MapStorage) LogCallCost(uuid, source string, cc *CallCost) error {
	result, err := ms.ms.Marshal(cc)
	ms.dict[LOG_CALL_COST_PREFIX+source+"_"+uuid] = result
	return err
}

func (ms *MapStorage) GetCallCostLog(uuid, source string) (cc *CallCost, err error) {
	if values, ok := ms.dict[LOG_CALL_COST_PREFIX+source+"_"+uuid]; ok {
		err = ms.ms.Unmarshal(values, &cc)
	} else {
		return nil, errors.New("not found")
	}
	return
}

func (ms *MapStorage) LogActionTrigger(ubId, source string, at *ActionTrigger, as Actions) (err error) {
	mat, err := ms.ms.Marshal(at)
	if err != nil {
		return
	}
	mas, err := ms.ms.Marshal(&as)
	if err != nil {
		return
	}
	ms.dict[LOG_ACTION_TRIGGER_PREFIX+source+"_"+time.Now().Format(time.RFC3339Nano)] = []byte(fmt.Sprintf("%s*%s*%s", ubId, string(mat), string(mas)))
	return
}

func (ms *MapStorage) LogActionTiming(source string, at *ActionTiming, as Actions) (err error) {
	mat, err := ms.ms.Marshal(at)
	if err != nil {
		return
	}
	mas, err := ms.ms.Marshal(&as)
	if err != nil {
		return
	}
	ms.dict[LOG_ACTION_TIMMING_PREFIX+source+"_"+time.Now().Format(time.RFC3339Nano)] = []byte(fmt.Sprintf("%s*%s", string(mat), string(mas)))
	return
}

func (ms *MapStorage) LogError(uuid, source, errstr string) (err error) {
	ms.dict[LOG_ERR+source+"_"+uuid] = []byte(errstr)
	return nil
}

func (ms *MapStorage) SetCdr(utils.CDR) error {
	return nil
}

func (ms *MapStorage) SetRatedCdr(utils.CDR, *CallCost) error {
	return nil
}

func (ms *MapStorage) GetAllRatedCdr() ([]utils.CDR, error) {
	return nil, nil
}

func (ms *MapStorage) GetTpDestinations(tpid, tag string) ([]*Destination, error) {
	return nil, nil
}

func (ms *MapStorage) GetTpRates(tpid, tag string) (map[string][]*Rate, error) {
	return nil, nil
}
func (ms *MapStorage) GetTpTimings(tpid, tag string) (map[string]*Timing, error) {
	return nil, nil
}
func (ms *MapStorage) GetTpRateTimings(tpid, tag string) ([]*RateTiming, error) {
	return nil, nil
}
func (ms *MapStorage) GetTpRatingProfiles(tpid, tag string) (map[string]*RatingProfile, error) {
	return nil, nil
}
func (ms *MapStorage) GetTpActions(tpid, tag string) (map[string][]*Action, error) {
	return nil, nil
}
func (ms *MapStorage) GetTpActionTimings(tpid, tag string) (map[string][]*ActionTiming, error) {
	return nil, nil
}
func (ms *MapStorage) GetTpActionTriggers(tpid, tag string) (map[string][]*ActionTrigger, error) {
	return nil, nil
}
func (ms *MapStorage) GetTpAccountActions(tpid, tag string) ([]*AccountAction, error) {
	return nil, nil
}
