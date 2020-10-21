package bdsql

import (
	"fmt"

	"database/sql"
)

type modificar struct {
	bd *BD
	tx *sql.Tx

	tabla   string
	campos  []string
	valores []interface{}

	condicion        string
	condicionValores []interface{}

	limite int

	senSQLExiste bool
	senSQLNombre string
	senSQL       string
}

// Tabla establece el nombre de la tabla a modificar.
func (o *modificar) Tabla(tabla string) *modificar {
	if o.senSQLExiste {
		return o
	}

	o.tabla = tabla
	return o
}

// Campos establece los nombres de campos a actualizar.
func (o *modificar) Campos(campos ...string) *modificar {
	o.campos = campos
	return o
}

// Valores establece los valores que recibirán los campos a actualizar.
func (o *modificar) Valores(valores ...interface{}) *modificar {
	o.valores = valores

	return o
}

// Condicion implementa la cláusula 'where' de la sentencia 'update'.
func (o *modificar) Condicion(condicion string, valores ...interface{}) *modificar {
	if !o.senSQLExiste {
		o.condicion = condicion
	}
	o.condicionValores = valores

	return o
}

// Limitar implementa la cláusula 'limit' de la sentencia 'update'.
func (o *modificar) Limitar(limite int) *modificar {
	o.limite = limite

	return o
}

// SQL devuelve la sentencia SQL.
func (o *modificar) SQL() (string, error) {
	return o.generarSQL()
}

// SentenciaPreparada devuelve una sentencia preparada para ser utilizada
// múltiples veces.
func (o *modificar) SentenciaPreparada() (*sentenciaPreparadaModificar, error) {
	var sentencia, err = o.generarSQL()
	if err != nil {
		return nil, err
	}

	var sp = &sentenciaPreparadaModificar{cantCampos: len(o.campos)}
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
func (o *modificar) Ejecutar() error {
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
		res, err = o.bd.db.Exec(sentencia, append(o.valores, o.condicionValores...)...)
	} else {
		// ejecución dentro de una transacción
		res, err = o.tx.Exec(sentencia, append(o.valores, o.condicionValores...)...)
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

func (o *modificar) generarSQL() (string, error) {
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
	// verificar que la condición no se encuentre vacía
	// no se permite realizar modificaciones sin condición
	if o.condicion == "" {
		err.asignarMotivoCondicionVacia()
	}
	if len(err.mensajes) != 0 {
		return "", err
	}

	// campos
	var campos string
	for _, v := range o.campos {
		if campos != "" {
			campos += ", "
		}
		campos += v + " = ?"
	}
	// sentencia con cláusula where
	var sentencia = fmt.Sprintf("update %v set %v where %v", o.tabla, campos, o.condicion)
	// limit
	if o.limite > 0 {
		sentencia += fmt.Sprintf(" limit %v", o.limite)
	}

	sentencia += ";"
	o.bd.guardarSentenciaSQL(o.senSQLNombre, sentencia)

	return sentencia, nil
}

// -----------------------------------------------------------------------------

type sentenciaPreparadaModificar struct {
	stmt *sql.Stmt

	cantCampos int
	valores    []interface{}
}

// Valores establece los valores que recibirán los campos a actualizar.
func (o *sentenciaPreparadaModificar) Valores(valores ...interface{}) *sentenciaPreparadaModificar {
	o.valores = valores

	return o
}

// Ejecutar ejecuta la sentencia SQL.
func (o *sentenciaPreparadaModificar) Ejecutar() error {
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

	if cant, err := res.RowsAffected(); err != nil {
		return errorNuevo().asignarMotivoObtencionDeRegistrosAfectados()
	} else if cant == 0 {
		return errorNuevo().asignarMotivoNingunRegistroAfectado()
	}

	return nil
}

// Cerrar cierra la sentencia preparada.
func (o *sentenciaPreparadaModificar) Cerrar() error {
	return resolverErrorMysql(o.stmt.Close())
}
