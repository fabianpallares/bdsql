package bdsql

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// EsError devuelve el error del paquete y un valor lógico que confirma el tipo.
func EsError(err error) (*errorPaquete, bool) {
	ep, ok := err.(*errorPaquete)
	return ep, ok
}

type errorPaquete struct {
	// origen (causa) del error (error original de la base de datos)
	origen error

	// mensajes de error
	mensajes []string

	// los diversos motivos (causas) del origen del error
	errorMotivos struct {
		// apertura y cierre de conexión
		esConexionAbrir  bool // no es posible conectase con el motor de la base de datos
		esConexionCerrar bool // no es posible cerrar la conexión con la base de datos

		// genéricos de validación del paquete (GENERAR SQL)
		esNombreDeTablaVacia    bool // el nombre de la tabla se encuentra vacía
		esNombresDeCamposVacios bool // los nombres de los campos se encuentran vacíos
		esCondicionVacia        bool // no se ha recibido la condición para ejecutar la sentencia. Se aplica a 'update' y 'delete'
		esValoresCondicionVacia bool // no se han recibido los valores de la condición para ejecutar la sentencia. Se aplica a 'update' y 'delete'

		// genéricos de validación del paquete (EJECUCION SQL)
		esValoresVacios                  bool // no se han recibido valores para poder ejecutar la sentencia
		esCamposValoresDiferenteCantidad bool // la cantidad de campos no coincide con la cantidad de valores recibidos

		// no atrapado
		esErrorNoAtrapado bool // No posible ejecutar la sentencia, debido a que se ha producido un error inesperado en la base de datos y el error no fue atrapado

		// causados por la base de datos, pero han sido atrapados
		esTablaInexistente              bool // el nombre de la tabla en la base de datos es inexistente
		esCampoDeTablaInexistente       bool // no es posible ejecutar la sentencia porque el nombre de campo es inexistente
		esEntradaDuplicada              bool // la tabla contiene un un campo con clave única y el valor recibido ya existe
		esTipoDeCampoIncorrecto         bool // se intenta guardar un valor en un campo de una tabla donde el tipo de valor es incorrecto
		esTipoDeCampoJSONIncorrecto     bool // se intenta guardar un valor en un campo JSON de una tabla donde el tipo de valor es incorrecto
		esCampoFueraDeRango             bool // no es posible ejecutar la sentencia porque hay al menos un valor que se desea guardar que supera el límite permitido po el campo
		esObtencionDeRegistrosAfectados bool // error al obtener la cantidad de registros afectados
		esNingunRegistroAfectado        bool // elemento inexistente o existen otros elementos con los mismos valores o no se ha cambiado ningún valor del elemento

		// insertar
		esObtencionDeID bool // No es posible obtener el id insertado

		// seleccionar
		esSeleccionarPunteroDeSlice        bool // El objeto recibido no es un puntero de slice de estructura
		esSeleccionarCamposSinRelacion     bool // No es posible ejecutar la sentencia porque los campos de la estructura del objeto recibido no tienen asignados la relación con los campos de la tabla de la base de datos
		esSeleccionarContieneEstructura    bool // No es posible recibir un objeto que contenga dentro otra estructura
		esSeleccionarTipoDeCampoIncorrecto bool // No es posible ejecutar la sentencia porque existe al menos un campo de la estructura que contiene un tipo erroneo (no se permiten punteros)
		esSeleccionarCamposFaltantes       bool // Los campos obtenidos de la consulta, no existen en su totalidad en la estructura
		esSeleccionarLecturaDeCampos       bool // No es posible leer los campos de la consulta
		esSeleccionarAsignacionDeCampos    bool // No es posible asignar los campos de la consulta de la base de datos a los campos de la estructura

		// sentencia preparada
		esSentenciaPreparadaCrear bool // No es posible crear la sentencia preparada

		// transacción
		esTxIniciar   bool // error al intentar iniciar una transacción
		esTxConfirmar bool // error al intentar confirmar la transacción
		esTxRevertir  bool // error al intentar revertir la transacción
	}
}

// Error devuelve el mensaje de error.
func (err *errorPaquete) Error() string {
	return strings.Join(err.mensajes, ". ")
}

// Origen devuelve el error de origen (error original).
func (err *errorPaquete) ObtenerOrigen() error {
	return err.origen
}

