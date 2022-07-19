module github.com/bitcapybara/geckod

go 1.18

// replace github.com/bitcapybara/geckod-proto => ../geckod-proto

require (
	github.com/bitcapybara/geckod-proto v0.0.0-20220714105617-6f7fdcf15669
	github.com/bits-and-blooms/bitset v1.2.2
	go.uber.org/atomic v1.9.0
)

require google.golang.org/protobuf v1.28.0 // indirect
