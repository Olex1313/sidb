package sidb

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

const testDbFilename = "test-db"

func executeCommands(commands []string) ([]string, error) {
	sidbCmd := exec.Command("./sidb", testDbFilename)
	sidbStdin, err := sidbCmd.StdinPipe()

	stdOut := new(bytes.Buffer)
	sidbCmd.Stdout = io.MultiWriter(os.Stdout, stdOut)

	sidbCmd.Start()
	for _, command := range commands {
		sidbStdin.Write([]byte(command + "\n"))
	}
	sidbCmd.Wait()
	return strings.Split(stdOut.String(), "\n"), err
}

func TestSiDb(t *testing.T) {
	tests := []struct {
		name           string
		commands       []string
		expectedOutput []string
	}{
		{
			name:     "insert should add row",
			commands: []string{"insert 1 b c", "select", ".exit"},
			expectedOutput: []string{
				"db > Executed",
				"db > 1, b, c",
				"Executed",
				"db > ",
			},
		},
		{
			name:     "insert should add row",
			commands: []string{"insert 1 b c", "select", ".exit"},
			expectedOutput: []string{
				"db > Executed",
				"db > 1, b, c",
				"Executed",
				"db > ",
			},
		},
		{
			name:     "insert should fail with negative id",
			commands: []string{"insert -1 b c", ".exit"},
			expectedOutput: []string{
				"db > ID must be positive",
				"db > ",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := executeCommands(test.commands)
			if err != nil {
				t.Error("Unexpected error:", err)
			}
			if !reflect.DeepEqual(output, test.expectedOutput) {
				t.Errorf("Expected\n%v\n to be equal to\n%v\n", output, test.expectedOutput)
			}
		})
		os.Remove(testDbFilename)
	}
}

func TestMaxStringLengthShouldBeEqualToMaxSize(t *testing.T) {
	longUsername := strings.Repeat("a", 32)
	longEmail := strings.Repeat("b", 255)
	commands := []string{
		fmt.Sprintf("insert 1 %s %s", longUsername, longEmail),
		"select",
		".exit",
	}
	output, err := executeCommands(commands)
	if err != nil {
		panic(err)
	}
	expectedOutput := []string{
		"db > Executed",
		fmt.Sprintf("db > 1, %s, %s", longUsername, longEmail),
		"Executed",
		"db > ",
	}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Test failed\n expected: %v\nbe equal to: %v", output, expectedOutput)
	}
	os.Remove(testDbFilename)
}

func TestErrorIsReturnedWhenStringIsTooLong(t *testing.T) {
	forbiddenUsername := strings.Repeat("a", 33)
	forbiddenEmail := strings.Repeat("b", 256)
	tests := []struct {
		name           string
		commands       []string
		expectedOutput []string
	}{
		{
			name:     "insert should forbid too long username",
			commands: []string{fmt.Sprintf("insert 1 %s c", forbiddenUsername), "select", ".exit"},
			expectedOutput: []string{
				"db > String too long",
				"db > Executed",
				"db > ",
			},
		},
		{
			name:     "insert should forbid too long email",
			commands: []string{fmt.Sprintf("insert 1 b %s", forbiddenEmail), "select", ".exit"},
			expectedOutput: []string{
				"db > String too long",
				"db > Executed",
				"db > ",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := executeCommands(test.commands)
			if err != nil {
				t.Error("Unexpected error:", err)
			}
			if !reflect.DeepEqual(output, test.expectedOutput) {
				t.Errorf("Test failed, expected %v be equal to %v", output, test.expectedOutput)
			}
		})

	}
	os.Remove(testDbFilename)
}

func TestInMemoryPersistence(t *testing.T) {
	firstCommands := []string{"insert 1 user1 user1@example.com", ".exit"}
	firstExpectedOutput := []string{"db > Executed", "db > "}
	secondCommands := []string{"select", ".exit"}
	secondExpectedOutput := []string{"db > 1, user1, user1@example.com", "Executed", "db > "}
	output, err := executeCommands(firstCommands)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(output, firstExpectedOutput) {
		t.Errorf("Expected\n%v\n to be equal to\n%v\n", output, firstExpectedOutput)
	}
	output, err = executeCommands(secondCommands)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if !reflect.DeepEqual(output, secondExpectedOutput) {
		t.Errorf("Expected\n%v\n to be equal to\n%v\n", output, secondExpectedOutput)
	}
	os.Remove(testDbFilename)
}
