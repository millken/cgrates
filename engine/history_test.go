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
	"strings"
	"testing"

	"github.com/cgrates/cgrates/history"
)

func TestHistoryRatinPlans(t *testing.T) {
	scribe := historyScribe.(*history.MockScribe)
	buf := scribe.BufMap[history.RATING_PROFILES_FN]
	if !strings.Contains(buf.String(), `{"Id":"*out:vdf:0:minu","RatingPlanActivations":[{"ActivationTime":"2012-01-01T00:00:00Z","RatingPlanId":"EVENING","FallbackKeys":null}]}`) {
		t.Error("Error in destination history content:", buf.String())
	}
}

func TestHistoryDestinations(t *testing.T) {
	scribe := historyScribe.(*history.MockScribe)
	buf := scribe.BufMap[history.DESTINATIONS_FN]
	expected := `[{"Id":"ALL","Prefixes":["49","41","43"]},
{"Id":"GERMANY","Prefixes":["49"]},
{"Id":"GERMANY_O2","Prefixes":["41"]},
{"Id":"GERMANY_PREMIUM","Prefixes":["43"]},
{"Id":"NAT","Prefixes":["0256","0257","0723"]},
{"Id":"PSTN_70","Prefixes":["+4970"]},
{"Id":"PSTN_71","Prefixes":["+4971"]},
{"Id":"PSTN_72","Prefixes":["+4972"]},
{"Id":"RET","Prefixes":["0723","0724"]},
{"Id":"nat","Prefixes":["0257","0256","0723"]}]`
	if buf.String() != expected {
		t.Error("Error in destination history content:", buf.String())
	}
}
