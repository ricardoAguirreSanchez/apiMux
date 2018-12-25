package googleDrive

import (
	"io/ioutil"
	"log"
	"strings"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"golang.org/x/net/context"
	"../utils"
)

/*
Este archivo tendra los dos metodos que usaran la API de GOOGLE DRIVE
PD: es necesario dejar los metodos con MAYUSCULA para que sean publicos
*/

type Documento struct{  //los atributos publicos
	Id string `json:"id"`
	Titulo string `json:"titulo"`
	Contenido string `json:"contenido"`
}

type ListDocument []DocumentoMeta

type DocumentoMeta struct{  //los atributos publicos
	Id string `json:"id"`
	Title string `json:"title"`
	MimeType string `json:"mimiType"`
	CreatedDate string `json:"createdDate"`
	WebContentLink string `json:"webContentLinkeatedDate"`
}

func init(){
}

// GET 
func List() ListDocument{
	log.Println("CONSULTO la API GOOGLE DRIVE para listar docs")
	var resultado ListDocument
	//----------------Busca las credenciales----------------------//
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
			log.Fatalf("No se pudo leer el archivo credentials.json : %v", err)
	}
	log.Println("Buscando credenciales")
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
			log.Fatalf("No se puede analizar el 'client secret file' para configurar: %v", err)
	}
	
	tokFile := "token.json"
	tok, err := utils.TokenFromFile(tokFile)
	client := config.Client(context.Background(), tok)
	//--------------------------------------------------------//

	//Busca el documento
	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
	}

	//Q("mimeType='application/pdf' and name contains 'myfile' and trashed=false")
	r, err := srv.Files.List().MaxResults(10).Q("trashed=false").Do()
	if err != nil {
			log.Fatalf("No se pudo recuperar los archivos: %v", err)
	}
	log.Println("Archivos:")
	if len(r.Items) == 0 {
		log.Println("No tiene archivos.")
	} else {
			for _, i := range r.Items {
				log.Printf("Id:%s Title:%s Mimetype:%s CreatedDate:%s WebContentLink:%s\n", i.Id, i.Title,i.MimeType,i.CreatedDate,i.WebContentLink)
				docu := DocumentoMeta{i.Id,i.Title,i.MimeType,i.CreatedDate,i.WebContentLink}
				resultado = append(resultado, docu) 
			}
	}
	return resultado
}

// GET 
func SerchInDocument(id string,word string) string{
	log.Println("CONSULTO la API GOOGLE DRIVE SI EL DOCUMENTO ID: " + id + " TIENE LA PALABRA: "+ word)

	//----------------Busca las credenciales----------------------//
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
			log.Fatalf("No se pudo leer el archivo credentials.json : %v", err)
	}
	log.Println("Buscando credenciales")
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
			log.Fatalf("No se puede analizar el 'client secret file' para configurar: %v", err)
	}
	
	tokFile := "token.json"
	tok, err := utils.TokenFromFile(tokFile)
	client := config.Client(context.Background(), tok)
	
	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
			return "No encontrado!"
	}
	log.Println("Drive recuperado")
	
	//-----------------Busca el documento-------------------//
	r, err := srv.Files.Get(id).Do()
	if err != nil {
			log.Fatalf("No se pudo recuperar el archivo buscado, Error : %v", err)
			return "No encontrado!"
	}

	//-----Descargamos y leemos el documento buscado---------------//
	round := client.Transport
	contenido,err := utils.DownloadAndReadFile(round , r)
	if err != nil {
		log.Printf("Error al descargar el archivo: %v\n", err)
	}else{
		log.Printf("Descarga y lectura realizado correctamente....")
	}

	//-----Controlamos si existe el word en el contenido
	log.Printf("Revisamos el contenido....")
	if strings.Contains(contenido, word){
		log.Println("Existe la palabra en el contenido!")
		return "Encontrado!"
	}else{
		log.Println("No existe la palabra en el contenido :(")
		return "No encontrado!"
	}
}

// POST (Obs:Crea un archivo .txt con el contenido mandado)
func CreateFile(documentoACrear Documento) Documento{
		
	log.Println("CONSULTO la API GOOGLE DRIVE PARA CREAR DOCUMENTO DE Titulo: " + documentoACrear.Titulo + " Y Contenido: " + documentoACrear.Contenido)
	
	//Busca las credenciales
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
			log.Fatalf("No se pudo leer el archivo credentials.json : %v", err)
	}
	log.Println("Buscando credenciales")
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
			log.Fatalf("No se puede analizar el 'client secret file' para configurar: %v", err)
	}

	tokFile := "token.json"
	tok, err := utils.TokenFromFile(tokFile)
	client := config.Client(context.Background(), tok)

	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
			return Documento{"DFEEWEFSEE34FF","default","default"}
	}
	log.Println("Drive recuperado")

	//Insertamos el doc.txt
	
	fileCreado,err := utils.InsertFileInDrive(srv,documentoACrear.Titulo,documentoACrear.Contenido,"","text/plain",documentoACrear.Titulo)

	if err != nil {
		log.Fatalf("No se pudo insertar el archivo en el drive del cliente: %v", err)
		return Documento{"DFEEWEFSEE34FF","default","default"}
	}
	
	documentoNuevo := Documento{fileCreado.Id,documentoACrear.Titulo,documentoACrear.Contenido}
	return documentoNuevo
}

