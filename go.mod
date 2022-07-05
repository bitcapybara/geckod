module github.com/bitcapybara/geckod

go 1.18

replace github.com/bitcapybara/geckod-proto => ../geckod-proto

require (
	github.com/bitcapybara/geckod-proto v0.0.0-20220704151304-47e7dad1832b
	go.uber.org/atomic v1.9.0
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
