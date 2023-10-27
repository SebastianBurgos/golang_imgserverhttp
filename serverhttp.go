package main

import (
	"encoding/base64" // Importamos el paquete para codificar y decodificar en base64
	"html/template"   // Importamos el paquete para manejar plantillas HTML

	// Importamos el paquete para leer archivos
	"math/rand"     // Importamos el paquete para generar números aleatorios
	"net/http"      // Importamos el paquete para manejar el protocolo HTTP
	"os"            // Importamos el paquete para interactuar con el sistema operativo
	"path/filepath" // Importamos el paquete para manipular rutas de archivos
)

// Definimos una estructura para los datos de la página
type PageData struct {
	Port       string      // Puerto del servidor
	HostName   string      // Nombre del host donde se ejecuta el servidor
	RandomPics []ImageData // Lista de datos de imágenes aleatorias
}

// Definimos una estructura para los datos de las imágenes
type ImageData struct {
	Name string // Nombre del archivo de la imagen
	Data string // Datos de la imagen codificados en base64
}

func main() {
	if len(os.Args) < 3 {
		println("Uso: go run main.go <puerto> <directorio_de_imagenes>")
		return
	}

	port := os.Args[1]     // Obtenemos el puerto de los argumentos de la línea de comandos
	imageDir := os.Args[2] // Obtenemos el directorio de imágenes de los argumentos

	// Configuramos el manejador para la ruta raíz del servidor
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, err := os.Hostname() // Obtenemos el nombre del host
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Leemos la lista de archivos de imágenes en el directorio
		imageFiles, err := os.ReadDir(imageDir)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var randomPics []ImageData
		for i := 0; i < 3 && i < len(imageFiles); i++ {
			randomIndex := rand.Intn(len(imageFiles))              // Generamos un índice aleatorio
			randomFile := imageFiles[randomIndex]                  // Obtenemos un archivo aleatorio
			filePath := filepath.Join(imageDir, randomFile.Name()) // Construimos la ruta del archivo

			fileData, err := os.ReadFile(filePath) // Leemos el contenido del archivo
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Codificamos los datos de la imagen en base64
			base64Data := base64.StdEncoding.EncodeToString(fileData)

			// Creamos una estructura con los datos de la imagen
			imageInfo := ImageData{
				Name: randomFile.Name(),
				Data: base64Data,
			}
			randomPics = append(randomPics, imageInfo) // Añadimos la imagen a la lista
		}

		// Creamos una estructura con los datos de la página
		data := PageData{
			Port:       port,
			HostName:   host,
			RandomPics: randomPics,
		}

		// Creamos una nueva plantilla a partir del HTML
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Ejecutamos la plantilla con los datos y escribimos la respuesta
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	println("Servidor corriendo en http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil) // Iniciamos el servidor en el puerto especificado
	if err != nil {
		println("Error al iniciar el servidor:", err)
	}
}
