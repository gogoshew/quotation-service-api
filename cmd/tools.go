package main

func panicIfErr(e error) {
	if e != nil {
		panic(e)
	}
}
