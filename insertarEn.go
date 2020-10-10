package bdsql

import (
	"database/sql"
	"fmt"
	"strings"
)

type insertarEn struct {
	bd *baseDeDatos

	tx *sql.Tx

	tabla   string
	campos  []string
	valores []interface{}

	sentenciaSQLNombre string
	sentenciaSQL       string
	sentenciaSQLExiste bool

	idPtr *int64 // puntero de la variable o campo de un objeto donde se guardará el valor del útlimo id insertado de la tabla
}

// Sentencia establece el nombre de la sentencia para ser utilizada ante
// una nueva llamada.
func (o *insertarEn) Sentencia(nombre string) *insertarEn {
	o.sentenciaSQLNombre = nombre
	o.sentenciaSQL, o.sentenciaSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// Campos establece los nombres de campos a actualizar.
func (o *insertarEn) Campos(campos ...string) *insertarEn {
	if o.sentenciaSQLExiste {
		return o
	}

	o.campos = campos
	return o
}

// Valores establece los valores que recibirán los campos a actualizar.
func (o *insertarEn) Valores(valores ...interface{}) *insertarEn {
	o.valores = valores

	return o
}

// ObtenerID obtiene el último id insertado de la tabla.
// Se debe utilizar para los casos en que la tabla contenga
// una clave principal (PK) del tipo autoincremental.
func (o *insertarEn) ObtenerID(varPtr *int64) *insertarEn {
	o.idPtr = varPtr

	return o
}

// SQL devuelve la sentencia SQL.
func (o *insertarEn) SQL() (string, error) {
	if o.sentenciaSQLExiste {
		return o.sentenciaSQL, nil
	}

	if err := o.validarSQL(); err != nil {
		return "", err
	}

	return o.generarSQL(), nil
}

// Ejecutar ejecuta la sentencia sql.
func (o *insertarEn) Ejecutar() error {
	if o.sentenciaSQLExiste {
		return o.ejecutar()
	}

	// validar la
	if err := o.validarSQL(); err != nil {
		return err
	}

	// verificar que los valores no se encuentren vacíos
	if len(o.valores) == 0 {
		return errorNuevo().asignarMotivoValoresVacios()
	}

	// verificar que la cantidad de campos coincida con la cantidad de valores recibidos
	if len(o.campos) != len(o.valores) {
		return errorNuevo().asignarMotivoCamposValores()
	}

	return o.ejecutar()
}

// validarSQL valida que la sentencia pueda ser creada correctamente.
func (o *insertarEn) validarSQL() error {
	err := errorNuevo()

	// verficar que el nombre de la tabla no se encuentre vacía
	if o.tabla == "" {
		err.asignarMotivoNombreDeTablaVacia()
	}
	// verificar que los nombres de campos a insertar no se encuentren vacíos
	if len(o.campos) == 0 {
		err.asignarMotivoNombresDeCamposVacios()
	}

	if len(err.mensajes) != 0 {
		return err
	}

	return nil
}

func (o *insertarEn) generarSQL() string {
	var s = fmt.Sprintf("insert into %v (%v) values (%v);",
		o.tabla,
		strings.Join(o.campos, ", "),
		strings.Repeat("?, ", len(o.campos)-1)+"?",
	)

	o.bd.guardarSentenciaSQL(o.sentenciaSQLNombre, s)

	return s
}

func (o *insertarEn) ejecutar() error {
	var err error
	var res sql.Result

	// obtener sentencia 'insert...'
	sentencia := o.generarSQL()

	if o.bd.db != nil {
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
