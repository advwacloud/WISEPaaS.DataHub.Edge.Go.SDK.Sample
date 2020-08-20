module github.com/advwacloud/WISEPaaS.DataHub.Edge.Go.SDK.Sample

go 1.13

replace github.com/advwacloud/WISEPaaS.DataHub.Edge.Go.SDK => ../WISEPaaS.DataHub.Edge.Go.SDK // To use the local SDK

require (
	github.com/advwacloud/WISEPaaS.DataHub.Edge.Go.SDK v0.0.0-20200807070017-dc6a6ab5cd9b
	github.com/mattn/go-sqlite3 v1.13.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
)
