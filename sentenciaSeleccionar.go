package bdsql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"database/sql"
)

type seleccionar struct {
	bd *BD
	tx *sql.Tx

	tabla            string
	campos           []string
	condicion        string
	condicionValores []interface{}
	ordenadoPor      []string

	agruparPor        []string
	teniendoCondicion string
	teniendoValores   []interface{}

	limite int
	salto  int

	objeto interface{} // puntero de slice de objeto para el método Resultado().

	juntaInternaTabla       []string
	juntaInternaCondicion   []string
	juntaIzquierdaTabla     []string
	juntaIzquierdaCondicion []string
	juntaDerechaTabla       []string
	juntaDerechaCondicion   []string
	juntaExternaTabla       []string
	juntaExternaCondicion   []string

	senSQLExiste bool
	senSQLNombre string
	senSQL       string
}

// Tabla establece el nombre de la tabla a seleccionar.
func (o *seleccionar) Tabla(tabla string) *seleccionar {
	if o.senSQLExiste {
		return o
	}

	o.tabla = tabla
	return o
}

// Campos establece los nombres de campos a seleccionar.
func (o *seleccionar) Campos(campos ...string) *seleccionar {
	o.campos = campos
	return o
}

// Condicion implementa la cláusula 'where' de la sentencia 'select'.
func (o *seleccionar) Condicion(condicion string, valores ...interface{}) *seleccionar {
	if !o.senSQLExiste {
		o.condicion = condicion
	}
	o.condicionValores = valores

	return o
}

// OrdenarPor implementa a la cláusula 'order by' de la sentencia 'select'.
func (o *seleccionar) OrdenarPor(valores ...string) *seleccionar {
	if !o.senSQLExiste {
		o.ordenadoPor = valores
	}

	return o
}

// AgruparPor implementa a la cláusula 'group by' de la sentencia 'select'.
func (o *seleccionar) AgruparPor(valores ...string) *seleccionar {
	if !o.senSQLExiste {
		o.agruparPor = valores
	}
	return o
}

// Teniendo implementa la cláusula 'having' de la sentencia 'select'.
func (o *seleccionar) Teniendo(condicion string, valores ...interface{}) *seleccionar {
	if !o.senSQLExiste {
		o.teniendoCondicion = condicion
	}
	o.teniendoValores = valores

	return o
}

// Limitar implementa la cláusula 'limit' de la sentencia 'select'.
func (o *seleccionar) Limitar(limite int) *seleccionar {
	o.limite = limite

	return o
}

// Saltar implementa la cláusula 'offset' de la sentencia 'select'.
func (o *seleccionar) Saltar(salto int) *seleccionar {
	o.salto = salto

	return o
}

// Juntar implementa la cláusula 'inner join' de la sentencia 'select'.
func (o *seleccionar) JuntarCon(tabla, condicion string) *seleccionar {
	if !o.senSQLExiste {
		o.juntaInternaTabla = append(o.juntaInternaTabla, tabla)
		o.juntaInternaCondicion = append(o.juntaInternaCondicion, condicion)
	}
	return o
}

// JuntarIzquierda implementa la cláusula 'left join' de la sentencia 'select'.
func (o *seleccionar) JuntarIzquierda(tabla, condicion string) *seleccionar {
	if !o.senSQLExiste {
		o.juntaIzquierdaTabla = append(o.juntaIzquierdaTabla, tabla)
		o.juntaIzquierdaCondicion = append(o.juntaIzquierdaCondicion, condicion)
	}
	return o
}

// JuntarDerecha implementa la cláusula 'right join' de la sentencia 'select'.
func (o *seleccionar) JuntarDerecha(tabla, condicion string) *seleccionar {
	if !o.senSQLExiste {
		o.juntaDerechaTabla = append(o.juntaDerechaTabla, tabla)
		o.juntaDerechaCondicion = append(o.juntaDerechaCondicion, condicion)
	}
	return o
}

// JuntarExterior implementa la cláusula 'outer join' de la sentencia 'select'.
func (o *seleccionar) JuntarExterior(tabla, condicion string) *seleccionar {
	if !o.senSQLExiste {
		o.juntaExternaTabla = append(o.juntaExternaTabla, tabla)
		o.juntaExternaCondicion = append(o.juntaExternaCondicion, condicion)
	}
	return o
}

// Resultado recibe el objeto donde se almacena el resultado de la consulta.
// Debe ser un puntero de slice de una estructura.
//	Ejemplo:
//	var elementos = []struct {
//		ID            int64  `bdsql:"id"`
//		Nombre        string `bdsql:"nombre"`
//		EsActivo      bool   `bdsql:"es_activo"`
//		Observaciones string `bdsql:"observaciones"`
//	}{}
//	var cant, err = bd.Seleccionar("consulta").
//					Tabla("elementos").
//					Campos("*").
//					Resultado(&elementos).
//					Ejecutar()
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	fmt.Println(cant, elementos)
//
func (o *seleccionar) Resultado(objeto interface{}) *seleccionar {
	o.objeto = objeto
	return o
}