func (err *errorPaquete) EsConexionAbrir() bool    { return err.errorMotivos.esConexionAbrir }
func (err *errorPaquete) EsConexionCerrar() bool   { return err.errorMotivos.esConexionCerrar }
func (err *errorPaquete) EsErrorNoAtrapado() bool  { return err.errorMotivos.esErrorNoAtrapado }
func (err *errorPaquete) EsEntradaDuplicada() bool { return err.errorMotivos.esEntradaDuplicada }

func (err *errorPaquete) EsNombreDeTablaVacia() bool { return err.errorMotivos.esNombreDeTablaVacia }
func (err *errorPaquete) EsNombresDeCamposVacios() bool {
	return err.errorMotivos.esNombresDeCamposVacios
}
func (err *errorPaquete) EsValoresVacios() bool { return err.errorMotivos.esValoresVacios }
func (err *errorPaquete) EsCamposValoresDiferenteCantidad() bool {
	return err.errorMotivos.esCamposValoresDiferenteCantidad
}
func (err *errorPaquete) EsCondicionVacia() bool { return err.errorMotivos.esCondicionVacia }
func (err *errorPaquete) EsValoresCondicionVacia() bool {
	return err.errorMotivos.esValoresCondicionVacia
}
func (err *errorPaquete) EsTablaInexistente() bool { return err.errorMotivos.esTablaInexistente }
func (err *errorPaquete) EsCampoDeTablaInexistente() bool {
	return err.errorMotivos.esCampoDeTablaInexistente
}
func (err *errorPaquete) EsTipoDeCampoIncorrecto() bool {
	return err.errorMotivos.esTipoDeCampoIncorrecto
}
func (err *errorPaquete) EsTipoDeCampoJSONIncorrecto() bool {
	return err.errorMotivos.esTipoDeCampoJSONIncorrecto
}
func (err *errorPaquete) EsObtencionDeRegistrosAfectados() bool {
	return err.errorMotivos.esObtencionDeRegistrosAfectados
}
func (err *errorPaquete) EsNingunRegistroAfectado() bool {
	return err.errorMotivos.esNingunRegistroAfectado
}
func (err *errorPaquete) EsCampoFueraDeRango() bool { return err.errorMotivos.esCampoFueraDeRango }
func (err *errorPaquete) EsObtencionDeID() bool     { return err.errorMotivos.esObtencionDeID }
func (err *errorPaquete) EsSeleccionarPunteroDeSlice() bool {
	return err.errorMotivos.esSeleccionarPunteroDeSlice
}
func (err *errorPaquete) EsSeleccionarCamposSinRelacion() bool {
	return err.errorMotivos.esSeleccionarCamposSinRelacion
}
func (err *errorPaquete) EsSeleccionarContieneEstructura() bool {
	return err.errorMotivos.esSeleccionarContieneEstructura
}
func (err *errorPaquete) EsSeleccionarTipoDeCampoIncorrecto() bool {
	return err.errorMotivos.esSeleccionarTipoDeCampoIncorrecto
}
func (err *errorPaquete) EsSeleccionarCamposFaltantes() bool {
	return err.errorMotivos.esSeleccionarCamposFaltantes
}
func (err *errorPaquete) EsSeleccionarLecturaDeCampos() bool {
	return err.errorMotivos.esSeleccionarLecturaDeCampos
}
func (err *errorPaquete) EsSeleccionarAsignacionDeCampos() bool {
	return err.errorMotivos.esSeleccionarAsignacionDeCampos
}
func (err *errorPaquete) EsSentenciaPreparadaCrear() bool {
	return err.errorMotivos.esSentenciaPreparadaCrear
}
func (err *errorPaquete) EsTxIniciar() bool {
	return err.errorMotivos.esTxIniciar
}
func (err *errorPaquete) EsTxConfirmar() bool {
	return err.errorMotivos.esTxConfirmar
}
func (err *errorPaquete) EsTxRevertir() bool {
	return err.errorMotivos.esTxRevertir
}

// -----------------------------------------------------------------------------

func (err *errorPaquete) asignarOrigen(origen error) *errorPaquete {
	err.origen = origen
	return err
}

