package checkerror

func CheckError(values ...error) {
	for _, err := range values {
		if err != nil {
			panic(err)
		}
	}
}