// SQL devuelve la sentencia SQL.
func (o *seleccionar) SQL() (string, error) {
	return o.generarSQL()
}

// Ejecutar ejecuta la sentencia SQL.
func (o *seleccionar) Ejecutar() (int, error) {
	var sentencia, err = o.generarSQL()
	if err != nil {
		return 0, err
	}

	// validar que el objeto de resultado, sea un puntero de slice de estructura
	if ok := reflect.TypeOf(o.objeto).Kind() == reflect.Ptr && reflect.TypeOf(o.objeto).Elem().Kind() == reflect.Slice && reflect.TypeOf(o.objeto).Elem().Elem().Kind() == reflect.Struct; !ok {
		return 0, errorNuevo().asignarMotivoSeleccionarPunteroDeSlice()
	}

	// slice de parámetros de la sentencia sql a ejecutar
	var parametros []interface{}
	// where
	if o.condicion != "" {
		parametros = append(parametros, o.condicionValores...)
	}
	// having
	if o.teniendoCondicion != "" {
		parametros = append(parametros, o.teniendoValores...)
	}

	// var err error
	var filas *sql.Rows
	if o.tx == nil {
		filas, err = o.bd.db.Query(sentencia, parametros...)
	} else {
		filas, err = o.tx.Query(sentencia, parametros...)
	}
	if err != nil {
		return 0, resolverErrorMysql(err)
	}
	defer filas.Close()

	cant, err := asignarAObjeto(filas, o.objeto)
	if err != nil {
		return 0, err
	}

	return cant, nil
}

func (o *seleccionar) generarSQL() (string, error) {
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
	if len(err.mensajes) != 0 {
		return "", err
	}
	// sentencia
	sentencia := fmt.Sprintf("select %v from %v", strings.Join(o.campos, ", "), o.tabla)
	// inner join
	for i := 0; i < len(o.juntaInternaTabla); i++ {
		sentencia += fmt.Sprintf(" inner join %v on %v", o.juntaInternaTabla[i], o.juntaInternaCondicion[i])
	}
	// left join
	for i := 0; i < len(o.juntaIzquierdaTabla); i++ {
		sentencia += fmt.Sprintf(" left join %v on %v", o.juntaIzquierdaTabla[i], o.juntaIzquierdaCondicion[i])
	}
	// right join
	for i := 0; i < len(o.juntaDerechaTabla); i++ {
		sentencia += fmt.Sprintf(" right join %v on %v", o.juntaDerechaTabla[i], o.juntaDerechaCondicion[i])
	}
	// outer join
	for i := 0; i < len(o.juntaExternaTabla); i++ {
		sentencia += fmt.Sprintf(" outer join %v on %v", o.juntaExternaTabla[i], o.juntaExternaCondicion[i])
	}
	// where
	if o.condicion != "" {
		sentencia += fmt.Sprintf(" where %v", o.condicion)
	}
	// group by
	if len(o.agruparPor) > 0 {
		sentencia += fmt.Sprintf(" group by %v", strings.Join(o.agruparPor, ", "))
	}
	// having
	if o.teniendoCondicion != "" {
		sentencia += fmt.Sprintf(" having %v", o.teniendoCondicion)
	}
	// order by
	if len(o.ordenadoPor) > 0 {
		sentencia += fmt.Sprintf(" order by %v", strings.Join(o.ordenadoPor, ", "))
	}
	// limit
	if o.limite > 0 {
		sentencia += fmt.Sprintf(" limit %v", o.limite)
	}
	// offset
	if o.salto > 0 {
		if o.limite == 0 {
			sentencia += fmt.Sprintf(" limit %v", 2100000000)
		}
		sentencia += fmt.Sprintf(" offset %v", o.salto)
	}

	return fmt.Sprintf("%v;", sentencia), nil
}

