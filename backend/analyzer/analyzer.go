package analyzer

import (
	commands "backend/commands"
	"errors"
	"fmt"
	"strings"
)

func Analyzer(input string) (string, error) {

	// Ignorar líneas de comentarios o vacías desde el principio
	if strings.TrimSpace(input) == "" || strings.HasPrefix(strings.TrimSpace(input), "#") {
		return "", nil
	}

	tokens := strings.Fields(input)

	if len(tokens) == 0 {
		return "", errors.New("no se proporcionó ningún comando")
	}

	switch tokens[0] {
	case "mkdisk":
		return commands.ParseMkdisk(tokens[1:])
	case "rmdisk":
		return commands.ParserRmdisk(tokens[1:])
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
	case "login":
		return commands.ParseLogin(tokens[1:])
	case "logout":
		return commands.ParseLogout(tokens[1:])
	// case "#":
	// 	return "", nil
	default:
		return "", fmt.Errorf("comando desconocido: %s", tokens[0])
	}
}
