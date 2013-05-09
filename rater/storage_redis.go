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

package rater

import (
	"errors"
	"fmt"
	"github.com/fzzy/radix/redis"
	"time"
)

type RedisStorage struct {
	dbNb int
	db   *redis.Client
	ms   Marshaler
}

func NewRedisStorage(address string, db int, pass string) (DataStorage, error) {
	ndb, err := redis.DialTimeout("tcp", address, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}
	if pass != "" {
		r := ndb.Cmd("auth", pass)
		if r.Err != nil {
			return nil, errors.New(r.Err.Error())
		}
	}
	if db > 0 {
		r := ndb.Cmd("select", db)
		if r.Err != nil {
			return nil, errors.New(r.Err.Error())
		}
	}
	ms := new(MyMarshaler)
	return &RedisStorage{db: ndb, dbNb: db, ms: ms}, nil
}

func (rs *RedisStorage) Close() {
	rs.db.Close()
}

func (rs *RedisStorage) Flush() (err error) {
	r := rs.db.Cmd("flushdb")
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}

func (rs *RedisStorage) GetRatingProfile(key string) (rp *RatingProfile, err error) {
	if values, err := rs.db.Cmd("get", RATING_PROFILE_PREFIX+key).Bytes(); err == nil {
		rp = new(RatingProfile)
		err = rs.ms.Unmarshal(values, rp)
	} else {
		return nil, err
	}
	return
}

func (rs *RedisStorage) SetRatingProfile(rp *RatingProfile) (err error) {
	result, err := rs.ms.Marshal(rp)
	r := rs.db.Cmd("set", RATING_PROFILE_PREFIX+rp.Id, result)
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}

func (rs *RedisStorage) GetDestination(key string) (dest *Destination, err error) {
	if values, err := rs.db.Cmd("get", DESTINATION_PREFIX+key).Bytes(); err == nil {
		dest = &Destination{Id: key}
		err = rs.ms.Unmarshal(values, dest)
	} else {
		return nil, err
	}
	return
}

func (rs *RedisStorage) SetDestination(dest *Destination) (err error) {
	result, err := rs.ms.Marshal(dest)
	if err != nil {
		return err
	}
	r := rs.db.Cmd("set", DESTINATION_PREFIX+dest.Id, result)
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}

func (rs *RedisStorage) GetActions(key string) (as []*Action, err error) {
	if values, err := rs.db.Cmd("get", ACTION_PREFIX+key).Bytes(); err == nil {
		err = rs.ms.Unmarshal(values, &as)
	} else {
		return nil, err
	}
	return
}

func (rs *RedisStorage) SetActions(key string, as []*Action) (err error) {
	result, err := rs.ms.Marshal(as)
	r := rs.db.Cmd("set", ACTION_PREFIX+key, result)
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}

func (rs *RedisStorage) GetUserBalance(key string) (ub *UserBalance, err error) {
	if values, err := rs.db.Cmd("get", USER_BALANCE_PREFIX+key).Bytes(); err == nil {
		ub = &UserBalance{Id: key}
		err = rs.ms.Unmarshal(values, ub)
	} else {
		return nil, err
	}
	return
}

func (rs *RedisStorage) SetUserBalance(ub *UserBalance) (err error) {
	result, err := rs.ms.Marshal(ub)
	r := rs.db.Cmd("set", USER_BALANCE_PREFIX+ub.Id, result)
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}

func (rs *RedisStorage) GetActionTimings(key string) (ats []*ActionTiming, err error) {
	if values, err := rs.db.Cmd("get", ACTION_TIMING_PREFIX+key).Bytes(); err == nil {
		err = rs.ms.Unmarshal(values, &ats)
	} else {
		return nil, err
	}
	return
}

func (rs *RedisStorage) SetActionTimings(key string, ats []*ActionTiming) (err error) {
	if len(ats) == 0 {
		// delete the key
		r := rs.db.Cmd("del", ACTION_TIMING_PREFIX+key)
		if r.Err != nil {
			return errors.New(r.Err.Error())
		}
		return
	}
	result, err := rs.ms.Marshal(ats)
	r := rs.db.Cmd("set", ACTION_TIMING_PREFIX+key, result)
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}

func (rs *RedisStorage) GetAllActionTimings() (ats map[string][]*ActionTiming, err error) {
	keys, err := rs.db.Cmd("keys", ACTION_TIMING_PREFIX+"*").List()
	if err != nil {
		return
	}
	ats = make(map[string][]*ActionTiming, len(keys))
	for _, key := range keys {
		values, err := rs.db.Cmd("get", key).Bytes()
		if err != nil {
			continue
		}
		var tempAts []*ActionTiming
		err = rs.ms.Unmarshal(values, &tempAts)
		ats[key[len(ACTION_TIMING_PREFIX):]] = tempAts
	}

	return
}

func (rs *RedisStorage) LogCallCost(uuid, source string, cc *CallCost) (err error) {
	result, err := rs.ms.Marshal(cc)
	if err != nil {
		return
	}
	r := rs.db.Cmd("set", LOG_CALL_COST_PREFIX+source+"_"+uuid, result)
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}

func (rs *RedisStorage) GetCallCostLog(uuid, source string) (cc *CallCost, err error) {
	if values, err := rs.db.Cmd("get", LOG_CALL_COST_PREFIX+source+"_"+uuid).Bytes(); err == nil {
		err = rs.ms.Unmarshal(values, cc)
	} else {
		return nil, err
	}
	return
}

func (rs *RedisStorage) LogActionTrigger(ubId, source string, at *ActionTrigger, as []*Action) (err error) {
	mat, err := rs.ms.Marshal(at)
	if err != nil {
		return
	}
	mas, err := rs.ms.Marshal(as)
	if err != nil {
		return
	}
	rs.db.Cmd("set", LOG_ACTION_TRIGGER_PREFIX+source+"_"+time.Now().Format(time.RFC3339Nano), []byte(fmt.Sprintf("%s*%s*%s", ubId, string(mat), string(mas))))
	return
}

func (rs *RedisStorage) LogActionTiming(source string, at *ActionTiming, as []*Action) (err error) {
	mat, err := rs.ms.Marshal(at)
	if err != nil {
		return
	}
	mas, err := rs.ms.Marshal(as)
	if err != nil {
		return
	}
	rs.db.Cmd("set", LOG_ACTION_TIMMING_PREFIX+source+"_"+time.Now().Format(time.RFC3339Nano), []byte(fmt.Sprintf("%s*%s", string(mat), string(mas))))
	return
}

func (rs *RedisStorage) LogError(uuid, source, errstr string) (err error) {
	r := rs.db.Cmd("set", LOG_ERR+source+"_"+uuid, errstr)
	if r.Err != nil {
		return errors.New(r.Err.Error())
	}
	return
}
