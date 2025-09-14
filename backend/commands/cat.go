package commands

import (
    "fmt"
    "os"
    stores "backend/stores"
)

// Cat lee el contenido de un archivo dentro de la partición logueada
func Cat(path string, id string) (string, error) {
    // 🔹 Validar que haya sesión activa
    if !stores.Auth.IsAuthenticated() {
        return "", fmt.Errorf("Error: no se ha iniciado sesión en ninguna partición")
    }

    // 🔹 Obtener la partición activa desde la sesión
    _, _, currentPartitionID := stores.Auth.GetCurrentUser()

    // 🔹 Verificar que el ID recibido coincida con la sesión activa
    if currentPartitionID != id {
        return "", fmt.Errorf("Error: la sesión actual no corresponde al ID %s", id)
    }

    // 🔹 Intentar leer el archivo
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("Error: no se puede leer el archivo %s", path)
    }

    // 🔹 Retornar el contenido
    return string(data), nil
}