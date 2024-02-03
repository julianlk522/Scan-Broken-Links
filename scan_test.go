package main

import (
	"os"
	"testing"
)

func TestCollectUserInput(t *testing.T) {

	// root_url
	url_tests := []struct {
		name string
		input string
		output string
	} {
		// Default root_url is "https://www.google.com"
		// so that should be expected output if no input
		// provided
		{"Test1", "", "https://www.google.com"},
		{"Test2", "https://www.netflix.com", "https://www.netflix.com"},
		{"Test3", "TestString", "TestString"},
	}

	for _, tt := range url_tests {
		t.Run(tt.name, func(t *testing.T) {
			// pipe pre-generated test cases to stdin
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdin = r
			// write test input to stdin
			w.Write([]byte(tt.input))
			// use default value for numbers of
			// consumer threads
			w.Write([]byte("\n"))
			w.Close()

			if got_url, _ := collect_user_input(); got_url != tt.output {
				t.Errorf("collect_user_input() = %v, want %v", got_url, tt.output)
			}
		})
	}

	// consumer threads
	ct_tests := []struct {
		name string
		input string
		output int
	} {
		// Default ct is 2 so that should be expected
		// output if no input provided
		{"Test4", "", 2},
		{"Test5", "2", 2},
		{"Test6", "4", 4},
	}

	for _, tt := range ct_tests {
		t.Run(tt.name, func(t *testing.T) {
			// pipe pre-generated test cases to stdin
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdin = r
			// use default value for root_url
			w.Write([]byte("\n"))

			// write test input to stdin
			w.Write([]byte(tt.input))
			w.Close()

			if _, got_ct := collect_user_input(); got_ct != tt.output {
				t.Errorf("collect_user_input() = %v, want %v", got_ct, tt.output)
			}
		})
	}
}