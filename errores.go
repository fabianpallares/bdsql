package bdsql

import (
	"fmt"
	"strings"
)

// motivos almacena los distintos motivos del paquete.
// Determina CUÁL es el motivo del error.
type motivos struct {
	esConexionAbrir  bool // No es posible conectase con el motor de la base de datos.
	esConexionCerrar bool // No es posible cerrar la conexión con la base de datos.

	// Genéricos de validación del paquete.
	esNombreDeTablaVacia    bool // El nombre de la tabla se encuentra vacía.
	esNombresDeCamposVacios bool // Los nombres de los campos se encuentran vacíos.
	esValoresVacios         bool // No se han recibido valores para poder ejecutar la sentencia.
	esCamposValores         bool // La cantidad de campos no coincide con la cantidad de valores recibidos.

	// No atrapado.
	esEjecucionFallida bool // No posible ejecutar la sentencia, debido a que se ha producido un error inesperado en la base de datos y no fue atrapado.

	// Causados por la base de datos, pero han sido atrapados.
	esTablaInexistente              bool // El nombre de la tabla en la base de datos es inexistente.
	esCampoDeTablaInexistente       bool // No es posible ejecutar la sentencia porque el nombre de campo es inexistente.
	esEntradaDuplicada              bool // La tabla contiene un un campo con clave única y el valor recibido ya existe.
	esTipoDeCampoIncorrecto         bool // Se intenta guardar un valor en un campo de una tabla donde el tipo de valor es incorrecto.
	esCampoFueraDeRango             bool // No es posible ejecutar la sentencia porque hay al menos un valor que se desea guardar que supera el límite permitido po el campo
	esObtencionDeRegistrosAfectados bool // Error al obtener la cantidad de registros afectados.
	esNingunRegistroAfectado        bool // Elemento inexistente o existen otros elementos con los mismos valores o no se ha cambiado ningún valor del elemento.

	// Insertar.
	esObtencionDeID bool // No es posible obtener el id insertado.

	// Seleccionar
	esSeleccionarPunteroDeSlice        bool // El objeto recibido no es un puntero de slice de estructura.
	esSeleccionarCamposSinRelacion     bool // No es posible ejecutar la sentencia porque los campos de la estructura del objeto recibido no tienen asignados la relación con los campos de la tabla de la base de datos.
	esSeleccionarContieneEstructura    bool // No es posible recibir un objeto que contenga dentro otra estructura.
	esSeleccionarTipoDeCampoIncorrecto bool // No es posible ejecutar la sentencia porque existe al menos un campo de la estructura que contiene un tipo erroneo (no se permiten punteros).
	esSeleccionarCamposFaltantes       bool // Los campos obtenidos de la consulta, no existen en su totalidad en la estructura.
	esSeleccionarLecturaDeCampos       bool // No es posible leer los campos de la consulta.
	esSeleccionarAsignacionDeCampos    bool // No es posible asignar los campos de la consulta de la base de datos a los campos de la estructura.

	// Transacción
	esTxIniciar   bool // Error al intentar iniciar una transacción.
	esTxConfirmar bool // Error al intentar confirmar la transacción.
	esTxRevertir  bool // Error al intentar revertir la transacción.
}

// errorPaquete es el error del paquete.
type errorPaquete struct {
	origen   error    // Origen (causa) del error (error original).
	mensajes []string // Mensajes de error.
	motivos           // Los diversos motivos (causas) del origen del error.
}

// asignarOrigen asigna el origen del error (envuelve el error) original.
func (err *errorPaquete) asignarOrigen(origen error) *errorPaquete {
	err.origen = origen
	return err
}

// -----------------------------------------------------------------------------
// Métodos internos de cambios de motivos de error (comportamiento).

func (err *errorPaquete) cambiarMotivo(mensaje string, motivo *bool) *errorPaquete {
	err.mensajes = append(err.mensajes, mensaje)
	*motivo = true

	return err
}

func (err *errorPaquete) motivoConexion() *errorPaquete {
	s := "Error al conectarse con la base de datos"
	return err.cambiarMotivo(s, &err.motivos.esConexionAbrir)
}
func (err *errorPaquete) motivoCerrar() *errorPaquete {
	s := "Error al cerrar la conexión con la base de datos"
	return err.cambiarMotivo(s, &err.motivos.esConexionCerrar)
}

// Validación de sentencias

func (err *errorPaquete) motivoNombreDeTablaVacia() *errorPaquete {
	s := "No es posible generar la sentencia porque el nombre de la tabla se encuentra vacía"
	return err.cambiarMotivo(s, &err.motivos.esNombreDeTablaVacia)
}
func (err *errorPaquete) motivoNombresDeCamposVacios() *errorPaquete {
	s := "No es posible generar la sentencia porque la lista de nombres de campos se encuentra vacía"
	return err.cambiarMotivo(s, &err.motivos.esNombresDeCamposVacios)
}
func (err *errorPaquete) motivoValoresVacios() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque la lista de valores de los campos se encuentra vacía"
	return err.cambiarMotivo(s, &err.motivos.esValoresVacios)
}
func (err *errorPaquete) motivoCamposValores() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque la cantidad de nombres de campos no coincide con la cantidad de valores recibidos"
	return err.cambiarMotivo(s, &err.motivos.esCamposValores)
}

