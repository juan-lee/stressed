// Package stressng provides stress configurations.
package stressng

const foobar string = `
verbose
metrics-brief
timeout 0
#vm 4
#vm-bytes 18G
#vm-keep
#vm-populate
aio 2
aiol 2
hdd 2
`

type Jobfile struct {
}

func (j *Jobfile) String() string {
	return foobar
}
