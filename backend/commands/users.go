package commands

import (
    "backend/stores"
    "backend/structures"
    "errors"
    "fmt"
    "regexp"
    "strings"
)

type MKGRP struct {
    name string
}

type RMGRP struct {
    name string
}

type MKUSR struct {
    user string
    pass string
    grp  string
}

type RMUSR struct {
    user string
}

type CHGRP struct {
    user string
    grp  string
}

// ------------------ FUNCIONES AUXILIARES ------------------

// Leer contenido de users.txt
func readUsersFile(sb *structures.SuperBlock, partitionPath string) (string, error) {
    data, err := sb.ReadFile(partitionPath, "/users.txt")
    if err != nil {
        return "", fmt.Errorf("error leyendo users.txt: %w", err)
    }
    return string(data), nil
}

// Guardar contenido en users.txt
func writeUsersFile(sb *structures.SuperBlock, partitionPath string, content string) error {
    return sb.UpdateFile(partitionPath, "/users.txt", []byte(content))
}

// Validar si usuario actual es root
func validateRoot() error {
    currentUser, _, _ := stores.Auth.GetCurrentUser()
    if currentUser != "root" {
        return errors.New("solo el usuario root puede ejecutar este comando")
    }
    return nil
}

// ------------------ MKGRP ------------------
func ParseMkgrp(tokens []string) (string, error) {
    cmd := &MKGRP{}
    re := regexp.MustCompile(`-name=[^\s]+`)
    matches := re.FindAllString(strings.Join(tokens, " "), -1)

    if len(matches) != 1 {
        return "", errors.New("uso: mkgrp -name=nombre_grupo")
    }

    cmd.name = strings.SplitN(matches[0], "=", 2)[1]
    return commandMkgrp(cmd)
}

func commandMkgrp(mkgrp *MKGRP) (string, error) {
    if !stores.Auth.IsAuthenticated() {
        return "", errors.New("debes iniciar sesión para usar mkgrp")
    }
    if err := validateRoot(); err != nil {
        return "", err
    }

    sb, _, path, err := stores.GetMountedPartitionSuperblock(stores.Auth.GetPartitionID())
    if err != nil {
        return "", err
    }

    content, _ := readUsersFile(sb, path)
    lines := strings.Split(content, "\n")

    // Verificar que no exista el grupo
    for _, line := range lines {
        fields := strings.Split(line, ",")
        if len(fields) >= 3 && fields[1] == "G" && fields[2] == mkgrp.name && fields[0] != "0" {
            return "", fmt.Errorf("el grupo '%s' ya existe", mkgrp.name)
        }
    }

    // Crear nuevo grupo
    newID := len(lines) + 1
    newLine := fmt.Sprintf("%d,G,%s", newID, mkgrp.name)
    if content != "" {
        content += "\n" + newLine
    } else {
        content = newLine
    }

    err = writeUsersFile(sb, path, content)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("MKGRP: Grupo '%s' creado correctamente", mkgrp.name), nil
}

// ------------------ RMGRP ------------------
func ParseRmgrp(tokens []string) (string, error) {
    cmd := &RMGRP{}
    re := regexp.MustCompile(`-name=[^\s]+`)
    matches := re.FindAllString(strings.Join(tokens, " "), -1)

    if len(matches) != 1 {
        return "", errors.New("uso: rmgrp -name=nombre_grupo")
    }

    cmd.name = strings.SplitN(matches[0], "=", 2)[1]
    return commandRmgrp(cmd)
}

func commandRmgrp(rmgrp *RMGRP) (string, error) {
    if !stores.Auth.IsAuthenticated() {
        return "", errors.New("debes iniciar sesión para usar rmgrp")
    }
    if err := validateRoot(); err != nil {
        return "", err
    }

    sb, _, path, err := stores.GetMountedPartitionSuperblock(stores.Auth.GetPartitionID())
    if err != nil {
        return "", err
    }

    content, _ := readUsersFile(sb, path)
    lines := strings.Split(content, "\n")
    updated := false

    for i, line := range lines {
        fields := strings.Split(line, ",")
        if len(fields) >= 3 && fields[1] == "G" && fields[2] == rmgrp.name && fields[0] != "0" {
            lines[i] = "0," + fields[1] + "," + fields[2] // Marcamos como eliminado
            updated = true
        }
    }

    if !updated {
        return "", fmt.Errorf("el grupo '%s' no existe", rmgrp.name)
    }

    err = writeUsersFile(sb, path, strings.Join(lines, "\n"))
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("RMGRP: Grupo '%s' eliminado correctamente", rmgrp.name), nil
}

// ------------------ MKUSR ------------------
func ParseMkusr(tokens []string) (string, error) {
    cmd := &MKUSR{}
    re := regexp.MustCompile(`-user=[^\s]+|-pass=[^\s]+|-grp=[^\s]+`)
    matches := re.FindAllString(strings.Join(tokens, " "), -1)

    if len(matches) != 3 {
        return "", errors.New("uso: mkusr -user=nombre -pass=clave -grp=grupo")
    }

    for _, match := range matches {
        kv := strings.SplitN(match, "=", 2)
        switch kv[0] {
        case "-user":
            cmd.user = kv[1]
        case "-pass":
            cmd.pass = kv[1]
        case "-grp":
            cmd.grp = kv[1]
        }
    }
    return commandMkusr(cmd)
}

