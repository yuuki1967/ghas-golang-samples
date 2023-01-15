module go_samples

go 1.18

require (
	github.com/sirupsen/logrus v1.9.0
	logmod.com/logmod v0.0.0-00010101000000-000000000000
)

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

replace logmod.com/logmod => ./logmod
