package discordcommands

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func TimeToTimestamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

func parseArgs(regex map[string]string, command *Command, args []string) (newArgs map[string]string, err error) {
	newArgs = make(map[string]string)
	if command == nil {
		return newArgs, nil
	}
	for index, argument := range command.Arguments {
		fmt.Println(index, argument)
		value, err := newArgValue(argument, args, index)
		if err != nil {
			return nil, err
		}
		newArgs[argument.Name] = value
		pattern, ok := regex[argument.Type]
		if value != "" && ok {
			match, _ := regexp.MatchString(pattern, value)
			if !match {
				return nil, UsageError{}
			}
		}
	}
	return
}

func newArgValue(argument *Argument, args []string, index int) (string, error) {
	if index >= len(args) && !argument.Optional {
		return "", UsageError{}
	} else if index >= len(args) && argument.Optional {
		return "", nil
	} else if argument.Infinite {
		return strings.Join(args[index:], " "), nil
	} else {
		return args[index], nil
	}
}
