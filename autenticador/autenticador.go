package autenticador

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "os"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/drive/v3"
)

// Recupera el token, lo guarda y devuelve el cliente generado
func getClient(config *oauth2.Config) *http.Client {
        // El archivo token.json almacena los token de ACCESO y de ACTUALIZACION del usuario, y es
        // creado automaticamente cuando el flujo de autorizacion es completado por primera vez
        tokFile := "token.json"
        tok, err := tokenFromFile(tokFile)
        if err != nil {
			//si hay error, entonces busca el token web
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
        }else{
			//el token ya lo tenemos!!!!!
		}
        return config.Client(context.Background(), tok)
}

// Solicita un token a la web y lo devuelve (esto lo puede hacer a partir del authorization code).
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
        authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
        fmt.Printf("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                log.Fatalf("No se pudo leer el codigo de autorizacion %v", err)
        }

        tok, err := config.Exchange(context.TODO(), authCode)
        if err != nil {
                log.Fatalf("No se pudo recuperar el token de la web %v", err)
        }
        return tok
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

// Guarda el token en el archivo path
func saveToken(path string, token *oauth2.Token) {
		log.Println("Guardando el token en el archivo:", path)
		
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        if err != nil {
                log.Fatalf("Unable to cache oauth token: %v", err)
        }
        defer f.Close()
        json.NewEncoder(f).Encode(token)
}

func logearme(){//Se escribe en mayuscula por ser publico
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
			log.Fatalf("No se pudo leer el archivo credentials.json : %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
			log.Fatalf("No se puede analizar el 'client secret file' para configurar: %v", err)
	}
	client := getClient(config) //este vera si tiene que redirigirlo a la web o tiene el token precargado

	
	//logica para usar el drive
	srv, err := drive.New(client)
	if err != nil {
			log.Fatalf("No se pudo recuperar el drive del cliente: %v", err)
	}

	//srv.Files *FilesService -> r *FilesListCall
	r, err := srv.Files.List().PageSize(10).
			Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
			log.Fatalf("No se pudo recuperar los archivos: %v", err)
	}
	fmt.Println("Archivos:")
	if len(r.Files) == 0 {
			fmt.Println("No tiene archivos.")
	} else {
			for _, i := range r.Files {
					fmt.Printf("%s (%s)\n", i.Name, i.Id)
			}
	}
}

//Funcion que me avisa si estoy autenticado o me da la url para hacerlo
func GetEstadoAutenticacion() string{//Se escribe en mayuscula por ser publico
        b, err := ioutil.ReadFile("credentials.json")
        if err != nil {
                log.Fatalf("No se pudo leer el archivo credentials.json : %v", err)
        }

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, drive.DriveScope)
        if err != nil {
                log.Fatalf("No se puede analizar el 'client secret file' para configurar: %v", err)
        }
		
		tokFile := "token.json"
		tok, err := tokenFromFile(tokFile)
		if err != nil {
			//hay q logear !!
			authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
			var url string
			url = authURL
			log.Println("La url para logear es: " + authURL)
			return url
		}else{
			//el token ya lo tenemos!!!!!
			log.Print("Ya tenemos el token ")
			log.Println(tok)
			return "AUTENTICADO"
		}
}

//Funcion que me autentica
func Autenticar(code string) string{//Se escribe en mayuscula por ser publico
	log.Println("Empezamos a autenticar...")
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
			log.Fatalf("No se pudo leer el archivo credentials.json : %v", err)
	}

	//BUsca la configuracion
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
			log.Fatalf("No se puede analizar el 'client secret file' para configurar: %v", err)
	}

	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		//Si no lo tengo, hay que cargar el token desde la web
		tokAux, errAux := config.Exchange(context.TODO(), code)
		if errAux != nil {
			log.Fatalf("No se pudo recuperar el token de la web %v", errAux)
			return "ERROR"
		}
		tok = tokAux
		saveToken(tokFile, tok)
		return "OK"
	}else{
		//el token ya lo tenemos!!!!!
		log.Fatalf("Esto no se deberia var ya que en teoria no estoy autenticado aun.")
		return "ERROR"
	}
	
}
