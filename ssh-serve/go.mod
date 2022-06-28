module ssh-serve

go 1.18

require (
	MINTCP v0.0.0
	golang.org/x/text v0.3.7
	minlib v0.0.0
)

replace minlib v0.0.0 => ../../../minlib

replace MINTCP v0.0.0 => ../../../MINTCP