// Error inesperado de la base de datos al ejecutar la sentencia.
func (err *errorPaquete) motivoEjecucionFallida() *errorPaquete {
	s := "Error al ejecutar la sentencia, se produjo un error inesperado"
	return err.cambiarMotivo(s, &err.motivos.esEjecucionFallida)
}

// Errores generados por la base de datos que han sido atrapados

func (err *errorPaquete) motivoTablaInexistente() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque el nombre de la tabla no existe en la base de datos"
	return err.cambiarMotivo(s, &err.motivos.esTablaInexistente)
}
func (err *errorPaquete) motivoCampoInexistente() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque hay al menos un nombre de campo que no existe en la tabla"
	return err.cambiarMotivo(s, &err.motivos.esCampoDeTablaInexistente)
}
func (err *errorPaquete) motivoEntradaDuplicada() *errorPaquete {
	s := "No es posible guardar los datos porque ya existe un campo que contiene el mismo valor que se ha recibido (entrada duplicada)"
	return err.cambiarMotivo(s, &err.motivos.esEntradaDuplicada)
}
func (err *errorPaquete) motivoTipoDeCampoIncorrecto() *errorPaquete {
	s := "No es posible guardar los datos porque existe al menos un campo de la tabla que está recibiendo un tipo de valor incorrecto"
	return err.cambiarMotivo(s, &err.motivos.esTipoDeCampoIncorrecto)
}
func (err *errorPaquete) motivoCampoFueraDeRango() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque existe al menos un valor recibido que supera el límite permitido por el campo de la tabla"
	return err.cambiarMotivo(s, &err.motivos.esCampoFueraDeRango)
}
func (err *errorPaquete) motivoObtencionDeRegistrosAfectados() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque se produjo un error al obtener la cantidad de registros afectados"
	return err.cambiarMotivo(s, &err.motivos.esObtencionDeRegistrosAfectados)
}
func (err *errorPaquete) motivoNingunRegistroAfectado() *errorPaquete {
	s := "La sentencia se ejecutó satisfactoriamente pero ningún registro de la tabla fue afectado. Los posibles motivos son: Elemento inexistente o no se ha cambiado ningún valor del registro"
	return err.cambiarMotivo(s, &err.motivos.esNingunRegistroAfectado)
}

func (err *errorPaquete) motivoObtencionDeID() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque se produjo un error al obtener el id insertado"
	return err.cambiarMotivo(s, &err.motivos.esObtencionDeID)
}

// Seleccionar

func (err *errorPaquete) motivoSeleccionarPunteroDeSlice() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque el objeto recibido no es un puntero de slice de estructura"
	return err.cambiarMotivo(s, &err.motivos.esSeleccionarPunteroDeSlice)
}
func (err *errorPaquete) motivoSeleccionarCamposSinRelacion() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque los campos de la estructura del objeto recibido no tienen asignados la relación de nombres con los campos de la tabla de la base de datos"
	return err.cambiarMotivo(s, &err.motivos.esSeleccionarCamposSinRelacion)
}
func (err *errorPaquete) motivoSeleccionarContieneEstrucutura() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque el objeto recibido contiene al menos una estructura"
	return err.cambiarMotivo(s, &err.motivos.esSeleccionarContieneEstructura)
}
func (err *errorPaquete) motivoSeleccionarTipoDeCampoIncorrecto() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque existe al menos un campo de la estructura que contiene un tipo erroneo (no se permiten punteros)"
	return err.cambiarMotivo(s, &err.motivos.esSeleccionarTipoDeCampoIncorrecto)
}
func (err *errorPaquete) motivoSeleccionarCamposFaltantes(camposFaltantes string) *errorPaquete {
	s := fmt.Sprintf("No es posible ejecutar la sentencia porque los campos obtenidos de la consulta no existen en su totalidad dentro de la estructura. Los campos faltantes son: %v", camposFaltantes)
	return err.cambiarMotivo(s, &err.motivos.esSeleccionarCamposFaltantes)
}
func (err *errorPaquete) motivoSeleccionarLecturaDeCampos() *errorPaquete {
	s := "No es posible ejecutar la sentencia porque se produjo un error al leer los campos de la consulta"
	return err.cambiarMotivo(s, &err.motivos.esSeleccionarLecturaDeCampos)
}
func (err *errorPaquete) motivoSeleccionarAsignacionDeCampos(mensaje string) *errorPaquete {
	s := fmt.Sprintf("No es posible ejecutar la sentencia porque se produjo un error al asignar los campos de la consulta a los campos de la estructura. %v", mensaje)
	return err.cambiarMotivo(s, &err.motivos.esSeleccionarAsignacionDeCampos)
}

