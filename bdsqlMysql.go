package bdsql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func errorMysql(err error) error {
	if err == nil {
		return nil
	}

	ep, ok := err.(*errorPaquete)
	if !ok {
		return err
	}

	if ep.origen == nil {
		return err
	}

	errMysql, ok := ep.origen.(*mysql.MySQLError)
	if !ok {
		return err
	}

	switch errMysql.Number {
	case 1146:
		// Nombre de tabla inexistente.
		return errorNuevo().asignarOrigen(errMysql).motivoTablaInexistente()
	case 1054:
		// Nombre de campo de la tabla inexistente.
		return errorNuevo().asignarOrigen(errMysql).motivoCampoInexistente()
	case 1062:
		// Entrada duplicada.
		return errorNuevo().asignarOrigen(errMysql).motivoEntradaDuplicada()
	case 1264:
		// Campo fuera de rango (se quiere guardar un valor superior a la capacidad del campo).
		return errorNuevo().asignarOrigen(errMysql).motivoCampoFueraDeRango()
	case 1366:
		// Tipo de campo incorrecto.
		return errorNuevo().asignarOrigen(errMysql).motivoTipoDeCampoIncorrecto()
	}

	return err
}

// ConectarConMysql crea y retorna una conección con el motor de base de datos Mysql.
func ConectarConMysql(dsn string, maxConAbiertas, maxConOciosas int) (Almacenador, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errorNuevo().asignarOrigen(err).motivoConexion()
	}
	err = db.Ping()
	if err != nil {
		return nil, errorNuevo().asignarOrigen(err).motivoConexion()
	}
	db.SetMaxOpenConns(maxConAbiertas)
	db.SetMaxIdleConns(maxConOciosas)

	return &almacenadorMysql{db: db}, nil
}

// ConectarConMariaDB crea y retorna una conección con el motor de base de datos MariaDB.
func ConectarConMariaDB(dsn string, maxConAbiertas, maxConOciosas int) (Almacenador, error) {
	return ConectarConMysql(dsn, maxConAbiertas, maxConOciosas)
}

// almacenadorMysql es la base de datos para Mysql.
// Debe implementar todas las funciones de la interface 'BaseDeDatos'.
type almacenadorMysql struct{ db *sql.DB }

// InsertarEn representa a la sentencia 'insert' de sql.
func (b *almacenadorMysql) InsertarEn(tabla string) *insertarEn {
	sen := &insertarEn{db: b.db, tabla: tabla}
	sen.ejecutorDeInsercion = &insertarEnMysql{sentenciaPtr: sen}
	return sen
}

// insertarEnMysql es el ejecutor de la sentencia 'insert' de sql
// para el motor de base de datos Mysql.
type insertarEnMysql struct{ sentenciaPtr *insertarEn }

// sql devuelve la sentencia sql.
func (o *insertarEnMysql) sql() string {
	return fmt.Sprintf("insert into %v (%v) values (%v);",
		o.sentenciaPtr.tabla,
		strings.Join(o.sentenciaPtr.campos, ", "),
		strings.Repeat("?, ", len(o.sentenciaPtr.campos)-1)+"?",
	)
}

// ejecutar ejecuta la sentencia sql.
func (o *insertarEnMysql) ejecutar() error {
	var err error
	var res sql.Result

	// obtener sentencia 'insert...'.
	sentencia := o.sql()

	if o.sentenciaPtr.db != nil {
		// Ejecución fuera de una transacción.
		res, err = o.sentenciaPtr.db.Exec(sentencia, o.sentenciaPtr.valores...)
	} else {
		// Ejecución dentro de una transacción.
		res, err = o.sentenciaPtr.tx.Exec(sentencia, o.sentenciaPtr.valores...)
	}
	if err != nil {
		return errorMysql(errorNuevo().asignarOrigen(err).motivoEjecucionFallida())
	}

	// Obtener el último id insertado.
	if o.sentenciaPtr.idPtr != nil {
		*o.sentenciaPtr.idPtr, err = res.LastInsertId()
		if err != nil {
			return errorMysql(errorNuevo().asignarOrigen(err).motivoObtencionDeID())
		}
	}

	return nil
}
