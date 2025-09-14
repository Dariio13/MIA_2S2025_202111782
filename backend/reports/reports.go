package reports

import (
    "fmt"
    "os"
    structures "backend/structures"
)

// Reporte general para DISK
func ReportDisk(sb *structures.SuperBlock, diskPath string, outputPath string) error {
    fmt.Println("Generando reporte DISK en:", outputPath)
    return os.WriteFile(outputPath, []byte("REPORTE DISK"), 0644)
}

// Reporte de BLOQUES
func ReportBlock(sb *structures.SuperBlock, diskPath string, outputPath string) error {
    fmt.Println("Generando reporte BLOCK en:", outputPath)
    return os.WriteFile(outputPath, []byte("REPORTE BLOCK"), 0644)
}

// Reporte de BITMAP de BLOQUES
func ReportBMBlock(sb *structures.SuperBlock, diskPath string, outputPath string) error {
    fmt.Println("Generando reporte BM_BLOCK en:", outputPath)
    return os.WriteFile(outputPath, []byte("REPORTE BM_BLOCK"), 0644)
}

// Reporte del SUPERBLOQUE
func ReportSB(sb *structures.SuperBlock, diskPath string, outputPath string) error {
    fmt.Println("Generando reporte SUPERBLOCK en:", outputPath)
    return os.WriteFile(outputPath, []byte("REPORTE SUPERBLOCK"), 0644)
}

// Reporte de un archivo específico
func ReportFile(sb *structures.SuperBlock, diskPath string, filePath string, outputPath string) error {
    fmt.Println("Generando reporte FILE en:", outputPath)
    return os.WriteFile(outputPath, []byte("REPORTE FILE"), 0644)
}

// Reporte tipo LS (lista de archivos en un directorio)
func ReportLS(sb *structures.SuperBlock, diskPath string, dirPath string, outputPath string) error {
    fmt.Println("Generando reporte LS en:", outputPath)
    return os.WriteFile(outputPath, []byte("REPORTE LS"), 0644)
}

// Reporte de la estructura de árbol
func GenerateDummyTree(sb *structures.SuperBlock, diskPath string, outputPath string) error {
    fmt.Println("Generando reporte TREE en:", outputPath)
    return os.WriteFile(outputPath, []byte("REPORTE TREE"), 0644)
}