func asignarAObjeto(filas *sql.Rows, objeto interface{}) (int, error) {
	// nombres de campos del resultado obtenido de la base de datos.
	camposFila, err := filas.Columns()
	if err != nil {
		return 0, resolverErrorMysql(err)
	}

	// obtener la lista de campos de la estructura que podrían ser asignados
	// desde la consulta obtenida.
	estructura := reflect.TypeOf(objeto).Elem().Elem()

	// mapa de estructura del objeto:
	// la clave es el valor de la etiqueta.
	// el valor es el nombre del campo de la estructura.
	camposEstructura := make(map[string]string)

	var encontrados int
	for i := 0; i < estructura.NumField(); i++ {
		campo := estructura.Field(i)
		campoTabla := strings.Trim(campo.Tag.Get("bdsql"), " ")

		switch campoTabla {
		case "-":
			// no asignar el valor
			continue
		case "":
			// en caso que el campo de la estructura no tenga asignado el
			// tag "bdsql", se asume que el nombre del campo de la la consulta
			// SQL, es el mismo (en minúsculas) que el nombre de campo de la
			// estructura.
			camposEstructura[strings.ToLower(campo.Name)] = campo.Name
		default:
			// en caso que el campo de la estructura tenga asignado el
			// tag "bdsql", se toma el tag para obtener el dato de la
			// consulta SQL.
			camposEstructura[campoTabla] = campo.Name
		}

		encontrados++
	}
	// podría pasar (raramente) que todos los campos de la estructura contengan
	// el tag `bdsql:"-"`. En ese caso, la estrucutra no permite que ninguno de
	// sus campos sean asignables por los valores recibidos de la consulta SQL.
	if encontrados == 0 {
		return 0, errorNuevo().asignarMotivoSeleccionarCamposSinRelacion()
	}

	// mapa de funciones:
	// se crea un slice el cuál mantiene la función que completa cada campo
	// de la estructura.
	funciones := make(map[string]func(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error), estructura.NumField())

	// verificar qe todos los campos de la consulta obtenida, puedan ser
	// ingresados en el objeto recibido.
	var camposFaltantes []string
	estructuraNueva := reflect.New(estructura).Elem()
	for _, campoFila := range camposFila {
		campoEstructura, ok := camposEstructura[campoFila]
		if !ok {
			camposFaltantes = append(camposFaltantes, campoFila)
			continue
		}

		// llenar el slice de funciones según cada tipo de campo de la estructura.
		switch estructuraNueva.FieldByName(campoEstructura).Kind() {
		case reflect.Int:
			funciones[campoEstructura] = valori
		case reflect.Int8:
			funciones[campoEstructura] = valori8
		case reflect.Int16:
			funciones[campoEstructura] = valori16
		case reflect.Int32:
			funciones[campoEstructura] = valori32
		case reflect.Int64:
			funciones[campoEstructura] = valori64
		case reflect.Uint:
			funciones[campoEstructura] = valorui
		case reflect.Uint8:
			funciones[campoEstructura] = valorui8
		case reflect.Uint16:
			funciones[campoEstructura] = valorui16
		case reflect.Uint32:
			funciones[campoEstructura] = valorui32
		case reflect.Uint64:
			funciones[campoEstructura] = valorui64
		case reflect.Float32:
			funciones[campoEstructura] = valorf32
		case reflect.Float64:
			funciones[campoEstructura] = valorf64
		case reflect.Bool:
			funciones[campoEstructura] = valorbool
		case reflect.String:
			funciones[campoEstructura] = valorstring
		case reflect.Struct:
			if estructuraNueva.FieldByName(campoEstructura).Type().String() == "time.Time" {
				funciones[campoEstructura] = valorFecha
			} else {
				return 0, errorNuevo().asignarMotivoSeleccionarContieneEstructura()
			}
		default:
			return 0, errorNuevo().asignarMotivoSeleccionarTipoDeCampoIncorrecto()
		}
	}
	if len(camposFaltantes) > 0 {
		return 0, errorNuevo().asignarMotivoSeleccionarCamposFaltantes(strings.Join(camposFaltantes, ","))
	}

	// obtener los valores de los campos de la base de datos.
	var valores = make([]interface{}, len(camposFila))
	for i := range valores {
		var ii interface{}
		valores[i] = &ii
	}

	// sabiendo que el objeto es un puntero de slice de una estructura se
	// asigna un elemento para ir incorporando las estructuras.
	sliceDeObjetos := reflect.ValueOf(objeto).Elem()

	// recorrer todas las filas e insertarlas en el objeto.
	var cant int
	for filas.Next() {
		cant++

		if err := filas.Scan(valores...); err != nil {
			return 0, errorNuevo().asignarOrigen(err).asignarMotivoSeleccionarLecturaDeCampos()
		}

		for i, campoFila := range camposFila {
			var valorCrudo = *(valores[i].(*interface{}))
			var tipoCrudo = reflect.TypeOf(valorCrudo)

			campoEstructura, ok := camposEstructura[campoFila]
			if !ok {
				return 0, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos(fmt.Sprintf("Inexistente: %v, Campo: %v", err, campoFila))
			}

			valor, err := funciones[campoEstructura](valorCrudo, tipoCrudo)
			if err != nil {
				return 0, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos(fmt.Sprintf("Error: %v, Campo: %v", err, campoFila))
			}
			estructuraNueva.FieldByName(campoEstructura).Set(valor)
		}

		sliceDeObjetos.Set(reflect.Append(sliceDeObjetos, estructuraNueva))
	}

	return cant, nil
}

// ---- Funciones de asignación de campos de la estructura ---------------------

