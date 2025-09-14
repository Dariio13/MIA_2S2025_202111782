package commands

import (
    stores "backend/stores"
    structures "backend/structures"
    utils "backend/utils"
    "errors"
    "fmt"
    "regexp"
    "strings"
)

// MOUNT estructura que representa el comando mount con sus parámetros
type MOUNT struct {
    path string // Ruta del archivo del disco
    name string // Nombre de la partición
}

/*
    Ejemplos:
    mount -path=/home/vboxuser/MIA_2S2025_202111782/disks/Disco1.mia -name=Part1
    mount -path=/home/vboxuser/MIA_2S2025_202111782/disks/Disco2.mia -name=Part2
*/

// ParseMount parsea el comando mount y devuelve una instancia de MOUNT
func ParseMount(tokens []string) (string, error) {
    cmd := &MOUNT{}

    // Unir tokens en una sola cadena
    args := strings.Join(tokens, " ")
    // Expresión regular para obtener parámetros válidos
    re := regexp.MustCompile(`-path="[^"]+"|-path=[^\s]+|-name="[^"]+"|-name=[^\s]+`)
    matches := re.FindAllString(args, -1)

    // Validar tokens inválidos
    if len(matches) != len(tokens) {
        for _, token := range tokens {
            if !re.MatchString(token) {
                return "", fmt.Errorf("parámetro inválido: %s", token)
            }
        }
    }

    // Extraer parámetros
    for _, match := range matches {
        kv := strings.SplitN(match, "=", 2)
        if len(kv) != 2 {
            return "", fmt.Errorf("formato de parámetro inválido: %s", match)
        }
        key, value := strings.ToLower(kv[0]), kv[1]

        // Limpiar comillas si existen
        if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
            value = strings.Trim(value, "\"")
        }

        switch key {
        case "-path":
            if value == "" {
                return "", errors.New("el path no puede estar vacío")
            }
            cmd.path = value
        case "-name":
            if value == "" {
                return "", errors.New("el nombre no puede estar vacío")
            }
            cmd.name = value
        default:
            return "", fmt.Errorf("parámetro desconocido: %s", key)
        }
    }

    // Validar parámetros requeridos
    if cmd.path == "" {
        return "", errors.New("faltan parámetros requeridos: -path")
    }
    if cmd.name == "" {
        return "", errors.New("faltan parámetros requeridos: -name")
    }

    // Montar la partición
    idPartition, err := commandMount(cmd)
    if err != nil {
        return "", err
    }

    // Respuesta exitosa
    return fmt.Sprintf("MOUNT: Partición montada exitosamente\n"+
        "-> Path: %s\n"+
        "-> Nombre: %s\n"+
        "-> ID: %s",
        cmd.path, cmd.name, idPartition), nil
}

func commandMount(mount *MOUNT) (string, error) {
    var mbr structures.MBR

    // Deserializar MBR desde el archivo de disco
    err := mbr.Deserialize(mount.path)
    if err != nil {
        return "", fmt.Errorf("error deserializando el MBR: %w", err)
    }

    // Buscar la partición
    partition, indexPartition := mbr.GetPartitionByName(mount.name)
    if partition == nil {
        return "", fmt.Errorf("error: la partición '%s' no existe en el disco", mount.name)
    }

    // --- NUEVA VALIDACIÓN: VERIFICAR SI YA ESTÁ MONTADA ---
    for id, path := range stores.MountedPartitions {
        partitionName := strings.TrimSpace(string(partition.Part_name[:]))
        if path == mount.path && strings.EqualFold(partitionName, mount.name) {
            return "", fmt.Errorf("error: la partición '%s' ya está montada con ID %s", mount.name, id)
        }
    }
    // ------------------------------------------------------

    // Generar ID único para la partición
    idPartition, partitionCorrelative, err := generatePartitionID(mount)
    if err != nil {
        return "", fmt.Errorf("error generando el ID de la partición: %w", err)
    }

    // Guardar la partición en la lista de montajes globales
    stores.MountedPartitions[idPartition] = mount.path

    // Modificar la partición para indicar que está montada
    partition.MountPartition(partitionCorrelative, idPartition)

    fmt.Println("\nPartición montada (modificada):")
    partition.PrintPartition()

    // Guardar cambios en el MBR
    mbr.Mbr_partitions[indexPartition] = *partition
    err = mbr.Serialize(mount.path)
    if err != nil {
        return "", fmt.Errorf("error serializando el MBR: %w", err)
    }

    return idPartition, nil
}

func generatePartitionID(mount *MOUNT) (string, int, error) {
    // Obtener la letra y correlativo para el ID
    letter, partitionCorrelative, err := utils.GetLetterAndPartitionCorrelative(mount.path)
    if err != nil {
        return "", 0, fmt.Errorf("error obteniendo la letra y el correlativo: %w", err)
    }

    // Crear ID con formato -> <CARNET><NUM><LETRA>
    idPartition := fmt.Sprintf("%s%d%s", stores.Carnet, partitionCorrelative, letter)

    return idPartition, partitionCorrelative, nil
}