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
	"fmt"
	"github.com/cgrates/cgrates/utils"
	"strconv"
	"strings"
	"time"
)

const (
	// Freswitch event property names
	FS_CDR_MAP      = "variables"
	FS_DIRECTION    = "direction"
	FS_SUBJECT      = "cgr_subject"
	FS_ACCOUNT      = "cgr_account"
	FS_DESTINATION  = "cgr_destination"
	FS_REQTYPE      = "cgr_reqtype" //prepaid or postpaid
	FS_TOR          = "cgr_tor"
	FS_UUID         = "uuid" // -Unique ID for this call leg
	FS_CSTMID       = "cgr_tenant"
	FS_CALL_DEST_NR = "dialed_extension"
	FS_PARK_TIME    = "start_epoch"
	FS_ANSWER_TIME  = "answer_epoch"
	FS_HANGUP_TIME  = "end_epoch"
	FS_DURATION     = "billsec"
	FS_USERNAME     = "user_name"
	FS_IP           = "sip_local_network_addr"
	FS_CDR_SOURCE   = "freeswitch_json"
)

type FSCdr map[string]string

func (fsCdr FSCdr) New(body []byte) (utils.RawCDR, error) {
	fsCdr = make(map[string]string)
	var tmp map[string]interface{}
	var err error
	if err = json.Unmarshal(body, &tmp); err == nil {
		if variables, ok := tmp[FS_CDR_MAP]; ok {
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
	return utils.FSCgrId(fsCdr[FS_UUID])
}
func (fsCdr FSCdr) GetAccId() string {
	return fsCdr[FS_UUID]
}
func (fsCdr FSCdr) GetCdrHost() string {
	return fsCdr[FS_IP]
}
func (fsCdr FSCdr) GetCdrSource() string {
	return FS_CDR_SOURCE
}
func (fsCdr FSCdr) GetDirection() string {
	//TODO: implement direction, not related to FS_DIRECTION but traffic towards or from subject/account
	return "*out"
}
func (fsCdr FSCdr) GetSubject() string {
	return utils.FirstNonEmpty(fsCdr[FS_SUBJECT], fsCdr[FS_USERNAME])
}
func (fsCdr FSCdr) GetAccount() string {
	return utils.FirstNonEmpty(fsCdr[FS_ACCOUNT], fsCdr[FS_USERNAME])
}

// Charging destination number
func (fsCdr FSCdr) GetDestination() string {
	return utils.FirstNonEmpty(fsCdr[FS_DESTINATION], fsCdr[FS_CALL_DEST_NR])
}

func (fsCdr FSCdr) GetTOR() string {
	return utils.FirstNonEmpty(fsCdr[FS_TOR], cfg.DefaultTOR)
}

func (fsCdr FSCdr) GetTenant() string {
	return utils.FirstNonEmpty(fsCdr[FS_CSTMID], cfg.DefaultTenant)
}
func (fsCdr FSCdr) GetReqType() string {
	return utils.FirstNonEmpty(fsCdr[FS_REQTYPE], cfg.DefaultReqType)
}
func (fsCdr FSCdr) GetExtraFields() map[string]string {
	extraFields := make(map[string]string, len(cfg.CDRSExtraFields))
	for _, field := range cfg.CDRSExtraFields {
		extraFields[field] = fsCdr[field]
	}
	return extraFields
}
func (fsCdr FSCdr) GetAnswerTime() (t time.Time, err error) {
	//ToDo: Make sure we work with UTC instead of local time
	at, err := strconv.ParseInt(fsCdr[FS_ANSWER_TIME], 0, 64)
	t = time.Unix(at, 0)
	return
}
func (fsCdr FSCdr) GetHangupTime() (t time.Time, err error) {
	hupt, err := strconv.ParseInt(fsCdr[FS_HANGUP_TIME], 0, 64)
	t = time.Unix(hupt, 0)
	return
}

// Extracts duration as considered by the telecom switch
func (fsCdr FSCdr) GetDuration() time.Duration {
	dur, _ := utils.ParseDurationWithSecs(fsCdr[FS_DURATION])
	return dur
}

func (fsCdr FSCdr) Store() (result string, err error) {
	result += fsCdr.GetCgrId() + "|"
	result += fsCdr.GetAccId() + "|"
	result += fsCdr.GetCdrHost() + "|"
	result += fsCdr.GetDirection() + "|"
	result += fsCdr.GetSubject() + "|"
	result += fsCdr.GetAccount() + "|"
	result += fsCdr.GetDestination() + "|"
	result += fsCdr.GetTOR() + "|"
	result += fsCdr.GetAccId() + "|"
	result += fsCdr.GetTenant() + "|"
	result += fsCdr.GetReqType() + "|"
	st, err := fsCdr.GetAnswerTime()
	if err != nil {
		return "", err
	}
	result += strconv.FormatInt(st.UnixNano(), 10) + "|"
	et, err := fsCdr.GetHangupTime()
	if err != nil {
		return "", err
	}
	result += strconv.FormatInt(et.UnixNano(), 10) + "|"
	result += strconv.FormatInt(int64(fsCdr.GetDuration().Seconds()), 10) + "|"
	return
}

func (fsCdr FSCdr) Restore(input string) error {
	return errors.New("Not implemented")
}

// Used in extra mediation
func (fsCdr FSCdr) AsStoredCdr(runId, reqTypeFld, directionFld, tenantFld, torFld, accountFld, subjectFld, destFld, answerTimeFld, durationFld string, extraFlds []string, fieldsMandatory bool) (*utils.StoredCdr, error) {
	if utils.IsSliceMember([]string{runId, reqTypeFld, directionFld, tenantFld, torFld, accountFld, subjectFld, destFld, answerTimeFld, durationFld}, "") {
		return nil, errors.New(fmt.Sprintf("%s:FieldName", utils.ERR_MANDATORY_IE_MISSING)) // All input field names are mandatory
	}
	var err error
	var hasKey bool
	var aTimeStr, durStr string
	rtCdr := new(utils.StoredCdr)
	rtCdr.MediationRunId = runId
	rtCdr.Cost = -1.0 // Default for non-rated CDR
	if rtCdr.AccId = fsCdr.GetAccId(); len(rtCdr.AccId) == 0 {
		if fieldsMandatory {
			return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, utils.ACCID))
		} else { // Not mandatory, need to generate here CgrId
			rtCdr.CgrId = utils.GenUUID()
		}
	} else { // hasKey, use it to generate cgrid
		rtCdr.CgrId = utils.FSCgrId(rtCdr.AccId)
	}
	if rtCdr.CdrHost = fsCdr.GetCdrHost(); len(rtCdr.CdrHost) == 0 && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, utils.CDRHOST))
	}
	if rtCdr.CdrSource = fsCdr.GetCdrSource(); len(rtCdr.CdrSource) == 0 && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, utils.CDRSOURCE))
	}
	if strings.HasPrefix(reqTypeFld, utils.STATIC_VALUE_PREFIX) { // Values starting with prefix are not dynamically populated
		rtCdr.ReqType = reqTypeFld[1:]
	} else if rtCdr.ReqType, hasKey = fsCdr[reqTypeFld]; !hasKey && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, reqTypeFld))
	}
	if strings.HasPrefix(directionFld, utils.STATIC_VALUE_PREFIX) {
		rtCdr.Direction = directionFld[1:]
	} else if rtCdr.Direction, hasKey = fsCdr[directionFld]; !hasKey && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, directionFld))
	}
	if strings.HasPrefix(tenantFld, utils.STATIC_VALUE_PREFIX) {
		rtCdr.Tenant = tenantFld[1:]
	} else if rtCdr.Tenant, hasKey = fsCdr[tenantFld]; !hasKey && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, tenantFld))
	}
	if strings.HasPrefix(torFld, utils.STATIC_VALUE_PREFIX) {
		rtCdr.TOR = torFld[1:]
	} else if rtCdr.TOR, hasKey = fsCdr[torFld]; !hasKey && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, torFld))
	}
	if strings.HasPrefix(accountFld, utils.STATIC_VALUE_PREFIX) {
		rtCdr.Account = accountFld[1:]
	} else if rtCdr.Account, hasKey = fsCdr[accountFld]; !hasKey && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, accountFld))
	}
	if strings.HasPrefix(subjectFld, utils.STATIC_VALUE_PREFIX) {
		rtCdr.Subject = subjectFld[1:]
	} else if rtCdr.Subject, hasKey = fsCdr[subjectFld]; !hasKey && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, subjectFld))
	}
	if strings.HasPrefix(destFld, utils.STATIC_VALUE_PREFIX) {
		rtCdr.Destination = destFld[1:]
	} else if rtCdr.Destination, hasKey = fsCdr[destFld]; !hasKey && fieldsMandatory {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, destFld))
	}
	if aTimeStr, hasKey = fsCdr[answerTimeFld]; !hasKey && fieldsMandatory && !strings.HasPrefix(answerTimeFld, utils.STATIC_VALUE_PREFIX) {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, answerTimeFld))
	} else {
		if strings.HasPrefix(answerTimeFld, utils.STATIC_VALUE_PREFIX) {
			aTimeStr = answerTimeFld[1:]
		}
		if rtCdr.AnswerTime, err = utils.ParseTimeDetectLayout(aTimeStr); err != nil && fieldsMandatory {
			return nil, err
		}
	}
	if durStr, hasKey = fsCdr[durationFld]; !hasKey && fieldsMandatory && !strings.HasPrefix(durationFld, utils.STATIC_VALUE_PREFIX) {
		return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, durationFld))
	} else {
		if strings.HasPrefix(durationFld, utils.STATIC_VALUE_PREFIX) {
			durStr = durationFld[1:]
		}
		if rtCdr.Duration, err = utils.ParseDurationWithSecs(durStr); err != nil && fieldsMandatory {
			return nil, err
		}
	}
	rtCdr.ExtraFields = make(map[string]string, len(extraFlds))
	for _, fldName := range extraFlds {
		if fldVal, hasKey := fsCdr[fldName]; !hasKey && fieldsMandatory {
			return nil, errors.New(fmt.Sprintf("%s:%s", utils.ERR_MANDATORY_IE_MISSING, fldName))
		} else {
			rtCdr.ExtraFields[fldName] = fldVal
		}
	}
	return rtCdr, nil
}
