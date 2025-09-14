package stores

// Estructura para manejar la sesión actual
type Session struct {
    IsActive    bool   // Si hay sesión iniciada
    User        string // Usuario actual
    PartitionID string // ID de la partición montada
}

// Variable global de sesión
var CurrentSession *Session

// Iniciar sesión
func StartSession(user string, partitionID string) {
    CurrentSession = &Session{
        IsActive:    true,
        User:        user,
        PartitionID: partitionID,
    }
}

// Cerrar sesión
func EndSession() {
    CurrentSession = nil
}

// Verificar si hay sesión activa
func IsAuthenticated() bool {
    return CurrentSession != nil && CurrentSession.IsActive
}