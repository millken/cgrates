#! /usr/bin/env sh

go test -i github.com/cgrates/cgrates/apier/v1
go test -i github.com/cgrates/cgrates/engine
go test -i github.com/cgrates/cgrates/sessionmanager
go test -i github.com/cgrates/cgrates/config
go test -i github.com/cgrates/cgrates/cmd/cgr-engine
go test -i github.com/cgrates/cgrates/mediator
go test -i github.com/cgrates/fsock
go test -i github.com/cgrates/cgrates/cache2go
go test -i github.com/cgrates/cgrates/cdrs
go test -i github.com/cgrates/cgrates/utils
go test -i github.com/cgrates/cgrates/history
go test -i github.com/cgrates/cgrates/cdrexporter

go test github.com/cgrates/cgrates/apier/v1 -local
ap=$?
go test github.com/cgrates/cgrates/engine -local
en=$?
go test github.com/cgrates/cgrates/sessionmanager
sm=$?
go test github.com/cgrates/cgrates/config
cfg=$?
go test github.com/cgrates/cgrates/cmd/cgr-engine
cr=$?
go test github.com/cgrates/cgrates/mediator
md=$?
go test github.com/cgrates/cgrates/cdrs
cdr=$?
go test github.com/cgrates/cgrates/utils
ut=$?
go test github.com/cgrates/fsock
fs=$?
go test github.com/cgrates/cgrates/history
hs=$?
go test github.com/cgrates/cgrates/cache2go
c2g=$?
go test github.com/cgrates/cgrates/cdrexporter
cdre=$?


exit $ap && $en && $sm && $cfg && $bl && $cr && $md && $cdr && $fs && $ut && $hs && $c2g && $cdre
