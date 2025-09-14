# Manual Técnico - Sistema de Gestión de Discos y Archivos

## Índice
1. [Introducción](#introducción)
2. [Arquitectura del Sistema](#arquitectura-del-sistema)
3. [Estructuras de Datos Principales](#estructuras-de-datos-principales)
   3.1 [MBR](#mbr)  
   3.2 [Partición](#partición)  
   3.3 [SuperBlock](#superblock)  
   3.4 [Inode](#inode)  
   3.5 [Bloques de Archivo y Carpeta](#bloques-de-archivo-y-carpeta)  
4. [Comandos y Funcionamiento Interno](#comandos-y-funcionamiento-interno)
   4.1 [Gestión de Discos](#gestión-de-discos)  
   4.2 [Gestión de Particiones](#gestión-de-particiones)  
   4.3 [Montaje de Particiones](#montaje-de-particiones)  
   4.4 [Sistemas de Archivos](#sistemas-de-archivos)  
   4.5 [Gestión de Usuarios y Grupos](#gestión-de-usuarios-y-grupos)  
   4.6 [Reportes](#reportes)  
5. [Seguridad y Manejo de Sesiones](#seguridad-y-manejo-de-sesiones)
6. [Persistencia y Serialización](#persistencia-y-serialización)
7. [Errores Comunes y Soluciones](#errores-comunes-y-soluciones)
8. [Recomendaciones Técnicas](#recomendaciones-técnicas)

---

## 1. Introducción
El sistema simula un **gestor de discos virtuales** que permite:
- Crear y eliminar discos y particiones.
- Montar particiones y asignar IDs únicos.
- Formatear particiones con sistemas de archivos EXT2 y EXT3.
- Gestionar archivos y directorios.
- Manejar usuarios y grupos.
- Generar reportes visuales del disco y sus estructuras internas.

Se diseñó para estudiantes y profesionales que desean aprender cómo funcionan los sistemas de archivos de manera didáctica.

---

## 2. Arquitectura del Sistema
El sistema sigue una arquitectura tipo **MVC simplificada** con los siguientes componentes:
- **Commands:** Implementa los comandos del sistema (`mkdisk`, `fdisk`, `mkfs`, `mkfile`, `mkdir`, `login`, `logout`, `rep`, etc.)
- **Structures:** Contiene las estructuras de datos como `SuperBlock`, `Inode`, `FileBlock`, `FolderBlock`, `Partition` y `MBR`.
- **Stores:** Administra las particiones montadas, la sesión activa y el acceso a las estructuras.
- **Analyzer:** Se encarga de recibir el input del usuario, parsear los comandos y enviarlos a la capa de `Commands`.
- **Reports:** Genera reportes gráficos en formato `.pdf` o `.dot` basados en la información de los discos y archivos.

---

## 3. Estructuras de Datos Principales

### 3.1 MBR
El **Master Boot Record** almacena:
- Tamaño total del disco.
- Fecha de creación.
- Fit del disco (FF/BF/WF).
- Tabla de particiones: hasta 4 entradas de particiones primarias o extendidas.

```go
type MBR struct {
    Mbr_tamano int32
    Mbr_fecha_creacion [20]byte
    Mbr_disk_signature int32
    Disk_fit byte
    Mbr_partitions [4]Partition
}
```

### 3.2 Partición
Cada partición tiene:
- Nombre, tamaño, tipo (P/E/L) y fit (FF/BF/WF).
- Inicio en bytes dentro del disco.
- Información sobre particiones extendidas y lógicas.

```go
type Partition struct {
    Part_status byte
    Part_type byte
    Part_fit byte
    Part_start int32
    Part_size int32
    Part_name [16]byte
}
```

### 3.3 SuperBlock
El **SuperBlock** controla la partición formateada:
- Tipo de sistema de archivos (2fs/3fs).
- Contadores de inodos y bloques libres.
- Punteros iniciales de inodos, bloques y bitmaps.
- Información de montaje y timestamps.

```go
type SuperBlock struct {
    S_filesystem_type int32
    S_inodes_count int32
    S_blocks_count int32
    S_free_inodes_count int32
    S_free_blocks_count int32
    S_mtime float32
    S_umtime float32
    S_mnt_count int32
    S_magic int32
    S_inode_size int32
    S_block_size int32
    S_first_ino int32
    S_first_blo int32
    S_bm_inode_start int32
    S_bm_block_start int32
    S_inode_start int32
    S_block_start int32
}
```

### 3.4 Inode
El **Inode** almacena metadatos de archivos y carpetas:
- UID y GID.
- Tamaño, tipo (archivo o carpeta) y permisos.
- Punteros a bloques de datos (directos, indirectos).

```go
type Inode struct {
    I_uid int32
    I_gid int32
    I_size int32
    I_atime float32
    I_ctime float32
    I_mtime float32
    I_block [15]int32
    I_type [1]byte
    I_perm [3]byte
}
```

### 3.5 Bloques de Archivo y Carpeta
- **FileBlock:** Contiene los datos de un archivo.  
- **FolderBlock:** Contiene referencias a subcarpetas y archivos.

```go
type FileBlock struct {
    B_content [64]byte
}

type FolderBlock struct {
    B_content [4]FolderContent
}

type FolderContent struct {
    B_name [16]byte
    B_inodo int32
}
```

---

## 4. Comandos y Funcionamiento Interno

### 4.1 Gestión de Discos
```
mkdisk -size=<tamaño> -path=<ruta> -fit=<FF/BF/WF>
rmdisk -path=<ruta>
```

### 4.2 Gestión de Particiones
```
fdisk -size=<tamaño> -path=<ruta> -name=<nombre> -type=<P/E/L> -fit=<FF/BF/WF>
fdisk -delete=<fast/full> -name=<nombre> -path=<ruta>
```

### 4.3 Montaje de Particiones
```
mount -path=<ruta_disco> -name=<nombre_partición>
mounted
```

### 4.4 Sistemas de Archivos
```
mkfs -id=<ID_partición> -type=full -fs=<2fs/3fs>
mkdir -path=<ruta> -r
mkfile -path=<ruta> -size=<tamaño> -cont=<contenido> -r
cat -path=<ruta>
```

### 4.5 Gestión de Usuarios y Grupos
```
login -user=<usuario> -pass=<contraseña> -id=<ID_partición>
logout
mkgrp -name=<nombre_grupo>
rmgrp -name=<nombre_grupo>
mkusr -user=<nombre> -pass=<contraseña> -grp=<grupo>
rmusr -user=<nombre_usuario>
chgrp -user=<nombre_usuario> -grp=<nuevo_grupo>
```

### 4.6 Reportes
```
rep -id=<ID_partición> -path=<ruta_reporte> -name=<tipo_reporte> [-path_file_ls=<ruta_ls>]
```
- Tipos: mbr, tree, disk, inode, block, bm_inode, bm_block, sb, file, ls  

---

## 5. Seguridad y Manejo de Sesiones
- Se requiere **login** para ejecutar comandos que modifiquen la partición.  
- **Logout** limpia la sesión y protege las operaciones críticas.  
- Los reportes dependen de tener la partición montada y la sesión activa.  

---

## 6. Persistencia y Serialización
- **Serialize:** Guarda estructuras en disco (SuperBlock, Inode, FileBlock).  
- **Deserialize:** Recupera estructuras desde disco.  
- **Bitmaps:** Controlan inodos y bloques libres.  
- **Users.txt:** Archivo dentro del FS que almacena usuarios y grupos.  

---

## 7. Errores Comunes y Soluciones
- **no se ha iniciado sesión:** Ejecutar `login` en una partición montada.  
- **partición no encontrada:** Verificar que la partición esté montada.  
- **comando desconocido:** Revisar sintaxis y nombre del comando.  

---

## 8. Recomendaciones Técnicas
- Montar la partición antes de formatear.  
- Mantener rutas absolutas.  
- Evitar caracteres especiales en nombres.  
