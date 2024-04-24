module github.com/producer

go 1.22.0

require (
	github.com/google/uuid v1.6.0
	github.com/models v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.47
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

replace github.com/models => ../models
