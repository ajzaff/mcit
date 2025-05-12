module github.com/ajzaff/mcit/internal/tools/tune_fastlog_constants

go 1.24.2

require (
	github.com/ajzaff/fastlog v0.0.0-20250504190535-7ae1b28d450a
	github.com/ajzaff/fastlog/suite v0.0.0-20250504190535-7ae1b28d450a
	github.com/ajzaff/mcit v0.0.0-20250512211224-504e08ac448b
)

require github.com/ajzaff/lazyq v0.2.0 // indirect

replace github.com/ajzaff/mcit => ../../..
