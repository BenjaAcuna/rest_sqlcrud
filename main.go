package main

import (
	"database/sql"
	"fmt"

	//"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {

	Driver := "mysql"
	Usuario := "root"
	Contrasenia := "Holabb13*"
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {

	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	fmt.Println("Servidor Corriendo...")
	http.ListenAndServe(":8080", nil)

}

func Borrar(a http.ResponseWriter, b *http.Request) {

	idEmpleado := b.URL.Query().Get("ID")
	fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)

	http.Redirect(a, b, "/", 301)

}

type Empleado struct {
	ID     int
	Nombre string
	Correo string
}

func Inicio(a http.ResponseWriter, b *http.Request) {

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}
	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {

		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.ID = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)

	}
	//fmt.Println(arregloEmpleado)

	plantillas.ExecuteTemplate(a, "inicio", arregloEmpleado)

}
func Editar(a http.ResponseWriter, b *http.Request) {
	idEmpleado := b.URL.Query().Get("ID")
	fmt.Println(idEmpleado)

	conexionEstablecida := conexionBD()
	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.ID = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}

	fmt.Println(empleado)
	plantillas.ExecuteTemplate(a, "editar", empleado)
}

func Crear(a http.ResponseWriter, b *http.Request) {

	plantillas.ExecuteTemplate(a, "crear", nil)

}
func Insertar(a http.ResponseWriter, b *http.Request) {
	if b.Method == "POST" {

		Nombre := b.FormValue("Nombre")
		Correo := b.FormValue("Correo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre,correo) VALUES(?,?)")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(Nombre, Correo)

		http.Redirect(a, b, "/", 301)
	}

}
func Actualizar(a http.ResponseWriter, b *http.Request) {
	if b.Method == "POST" {
		ID := b.FormValue("ID")
		Nombre := b.FormValue("Nombre")
		Correo := b.FormValue("Correo")

		conexionEstablecida := conexionBD()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET Nombre=?,Correo=?,ID=?")

		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(Nombre, Correo, ID)

		http.Redirect(a, b, "/", 301)
	}

}
