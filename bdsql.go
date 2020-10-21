/*
Package bdsql gestiona de manera simple, rápida y eficiente; las sentencias
que se realizan con el motor de base de datos Mysql/MariaDB.*/
package bdsql

import (
	"database/sql"
	"sync"
)

// Conectar crea una conección con el motor de base de datos Mysql/MariaDB.
func Conectar(dsn string, maxConAbiertas, maxConOciosas int) (*BD, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errorNuevo().asignarOrigen(err).asignarMotivoConexionAbrir()
	}
	if err = db.Ping(); err != nil {
		return nil, errorNuevo().asignarOrigen(err).asignarMotivoConexionAbrir()
	}
	db.SetMaxOpenConns(maxConAbiertas)
	db.SetMaxIdleConns(maxConOciosas)

	var bd = &BD{db: db}
	bd.setencias = make(map[string]string)

	return bd, nil
}

// BD representa el pool de conexiones con la base de datos.
type BD struct {
	db  *sql.DB    // manejador de la base de datos
	mux sync.Mutex // bloqueador de exclusión mutua

	// sentencias almacena sentencias SQL para que no vuelvan a
	// ser generadas por cada llamada
	setencias map[string]string
}

// Insertar representa la sentencia 'insert' de SQL.
func (bd *BD) Insertar(nombre string) *insertar {
	var o = &insertar{bd: bd, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// Modificar representa la sentencia 'update' de SQL.
func (bd *BD) Modificar(nombre string) *modificar {
	var o = &modificar{bd: bd, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// Eliminar representa la sentencia 'delete' de SQL.
func (bd *BD) Eliminar(nombre string) *eliminar {
	var o = &eliminar{bd: bd, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// Seleccionar representa la sentencia 'select' de SQL.
func (bd *BD) Seleccionar(nombre string) *seleccionar {
	var o = &seleccionar{bd: bd, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// TxIniciar inicia una nueva transacción.
// Representa a la sentencia 'Begin' de SQL.
func (bd *BD) TxIniciar() (*TX, error) {
	txBD, err := bd.db.Begin()
	if err != nil {
		return nil, errorNuevo().asignarOrigen(err).asignarMotivoTxIniciar()
	}

	return &TX{bd: bd, tx: txBD}, nil
}

// Cerrar cierra la conexión con la base de datos.
func (bd *BD) Cerrar() error {
	if err := bd.db.Close(); err != nil {
		return errorNuevo().asignarOrigen(err).asignarMotivoConexionCerrar()
	}

	return nil
}

// -----------------------------------------------------------------------------

// TX representa a una transacción de la base de datos.
// Es quien tiene la resposabilidad de mantener y otorgar
// las sentencias a ejecutar de SQL dentro de una transacción.
type TX struct {
	bd *BD
	tx *sql.Tx
}

// Insertar representa la sentencia 'insert' de SQL.
func (tx *TX) Insertar(nombre string) *insertar {
	var o = &insertar{bd: tx.bd, tx: tx.tx, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// Modificar representa la sentencia 'update' de SQL.
func (tx *TX) Modificar(nombre string) *modificar {
	var o = &modificar{bd: tx.bd, tx: tx.tx, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// Eliminar representa la sentencia 'delete' de SQL.
func (tx *TX) Eliminar(nombre string) *eliminar {
	var o = &eliminar{bd: tx.bd, tx: tx.tx, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// Seleccionar representa la sentencia 'select' de SQL.
func (tx *TX) Seleccionar(nombre string) *seleccionar {
	var o = &seleccionar{bd: tx.bd, tx: tx.tx, senSQLNombre: nombre}
	o.senSQL, o.senSQLExiste = o.bd.obtenerSentenciaSQL(nombre)

	return o
}

// SeleccionarSql(sentencia string, valores ...interface{}) *seleccionarSql

// TxConfirmar representa a la sentencia 'commit' de SQL.
func (tx *TX) TxConfirmar() error {
	if err := tx.tx.Commit(); err != nil {
		return errorNuevo().asignarOrigen(err).asignarMotivoTxConfirmar()
	}

	return nil
}

// TxRevertir representa a la sentencia 'rollback' de SQL.
func (tx *TX) TxRevertir() error {
	if err := tx.tx.Rollback(); err != nil {
		return errorNuevo().asignarOrigen(err).asignarMotivoTxRevertir()
	}

	return nil
}

// TxFinalizar decide si confirma o revierte la transacción
// tomando en cuenta el error recibido.
// Si es nulo, confirma la transacción (realiza 'commit').
// En caso de recibir un error, revierte la transacción (realiza 'rollback').
func (tx *TX) TxFinalizar(err error) error {
	if err != nil {
		return tx.TxRevertir()
	}

	return tx.TxConfirmar()
}

// -----------------------------------------------------------------------------
// funciones internas

func (bd *BD) obtenerSentenciaSQL(nombre string) (string, bool) {
	bd.mux.Lock()
	s, ok := bd.setencias[nombre]
	bd.mux.Unlock()

	return s, ok
}

func (bd *BD) guardarSentenciaSQL(nombre, sentenciaSQL string) {
	bd.mux.Lock()
	bd.setencias[nombre] = sentenciaSQL
	bd.mux.Unlock()
}
