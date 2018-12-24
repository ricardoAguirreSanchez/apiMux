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
	"net/http"
	"os"
		
)

/*
Este archivo tendra los dos metodos que usaran la API de GOOGLE DRIVE
PD: es necesario dejar los metodos con MAYUSCULA para que sean publicos
*/

type Documento struct{  //los atributos publicos
	Id string `json:"id"`
	Titulo string `json:"titulo"`
	Descripcion string `json:"descripcion"`
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
				log.Printf("Id:%s Title:%s Mimetype:%s CreatedDate:%s WebContentLink:%s\n", i.Id, i.Title,i.MimeType,i.CreatedDate,i.WebContentLink)
				docu := DocumentoMeta{i.Id,i.Title,i.MimeType,i.CreatedDate,i.WebContentLink}
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
	contenido,err := DownloadAndReadFile(round , r)
	if err != nil {
		log.Printf("Error al descargar el archivo: %v\n", err)
	}else{
		log.Printf("Descargado y lectura realizado correctamente....")
	}

	//-----Controlamos si existe el word en el contenido
	if strings.Contains(contenido, word){
		log.Println("Existe la palabra en el contenido!")
		return "Encontrado!"
	}else{
		log.Println("No existe la palabra en el contenido :(")
		return "No encontrado!"
	}
}


//Permite descargar y leer archivos .txt y .pdf
func DownloadAndReadFile(t http.RoundTripper, f *drive.File) (string, error) {
	
	downloadUrl := f.DownloadUrl
	if downloadUrl == "" {
	  // If there is no downloadUrl, there is no body
	  log.Printf("An error occurred: File is not downloadable")
	  return "", nil
	}
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
	  return "", err
	}
	resp, err := t.RoundTrip(req)
	// Si o si, luego del return o que explote, con defer ejecuta close siempre
	defer resp.Body.Close()
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
	  return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)
	  return "", err
	}

	if strings.Contains(f.MimeType,"text/plain"){
		return string(body), nil
	}else{
		//hay que cambiar el binario a texto ya que es .pdf
		return "", nil
	}
}


//https://gist.github.com/atotto/86fa30668473b41eeac7d750e5ad5f5c
//https://stackoverflow.com/questions/46334646/google-drive-api-v3-create-and-upload-file
//Crea un archivo .txt con el contenido mandabo
func CreateFile(documentoACrear Documento) Documento{
		
	log.Println("CONSULTO la API GOOGLE DRIVE PARA CREAR DOCUMENTO DE Titulo: " + documentoACrear.Titulo + " Y Contenido: " + documentoACrear.Descripcion)
	
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

	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
			return Documento{"DFEEWEFSEE34FF","default","default"}
	}
	log.Println("Drive recuperado")

	//----------------------Insertamos el doc.txt------------------------------------//
	
	fileCreado,err := InsertFile(srv,documentoACrear.Titulo,documentoACrear.Descripcion,"","text/plain",documentoACrear.Titulo)

	if err != nil {
		log.Fatalf("No se pudo insertar el archivo en el drive del cliente: %v", err)
		return Documento{"DFEEWEFSEE34FF","default","default"}
	}
	
	documentoNuevo := Documento{fileCreado.Id,documentoACrear.Titulo,documentoACrear.Descripcion}
	return documentoNuevo
}

func InsertFile(d *drive.Service, title string, description string,parentId string, mimeType string, filename string) (*drive.File, error) {
	//creo un .txt con ese contenido 


	//inserto el .txt en drive

	
	//borro el .txt local

	//   m, err := os.Open(filename)
//   if err != nil {
//     log.Printf("An error occurred: %v\n", err)
//     return nil, err
//   }
	f := &drive.File{Title: title, Description: description, MimeType: mimeType}
  	if parentId != "" {
    	p := &drive.ParentReference{Id: parentId}
    	f.Parents = []*drive.ParentReference{p}
  	}
//   r, err := d.Files.Insert(f).Media(m).Do()
  	r, err := d.Files.Insert(f).Do()
  	if err != nil {
    	log.Printf("Error tratando de insertar el archivo: %v\n", err)
    	return nil, err
  	}
  	return r, nil
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