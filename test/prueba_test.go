package test

//ACLARACIONES:
//El test se prueba parandose en el package y ejecutando "go test -v"
//Para poder realizar correctamente los test es que se agrego el archivo de credenciales y el archivoTest.pdf

import (
	"testing"
	"strings"
	"../utils"
	"../autenticador"
)

//Test sobre una autenticacion con un codigo ficticio, es correcto que de error
func Test01(t * testing.T){
	resultado := autenticador.Autenticar("111111")

	if resultado == "ERROR"{
		t.Log("[Test - 01] Respuesta esperada - OK ")
	}else{
		t.Log("[Test - 01] Respuesta diferente de la esperada - Posible error")
		t.Fail()
	}
}

//Test sobre la lectura del pdf, es correcto que si encuentre la palabra en el pdf
func Test02(t * testing.T){

	var nombreFile = "archivoTest.pdf"

	//buscamos la palabra "hace" dentro de "archivoTest.pdf" que tiene "hola que hace"
	contenido,err := utils.ReadPdf(nombreFile)

	if err != nil{
		t.Log("Error al quere leer el pdf")
		t.Fail()
	}
	//controlamos
	if strings.Contains(contenido,"hace"){
		t.Log("[Test - 02] Respuesta esperada - OK ")
	}else{
		t.Log("[Test - 02] Respuesta diferente de la esperada - Posible error")
		t.Fail()
	}

}


