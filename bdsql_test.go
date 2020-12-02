package bdsql

import (
	"fmt"
	"testing"
)

const (
	dsn                = "root:mysql@tcp(localhost:3306)/bdsql?charset=utf8&parseTime=true&clientFoundRows=true"
	parametrosPostgres = "user=postgres password=postgres dbname=pruebas sslmode=disable"
)

func TestSeleccionar(t *testing.T) {
	t.Skip()
	bd, err := Conectar(dsn, 10, 0)
	if err != nil {
		t.Error("No es posible conectar con la bd:", err)
	}

	seleccionarEnCosas(bd)
}
func TestModificar(t *testing.T) {
	t.Skip()
	bd, err := Conectar(dsn, 10, 0)
	if err != nil {
		t.Error("No es posible conectar con la bd:", err)
	}
	if err := modificarEnCosas(bd, 11, "doce modificado otra vez", false, "obs modi"); err != nil {
		fmt.Println("-1-----------------------------")
		if errBdsql, ok := EsError(err); ok {
			fmt.Println("Error: ", err)
			fmt.Println("Origen:", errBdsql.ObtenerOrigen())
		}
		fmt.Println("-------------------------------")
	} else {
		fmt.Println("Modificado!")
	}

	if err := modificarEnCosas(bd, 12, "doce modificado otra vez", false, "obs modi"); err != nil {
		fmt.Println("-2-----------------------------")
		if errBdsql, ok := EsError(err); ok {
			fmt.Println("Error: ", err)
			fmt.Println("Origen:", errBdsql.ObtenerOrigen())
		}
		fmt.Println("-------------------------------")
	} else {
		fmt.Println("Modificado!")
	}

}
func TestInsertarEnA(t *testing.T) {
	t.Skip()
	bd, err := Conectar(dsn, 10, 0)
	if err != nil {
		t.Error("No es posible conectar con la bd:", err)
	}
	id, err := insertarEnCosas(bd, "un nombre efg10", "{\"id\": 13, \"nombre\": \"un nombre\", \"unArray\": [{\"id\": 1}, {\"id\": 2}]", true, "observaciones locas")
	if err != nil {
		fmt.Println("-1-----------------------------")
		fmt.Printf("No se insertó nada\n%v\n", err)
		if errBdsql, ok := EsError(err); ok {
			fmt.Println("Origen:", errBdsql.ObtenerOrigen())
			fmt.Println("Es entrada duplicada:", errBdsql.EsEntradaDuplicada())
		}
		fmt.Println("-------------------------------")
	} else {
		fmt.Println("Primer id insertado:", id)
	}

	// id, err = insertarEnCosas(bd, "un nombre efg", "{\"id\": 13, \"nombre\": \"un nombre\", \"unArray\": [{\"id\": 1}, {\"id\": 2}]}", true, "observaciones locas")
	// if err != nil {
	// 	fmt.Println("-2-----------------------------")
	// 	fmt.Printf("No se insertó nada\n%v\n", err)
	// 	if errBdsql, ok := EsError(err); ok {
	// 		fmt.Println("Origen:", errBdsql.ObtenerOrigen())
	// 		fmt.Println("Es entrada duplicada:", errBdsql.EsEntradaDuplicada())
	// 	}
	// 	fmt.Println("-------------------------------")
	// } else {
	// 	fmt.Println("Segundo id insertado:", id)
	// }
}
func TestInsertarEnSentenciaPreparada(t *testing.T) {
	// t.Skip()
	bd, err := Conectar(dsn, 10, 0)
	if err != nil {
		t.Error("No es posible conectar con la bd:", err)
	}
	sp, err := bd.
		Insertar("").
		Tabla("cosas").
		Campos("nombre", "datos", "es_activo", "observaciones").
		SentenciaPreparada()

	defer sp.Cerrar()

	var id int64
	for i := 0; i < 3; i++ {
		if err := sp.Valores(fmt.Sprintf("nombre-%v", i), "{}", true, fmt.Sprintf("observaciones-%v", i)).ObtenerID(&id).Ejecutar(); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Insertado:", id)
	}
}

func insertarEnCosas(bd *BD, nombre, datos string, esActivo bool, observaciones string) (int64, error) {
	var id int64
	err := bd.
		Insertar("cosasInsertar").
		Tabla("cosas").
		Campos("nombre", "datos", "es_activo", "observaciones").
		Valores(nombre, datos, esActivo, observaciones).
		ObtenerID(&id).
		Ejecutar()

	return id, err
}

func modificarEnCosas(bd *BD, id int64, nombre string, esActivo bool, observaciones string) error {
	mod := bd.
		Modificar("cosasModificar").
		Tabla("cosas").
		Campos("nombre", "es_activo", "observaciones").
		Valores(nombre, esActivo, observaciones).
		Condicion("id = ? ", id)

	// s, err := mod.SQL()
	// fmt.Println("SQL:", s, err)
	return mod.Ejecutar()

	// return err
}

func seleccionarEnCosas(bd *BD) error {
	var cosas = []struct {
		ID            int64  `bdsql:"id"`
		Nombre        string `bdsql:"nombre"`
		Datos         string
		EsActivo      bool   `bdsql:"es_activo"`
		Observaciones string `bdsql:"observaciones"`
	}{}

	var cant, err = bd.Seleccionar("seleccionarEnCosas").Tabla("cosas").Campos("*").Resultado(&cosas).Ejecutar()
	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}
	fmt.Println("Cantidad:", cant)
	for _, cosa := range cosas {
		fmt.Println(cosa)
	}

	return nil
}
