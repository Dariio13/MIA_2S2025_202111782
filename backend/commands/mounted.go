package commands

import (
	"fmt"
	"backend/stores"
)

func Mounted() string {
	if len(stores.MountedPartitions) == 0 {
		return "No hay particiones montadas"
	}

	result := "=== PARTICIONES MONTADAS ===\n"
	for id, path := range stores.MountedPartitions {
		result += fmt.Sprintf("ID: %s | Disco: %s\n", id, path)
	}
	return result
}