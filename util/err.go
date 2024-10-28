package util

import "log"

func FatalErr(str string, err error) {
	if err != nil {
		log.Printf("%s: %v\n", str, err)
		log.Fatalln("=====DNS-Shift FATAL=====")
	}
}

func PrintErr(str string, err error) bool {
	if err != nil {
		log.Printf("%s: %v\n", str, err)
		return true
	}
	return false
}
