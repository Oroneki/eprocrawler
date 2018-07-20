package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type EprocessoData struct {
	Data map[string]map[string]string
	Port string
	// str  string
}

type PayloadDeleteResponse struct {
	data map[string][]string
}

func withData(data *EprocessoData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		t := template.Must(template.ParseFiles("front_build/index.html"))
		// Trace.Printf("\n\n\n\nTEMPLATE:\n%#v\n\n<- do template\n", data)
		e := t.Execute(w, data)
		if e != nil {
			Info.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error)
			Trace.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error)
		}
	}
}

func handlerWithInitialTemplate(api *apiConn, pasta_down string, port string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Trace.Printf("-")
		var data map[string]map[string]string
		Trace.Printf("\n-x pegar dados do api pra o servidor")
		api.patchWinPrincipal()
		Trace.Printf("-")
		strData := api.getInitialJSONData()
		// Trace.Printf("\nString Retornada do api\n%#v\n\n\n", strData)
		// Trace.Printf("\nmap antes\n%#v\n\n\n", data)
		dataBytes := []byte(strData)
		// Trace.Printf("\ndata bytes\n%#v\n\n\n", dataBytes)
		err := json.Unmarshal(dataBytes, &data)
		if err != nil {
			Trace.Println("Erro no Unmarshall", err)
			Trace.Println(data)
		}
		Trace.Printf("-")

		data["__META__"]["pasta_download"] = pasta_down

		Trace.Printf("\n-x")
		t := template.Must(template.ParseFiles("front_build/index.html"))
		Trace.Printf("-")

		// Trace.Printf("\n\n\n\nTEMPLATE:\n%#v\n\n<- do template\n", data)

		e := t.Execute(w, &EprocessoData{data, port})
		Trace.Printf("-")

		if e != nil {
			Info.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error)
			Trace.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error)
		}
		Trace.Printf("-")

	}
}

func handlerWithInitialJson(api *apiConn, pasta_down string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Trace.Printf("-")
		var data map[string]map[string]string
		Trace.Printf("-")

		Trace.Printf("\n-x pegar dados do api pra o servidor")
		api.patchWinPrincipal()
		Trace.Printf("-")
		strData := api.getInitialJSONData()
		Trace.Printf("-")
		json.Unmarshal([]byte(strData), &data)
		Trace.Printf("-")

		data["__META__"]["pasta_download"] = pasta_down
		Trace.Printf("-")

		bytejson, _ := json.Marshal(data)
		Trace.Printf("-")

		// Trace.Printf("\n%#v", strData)
		w.Header().Set("Content-Type", "application/json")
		Trace.Printf("-")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		Trace.Printf("-")

		w.Write(bytejson)
		Trace.Printf("-")
		Trace.Printf("\n-x")
	}
}

func deleteFilesHandler(pasta_down string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			ok, errado := apagarArquivos(pasta_down, string(body))
			Trace.Println("Apagar arquivos handler... ok?", ok)
			var mapa = make(map[string][]string)
			mapa["certo"] = ok
			mapa["errado"] = errado
			Trace.Println("MAPA\n", mapa)
			arrBytes, err := json.Marshal(mapa)
			if err != nil {
				Trace.Println("nao fez o []bytes :(")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(arrBytes)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func serveHttp(api *apiConn, pdfPath string, port string) {
	Trace.Printf("\n-x")
	fs := http.FileServer(http.Dir("front_build/static/"))
	Trace.Printf("\n-x")
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	Trace.Printf("\n-File Server")
	fsPdf := http.FileServer(http.Dir(pdfPath))
	Trace.Printf("\n-x")
	http.Handle("/pdf/", corsMiddleware(http.StripPrefix("/pdf/", fsPdf)))
	Trace.Printf("\n-x")
	// Trace.Printf("\n\nDATA (serve api)%#v\n---\n-< do serveApi", data)
	http.HandleFunc("/", handlerWithInitialTemplate(api, pdfPath, port)) // set router
	http.HandleFunc("/json", handlerWithInitialJson(api, pdfPath))       // set router
	http.HandleFunc("/deletefiles", deleteFilesHandler(pdfPath))         // set router
	Info.Printf("\nServir na porta " + port + "... Visite http://localhost:" + port + " no Chrome (ou Firefox se tiver atualizado)")
	Trace.Printf("\n-x")
	err := http.ListenAndServe(":"+port, nil) // set listen port
	Trace.Printf("\n-x")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
