/*
Package bdsql gestiona de manera simple, rápida y eficiente; las sentencias
que se realizan con el motor de base de datos Mysql/MariaDB.*/
package bdsql

import (
	"database/sql"
	"sync"
)

// Conectar crea una conección con el motor de base de datos Mysql/MariaDB.
func Conectar(dsn string, maxConAbiertas, maxConOciosas int) (*baseDeDatos, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errorNuevo().asignarOrigen(err).asignarMotivoConexionAbrir()
	}
	err = db.Ping()
	if err != nil {
		return nil, errorNuevo().asignarOrigen(err).asignarMotivoConexionAbrir()
	}
	db.SetMaxOpenConns(maxConAbiertas)
	db.SetMaxIdleConns(maxConOciosas)

	var bd = &baseDeDatos{db: db}
	bd.setencias = make(map[string]string)

	return bd, nil
}

type baseDeDatos struct {
	db   *sql.DB    // manejador de la base de datos
	exmu sync.Mutex // manejador de exclusión mutua

	// sentencias almacena sentencias SQL para que no vuelvan a
	// ser generadas por cada llamada
	setencias map[string]string
}

// InsertarEn representa a la sentencia 'insert' de sql.
func (bd *baseDeDatos) InsertarEn(tabla string) *insertarEn {
	return &insertarEn{bd: bd, tabla: tabla}
}

func (bd *baseDeDatos) obtenerSentenciaSQL(nombre string) (string, bool) {
	bd.exmu.Lock()
	s, ok := bd.setencias[nombre]
	bd.exmu.Unlock()

	return s, ok
}

func (bd *baseDeDatos) guardarSentenciaSQL(nombre, sentenciaSQL string) {
	bd.exmu.Lock()
	bd.setencias[nombre] = sentenciaSQL
	bd.exmu.Unlock()
}
