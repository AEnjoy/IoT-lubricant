package main

/*
This file is used at E2E test.
E2E Test working-directory is at test/
*/
func main() {
	if err := Execute(); err != nil {
		println(err)
	}
}
