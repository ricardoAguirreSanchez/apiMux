package googleDrive

import (
	"io/ioutil"
	"log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"strings"
	"encoding/json"
	"os"
		
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
}

func init(){
}

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
	tok, err := tokenFromFile(tokFile)
	client := config.Client(context.Background(), tok)
	//--------------------------------------------------------//

	//Busca el documento
	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
	}

	//srv.Files *FilesService -> r *FilesListCall
	r, err := srv.Files.List().MaxResults(10).Do()
	if err != nil {
			log.Fatalf("No se pudo recuperar los archivos: %v", err)
	}
	log.Println("Archivos:")
	if len(r.Items) == 0 {
		log.Println("No tiene archivos.")
	} else {
			for _, i := range r.Items {
				log.Printf("Id:%s Title:%s Mimetype:%s CreatedDate:%s\n", i.Id, i.Title,i.MimeType,i.CreatedDate)

				// docu := new(DocumentoMeta)
				// docu.Id = i.Id
				// docu.Description = i.Description
				// docu.MimeType = i.MimeType

				docu := DocumentoMeta{i.Id,i.Title,i.MimeType,i.CreatedDate}


				resultado = append(resultado, docu) 
			}
	}
	return resultado
}

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
	tok, err := tokenFromFile(tokFile)
	client := config.Client(context.Background(), tok)
	//--------------------------------------------------------//

	//Busca el documento
	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
			return "No encontrado!"
	}
	log.Println("Drive recuperado")
	
	r, err := srv.Files.Get(id).Do()
	if err != nil {
			log.Fatalf("No se pudo recuperar el archivo error : %v", err)
			return "No encontrado!"
	}

	//por ahora el word lo busca en el mimetype
	if strings.Contains(r.Title, word){
		log.Println("Existe la palabra en el titulo!")
		return "Encontrado!"
	}else{
		log.Println("No existe la palabra en el titulo :(")
		return "No encontrado!"
	}
	
	//--parece que para pedir el archivo hay que mandar un request


}

//https://gist.github.com/atotto/86fa30668473b41eeac7d750e5ad5f5c
//https://stackoverflow.com/questions/46334646/google-drive-api-v3-create-and-upload-file
func CreateFile(documentoACrear Documento) Documento{
		
	log.Println("CONSULTO la API GOOGLE DRIVE PARA CREAR DOCUMENTO DE Titulo: " + documentoACrear.Titulo + " Y Contenido: " + documentoACrear.Contenido)
	
	//----------------Busca las credenciales----------------------//
	// b, err := ioutil.ReadFile("credentials.json")
	// if err != nil {
	// 		log.Fatalf("No se pudo leer el archivo credentials.json : %v", err)
	// }
	// log.Println("Buscando credenciales")
	// config, err := google.ConfigFromJSON(b, drive.DriveScope)
	// if err != nil {
	// 		log.Fatalf("No se puede analizar el 'client secret file' para configurar: %v", err)
	// }

	// tokFile := "token.json"
	// tok, err := tokenFromFile(tokFile)
	// client := config.Client(context.Background(), tok)
	//--------------------------------------------------------//

	documentoNuevo := Documento{"DFEEWEFSEE34FF",documentoACrear.Titulo,documentoACrear.Contenido}
	
	// srv, err := drive.New(client)
	// if err != nil {
	// 		log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
	// 		return documentoNuevo
	// }
	// log.Println("Drive recuperado")

	// //1.- generamos un Id
	// genereateId, err := srv.Files.GenerateIds().Do()
	// if err != nil {
	// 		log.Fatalf("No se pudo generar el id %v",err)
	// 		return documentoNuevo
	// }
	// //2.- armamos el *file
	// id := genereateId.Ids[0]
	// log.Println("Generamos el id: %s",id)
	// //3.- guardamos el *file

	return documentoNuevo
}



// Busca el token del archivo
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
			return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}