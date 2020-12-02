# Cambios del aplicativo

## [0.1.0] 2020-12-02
### Agregados
* Permite nombres de sentencias con caracter "-": no almacena sentencias generadas, por lo tanto; por cada invocación se genera una sentencia SQL nativa.

## [0.9.1] 2020-09-31
### Modificaciones
* Permitir nombres de sentencias con caracter "-" para que no almacene en en el mapa de sentencias generadas.

## [0.9.0] 2020-09-21
### Agregados
* Primera versión del paquete.

### Pendientes
* Las instrucciones select, intentar imitar scan: pasando cada valor de un puntero del slice; así no se utilizaría la reflexión.
* Utilizar contextos, saber para que se utilizan y si es aconsejable.
* Realizar select directo, sin generador de sentencia.
* Verificar que el slice que se va llenando en el select NO utilice append (debe estar inicializado con la cantidad filas obtenidas).