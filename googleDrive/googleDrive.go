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
	
	"bytes"
	"github.com/ledongthuc/pdf"
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
	tok, err := tokenFromFile(tokFile)
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
	contenido,err := downloadAndReadFile(round , r)
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
	tok, err := tokenFromFile(tokFile)
	client := config.Client(context.Background(), tok)

	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
			return Documento{"DFEEWEFSEE34FF","default","default"}
	}
	log.Println("Drive recuperado")

	//Insertamos el doc.txt
	
	fileCreado,err := insertFileInDrive(srv,documentoACrear.Titulo,documentoACrear.Contenido,"","text/plain",documentoACrear.Titulo)

	if err != nil {
		log.Fatalf("No se pudo insertar el archivo en el drive del cliente: %v", err)
		return Documento{"DFEEWEFSEE34FF","default","default"}
	}
	
	documentoNuevo := Documento{fileCreado.Id,documentoACrear.Titulo,documentoACrear.Contenido}
	return documentoNuevo
}

//------------------------------------Privados--------------------------------//

//Permite descargar y leer archivos .txt y .pdf
func downloadAndReadFile(t http.RoundTripper, f *drive.File) (string, error) {
	
	downloadUrl := f.DownloadUrl
	if downloadUrl == "" {
	  log.Printf("downloadUrl esta vacio!")
	  return "", nil
	}
	log.Printf("downloadUrl: %s",downloadUrl)
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		log.Printf("Error al ejecuar NewRequest: %v\n", err)
	  return "", err
	}
	resp, err := t.RoundTrip(req)
	// Si o si, luego del return o que explote, con defer ejecuta close siempre
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Error al ejecutar RoundTrip: %v\n", err)
	  return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error al ejecutar ReadAll: %v\n", err)
	  return "", err
	}

	if strings.Contains(f.MimeType,"text/plain"){
		log.Printf("El doc a analizar es un text/plain...")
		return string(body), nil
	}else if strings.Contains(f.MimeType,"application/pdf"){
		log.Printf("El doc a analizar es un application/pdf...")
		path := f.Title
		stri := string(body)
		
		createAndWriteFile(path,stri)
		
		content,err := readPdf(path)//hay que cambiar el binario a texto ya que es .pdf
		if err != nil{
			log.Println("Error al ejecutar readPdf: ",err)
		}
		log.Println("Documento pdf leido correctamente con readPdf")

		deleteFile(path)
		return content, nil
	}else{
		//es otro tipo de archivo, no lo podemos analizar
		log.Println("El doc a analizar no es ni un text ni pdf")
		return "", nil
	}
}

//libreria externa para poder leer archivos .pdf
func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
    b, err := r.GetPlainText()
    if err != nil {
        return "", err
    }
    buf.ReadFrom(b)
	return buf.String(), nil
}

//realiza todo el flujo para poder insertar un doc en drive
func insertFileInDrive(d *drive.Service, title string, contenido string,parentId string, mimeType string, filename string) (*drive.File, error) {
	
	//path := title + ".txt"

	//creo un .txt con ese contenido 
	createAndWriteFile(title,contenido)

	//inserto el .txt en drive
	r := subirInDrive(d, title, contenido,parentId, mimeType , title)

	//borro el .txt local
	deleteFile(title)
	
	return r,nil
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) { return }

	log.Println("==> se borro el archivo creado temporalmente")
}

func createAndWriteFile(path string,contenido string) {
	
	// detect if file exists
	var _, errFileCreado = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(errFileCreado) {
		var fileCreado, errFileCreado = os.Create(path)
		if isError(errFileCreado) { return }
		defer fileCreado.Close()
	}

	log.Println("==> se creo el archivo ", path)

	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) { return }
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contenido)
	if isError(err) { return }

	// save changes
	err = file.Sync()
	if isError(err) { return }

	log.Println("==> se escribio el contenido en el archivo creado")
}

//Subir al drive
func subirInDrive(d *drive.Service, title string, description string,parentId string, mimeType string, filename string) (*drive.File){
	m, err := os.Open(filename)
	if err != nil {
		log.Println("Error al abrir el archivo", filename)
		return nil
	}
	f := &drive.File{Title: title, Description: description, MimeType: mimeType}
  	if parentId != "" {
    	p := &drive.ParentReference{Id: parentId}
    	f.Parents = []*drive.ParentReference{p}
  	}
	r, err := d.Files.Insert(f).Media(m).Do()
  	if err != nil {
    	log.Printf("Error tratando de insertar el archivo: %v\n", err)
    	return nil
	}
	
	log.Printf("==> Se subio el archivo correctamente!")
	return r
}

func isError(err error) bool {
	if err != nil {
		log.Println(err.Error())
	}

	return (err != nil)
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