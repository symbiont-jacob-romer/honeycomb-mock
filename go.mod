module github.com/symbiont-io/assembly

go 1.12

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0

require (
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/facebookgo/limitgroup v0.0.0-20150612190941-6abd8d71ec01 // indirect
	github.com/facebookgo/muster v0.0.0-20150708232844-fd3d7953fd52 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/honeycombio/libhoney-go v1.12.1
	github.com/klauspost/compress v1.8.5
	github.com/stretchr/testify v1.4.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	google.golang.org/appengine v1.6.4 // indirect
	gopkg.in/alexcesaro/statsd.v2 v2.0.0 // indirect
)
