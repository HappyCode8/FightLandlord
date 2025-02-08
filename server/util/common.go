package util

func ChooseIf[T any](cond bool, onTrue, onFalse T) T {
	if cond {
		return onTrue
	} else {
		return onFalse
	}
}
