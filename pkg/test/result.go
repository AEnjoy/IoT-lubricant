package test

import (
	"fmt"
	"os"
)

type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r *Result) CheckResult(abort bool) bool {
	if !r.Success {
		_, _ = fmt.Fprintln(os.Stdout, "Test_Failed_Message:", r.Message)
		if abort {
			fmt.Println("According to the policy, subsequent tests have been terminated")
			os.Exit(12)
		}
	}
	return true
}
