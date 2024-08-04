module mainmodule

go 1.19

replace mainmodule/helper => ./helper

require (
	mainmodule/helper v0.0.0-00010101000000-000000000000
	mainmodule/service v0.0.0-00010101000000-000000000000
)

replace mainmodule/service => ./service
