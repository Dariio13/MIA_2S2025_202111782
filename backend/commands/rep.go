package commands

import (
    reports "backend/reports"
    stores "backend/stores"
    "errors"
    "fmt"
    "regexp"
    "strings"
)

// REP estructura que representa el comando rep con sus parámetros
type REP struct {
    id           string // ID del disco
    path         string // Ruta donde guardar el reporte
    name         string // Tipo de reporte
    path_file_ls string // Ruta opcional para reportes tipo ls/file
}

// ParseRep parsea el comando rep y devuelve la respuesta o un error
func ParseRep(tokens []string) (string, error) {
    cmd := &REP{}

    args := strings.Join(tokens, " ")
    re := regexp.MustCompile(`-id=[^\s]+|-path="[^"]+"|-path=[^\s]+|-name=[^\s]+|-path_file_ls="[^"]+"|-path_file_ls=[^\s]+`)
    matches := re.FindAllString(args, -1)

    // Validar tokens inválidos
    if len(matches) != len(tokens) {
        for _, token := range tokens {
            if !re.MatchString(token) {
                return "", fmt.Errorf("parámetro inválido: %s", token)
            }
        }
    }

    // Procesar parámetros reconocidos
    for _, match := range matches {
        kv := strings.SplitN(match, "=", 2)
        if len(kv) != 2 {
            return "", fmt.Errorf("formato de parámetro inválido: %s", match)
        }
        key, value := strings.ToLower(kv[0]), kv[1]

        if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
            value = strings.Trim(value, "\"")
        }

        switch key {
        case "-id":
            if value == "" {
                return "", errors.New("el id no puede estar vacío")
            }
            cmd.id = value
        case "-path":
            if value == "" {
                return "", errors.New("el path no puede estar vacío")
            }
            cmd.path = value
        case "-name":
            validNames := []string{"mbr", "tree", "disk", "inode", "block", "bm_inode", "bm_block", "sb", "file", "ls"}
            if !contains(validNames, value) {
                return "", errors.New("nombre inválido, debe ser uno de los siguientes: mbr, tree, disk, inode, block, bm_inode, bm_block, sb, file, ls")
            }
            cmd.name = value
        case "-path_file_ls":
            cmd.path_file_ls = value
        default:
            return "", fmt.Errorf("parámetro desconocido: %s", key)
        }
    }

    // Verificar parámetros obligatorios
    if cmd.id == "" || cmd.path == "" || cmd.name == "" {
        return "", errors.New("faltan parámetros requeridos: -id, -path, -name")
    }

    // Ejecutar el comando
    err := commandRep(cmd)
    if err != nil {
        return "", err
    }

    return fmt.Sprintf("REP: Reporte generado exitosamente\n"+
        "-> ID: %s\n"+
        "-> Path: %s\n"+
        "-> Tipo: %s%s",
        cmd.id,
        cmd.path,
        cmd.name,
        func() string {
            if cmd.path_file_ls != "" {
                return fmt.Sprintf("\n-> Path LS: %s", cmd.path_file_ls)
            }
            return ""
        }()), nil
}

// Función auxiliar para verificar si un valor está en una lista
func contains(list []string, value string) bool {
    for _, v := range list {
        if v == value {
            return true
        }
    }
    return false
}

// commandRep genera los reportes solicitados
func commandRep(rep *REP) error {
    // Validar sesión activa
    if !stores.Auth.IsAuthenticated() {
        return fmt.Errorf("Error: no se ha iniciado sesión en ninguna partición")
    }

    // Verificar que el ID solicitado corresponda al de la sesión activa
    _, _, currentPartitionID := stores.Auth.GetCurrentUser()
    if currentPartitionID != rep.id {
        return fmt.Errorf("Error: la sesión actual no corresponde al ID %s", rep.id)
    }

    // Obtener la partición montada para generar reportes
    mountedMbr, mountedSb, mountedDiskPath, err := stores.GetMountedPartitionRep(rep.id)
    if err != nil {
        return fmt.Errorf("error al obtener la partición montada: %w", err)
    }

    // Manejo de reportes según el tipo solicitado
    switch rep.name {
    case "mbr":
        return reports.ReportMBR(mountedMbr, rep.path)
    case "disk":
        return reports.ReportDisk(mountedSb, mountedDiskPath, rep.path)
    case "inode":
        return reports.ReportInode(mountedSb, mountedDiskPath, rep.path)
    case "block":
        return reports.ReportBlock(mountedSb, mountedDiskPath, rep.path)
    case "bm_inode":
        return reports.ReportBMInode(mountedSb, mountedDiskPath, rep.path)
    case "bm_block":
        return reports.ReportBMBlock(mountedSb, mountedDiskPath, rep.path)
    case "sb":
        return reports.ReportSB(mountedSb, mountedDiskPath, rep.path)
    case "file":
        if rep.path_file_ls == "" {
            return fmt.Errorf("error: para generar un reporte tipo 'file' debes especificar -path_file_ls")
        }
        return reports.ReportFile(mountedSb, mountedDiskPath, rep.path_file_ls, rep.path)
    case "ls":
        if rep.path_file_ls == "" {
            return fmt.Errorf("error: para generar un reporte tipo 'ls' debes especificar -path_file_ls")
        }
        return reports.ReportLS(mountedSb, mountedDiskPath, rep.path_file_ls, rep.path)
    case "tree":
        // 🔹 Si aún no implementas el árbol, devolvemos una imagen vacía
        return reports.GenerateDummyTree(mountedSb, mountedDiskPath, rep.path)
    default:
        return fmt.Errorf("reporte '%s' no implementado", rep.name)
    }
}