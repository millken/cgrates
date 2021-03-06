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

package utils

import (
	"net/url"
	"strconv"
	"time"
)

func NewStoredCdrFromRawCDR(rawcdr RawCDR) (*StoredCdr, error) {
	var err error
	rtCdr := new(StoredCdr)
	rtCdr.CgrId = rawcdr.GetCgrId()
	rtCdr.AccId = rawcdr.GetAccId()
	rtCdr.CdrHost = rawcdr.GetCdrHost()
	rtCdr.CdrSource = rawcdr.GetCdrSource()
	rtCdr.ReqType = rawcdr.GetReqType()
	rtCdr.Direction = rawcdr.GetDirection()
	rtCdr.Tenant = rawcdr.GetTenant()
	rtCdr.TOR = rawcdr.GetTOR()
	rtCdr.Account = rawcdr.GetAccount()
	rtCdr.Subject = rawcdr.GetSubject()
	rtCdr.Destination = rawcdr.GetDestination()
	if rtCdr.AnswerTime, err = rawcdr.GetAnswerTime(); err != nil {
		return nil, err
	}
	rtCdr.Duration = rawcdr.GetDuration()
	rtCdr.ExtraFields = rawcdr.GetExtraFields()
	rtCdr.MediationRunId = DEFAULT_RUNID
	rtCdr.Cost = -1
	return rtCdr, nil
}

// Rated CDR as extracted from StorDb. Kinda standard of internal CDR, complies to CDR interface also
type StoredCdr struct {
	CgrId          string
	AccId          string
	CdrHost        string
	CdrSource      string
	ReqType        string
	Direction      string
	Tenant         string
	TOR            string
	Account        string
	Subject        string
	Destination    string
	AnswerTime     time.Time
	Duration       time.Duration
	ExtraFields    map[string]string
	MediationRunId string
	Cost           float64
}

// Methods maintaining RawCDR interface

func (storedCdr *StoredCdr) GetCgrId() string {
	return storedCdr.CgrId
}

func (storedCdr *StoredCdr) GetAccId() string {
	return storedCdr.AccId
}

func (storedCdr *StoredCdr) GetCdrHost() string {
	return storedCdr.CdrHost
}

func (storedCdr *StoredCdr) GetCdrSource() string {
	return storedCdr.CdrSource
}

func (storedCdr *StoredCdr) GetDirection() string {
	return storedCdr.Direction
}

func (storedCdr *StoredCdr) GetSubject() string {
	return storedCdr.Subject
}

func (storedCdr *StoredCdr) GetAccount() string {
	return storedCdr.Account
}

func (storedCdr *StoredCdr) GetDestination() string {
	return storedCdr.Destination
}

func (storedCdr *StoredCdr) GetTOR() string {
	return storedCdr.TOR
}

func (storedCdr *StoredCdr) GetTenant() string {
	return storedCdr.Tenant
}

func (storedCdr *StoredCdr) GetReqType() string {
	return storedCdr.ReqType
}

func (storedCdr *StoredCdr) GetAnswerTime() (time.Time, error) {
	return storedCdr.AnswerTime, nil
}

func (storedCdr *StoredCdr) GetDuration() time.Duration {
	return storedCdr.Duration
}

func (storedCdr *StoredCdr) GetExtraFields() map[string]string {
	return storedCdr.ExtraFields
}

func (storedCdr *StoredCdr) AsStoredCdr(runId, reqTypeFld, directionFld, tenantFld, torFld, accountFld, subjectFld, destFld, answerTimeFld, durationFld string, extraFlds []string, fieldsMandatory bool) (*StoredCdr, error) {
	return storedCdr, nil
}

// Converts part of the rated Cdr as httpForm used to post remotely to CDRS
func (storedCdr *StoredCdr) AsRawCdrHttpForm() url.Values {
	v := url.Values{}
	v.Set(ACCID, storedCdr.AccId)
	v.Set(CDRHOST, storedCdr.CdrHost)
	v.Set(CDRSOURCE, storedCdr.CdrSource)
	v.Set(REQTYPE, storedCdr.ReqType)
	v.Set(DIRECTION, storedCdr.Direction)
	v.Set(TENANT, storedCdr.Tenant)
	v.Set(TOR, storedCdr.TOR)
	v.Set(ACCOUNT, storedCdr.Account)
	v.Set(SUBJECT, storedCdr.Subject)
	v.Set(DESTINATION, storedCdr.Destination)
	v.Set(ANSWER_TIME, storedCdr.AnswerTime.String())
	v.Set(DURATION, strconv.FormatFloat(storedCdr.Duration.Seconds(), 'f', -1, 64))
	for fld, val := range storedCdr.ExtraFields {
		v.Set(fld, val)
	}
	return v
}
