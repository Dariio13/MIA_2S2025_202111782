# Manual de Usuario - Sistema de Gestión de Discos y Archivos

## Índice
1. [Introducción](#introducción)  
2. [Instalación](#instalación)  
3. [Comandos del Sistema](#comandos-del-sistema)  
   3.1 [Gestión de Discos](#gestión-de-discos)  
   3.2 [Gestión de Particiones](#gestión-de-particiones)  
   3.3 [Gestión de Montaje](#gestión-de-montaje)  
   3.4 [Sistemas de Archivos](#sistemas-de-archivos)  
   3.5 [Gestión de Usuarios y Grupos](#gestión-de-usuarios-y-grupos)  
   3.6 [Reportes](#reportes)  
4. [Ejemplos de Uso](#ejemplos-de-uso)  
5. [Notas](#notas)  
6. [Recomendaciones](#recomendaciones)  

---

## 1. Introducción
Este sistema permite la creación y gestión de discos virtuales, particiones, sistemas de archivos (EXT2/EXT3), usuarios y grupos, así como la ejecución de comandos para manipular archivos y generar reportes. Está pensado para estudiantes de Ingeniería y administración de sistemas, permitiendo simular la gestión de un sistema de archivos y disco desde consola.

---

## 2. Instalación
1. Asegúrate de tener instalado **Go 1.20+**.  
2. Clona el repositorio en tu máquina:
   ```bash
   git clone <URL_REPO>
   cd backend
   ```

---

## 3. Comandos del Sistema

### 3.1 Gestión de Discos
**Crear disco**
```
mkdisk -size=<tamaño> -path=<ruta> -fit=<FF/BF/WF>
```
- `-size` → tamaño del disco en MB o KB.  
- `-path` → ruta completa donde se guardará el disco.  
- `-fit` → tipo de ajuste: FF (First Fit), BF (Best Fit), WF (Worst Fit).  

**Eliminar disco**
```
rmdisk -path=<ruta>
```
Elimina el disco ubicado en la ruta especificada.

---

### 3.2 Gestión de Particiones
**Crear partición**
```
fdisk -size=<tamaño> -path=<ruta> -name=<nombre> -type=<P/E/L> -fit=<FF/BF/WF>
```
- `-size` → tamaño de la partición.  
- `-type` → tipo: P (Primaria), E (Extendida), L (Lógica).  
- `-fit` → ajuste de partición: FF, BF, WF.  

**Eliminar partición**
```
fdisk -delete=<fast/full> -name=<nombre> -path=<ruta>
```

---

### 3.3 Gestión de Montaje
**Montar partición**
```
mount -path=<ruta_disco> -name=<nombre_partición>
```
Retorna un ID único de la partición montada.

**Ver particiones montadas**
```
mounted
```

---

### 3.4 Sistemas de Archivos
**Formatear partición**
```
mkfs -id=<ID_partición> -type=full -fs=<2fs/3fs>
```
- `-type` → siempre full.  
- `-fs` → sistema de archivos: 2fs (EXT2) o 3fs (EXT3).  

**Crear carpeta**
```
mkdir -path=<ruta> -r
```
- `-r` → opcional, crea recursivamente todos los directorios padres.  

**Crear archivo**
```
mkfile -path=<ruta> -size=<tamaño> -cont=<contenido> -r
```

**Ver contenido de archivo**
```
cat -path=<ruta>
```

---

### 3.5 Gestión de Usuarios y Grupos
**Login**
```
login -user=<usuario> -pass=<contraseña> -id=<ID_partición>
```

**Logout**
```
logout
```

**Crear grupo**
```
mkgrp -name=<nombre_grupo>
```

**Eliminar grupo**
```
rmgrp -name=<nombre_grupo>
```

**Crear usuario**
```
mkusr -user=<nombre> -pass=<contraseña> -grp=<grupo>
```

**Eliminar usuario**
```
rmusr -user=<nombre_usuario>
```

**Cambiar grupo de usuario**
```
chgrp -user=<nombre_usuario> -grp=<nuevo_grupo>
```

---

### 3.6 Reportes
**Generar reportes**
```
rep -id=<ID_partición> -path=<ruta_reporte> -name=<tipo_reporte> [-path_file_ls=<ruta_ls>]
```
- Tipos: mbr, tree, disk, inode, block, bm_inode, bm_block, sb, file, ls.  
- `-path_file_ls` → opcional para reportes tipo ls.  

---

## 4. Ejemplos de Uso
**Crear un disco de 50MB:**
```
mkdisk -size=50M -path="/home/user/disks/disco1.mia" -fit=FF
```

**Crear partición primaria de 10MB:**
```
fdisk -size=10M -path="/home/user/disks/disco1.mia" -name=Part1 -type=P -fit=BF
```

**Formatear partición montada:**
```
mkfs -id=821A -type=full -fs=3fs
```

**Crear usuario root:**
```
mkusr -user=root -pass=123 -grp=admin
```

**Crear carpeta en ext3:**
```
mkdir -path="/carpeta1/subcarpeta2" -r
```

**Crear archivo con contenido:**
```
mkfile -path="/carpeta1/archivo.txt" -size=50 -cont="Hola Mundo"
```

---

## 5. Notas
- Todos los comandos requieren sesión activa en la partición salvo creación de disco y particiones.  
- Los reportes no se pueden generar sin montar la partición correspondiente.  
- Usuarios y grupos se gestionan mediante el archivo `users.txt` dentro del sistema de archivos.  

---

## 6. Recomendaciones
- Siempre montar la partición antes de formatear.  
- Mantener la estructura de carpetas en ruta absoluta `/`.  
- Evitar usar caracteres especiales en nombres de archivos y carpetas.  
