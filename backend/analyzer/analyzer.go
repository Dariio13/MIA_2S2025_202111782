package analyzer

import (
    commands "backend/commands"
    "backend/stores"
    "bufio"
    "errors"
    "fmt"
    "strings"
)

// Analyzer recibe un comando completo, lo tokeniza y lo envía al parser correcto
func Analyzer(cmd string) (string, error) {
    var tokens []string

    // Usamos un scanner para dividir palabras respetando espacios dentro de comillas
    scanner := bufio.NewScanner(strings.NewReader(cmd))
    scanner.Split(bufio.ScanWords)

    for scanner.Scan() {
        token := scanner.Text()

        // Si detectamos que el token empieza con comillas, concatenamos hasta cerrarlas
        if strings.HasPrefix(token, "\"") {
            for !strings.HasSuffix(token, "\"") && scanner.Scan() {
                token += " " + scanner.Text()
            }
            // Conservamos las comillas para que los parsers manejen las rutas correctamente
        }

        tokens = append(tokens, token)
    }

    // Si no hay tokens, error
    if len(tokens) == 0 {
        return "", errors.New("no se proporcionó ningún comando")
    }

    // Comando principal
    cmdName := strings.ToLower(tokens[0])

    switch cmdName {
    // ------------- PARTE 1: MANEJO DE DISCOS -------------
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

    // ------------- PARTE 2: ARCHIVOS Y DIRECTORIOS -------------
    case "mkdir":
        return commands.ParseMkdir(tokens[1:])
    case "mkfile":
        return commands.ParserMkfile(tokens[1:])
    case "cat":
        var path string
        for _, token := range tokens[1:] {
            if strings.HasPrefix(token, "-path=") {
                path = strings.TrimPrefix(token, "-path=")
            }
        }
        if path == "" {
            return "", fmt.Errorf("Error: falta el parámetro -path")
        }
        return commands.Cat(path, stores.CurrentSession.PartitionID)

    // ------------- PARTE 3: SESIONES Y USUARIOS -------------
    /*case "login":
        return commands.ParseLogin(tokens[1:])
    case "logout":
        return commands.ParseLogout(tokens[1:])
    case "mkgrp":
        return commands.ParseMkgrp(tokens[1:])
    case "rmgrp":
        return commands.ParseRmgrp(tokens[1:])
    case "mkusr":
        return commands.ParseMkusr(tokens[1:])
    case "rmusr":
        return commands.ParseRmusr(tokens[1:])
    case "chgrp":
        return commands.ParseChgrp(tokens[1:])

	*/	

    // ------------- PARTE 4: LISTAR PARTICIONES MONTADAS -------------
    case "mounted":
        return commands.Mounted(), nil

    // ------------- COMENTARIOS -------------
    case "#":
        return "", nil // Ignoramos las líneas que empiezan con "#"

    // ------------- ERROR: COMANDO NO RECONOCIDO -------------

    default:
        return "", fmt.Errorf("comando desconocido: %s", tokens[0])
    }
}