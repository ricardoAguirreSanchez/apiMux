package main

import (
	"github.com/gorilla/mux" //go get -v -u github.com/gorilla/mux
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"./googleDrive"
	"./autenticador"
	"os"

    // "golang.org/x/net/context"
    // "golang.org/x/oauth2"
    // "golang.org/x/oauth2/google"
    // "google.golang.org/api/drive/v3"
)

// INDEX
func IndexHandler(w http.ResponseWriter , r *http.Request){
	
	estado := autenticador.GetEstadoAutenticacion()
	if(estado == "AUTENTICADO"){//significa que me autentique bien
		//le aviso que puede usar la api sin problema
		fmt.Fprintf(w,"Podes usar la api sin problema")
	}else{
		//le digo que tiene que autenticarse
		fmt.Fprintf(w, "<a href= %s>Logeate aqui</a>",estado)
	}
}

// OAUTH
func OauthHandler(w http.ResponseWriter , r *http.Request){
	codes, ok := r.URL.Query()["code"]
	if !ok || len(codes[0]) < 1 {
		log.Fatalf("No tenes el parametro code en la url.")
		fmt.Fprintf(w,"Falta el parametro code.")
	}else{
		code := codes[0]
		resultado := autenticador.Autenticar(code)
		if resultado == "OK"{
			fmt.Fprintf(w,"Logeo exitoso, ahora podes usar los endpoints") 
		}else{
			fmt.Fprintf(w,"Porfavor vaya al localhost:8080/ para que se autentique")
		}
	}
}


// GET 
func GetListHandler(w http.ResponseWriter , r *http.Request){
	estado := autenticador.GetEstadoAutenticacion()
	if estado != "AUTENTICADO"{
		//Significa que tengo que autenticarme
		log.Println("No estas autenticado, tenes que logearte")
		fmt.Fprintf(w,"Porfavor vaya al localhost:8080/ para que se autentique")
	}else{
		log.Println("Estas logeado!")
		log.Println("[GET] - Solicitando servicio list")
	
		resultado := googleDrive.List()
		json.NewEncoder(w).Encode(resultado)
	}
}



// GET 
func GetSearchInDocHandler(w http.ResponseWriter , r *http.Request){
	estado := autenticador.GetEstadoAutenticacion()
	if estado != "AUTENTICADO"{
		//Significa que tengo que autenticarme
		log.Println("No estas autenticado, tenes que logearte")
		fmt.Fprintf(w,"Porfavor vaya al localhost:8080/ para que se autentique")
	}else{
		log.Println("Estas logeado!")
		log.Println("[GET] - Solicitando servicio searchInDoc")
	
		//Para controlar la variable id
		vars := mux.Vars(r)
		id := vars["id"]
	
		//Para controlar el parametro word
		words, ok := r.URL.Query()["word"]
		if !ok || len(words[0]) < 1 {
			log.Fatalf("No tenes el parametro word en la url.")
			fmt.Fprintf(w,"Falta el parametro word.")
		}else{
			word := words[0]
			resultado := googleDrive.SerchInDocument(id,word)
			fmt.Fprintf(w,resultado)
		}
	}
}

// POST
func PostCreatFileHandler(w http.ResponseWriter , r *http.Request){
	estado := autenticador.GetEstadoAutenticacion()
	if estado != "AUTENTICADO"{
		//Significa que tengo que autenticarme
		fmt.Fprintf(w,"Porfavor vaya al localhost:8080/ para que se autentique")
	}else{
		log.Println("[POST] - Solicitando servicio file")
		var documento googleDrive.Documento
		
		//decodificamos el json recibio a un objeto documento
		error := json.NewDecoder(r.Body).Decode(&documento)
		if error != nil {
			log.Fatalf("[POST] - Error decodificando json en Solicitando servicio file")
			fmt.Fprintf(w,"Error al tratar de decodificar el json")
		}else{
			resultado := googleDrive.CreateFile(documento)
			json.NewEncoder(w).Encode(resultado)
		}
	}
	
	
}

func main() {
	//Borramos token.json si es que existe, asi limpiamos los token que allan
	var err = os.Remove("token.json")
	if err != nil {
        fmt.Println("No existe token.json")
    }else{
		fmt.Println("token.json borrado")
	}
	

	//Creamos un enroutador
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/oauth", OauthHandler).Methods("GET")
	r.HandleFunc("/search-in-doc/{id}", GetSearchInDocHandler).Methods("GET")
	r.HandleFunc("/file", PostCreatFileHandler).Methods("POST")
	r.HandleFunc("/list", GetListHandler).Methods("GET")

	//Podemos crear nuestro servidor a mano, asi lo podemos customizar mejor
	server := &http.Server{
			Addr:			":8080",
			Handler:		r,
			ReadTimeout:	10 * time.Second,
			WriteTimeout:	10 * time.Second,
			MaxHeaderBytes:	1 << 20,
	}
	log.Println("Escuchando....")
	server.ListenAndServe()
}


