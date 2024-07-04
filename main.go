package main

import (
	"database/sql"
	"fmt"
	"strings"

	"log"
	"net/http"
	"regexp"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Password := ""
	Nombre := "sistema"
	conexion, err := sql.Open(Driver, Usuario+":"+Password+"@tcp(127.0.0.1)/"+Nombre)

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

	log.Println("Servidor Corriendo...")
	http.ListenAndServe(":8080", nil)

}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionDB()
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
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}
	//fmt.Println(arregloEmpleado)

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func Crear(w http.ResponseWriter, r *http.Request) {

	plantillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conexionEstablecida := conexionDB()
		isValid := validateNombre(nombre)
		if !isValid {
			fmt.Println("El usuario solo debe agregar letras en el campo de nombre")
			return
		}
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre, correo) VALUES(?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}

}
func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	//fmt.Println(idEmpleado)
	conexionEstablecida := conexionDB()
	borrarRegistros, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	borrarRegistros.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)

}
func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	//fmt.Println(idEmpleado)
	conexionEstablecida := conexionDB()

	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)

	empleado := Empleado{}
	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	//fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)
}
func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")
		conexionEstablecida := conexionDB()
		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		modificarRegistros.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", 301)
	}

}
func validateNombre(nombre string) bool {

	nombre = strings.TrimSpace(nombre)

	re := regexp.MustCompile(`^[A-Za-z]+$`)
	return re.MatchString(nombre)
}
