package bdsql

import (
	"database/sql"
)

// insertarEn es la estructura que mantiene los datos de la sentencia
// 'insert' de sql.
type insertarEn struct {
	db *sql.DB
	tx *sql.Tx

	tabla   string
	campos  []string
	valores []interface{}

	idCampo string // nombre del campo PK del cual se obtiene el último id insertado de la tabla.
	idPtr   *int64 // puntero de la variable donde se guardará el valor del útlimo id insertado de la tabla.

	ejecutorDeInsercion
}

// Campos son los nombres de campos (las columnas) a actualizar.
func (o *insertarEn) Campos(campos ...string) *insertarEn {
	o.campos = campos

	return o
}

// Valores son los valores que recibirán los campos (las columnas).
func (o *insertarEn) Valores(valores ...interface{}) *insertarEn {
	o.valores = valores

	return o
}

// ObtenerID obtiene el último id insertado.
func (o *insertarEn) ObtenerID(campo string, varPtr *int64) *insertarEn {
	o.idCampo, o.idPtr = campo, varPtr

	return o
}

// SQL devuelve la sentencia sQL.
func (o *insertarEn) SQL() (string, error) {
	if err := o.validarSQL(); err != nil {
		return "", err
	}

	return o.sql(), nil
}

// Ejecutar ejecuta la sentencia sql.
func (o *insertarEn) Ejecutar() error {
	if err := o.validarSQL(); err != nil {
		return err
	}

	// Verificar que los valores no se encuentren vacíos.
	if len(o.valores) == 0 {
		return errorNuevo().motivoValoresVacios()
	}

	// Verificar que la cantidad de campos coincida con la cantidad de valores recibidos.
	if len(o.campos) != len(o.valores) {
		return errorNuevo().motivoCamposValores()
	}

	return o.ejecutar()
}

// validarSQL valida que la sentencia pueda ser creada correctamente.
func (o *insertarEn) validarSQL() error {
	err := errorNuevo()

	// Verficar que el nombre de la tabla no se encuentre vacía.
	if o.tabla == "" {
		err.motivoNombreDeTablaVacia()
	}

	// Verificar que los nombres de campos a insertar no se encuentren vacíos.
	if len(o.campos) == 0 {
		err.motivoNombresDeCamposVacios()
	}

	if len(err.mensajes) != 0 {
		return err
	}

	return nil
}