func commandMkusr(mkusr *MKUSR) (string, error) {
    if !stores.Auth.IsAuthenticated() {
        return "", errors.New("debes iniciar sesión para usar mkusr")
    }
    if err := validateRoot(); err != nil {
        return "", err
    }

    sb, _, path, err := stores.GetMountedPartitionSuperblock(stores.Auth.GetPartitionID())
    if err != nil {
        return "", err
    }

    content, _ := readUsersFile(sb, path)
    lines := strings.Split(content, "\n")

    // Verificar si usuario ya existe
    for _, line := range lines {
        fields := strings.Split(line, ",")
        if len(fields) >= 5 && fields[1] == "U" && fields[3] == mkusr.user && fields[0] != "0" {
            return "", fmt.Errorf("el usuario '%s' ya existe", mkusr.user)
        }
    }

    // Crear usuario
    newID := len(lines) + 1
    newLine := fmt.Sprintf("%d,U,%s,%s,%s", newID, mkusr.grp, mkusr.user, mkusr.pass)
    if content != "" {
        content += "\n" + newLine
    } else {
        content = newLine
    }

    err = writeUsersFile(sb, path, content)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("MKUSR: Usuario '%s' creado correctamente", mkusr.user), nil
}

// ------------------ RMUSR ------------------
func ParseRmusr(tokens []string) (string, error) {
    cmd := &RMUSR{}
    re := regexp.MustCompile(`-user=[^\s]+`)
    matches := re.FindAllString(strings.Join(tokens, " "), -1)

    if len(matches) != 1 {
        return "", errors.New("uso: rmusr -user=nombre_usuario")
    }

    cmd.user = strings.SplitN(matches[0], "=", 2)[1]
    return commandRmusr(cmd)
}

func commandRmusr(rmusr *RMUSR) (string, error) {
    if !stores.Auth.IsAuthenticated() {
        return "", errors.New("debes iniciar sesión para usar rmusr")
    }
    if err := validateRoot(); err != nil {
        return "", err
    }

    sb, _, path, err := stores.GetMountedPartitionSuperblock(stores.Auth.GetPartitionID())
    if err != nil {
        return "", err
    }

    content, _ := readUsersFile(sb, path)
    lines := strings.Split(content, "\n")
    updated := false

    for i, line := range lines {
        fields := strings.Split(line, ",")
        if len(fields) >= 5 && fields[1] == "U" && fields[3] == rmusr.user && fields[0] != "0" {
            lines[i] = "0," + fields[1] + "," + fields[2] + "," + fields[3] + "," + fields[4]
            updated = true
        }
    }

    if !updated {
        return "", fmt.Errorf("el usuario '%s' no existe", rmusr.user)
    }

    err = writeUsersFile(sb, path, strings.Join(lines, "\n"))
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("RMUSR: Usuario '%s' eliminado correctamente", rmusr.user), nil
}

// ------------------ CHGRP ------------------
func ParseChgrp(tokens []string) (string, error) {
    cmd := &CHGRP{}
    re := regexp.MustCompile(`-user=[^\s]+|-grp=[^\s]+`)
    matches := re.FindAllString(strings.Join(tokens, " "), -1)

    if len(matches) != 2 {
        return "", errors.New("uso: chgrp -user=nombre_usuario -grp=grupo_nuevo")
    }

    for _, match := range matches {
        kv := strings.SplitN(match, "=", 2)
        switch kv[0] {
        case "-user":
            cmd.user = kv[1]
        case "-grp":
            cmd.grp = kv[1]
        }
    }
    return commandChgrp(cmd)
}

func commandChgrp(chgrp *CHGRP) (string, error) {
    if !stores.Auth.IsAuthenticated() {
        return "", errors.New("debes iniciar sesión para usar chgrp")
    }
    if err := validateRoot(); err != nil {
        return "", err
    }

    sb, _, path, err := stores.GetMountedPartitionSuperblock(stores.Auth.GetPartitionID())
    if err != nil {
        return "", err
    }

    content, _ := readUsersFile(sb, path)
    lines := strings.Split(content, "\n")
    updated := false

    for i, line := range lines {
        fields := strings.Split(line, ",")
        if len(fields) >= 5 && fields[1] == "U" && fields[3] == chgrp.user && fields[0] != "0" {
            lines[i] = fmt.Sprintf("%s,U,%s,%s,%s", fields[0], chgrp.grp, fields[3], fields[4])
            updated = true
        }
    }

    if !updated {
        return "", fmt.Errorf("el usuario '%s' no existe", chgrp.user)
    }

    err = writeUsersFile(sb, path, strings.Join(lines, "\n"))
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("CHGRP: Usuario '%s' movido al grupo '%s'", chgrp.user, chgrp.grp), nil
}

/* ------------------ MOUNTED ------------------
func Mounted() string {
    if len(stores.MountedPartitions) == 0 {
        return "No hay particiones montadas"
    }

    result := "Particiones montadas:\n"
    for id, path := range stores.MountedPartitions {
        result += fmt.Sprintf("ID: %s | Path: %s\n", id, path)
    }
    return result
}*/