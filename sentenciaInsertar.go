package bdsql

import (
	"fmt"
	"strings"

	"database/sql"
)

type insertar struct {
	bd *BD
	tx *sql.Tx

	tabla   string
	campos  []string
	valores []interface{}

	idPtr *int64 // puntero de la variable o campo de un objeto donde se guardará el valor del útlimo id insertado de la tabla

	senSQLExiste bool
	senSQLNombre string
	senSQL       string
}

// Tabla establece el nombre de la tabla a insertar.
func (o *insertar) Tabla(tabla string) *insertar {
	if o.senSQLExiste {
		return o
	}

	o.tabla = tabla
	return o
}

// Campos establece los nombres de campos a actualizar.
func (o *insertar) Campos(campos ...string) *insertar {
	o.campos = campos
	return o
}

// Valores establece los valores que recibirán los campos a actualizar.
func (o *insertar) Valores(valores ...interface{}) *insertar {
	o.valores = valores

	return o
}

// ObtenerID obtiene el último id insertado de la tabla.
// Se debe utilizar para los casos en que la tabla contenga
// una clave principal (PK) del tipo autoincremental.
func (o *insertar) ObtenerID(varPtr *int64) *insertar {
	o.idPtr = varPtr

	return o
}

// SQL devuelve la sentencia SQL.
func (o *insertar) SQL() (string, error) {
	return o.generarSQL()
}

// SentenciaPreparada devuelve una sentencia preparada para ser utilizada
// múltiples veces.
func (o *insertar) SentenciaPreparada() (*sentenciaPreparadaInsertar, error) {
	var sentencia, err = o.generarSQL()
	if err != nil {
		return nil, err
	}

	var sp = &sentenciaPreparadaInsertar{cantCampos: len(o.campos)}
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
func (o *insertar) Ejecutar() error {
	var sentencia, err = o.generarSQL()
	if err != nil {
		return err
	}

	var errEjec = errorNuevo()
	// verificar que los valores no se encuentren vacíos
	if len(o.valores) == 0 {
		errEjec.asignarMotivoValoresVacios()
	}
	// verificar que la cantidad de campos coincida con la cantidad de valores recibidos
	if len(o.campos) != len(o.valores) {
		errEjec.asignarMotivoCamposValoresDiferenteCantidad()
	}
	if len(errEjec.mensajes) != 0 {
		return errEjec
	}

	var res sql.Result
	if o.tx == nil {
		// ejecución fuera de una transacción
		res, err = o.bd.db.Exec(sentencia, o.valores...)
	} else {
		// ejecución dentro de una transacción
		res, err = o.tx.Exec(sentencia, o.valores...)
	}
	if err != nil {
		return resolverErrorMysql(err)
	}

	// obtener el último id insertado
	if o.idPtr != nil {
		if *o.idPtr, err = res.LastInsertId(); err != nil {
			return resolverErrorMysql(err)
		}
	}

	return nil
}

func (o *insertar) generarSQL() (string, error) {
	if o.senSQLExiste {
		return o.senSQL, nil
	}

	var err = errorNuevo()
	// verificar que el nombre de la tabla no se encuentre vacía
	if o.tabla == "" {
		err.asignarMotivoNombreDeTablaVacia()
	}
	// verificar que los nombres de campos no se encuentren vacíos
	if len(o.campos) == 0 {
		err.asignarMotivoNombresDeCamposVacios()
	}
	if len(err.mensajes) != 0 {
		return "", err
	}

	var sentencia = fmt.Sprintf("insert into %v (%v) values (%v);",
		o.tabla,
		strings.Join(o.campos, ", "),
		strings.Repeat("?, ", len(o.campos)-1)+"?",
	)
	o.bd.guardarSentenciaSQL(o.senSQLNombre, sentencia)

	return sentencia, nil
}

// -----------------------------------------------------------------------------

type sentenciaPreparadaInsertar struct {
	stmt *sql.Stmt

	cantCampos int
	valores    []interface{}

	idPtr *int64 // puntero de la variable o campo de un objeto donde se guardará el valor del útlimo id insertado de la tabla
}

// Valores establece los valores que recibirán los campos a actualizar.
func (o *sentenciaPreparadaInsertar) Valores(valores ...interface{}) *sentenciaPreparadaInsertar {
	o.valores = valores

	return o
}

// ObtenerID obtiene el último id insertado de la tabla.
// Se debe utilizar para los casos en que la tabla contenga
// una clave principal (PK) del tipo autoincremental.
func (o *sentenciaPreparadaInsertar) ObtenerID(varPtr *int64) *sentenciaPreparadaInsertar {
	o.idPtr = varPtr

	return o
}

// Ejecutar ejecuta la sentencia SQL.
func (o *sentenciaPreparadaInsertar) Ejecutar() error {
	var errEjec = errorNuevo()
	// verificar que los valores no se encuentren vacíos
	if len(o.valores) == 0 {
		return errEjec.asignarMotivoValoresVacios()
	}
	// verificar que la cantidad de campos coincida con la cantidad de valores recibidos
	if o.cantCampos != len(o.valores) {
		errEjec.asignarMotivoCamposValoresDiferenteCantidad()
	}
	if len(errEjec.mensajes) != 0 {
		return errEjec
	}

	res, err := o.stmt.Exec(o.valores...)
	if err != nil {
		return resolverErrorMysql(err)
	}

	// obtener el último id insertado
	if o.idPtr != nil {
		if *o.idPtr, err = res.LastInsertId(); err != nil {
			return resolverErrorMysql(err)
		}
	}

	return nil
}

// Cerrar cierra la sentencia preparada.
func (o *sentenciaPreparadaInsertar) Cerrar() error {
	return resolverErrorMysql(o.stmt.Close())
}
