# bdsql: Manejador de bases de datos SQL para Go/Golang (Mysql y Mariadb)

[![Go Report Card](https://goreportcard.com/badge/github.com/fabianpallares/bdsql)](https://goreportcard.com/report/github.com/fabianpallares/bdsql) [![GoDoc](https://godoc.org/github.com/fabianpallares/bdsql?status.svg)](https://godoc.org/github.com/fabianpallares/bdsql)

Escribir y programar las funciones de persistencia de datos, es un trabajo un tanto tedioso y reiterativo.
Son muchas líneas de código repetidas que hay que escribir para cada entidad de negocio. El paquete **bdsql** intenta hacer esta tarea más simple, agradable y automática, pudiendo escribir instrucciones SQL en español.

Existen muchos paquetes que simplifican el código para trabajar contra la base de datos; la mayoría de ellos son implementaciones de ORM.

_**bdsql no es un ORM.**_<br />
_**bdsql no es un mapeador de objetos de bases de datos relacionales.**_

Con el paquete **bdsql**, se ecriben (en español) instrucciones SQL estandar sin importar el motor de base de datos.

En ningún momento existe de manera automática el mapeo de los datos de las estructuras con los datos de las tablas de la base de datos.

## Instalación:
Para instalar el paquete utilice la siguiente sentencia:
```
go get -u github.com/fabianpallares/bdsql
```

## Conección con Mysql y Mariadb:
Para conectarse con el motor de base de datos Mysql/Mariadb, **bdsql** dispone de la
siguiente funcion:

```GO
package main

import (
	"github.com/fabianpallares/bdsql"
)

func main() {
	dsn := "root:mysql@tcp(localhost:3306)/pruebas?charset=utf8&parseTime=true&clientFoundRows=true"
	maxConAbiertas, maxConOciosas := 10, 0

	bd, err = bdsql.Conectar(dsn, maxConAbiertas, maxConOciosas)
	if err != nil {
		// No se ha podido conectar, tratar el error.
	}

	// Aquí está viva la base de datos a través de la variable bd.
}
```

## Insertando datos:
Una vez que se dispone de una conección con la base de datos; estamos en condiciones de trabajar con ella. Comencemos con la sentencia 'insert'.

```GO
err := bd.
	Insertar("personasInsertar").
	Tabla("personas").
	Campos("apellidos", "nombres", "activo").
	Valores("Un apellido", "Un nombre", true).
	Ejecutar()
if err != nil {
	// No es posible insertar, tratar el error.
}
```

También es posible obtener el último id insertado:
```GO
var id int64
err := bd.
	Insertar("personasInsertar").
	Tabla("personas").
	Campos("apellidos", "nombres", "activo").
	Valores("Un apellido", "Un nombre", true).
	ObtenerID(&id).
	Ejecutar()
if err != nil {
	// No es posible insertar, tratar el error.
}
fmt.Println("Id insertado:", id)
```

## Modificando datos:
La sentencia 'update' se utiliza de la siguiente manera:
```GO
err := bd.
	Modificar("personasModificar").
	Tabla("personas").
	Campos("apellidos", "nombres", "activo").
	Valores("Un apellido", "Un nombre", true).
	Condicion("id = ?", 2).
	Ejecutar()
if err != nil {
	// No es posible modificar, tratar el error.
}
```

## Eliminando datos:
La sentencia 'delete' se utiliza de la siguiente manera:
```GO
err = bd.
	Eliminar("personasEliminar").
	Tabla("personas").
	Condicion("id = ?", 2).
	Ejecutar()
if err != nil {
	// No es posible eliminar, tratar el error.
}

```

## Obteniendo la sentencia SQL generada:
Para conocer la sentencia SQL que genera el paquete **bdsql**, se utilizará el método SQL() de cada sentencia:

```GO
sql, err := bd.
	Modificar("personasInsertar").
	Tabla("personas").
	Campos("apellidos", "nombres", "activo").
	SQL()
if err != nil {
	// Instrucción mal escrita, tratar el error.
}

fmt.Println(sql)

// Genera la sentencia SQL:
insert into personas (apellidos, nombres, activo) values (?, ?, ?);
```

## Seleccionando datos:
La sentencia 'select' se utiliza de la siguiente manera:

```GO
// Obtener todos los datos de la tabla personas:
var personas []struct {
	Apellidos string `bdsql:"apellidos"`
	Nombres   string `bdsql:"nombres"`
}{}

cant, err := bd.
	Seleccionar("personasSeleccionar").
	Tabla("personas").
	Campos("*").
	Recibir(&personas).
	Ejecutar()
if err != nil {
	// No es posible obtener datos, tratar el error.
}
fmt.Println("Filas obtenidas:", cant)

// Operaciones de SQL:
var max = []struct{
	ID int bdsql:"id"
}{}

cant, err := bd.
	Seleccionar("personasSeleccionar").
	Tabla("personas").
	Campos("max(id) as id").
	Recibir(&max).
	Ejecutar()
if err != nil {
	// No es posible obtener datos, tratar el error.
}
if cant == 0 {
	fmt.Println("tabla vacía")
	return
}
fmt.Println("Id máximo:", max[0].ID)

// Ordenamiento:
cant, err := bd.
	Seleccionar("personasSeleccionar").
	Tabla("personas").
	Campos("id", "apellidos", "nombres", "activo").
	OrdenarPor("apellidos desc", "id").
	Recibir(<objeto>).
	Ejecutar()

// Agrupamiento:
cant, err := bd.
	Seleccionar("personasSeleccionar").
	Tabla("personas").
	Campos("apellidos", "count(apellidos) as cantidad").
	AgruparPor("apellidos").
	OrdenarPor("apellidos desc", "cantidad").
	Recibir(<objeto>).
	Ejecutar()

// Operaciones mas complejas:
// Realizar inner join con tabla de teléfonos, obtener los teléfonos de las
// personas que se encuentren en determinada zona y con un estado activo.
// Obtener la segunda página de 100 elementos por página.
cant, filas, err := bd.
	Seleccionar("personasSeleccionar").
	Tabla("personas").
	Campos("p.apellidos", "p.nombres", "t.telefono").
	JuntarCon("telefonos t", "t.persona_id = p.id").
	Condicion("t.zona = ? and t.activo = ?", "centro", true).
	OrdenarPor("apllidos asc", "telefono desc").
	Limitar(100).
	Saltar(100).
	Recibir(<objeto>).
	Ejecutar()
```

## Transacciones:
Las transacciones son muy simples de utilizar con el paquete **bdsql**.
Una vez que nos hemos conectado con el motor, lo primero que haremos es crear
una transacción:
```GO
tx, err := bd.TxIniciar()
if err != nil {
	// No es posible crear la transacción, tratar el error.
}
```

En la variable 'tx' se obtiene el manejador de la transacción.
A partir de ahora, todas las acciones que se realicen con la base de datos, se deben hacer con la transacción ('tx').
```GO
func insertarEnTransaccion(bd bdsql.BD) {
	tx, err := bd.TxIniciar()
	if err != nil {
		// No es posible crear la transacción, tratar el error
		return
	}

	// Crear una variable de error, la cual guardará el resultado de la
	// última acción. Según el resultado del error, se realizará
	// Commit o Rollback de la transacción.
	var errTx error

	// Cuando finalice la función se ejecutará Commit o Rollback según el resutado en errTx...
	defer func() {
		if errTx != nil {
			tx.TxRevertir()
			return
		}
		tx.TxConfirmar()
	}()

	var id int64
	errTx = tx.
		Insertar("personasInsertar").
		Tabla("personas").
		Campos("apellidos", "nombres", "activo").
		Valores("Un apellido", "Un nombre", true).
		ObtenerID(&id).
		Ejecutar()
	if errTx != nil {
		// No es posible insertar a la persona, tratar el error.
		return
	}

	fmt.Println("Persona insertada:", id)

	telefonos := []string{"1234", "5678", "2222", "3333"}

	var telId int64
	for _, tel := range telefonos {
		errTx = tx.
			Insertar("telefonosInsertar").
			Tabla("telefonos").
			Campos("persona_id", "telefono").
			Valores(id, tel).
			ObtenerID(&telId).
			Ejecutar()
		if errTx != nil {
			// No es posible insertar el teléfono de la persona, tratar el error.
			return
		}
		fmt.Println("Teléfono insertado:", telId)
	}
}
```

Las mismas operaciones que pueden hacerse con la base de datos, pueden realizarse dentro de una transacción (Insertar, Modificar, Eliminar y Seleccionar).

## Sentencias preparadas:
Las sentencias preparadas agilizan la ejecución cuando hay que realizar repetidamente la misma acción.
Son ideales para ser utilizadas dentro de una transacción. Cada sentencia de insersión, modificación y eliminación poseen la generación de sentencias preparadas.
Una vez que se obtiene la sentencia preparada, solo hay que pasarle los parámetros para que realice la instrucción ya preparada con los datos que recibe.

```GO
func insertarEnTransaccionConSentenciasPreparadas(bd bdsql.BD) {
	tx, err := bd.TxIniciar()
	if err != nil {
		// No es posible crear la transacción, tratar el error
		return
	}

	var errTx error
	defer func() {
		// TxFinalizar determina si realiza Commit o Rollback basándose en el error.
		err := tx.TxFinalizar(errTx)
		if err != nil {
			// Hubo un error en la transacción, tratar el error
			return
		}
		// Todo perfecto
	}()

	var id int64
	errTx = tx.
		Insertar("personasInsertar").
		Tabla("personas").
		Campos("apellidos", "nombres", "activo").
		Valores("Un apellido", "Un nombre", true).
		ObtenerID(&id).
		Ejecutar()
	if errTx != nil {
		fmt.Println("error al insertar una persona")
		return
	}

	fmt.Println("Persona insertada:", id)

	telefonos := []string{"1234", "5678", "2222", "3333"}

	var telId int64
	stmt, errTx := tx.
		Insertar("telefonosInsertar").
		Insertar("telefonos").
		Campos("persona_id", "telefono").
		ObtenerID(&telId).
		SentenciaPreparada()
	if errTx != nil {
		// No es posible generar la sentencia preparada, tratar el error.
		return
	}

	for _, tel := range telefonos {
		// Guardar los teléfonos utilizando la sentencia preparada.
		errTx = stmt.Parametros(id, tel).Ejecutar()
		if errTx != nil {
			// No es posible realizar la operación, tratar el error.
			return
		}
		fmt.Println("Teléfono insertado:", telId)
	}
}
```

## Manejando errores:
En todo momento puede conocerse que sucedió exactamente con el error.
Para esto, el paquete **bdsql** cuenta con un método el cual obtiene el tipo de 
error del paquete:
```GO
// Si el error obtenido fue probocado por el paquete, se obtiene la información.
if bdError, ok := bdsql.EsError(err); ok {
	// Es un error no encontrado ?
	if bdError.EsErrorNoEncontrado() {
		// La sentencia llegó al motor de la base de datos y no puede ser
		// ejecutada porque el registro que desea modificar o eliminar no ha
		// sido encontrado.
	}
}
```
#### Documentación:
[Documentación en godoc](https://godoc.org/github.com/fabianpallares/bdsql)