func valori(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v int

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.Atoi(string(s))
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo int")
			}
			v = vTemp
		} else {
			vTemp, ok := valorCrudo.(int64)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo int64")
			}
			v = int(vTemp)
		}
	}

	return reflect.ValueOf(v), nil
}

func valori8(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v int8

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseInt(string(s), 10, 8)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo int8")
			}
			v = int8(vTemp)
		} else {
			vTemp, ok := valorCrudo.(int64)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo int64")
			}
			v = int8(vTemp)
		}
	}

	return reflect.ValueOf(v), nil
}

func valori16(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v int16

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseInt(string(s), 10, 16)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo int16")
			}
			v = int16(vTemp)
		} else {
			vTemp, ok := valorCrudo.(int64)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo int64")
			}
			v = int16(vTemp)
		}
	}

	return reflect.ValueOf(v), nil
}

func valori32(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v int32

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseInt(string(s), 10, 32)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo int32")
			}
			v = int32(vTemp)
		} else {
			vTemp, ok := valorCrudo.(int64)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo int64")
			}
			v = int32(vTemp)
		}
	}

	return reflect.ValueOf(v), nil
}

func valori64(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v int64

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseInt(string(s), 10, 64)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo int64")
			}
			v = vTemp
		} else {
			vTemp, ok := valorCrudo.(int64)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo int64")
			}
			v = vTemp
		}
	}

	return reflect.ValueOf(v), nil
}

func valorui(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v uint

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseUint(string(s), 10, 64)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo uint")
			}
			v = uint(vTemp)
		} else {
			vTemp, ok := valorCrudo.(uint)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo uint")
			}
			v = uint(vTemp)
		}
	}

	return reflect.ValueOf(v), nil
}

func valorui8(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v uint8

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseUint(string(s), 10, 8)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo uint8")
			}
			v = uint8(vTemp)
		} else {
			vTemp, ok := valorCrudo.(uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo uint8")
			}
			v = vTemp
		}
	}

	return reflect.ValueOf(v), nil
}

func valorui16(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v uint16

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseUint(string(s), 10, 16)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo uint16")
			}
			v = uint16(vTemp)
		} else {
			vTemp, ok := valorCrudo.(uint16)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo uint16")
			}
			v = vTemp
		}
	}

	return reflect.ValueOf(v), nil
}

func valorui32(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v uint32

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseUint(string(s), 10, 32)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo uint32")
			}
			v = uint32(vTemp)
		} else {
			vTemp, ok := valorCrudo.(uint32)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo uint32")
			}
			v = vTemp
		}
	}

	return reflect.ValueOf(v), nil
}

func valorui64(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v uint64

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseUint(string(s), 10, 64)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo uint64")
			}
			v = vTemp
		} else {
			vTemp, ok := valorCrudo.(uint64)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo uint64")
			}
			v = vTemp
		}
	}

	return reflect.ValueOf(v), nil
}

func valorf32(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v float32

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseFloat(string(s), 32)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo float32")
			}
			v = float32(vTemp)
		} else {
			vTemp, ok := valorCrudo.(float32)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo float32")
			}
			v = float32(vTemp)
		}
	}

	return reflect.ValueOf(v), nil
}

func valorf64(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v float64

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseFloat(string(s), 64)
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo float64")
			}
			v = vTemp
		} else {
			vTemp, ok := valorCrudo.(float64)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es del tipo float64")
			}
			v = vTemp
		}
	}

	return reflect.ValueOf(v), nil
}

func valorbool(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v bool

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			vTemp, err := strconv.ParseBool(string(s))
			if err != nil {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo bool")
			}
			v = vTemp
		} else {
			switch valorCrudo.(type) {
			case bool:
				v = valorCrudo.(bool)
			case int64:
				v = valorCrudo == int64(1)
			default:
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible asignar a tipo bool")
			}
		}
	}

	return reflect.ValueOf(v), nil
}

func valorstring(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v string

	if tipoCrudo != nil {
		if tipoCrudo.String() == "[]uint8" {
			s, ok := valorCrudo.([]uint8)
			if !ok {
				return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es de tipo uint8[]")
			}
			v = fmt.Sprintf("%s", s)
		} else {
			v = fmt.Sprintf("%s", valorCrudo)
		}
	}

	return reflect.ValueOf(v), nil
}

func valorFecha(valorCrudo interface{}, tipoCrudo reflect.Type) (reflect.Value, error) {
	var vacio reflect.Value
	var v time.Time = time.Time{}

	if valorCrudo != nil {
		fh, ok := valorCrudo.(time.Time)
		if !ok {
			return vacio, errorNuevo().asignarMotivoSeleccionarAsignacionDeCampos("No es posible convertir al tipo time.Time")
		}
		v = fh
	}

	return reflect.ValueOf(v), nil
}
