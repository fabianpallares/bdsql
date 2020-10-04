/*
Package bdsql gestiona de manera simple, rápida y eficiente; todas las acciones
que se realizan con un motor de base de datos SQL.

Los motores de bases de datos que pueden ser utilizados son:
	mysql
	mariadb
	postgres

El paquete ofrece la misma funcionalidad y uso, indistintamente
al motor de base de datos que se haya conectado.*/
package bdsql

// Almacenador representa a la base de datos.
// Se encarga de ejecutar sentencias SQL en la base de datos.
type Almacenador interface {
	InsertarEn(tabla string) *insertarEn
	// ModificarEn(tabla string) *modificarEn
	// EliminarEn(tabla string) *eliminarEn
	// SeleccionarDe(tabla string) *seleccionarDe
	// SeleccionarSql(sentencia string, valores ...interface{}) *seleccionarSql
	// TxIniciar() (transaccionador, error)
}

// Transaccionador representa a una transacción de la base de datos.
// Se encarga de ejecutar sentencias SQL dentro de una transacción.
type Transaccionador interface {
	InsertarEn(tabla string) *insertarEn
	// ModificarEn(tabla string) *modificarEn
	// EliminarEn(tabla string) *eliminarEn
	// SeleccionarDe(tabla string) *seleccionarDe
	// SeleccionarSql(sentencia string, valores ...interface{}) *seleccionarSql
	// TxConfirmar() error
	// TxRevertir() error
	// TxFinalizar(err error) error
}

// ejecutorDeInsercion establece la funcionalidad de la sentencia
// 'insert' de SQL.
type ejecutorDeInsercion interface {
	// crearSentenciaPreparada() (*sentenciaPreparadaInsertar, error)
	sql() string
	ejecutar() error
}