func (err *errorPaquete) asignarMotivoConexionAbrir() *errorPaquete {
	err.mensajes = append(err.mensajes, "Error al conectarse con la base de datos")
	err.errorMotivos.esConexionAbrir = true
	return err
}
func (err *errorPaquete) asignarMotivoConexionCerrar() *errorPaquete {
	err.mensajes = append(err.mensajes, "Error al cerrar la conexión con la base de datos")
	err.errorMotivos.esConexionCerrar = true
	return err
}
func (err *errorPaquete) asignarMotivoNombreDeTablaVacia() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible generar la sentencia SQL. El nombre de la tabla se encuentra vacía")
	err.errorMotivos.esNombreDeTablaVacia = true
	return err
}
func (err *errorPaquete) asignarMotivoNombresDeCamposVacios() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible generar la sentencia SQL. La lista de nombres de campos se encuentra vacía")
	err.errorMotivos.esNombresDeCamposVacios = true
	return err
}
func (err *errorPaquete) asignarMotivoCondicionVacia() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible generar la sentencia SQL. No se permite que la condición (cláusula where) se encuentra vacía")
	err.errorMotivos.esCondicionVacia = true
	return err
}
func (err *errorPaquete) asignarMotivoValoresCondicionVacia() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible generar la sentencia SQL. No se han recibido los valores de la condición")
	err.errorMotivos.esValoresCondicionVacia = true
	return err
}
func (err *errorPaquete) asignarMotivoValoresVacios() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. La lista de valores de los campos no han sido asignados")
	err.errorMotivos.esValoresVacios = true
	return err
}
func (err *errorPaquete) asignarMotivoCamposValoresDiferenteCantidad() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. La cantidad de nombres de campos no coincide con la cantidad de valores recibidos")
	err.errorMotivos.esCamposValoresDiferenteCantidad = true
	return err
}
func (err *errorPaquete) asignarMotivoErrorNoAtrapado() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Se produjo un error inesperado")
	err.errorMotivos.esErrorNoAtrapado = true
	return err
}
func (err *errorPaquete) asignarMotivoTablaInexistente() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. El nombre de la tabla no existe en la base de datos")
	err.errorMotivos.esTablaInexistente = true
	return err
}
func (err *errorPaquete) asignarMotivoCampoInexistente() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Hay al menos un nombre de campo que no existe en la tabla")
	err.errorMotivos.esCampoDeTablaInexistente = true
	return err
}
func (err *errorPaquete) asignarMotivoEntradaDuplicada() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible guardar los datos. Ya existe un campo que contiene el mismo valor que se ha recibido (entrada duplicada)")
	err.errorMotivos.esEntradaDuplicada = true
	return err
}
func (err *errorPaquete) asignarMotivoTipoDeCampoIncorrecto() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible guardar los datos. Existe al menos un campo de la tabla que está recibiendo un tipo de valor incorrecto")
	err.errorMotivos.esTipoDeCampoIncorrecto = true
	return err
}
func (err *errorPaquete) asignarMotivoTipoDeCampoJSONIncorrecto() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible guardar los datos. Existe al menos un campo JSON de la tabla que está recibiendo un tipo de valor incorrecto")
	err.errorMotivos.esTipoDeCampoJSONIncorrecto = true
	return err
}
func (err *errorPaquete) asignarMotivoCampoFueraDeRango() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Existe al menos un valor recibido que supera el límite permitido por el campo de la tabla")
	err.errorMotivos.esCampoFueraDeRango = true
	return err
}
func (err *errorPaquete) asignarMotivoObtencionDeRegistrosAfectados() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Se produjo un error al obtener la cantidad de registros afectados")
	err.errorMotivos.esObtencionDeRegistrosAfectados = true
	return err
}
func (err *errorPaquete) asignarMotivoNingunRegistroAfectado() *errorPaquete {
	err.mensajes = append(err.mensajes, "La sentencia se ejecutó satisfactoriamente pero ningún registro de la tabla fue afectado. Los posibles motivos son: Elemento inexistente o no se ha cambiado ningún valor del registro")
	err.errorMotivos.esNingunRegistroAfectado = true
	return err
}
func (err *errorPaquete) asignarMotivoObtencionDeID() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Se produjo un error al obtener el identificador insertado")
	err.errorMotivos.esObtencionDeID = true
	return err
}
func (err *errorPaquete) asignarMotivoSeleccionarPunteroDeSlice() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. El objeto recibido no es un puntero de slice de estructura")
	err.errorMotivos.esSeleccionarPunteroDeSlice = true
	return err
}
func (err *errorPaquete) asignarMotivoSeleccionarCamposSinRelacion() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Los campos de la estructura del objeto recibido no tienen asignados la relación de nombres con los campos de la tabla de la base de datos")
	err.errorMotivos.esSeleccionarCamposSinRelacion = true
	return err
}
func (err *errorPaquete) asignarMotivoSeleccionarContieneEstructura() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. El objeto recibido contiene al menos una estructura")
	err.errorMotivos.esSeleccionarContieneEstructura = true
	return err
}
func (err *errorPaquete) asignarMotivoSeleccionarTipoDeCampoIncorrecto() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Existe al menos un campo de la estructura que contiene un tipo erroneo (no se permiten punteros)")
	err.errorMotivos.esSeleccionarTipoDeCampoIncorrecto = true
	return err
}
func (err *errorPaquete) asignarMotivoSeleccionarCamposFaltantes(camposFaltantes string) *errorPaquete {
	err.mensajes = append(err.mensajes, fmt.Sprintf("No es posible ejecutar la sentencia SQL. Los campos obtenidos de la consulta no existen en su totalidad dentro de la estructura. Los campos faltantes son: %v", camposFaltantes))
	err.errorMotivos.esSeleccionarCamposFaltantes = true
	return err
}
func (err *errorPaquete) asignarMotivoSeleccionarLecturaDeCampos() *errorPaquete {
	err.mensajes = append(err.mensajes, "No es posible ejecutar la sentencia SQL. Se produjo un error al leer los campos de la consulta")
	err.errorMotivos.esSeleccionarLecturaDeCampos = true
	return err
}
func (err *errorPaquete) asignarMotivoSeleccionarAsignacionDeCampos(mensaje string) *errorPaquete {
	err.mensajes = append(err.mensajes, fmt.Sprintf("No es posible ejecutar la sentencia SQL. Se produjo un error al asignar los campos de la consulta a los campos de la estructura. %v", mensaje))
	err.errorMotivos.esSeleccionarAsignacionDeCampos = true
	return err
}
func (err *errorPaquete) asignarMotivoSentenciaPreparadaCrear() *errorPaquete {
	err.mensajes = append(err.mensajes, "Error al crear la sentencia preparada")
	err.errorMotivos.esSentenciaPreparadaCrear = true
	return err
}
func (err *errorPaquete) asignarMotivoTxIniciar() *errorPaquete {
	err.mensajes = append(err.mensajes, "Error al intentar iniciar una transacción")
	err.errorMotivos.esTxIniciar = true
	return err
}
func (err *errorPaquete) asignarMotivoTxConfirmar() *errorPaquete {
	err.mensajes = append(err.mensajes, "Error al intentar confirmar la transacción")
	err.errorMotivos.esTxConfirmar = true
	return err
}
func (err *errorPaquete) asignarMotivoTxRevertir() *errorPaquete {
	err.mensajes = append(err.mensajes, "Error al intentar revertir la transacción")
	err.errorMotivos.esTxRevertir = true
	return err
}