// Transacciones

func (err *errorPaquete) motivoTxIniciar() *errorPaquete {
	s := "Error al intentar iniciar una transacción"
	return err.cambiarMotivo(s, &err.motivos.esTxIniciar)
}
func (err *errorPaquete) motivoTxConfirmar() *errorPaquete {
	s := "Error al intentar confirmar la transacción"
	return err.cambiarMotivo(s, &err.motivos.esTxConfirmar)
}
func (err *errorPaquete) motivoTxRevertir() *errorPaquete {
	s := "Error al intentar revertir la transacción"
	return err.cambiarMotivo(s, &err.motivos.esTxRevertir)
}

// -----------------------------------------------------------------------------
// Métodos públicos.

// Error implementa la interface error.
func (err *errorPaquete) Error() string {
	return strings.Join(err.mensajes, ". ")
}

// Origen devuelve el error de origen (error original).
func (err *errorPaquete) ObtenerOrigen() error {
	return err.origen
}

// Métodos públicos del error (motivos/comportamientos/causas de error).

func (err *errorPaquete) EsConexionAbrir() bool    { return err.motivos.esConexionAbrir }
func (err *errorPaquete) EsConexionCerrar() bool   { return err.motivos.esConexionCerrar }
func (err *errorPaquete) EsEjecucionFallida() bool { return err.motivos.esEjecucionFallida }
func (err *errorPaquete) EsEntradaDuplicada() bool { return err.motivos.esEntradaDuplicada }

func (err *errorPaquete) EsNombreDeTablaVacia() bool             { return err.motivos.esNombreDeTablaVacia }
func (err *errorPaquete) EsNombresDeCamposVacios() bool          { return err.motivos.esNombresDeCamposVacios }
func (err *errorPaquete) EsValoresVacios() bool                  { return err.motivos.esValoresVacios }
func (err *errorPaquete) EsCamposValoresDiferenteCantidad() bool { return err.motivos.esCamposValores }
func (err *errorPaquete) EsTablaInexistente() bool               { return err.motivos.esTablaInexistente }
func (err *errorPaquete) EsCampoDeTablaInexistente() bool {
	return err.motivos.esCampoDeTablaInexistente
}
func (err *errorPaquete) EsTipoDeCampoIncorrecto() bool { return err.motivos.esTipoDeCampoIncorrecto }
func (err *errorPaquete) EsObtencionDeRegistrosAfectados() bool {
	return err.motivos.esObtencionDeRegistrosAfectados
}
func (err *errorPaquete) EsNingunRegistroAfectado() bool { return err.motivos.esNingunRegistroAfectado }
func (err *errorPaquete) EsCampoFueraDeRango() bool      { return err.motivos.esCampoFueraDeRango }
func (err *errorPaquete) EsObtencionDeID() bool          { return err.motivos.esObtencionDeID }

// Seleccionar -----------------------------------------------------------------

func (err *errorPaquete) EsSeleccionarPunteroDeSlice() bool {
	return err.motivos.esSeleccionarPunteroDeSlice
}
func (err *errorPaquete) EsSeleccionarCamposSinRelacion() bool {
	return err.motivos.esSeleccionarCamposSinRelacion
}
func (err *errorPaquete) EsSeleccionarContieneEstructura() bool {
	return err.motivos.esSeleccionarContieneEstructura
}
func (err *errorPaquete) EsSeleccionarTipoDeCampoIncorrecto() bool {
	return err.motivos.esSeleccionarTipoDeCampoIncorrecto
}
func (err *errorPaquete) EsSeleccionarCamposFaltantes() bool {
	return err.motivos.esSeleccionarCamposFaltantes
}
func (err *errorPaquete) EsSeleccionarLecturaDeCampos() bool {
	return err.motivos.esSeleccionarLecturaDeCampos
}
func (err *errorPaquete) EsSeleccionarAsignacionDeCampos() bool {
	return err.motivos.esSeleccionarAsignacionDeCampos
}

// Transacciones ---------------------------------------------------------------

func (err *errorPaquete) EsTxIniciar() bool {
	return err.motivos.esTxIniciar
}
func (err *errorPaquete) EsTxConfirmar() bool {
	return err.motivos.esTxConfirmar
}
func (err *errorPaquete) EsTxRevertir() bool {
	return err.motivos.esTxRevertir
}

// -----------------------------------------------------------------------------

// errorNuevo crea un nuevo error del paquete.
func errorNuevo() *errorPaquete {
	return &errorPaquete{}
}

// EsError devuelve el error del paquete y un valor lógico que confirma el tipo.
// La única manera de obtener el error del paquete es a traves de esta función.
func EsError(err error) (*errorPaquete, bool) {
	ep, ok := err.(*errorPaquete)
	return ep, ok
}
