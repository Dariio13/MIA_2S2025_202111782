package analyzer

import (
	commands "backend/commands"
	"errors"
	"fmt"
	"strings"
)

func Analyzer(input string) (string, error) {
	tokens := strings.Fields(input)

	if len(tokens) == 0 {
		return "", errors.New("no se proporcionó ningún comando")
	}

	switch tokens[0] {
	case "mkdisk":
		return commands.ParseMkdisk(tokens[1:])
	case "fdisk":
		return commands.ParseFdisk(tokens[1:])
	case "mount":
		return commands.ParseMount(tokens[1:])
	case "mkfs":
		return commands.ParseMkfs(tokens[1:])
	case "rep":
		return commands.ParseRep(tokens[1:])
	case "mkdir":
		return commands.ParseMkdir(tokens[1:])
	case "mkfile":
		return commands.ParserMkfile(tokens[1:])
	default:
		return "", fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}