// -----------------------------------------------------------------------------

func errorNuevo() *errorPaquete {
	return &errorPaquete{}
}

func resolverErrorMysql(err error) error {
	if err == nil {
		return nil
	}

	errMysql, ok := err.(*mysql.MySQLError)
	if !ok {
		return errorNuevo().asignarOrigen(err).asignarMotivoErrorNoAtrapado()
	}

	switch errMysql.Number {
	case 1146:
		// Nombre de tabla inexistente.
		return errorNuevo().asignarOrigen(errMysql).asignarMotivoTablaInexistente()
	case 1054:
		// Nombre de campo de la tabla inexistente.
		return errorNuevo().asignarOrigen(errMysql).asignarMotivoCampoInexistente()
	case 1062:
		// Entrada duplicada.
		return errorNuevo().asignarOrigen(errMysql).asignarMotivoEntradaDuplicada()
	case 1264:
		// Campo fuera de rango (se quiere guardar un valor superior a la capacidad del campo).
		return errorNuevo().asignarOrigen(errMysql).asignarMotivoCampoFueraDeRango()
	case 1366:
		// Tipo de campo incorrecto.
		return errorNuevo().asignarOrigen(errMysql).asignarMotivoTipoDeCampoIncorrecto()
	case 3140:
		// Tipo de campo incorrecto (JSON inválido).
		return errorNuevo().asignarOrigen(errMysql).asignarMotivoTipoDeCampoJSONIncorrecto()
	default:
		// No atrapado.
		return errorNuevo().asignarOrigen(errMysql).asignarMotivoErrorNoAtrapado()
	}
}
