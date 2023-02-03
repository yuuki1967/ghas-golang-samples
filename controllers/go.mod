module example.com/controllers

go 1.18

replace example.com/logger => ../../util/logger

require example.com/logger v0.0.0-00010101000000-000000000000

require (
	github.com/magefile/mage v1.9.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	go.elastic.co/ecslogrus v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
