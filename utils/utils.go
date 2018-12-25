package utils

//Este paquete se hizo para poder centralizar algunas funciones mas generales y tener el codigo más ordenado

import (
	"golang.org/x/oauth2"
	"strings"
	"encoding/json"
	"net/http"
	"os"
	"io/ioutil"
	"log"
	"google.golang.org/api/drive/v2"
	"bytes"
	"github.com/ledongthuc/pdf"


)
//Permite descargar y leer archivos .txt y .pdf
func DownloadAndReadFile(t http.RoundTripper, f *drive.File) (string, error) {
	
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
		
		CreateAndWriteFile(path,stri)
		
		content,err := ReadPdf(path)//hay que cambiar el binario a texto ya que es .pdf
		if err != nil{
			log.Println("Error al ejecutar readPdf: ",err)
		}
		log.Println("Documento pdf leido correctamente con readPdf")

		DeleteFile(path)
		return content, nil
	}else{
		//es otro tipo de archivo, no lo podemos analizar
		log.Println("El doc a analizar no es ni un text ni pdf")
		return "", nil
	}
}

//libreria externa para poder leer archivos .pdf
func ReadPdf(path string) (string, error) {
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
func InsertFileInDrive(d *drive.Service, title string, contenido string,parentId string, mimeType string, filename string) (*drive.File, error) {
	
	//path := title + ".txt"

	//creo un .txt con ese contenido 
	CreateAndWriteFile(title,contenido)

	//inserto el .txt en drive
	r := SubirInDrive(d, title, contenido,parentId, mimeType , title)

	//borro el .txt local
	DeleteFile(title)
	
	return r,nil
}

func DeleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if IsError(err) { return }

	log.Println("==> se borro el archivo creado temporalmente")
}

func CreateAndWriteFile(path string,contenido string) {
	
	// detect if file exists
	var _, errFileCreado = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(errFileCreado) {
		var fileCreado, errFileCreado = os.Create(path)
		if IsError(errFileCreado) { return }
		defer fileCreado.Close()
	}

	log.Println("==> se creo el archivo ", path)

	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if IsError(err) { return }
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contenido)
	if IsError(err) { return }

	// save changes
	err = file.Sync()
	if IsError(err) { return }

	log.Println("==> se escribio el contenido en el archivo creado")
}

//Subir al drive
func SubirInDrive(d *drive.Service, title string, description string,parentId string, mimeType string, filename string) (*drive.File){
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

func IsError(err error) bool {
	if err != nil {
		log.Println(err.Error())
	}

	return (err != nil)
}


// Busca el token del archivo
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
			return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}