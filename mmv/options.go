package mmv

import "regexp"

func OptionDryRun() func(*MultiMover) {
	return func(mover *MultiMover) {
		mover.dryRun = true
	}
}

func OptionInclude(includes *regexp.Regexp) func(*MultiMover) {
	return func(mover *MultiMover) {
		mover.includes = includes
	}
}

func OptionExclude(excludes *regexp.Regexp) func(*MultiMover) {
	return func(mover *MultiMover) {
		mover.excludes = excludes
	}
}
