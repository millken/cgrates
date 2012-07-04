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
	"github.com/simonz05/godis"
)

const (
	ACTION_TIMING_PREFIX = "acttmg"
)

type RedisStorage struct {
	dbNb int
	db   *godis.Client
	ms   *MarshalStrategy
}

func NewRedisStorage(address string, db int) (*RedisStorage, error) {
	ndb := godis.New(address, db, "")
	ms := &MarshalStrategy{}
	ms.SetMarshaler(&MyMarshaler{})
	return &RedisStorage{db: ndb, dbNb: db, ms: ms}, nil
}

func (rs *RedisStorage) Close() {
	rs.db.Quit()
}

func (rs *RedisStorage) Flush() error {
	return rs.db.Flushdb()
}

func (rs *RedisStorage) GetActivationPeriodsOrFallback(key string) (aps []*ActivationPeriod, fallbackKey string, err error) {
	elem, err := rs.db.Get(key)
	if err != nil {
		return
	}
	err = rs.ms.Unmarshal(elem, &aps)
	if err != nil {
		err = rs.ms.Unmarshal(elem, &fallbackKey)
	}
	return
}

func (rs *RedisStorage) SetActivationPeriodsOrFallback(key string, aps []*ActivationPeriod, fallbackKey string) (err error) {
	var result []byte
	if len(aps) > 0 {
		result, err = rs.ms.Marshal(aps)
	} else {
		result, err = rs.ms.Marshal(fallbackKey)
	}
	return rs.db.Set(key, result)
}

func (rs *RedisStorage) GetDestination(key string) (dest *Destination, err error) {
	if values, err := rs.db.Get(key); err == nil {
		dest = &Destination{Id: key}
		err = rs.ms.Unmarshal(values, dest)
	}
	return
}
func (rs *RedisStorage) SetDestination(dest *Destination) (err error) {
	result, err := rs.ms.Marshal(dest)
	return rs.db.Set(dest.Id, result)
}

func (rs *RedisStorage) GetActions(key string) (as []*Action, err error) {
	if values, err := rs.db.Get(key); err == nil {
		err = rs.ms.Unmarshal(values, &as)
	}
	return
}

func (rs *RedisStorage) SetActions(key string, as []*Action) (err error) {
	result, err := rs.ms.Marshal(as)
	return rs.db.Set(key, result)
}

func (rs *RedisStorage) GetUserBalance(key string) (ub *UserBalance, err error) {
	if values, err := rs.db.Get(key); err == nil {
		ub = &UserBalance{Id: key}
		err = rs.ms.Unmarshal(values, ub)
	}
	return
}

func (rs *RedisStorage) SetUserBalance(ub *UserBalance) (err error) {
	result, err := rs.ms.Marshal(ub)
	return rs.db.Set(ub.Id, result)
}

func (rs *RedisStorage) GetActionTimings(key string) (ats []*ActionTiming, err error) {
	if values, err := rs.db.Get(key); err == nil {
		err = rs.ms.Unmarshal(values, &ats)
	}
	return
}

func (rs *RedisStorage) SetActionTimings(key string, ats []*ActionTiming) (err error) {
	result, err := rs.ms.Marshal(ats)
	return rs.db.Set(key, result)
}

func (rs *RedisStorage) GetAllActionTimings() (ats []*ActionTiming, err error) {
	keys, err := rs.db.Keys(ACTION_TIMING_PREFIX + "*")
	if err != nil {
		return
	}
	values, err := rs.db.Mget(keys...)
	if err != nil {
		return
	}
	for _, v := range values.BytesArray() {
		var tempAts []*ActionTiming
		err = rs.ms.Unmarshal(v, &tempAts)
		ats = append(ats, tempAts...)
	}
	return
}
