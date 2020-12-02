package bdsql

import (
	"fmt"

	"database/sql"
)

type eliminar struct {
	bd *BD
	tx *sql.Tx

	tabla string

	condicion        string
	condicionValores []interface{}

	limite int

	senSQLExiste bool
	senSQLNombre string
	senSQL       string
}

// Tabla establece el nombre de la tabla donde se eliminarán los registros.
func (o *eliminar) Tabla(tabla string) *eliminar {
	if o.senSQLExiste {
		return o
	}

	o.tabla = tabla
	return o
}

// Condicion implementa la cláusula 'where' de la sentencia 'delete'.
func (o *eliminar) Condicion(condicion string, valores ...interface{}) *eliminar {
	if !o.senSQLExiste {
		o.condicion = condicion
	}
	o.condicionValores = valores

	return o
}

// Limitar implementa la cláusula 'limit' de la sentencia 'delete'.
func (o *eliminar) Limitar(limite int) *eliminar {
	o.limite = limite

	return o
}

// SQL devuelve la sentencia SQL.
func (o *eliminar) SQL() (string, error) {
	return o.generarSQL()
}

// SentenciaPreparada devuelve una sentencia preparada para ser utilizada
// múltiples veces.
func (o *eliminar) SentenciaPreparada() (*sentenciaPreparadaEliminar, error) {
	var sentencia, err = o.generarSQL()
	if err != nil {
		return nil, err
	}

	var sp = new(sentenciaPreparadaEliminar)
	if o.tx == nil {
		// ejecución fuera de una transacción
		sp.stmt, err = o.bd.db.Prepare(sentencia)
	} else {
		// ejecución dentro de una transacción
		sp.stmt, err = o.tx.Prepare(sentencia)
	}
	if err != nil {
		return nil, errorNuevo().asignarOrigen(err).asignarMotivoSentenciaPreparadaCrear()
	}

	return sp, nil
}

// Ejecutar ejecuta la sentencia SQL.
func (o *eliminar) Ejecutar() error {
	var sentencia, err = o.generarSQL()
	if err != nil {
		return err
	}

	var errEjec = errorNuevo()
	// verificar que los valores de la condición no se encuentren vacíos
	if len(o.condicionValores) == 0 {
		errEjec.asignarMotivoValoresCondicionVacia()
	}
	if len(errEjec.mensajes) != 0 {
		return errEjec
	}

	var res sql.Result
	if o.tx == nil {
		// ejecución fuera de una transacción
		res, err = o.bd.db.Exec(sentencia, o.condicionValores...)
	} else {
		// ejecución dentro de una transacción
		res, err = o.tx.Exec(sentencia, o.condicionValores...)
	}
	if err != nil {
		return resolverErrorMysql(err)
	}

	if cant, err := res.RowsAffected(); err != nil {
		return errorNuevo().asignarMotivoObtencionDeRegistrosAfectados()
	} else if cant == 0 {
		return errorNuevo().asignarMotivoNingunRegistroAfectado()
	}

	return nil
}

func (o *eliminar) generarSQL() (string, error) {
	if o.senSQLExiste {
		return o.senSQL, nil
	}

	var err = errorNuevo()
	// verificar que el nombre de la tabla no se encuentre vacía
	if o.tabla == "" {
		err.asignarMotivoNombreDeTablaVacia()
	}
	// verificar que la condición no se encuentre vacía
	// no se permite realizar eliminaciones sin condición
	if o.condicion == "" {
		err.asignarMotivoCondicionVacia()
	}
	// verificar que los valores de la condición no se encuentren vacíos
	if len(o.condicionValores) == 0 {
		err.asignarMotivoValoresCondicionVacia()
	}
	if len(err.mensajes) != 0 {
		return "", err
	}

	// sentencia con cláusula where
	var sentencia = fmt.Sprintf("delete from %v where %v", o.tabla, o.condicion)
	// limit
	if o.limite > 0 {
		sentencia += fmt.Sprintf(" limit %v", o.limite)
	}

	sentencia += ";"
	o.bd.guardarSentenciaSQL(o.senSQLNombre, sentencia)

	return sentencia, nil
}

// -----------------------------------------------------------------------------

type sentenciaPreparadaEliminar struct {
	stmt *sql.Stmt

	valores []interface{}
}

// Valores establece los valores que recibirán los campos a actualizar.
func (o *sentenciaPreparadaEliminar) Valores(valores ...interface{}) *sentenciaPreparadaEliminar {
	o.valores = valores

	return o
}

// Ejecutar ejecuta la sentencia SQL.
func (o *sentenciaPreparadaEliminar) Ejecutar() error {
	// verificar que los valores no se encuentren vacíos
	if len(o.valores) == 0 {
		return errorNuevo().asignarMotivoValoresVacios()
	}

	res, err := o.stmt.Exec(o.valores...)
	if err != nil {
		return resolverErrorMysql(err)
	}

	if cant, err := res.RowsAffected(); err != nil {
		return errorNuevo().asignarMotivoObtencionDeRegistrosAfectados()
	} else if cant == 0 {
		return errorNuevo().asignarMotivoNingunRegistroAfectado()
	}

	return nil
}

// Cerrar cierra la sentencia preparada.
func (o *sentenciaPreparadaEliminar) Cerrar() error {
	return resolverErrorMysql(o.stmt.Close())
}
