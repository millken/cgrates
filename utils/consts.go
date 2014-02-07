package utils

const (
	VERSION                    = "0.9.1c4"
	POSTGRES                   = "postgres"
	MYSQL                      = "mysql"
	MONGO                      = "mongo"
	REDIS                      = "redis"
	LOCALHOST                  = "127.0.0.1"
	FSCDR_FILE_CSV             = "freeswitch_file_csv"
	FSCDR_HTTP_JSON            = "freeswitch_http_json"
	NOT_IMPLEMENTED            = "not implemented"
	PREPAID                    = "prepaid"
	POSTPAID                   = "postpaid"
	PSEUDOPREPAID              = "pseudoprepaid"
	RATED                      = "rated"
	ERR_NOT_IMPLEMENTED        = "NOT_IMPLEMENTED"
	ERR_SERVER_ERROR           = "SERVER_ERROR"
	ERR_NOT_FOUND              = "NOT_FOUND"
	ERR_MANDATORY_IE_MISSING   = "MANDATORY_IE_MISSING"
	ERR_EXISTS                 = "EXISTS"
	ERR_BROKEN_REFERENCE       = "BROKEN_REFERENCE"
	TBL_TP_TIMINGS             = "tp_timings"
	TBL_TP_DESTINATIONS        = "tp_destinations"
	TBL_TP_RATES               = "tp_rates"
	TBL_TP_DESTINATION_RATES   = "tp_destination_rates"
	TBL_TP_RATING_PLANS        = "tp_rating_plans"
	TBL_TP_RATE_PROFILES       = "tp_rating_profiles"
	TBL_TP_ACTIONS             = "tp_actions"
	TBL_TP_ACTION_PLANS        = "tp_action_plans"
	TBL_TP_ACTION_TRIGGERS     = "tp_action_triggers"
	TBL_TP_ACCOUNT_ACTIONS     = "tp_account_actions"
	TBL_CDRS_PRIMARY           = "cdrs_primary"
	TBL_CDRS_EXTRA             = "cdrs_extra"
	TBL_COST_DETAILS           = "cost_details"
	TBL_RATED_CDRS             = "rated_cdrs"
	TIMINGS_CSV                = "Timings.csv"
	DESTINATIONS_CSV           = "Destinations.csv"
	RATES_CSV                  = "Rates.csv"
	DESTINATION_RATES_CSV      = "DestinationRates.csv"
	RATING_PLANS_CSV           = "RatingPlans.csv"
	RATING_PROFILES_CSV        = "RatingProfiles.csv"
	SHARED_GROUPS_CSV          = "SharedGroups.csv"
	ACTIONS_CSV                = "Actions.csv"
	ACTION_PLANS_CSV           = "ActionPlans.csv"
	ACTION_TRIGGERS_CSV        = "ActionTriggers.csv"
	ACCOUNT_ACTIONS_CSV        = "AccountActions.csv"
	TIMINGS_NRCOLS             = 6
	DESTINATIONS_NRCOLS        = 2
	RATES_NRCOLS               = 8
	DESTINATION_RATES_NRCOLS   = 3
	DESTRATE_TIMINGS_NRCOLS    = 4
	RATE_PROFILES_NRCOLS       = 7
	SHARED_GROUPS_NRCOLS       = 5
	ACTIONS_NRCOLS             = 12
	ACTION_PLANS_NRCOLS        = 4
	ACTION_TRIGGERS_NRCOLS     = 8
	ACCOUNT_ACTIONS_NRCOLS     = 5
	ROUNDING_UP                = "*up"
	ROUNDING_MIDDLE            = "*middle"
	ROUNDING_DOWN              = "*down"
	ANY                        = "*any"
	COMMENT_CHAR               = '#'
	CSV_SEP                    = ','
	FALLBACK_SEP               = ';'
	JSON                       = "json"
	MSGPACK                    = "msgpack"
	CSV_LOAD                   = "CSVLOAD"
	CGRID                      = "cgrid"
	ACCID                      = "accid"
	CDRHOST                    = "cdrhost"
	CDRSOURCE                  = "cdrsource"
	REQTYPE                    = "reqtype"
	DIRECTION                  = "direction"
	TENANT                     = "tenant"
	TOR                        = "tor"
	ACCOUNT                    = "account"
	SUBJECT                    = "subject"
	DESTINATION                = "destination"
	ANSWER_TIME                = "answer_time"
	DURATION                   = "duration"
	DEFAULT_RUNID              = "default"
	STATIC_VALUE_PREFIX        = "^"
	CDRE_CSV                   = "csv"
	CDRE_DRYRUN                = "dry_run"
	INTERNAL                   = "internal"
	ZERO_RATING_SUBJECT_PREFIX = "*zero"
)

var (
	CdreCdrFormats = []string{CDRE_CSV, CDRE_DRYRUN}
)
