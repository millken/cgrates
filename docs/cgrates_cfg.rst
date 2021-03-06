cgr-engine configuration file
=============================

Organized into configuration sections. All configuration options come with defaults and we have tried our best to choose the best ones for a minimum of efforts necessary when running.

Bellow is the default configuration file which comes hardcoded into cgr-engine, most of them being explained and exemplified there.

::

 [global]
 # ratingdb_type = redis 			# Rating subsystem database: <redis>.
 # ratingdb_host = 127.0.0.1 			# Rating subsystem database host address.
 # ratingdb_port = 6379 				# Rating subsystem port to reach the database.
 # ratingdb_name = 10 				# Rating subsystem database name to connect to.
 # ratingdb_user =		 		# Rating subsystem username to use when connecting to database.
 # ratingdb_passwd =				# Rating subsystem password to use when connecting to database.
 # accountdb_type = redis 			# Accounting subsystem database: <redis>.
 # accountdb_host = 127.0.0.1 			# Accounting subsystem database host address.
 # accountdb_port = 6379 			# Accounting subsystem port to reach the database.
 # accountdb_name = 11				# Accounting subsystem database name to connect to.
 # accountdb_user =		 		# Accounting subsystem username to use when connecting to database.
 # accountdb_passwd =				# Accounting subsystem password to use when connecting to database.
 # stordb_type = mysql				# Stor database type to use: <mysql>
 # stordb_host = 127.0.0.1 			# The host to connect to. Values that start with / are for UNIX domain sockets.
 # stordb_port = 3306				# The port to reach the logdb.
 # stordb_name = cgrates 			# The name of the log database to connect to.
 # stordb_user = cgrates	 			# Username to use when connecting to stordb.
 # stordb_passwd = CGRateS.org			# Password to use when connecting to stordb.
 # dbdata_encoding = msgpack			# The encoding used to store object data in strings: <msgpack|json>
 # rpc_json_listen = 127.0.0.1:2012		# RPC JSON listening address
 # rpc_gob_listen = 127.0.0.1:2013		# RPC GOB listening address
 # http_listen = 127.0.0.1:2080			# HTTP listening address
 # default_reqtype = rated			# Default request type to consider when missing from requests: <""|prepaid|postpaid|pseudoprepaid|rated>.
 # default_tor = call				# Default Type of Record to consider when missing from requests.
 # default_tenant = cgrates.org			# Default Tenant to consider when missing from requests.
 # default_subject = cgrates			# Default rating Subject to consider when missing from requests.
 # rounding_method = *middle			# Rounding method for floats/costs: <*up|*middle|*down>
 # rounding_decimals = 4				# Number of decimals to round float/costs at

 [balancer]
 # enabled = false 				# Start Balancer service: <true|false>.

 [rater]
 # enabled = false				# Enable RaterCDRSExportPath service: <true|false>.
 # balancer =  					# Register to Balancer as worker: <""|internal|127.0.0.1:2013>.

 [scheduler]
 # enabled = false				# Starts Scheduler service: <true|false>.

 [cdrs]
 # enabled = false				# Start the CDR Server service:  <true|false>.
 # extra_fields = 				# Extra fields to store in CDRs
 # mediator = 					# Address where to reach the Mediator. Empty for disabling mediation. <""|internal>

 [cdre]
 # cdr_format = csv					# Exported CDRs format <csv>
 # extra_fields = 					# List of extra fields to be exported out in CDRs
 # export_dir = /var/log/cgrates/cdr/cdrexport/csv	# Path where the exported CDRs will be placed

 [cdrc]
 # enabled = false				# Enable CDR client functionality
 # cdrs = internal				# Address where to reach CDR server. <internal|127.0.0.1:2080>
 # cdrs_method = http_cgr			# Mechanism to use when posting CDRs on server  <http_cgr>
 # run_delay = 0					# Sleep interval in seconds between consecutive runs, 0 to use automation via inotify
 # cdr_type = csv				# CDR file format <csv|freeswitch_csv>.
 # cdr_in_dir = /var/log/cgrates/cdr/cdrc/in 	# Absolute path towards the directory where the CDRs are stored.
 # cdr_out_dir =	/var/log/cgrates/cdr/cdrc/out	# Absolute path towards the directory where processed CDRs will be moved.
 # cdr_source_id = freeswitch_csv		# Free form field, tag identifying the source of the CDRs within CGRS database.
 # accid_field = 0				# Accounting id field identifier. Use index number in case of .csv cdrs.
 # reqtype_field = 1				# Request type field identifier. Use index number in case of .csv cdrs.
 # direction_field = 2				# Direction field identifier. Use index numbers in case of .csv cdrs.
 # tenant_field = 3				# Tenant field identifier. Use index numbers in case of .csv cdrs.
 # tor_field = 4					# Type of Record field identifier. Use index numbers in case of .csv cdrs.
 # account_field = 5				# Account field identifier. Use index numbers in case of .csv cdrs.
 # subject_field = 6				# Subject field identifier. Use index numbers in case of .csv CDRs.
 # destination_field = 7				# Destination field identifier. Use index numbers in case of .csv cdrs.
 # answer_time_field = 8				# Answer time field identifier. Use index numbers in case of .csv cdrs.
 # duration_field = 9				# Duration field identifier. Use index numbers in case of .csv cdrs.
 # extra_fields = 				# Extra fields identifiers. For .csv, format: <label_extrafield_1>:<index_extrafield_1>[...,<label_extrafield_n>:<index_extrafield_n>]

 [mediator]
 # enabled = false				# Starts Mediator service: <true|false>.
 # rater = internal				# Address where to reach the Rater: <internal|x.y.z.y:1234>
 # rater_reconnects = 3				# Number of reconnects to rater before giving up.
 # run_ids = 					# Identifiers of each extra mediation to run on CDRs
 # reqtype_fields = 				# Name of request type fields to be used during extra mediation. Use index number in case of .csv cdrs.
 # direction_fields = 				# Name of direction fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 # tenant_fields = 				# Name of tenant fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 # tor_fields = 					# Name of tor fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 # account_fields = 				# Name of account fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 # subject_fields = 				# Name of fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 # destination_fields = 				# Name of destination fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 # answer_time_fields = 				# Name of time_answer fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 # duration_fields = 				# Name of duration fields to be used during extra mediation. Use index numbers in case of .csv cdrs.
 
 [session_manager]
 # enabled = false				# Starts SessionManager service: <true|false>.
 # switch_type = freeswitch			# Defines the type of switch behind: <freeswitch>.
 # rater = internal				# Address where to reach the Rater.
 # rater_reconnects = 3				# Number of reconnects to rater before giving up.
 # debit_interval = 10				# Interval to perform debits on.
 # max_call_duration = 3h			# Maximum call duration a prepaid call can last

 [freeswitch]
 # server = 127.0.0.1:8021			# Adress where to connect to FreeSWITCH socket.
 # passwd = ClueCon				# FreeSWITCH socket password.
 # reconnects = 5				# Number of attempts on connect failure.

 [history_server]
 # enabled = false				# Starts History service: <true|false>.
 # history_dir = /var/log/cgrates/history	# Location on disk where to store history files.
 # save_interval = 1s				# Interval to save changed cache into .git archive

 [history_agent]
 # enabled = false				# Starts History as a client: <true|false>.
 # server = internal				# Address where to reach the master history server: <internal|x.y.z.y:1234>

 [mailer]
 # server = localhost					# The server to use when sending emails out
 # auth_user = cgrates					# Authenticate to email server using this user
 # auth_passwd = CGRateS.org				# Authenticate to email server with this password
 # from_address = cgr-mailer@localhost.localdomain	# From address used when sending emails out
