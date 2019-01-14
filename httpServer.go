package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type EprocessoData struct {
	Data map[string]map[string]string
	Port string
	// str  string
}

type PayloadDeleteResponse struct {
	data map[string][]string
}

type WebSocketMessage struct {
	Tipo    string `json:"tipo"`
	Payload string `json:"payload"`
}

var upgrader = websocket.Upgrader{}

func WebSocketHandle(api *apiConn) http.HandlerFunc {
	Trace.Print("WS FUNC: montou")
	return func(w http.ResponseWriter, r *http.Request) {
		Trace.Print("WS FUNC: entrou")
		// (w).Header().Set("Access-Control-Allow-Origin", "*")
		// (w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// (w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Trace.Print("err upgrade:", err)
			return
		}
		go websocketGoroutine(api, c)
	}
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

func initSidaHandler(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			api.SIDAInit()
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write([]byte("ok"))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func abreSidaWindow(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Trace.Println("init Sida apenas")
		resp := api.SIDAInit()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if resp == true {
			w.Write([]byte("ok"))
		} else {
			w.Write([]byte("erro"))
		}
	}
}

func pesquisaSidaProcesso(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			processo := string(body)
			Trace.Println("processo no handler (body)\n", processo)
			resp := api.SIDAConsultaProcesso(processo)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write([]byte(resp))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func handleInject(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			Trace.Println("handleInject request info: \n%s\n", &r)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			Trace.Println("handleInject")

			res := api.evalOnWindow(body)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(res)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func handleInjectSidaWindow(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			Trace.Println("handleInject request info: \n%s\n", &r)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			Trace.Println("handleInject")

			res := api.evalOnSidaWindow(body)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(res)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func pesquisaSidaVariosProcessos(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Trace.Printf("sida ...")

		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "POST" {
			Trace.Printf("post ...")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				Trace.Printf("body ... %o", err)
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			processos := string(body)
			Trace.Println("processos no handler (body)\n", processos)
			arrProcs := strings.Split(processos, "|")
			respostaFinal := ""
			api.SIDAInit()
			for _, processo := range arrProcs {
				// element is the element from someSlice for where we are
				resp := api.SIDAConsultaProcesso(processo)

				// TEMP --------------------------------------------------
				if strings.Contains(resp, "FORAM LOCALIZADAS") {
					Trace.Printf("processo %s tem mais de uma inscrição.", processo)
					api.waitForCondition(SIDA_WINDOW, `(function () {
						return document.getElementsByTagName('img').length > 0;
						})();`)
					Trace.Printf("processo %s [ 0]", processo)
					api.evalOnSidaWindow([]byte("function abrirJanela(href) {window.navigate(href);}"))
					Trace.Printf("processo %s [ 1]", processo)
					api.evalOnSidaWindow([]byte(jsPolyfills))
					Trace.Printf("processo %s [ 2] injetou polyfillss", processo)
					api.evalOnSidaWindow([]byte(`var arraYImages = document.querySelectorAll('img');
					arraYImages[arraYImages.length - 2].click();`))
					Trace.Printf("processo %s [ 3] clickou ?", processo)
					api.waitForCondition(SIDA_WINDOW, "document.getElementById('formatoHtml').checked === true")
					Trace.Printf("processo %s [ 4] avaliou true a condição ?", processo)
					api.evalOnSidaWindow([]byte(`window.print = function () {
						return undefined;
						};`))
					api.evalOnSidaWindow([]byte("document.getElementsByTagName('img')['ok'].click();"))
					api.waitForCondition(SIDA_WINDOW, `(function () {
							window.print = function () {
								return undefined;
								};
								
								var tables = document.getElementsByTagName('table');
								var i__;
								
								for (i__ = 0; i__ < tables.length; i__++) {
									if (tables[i__].className === "Cabecalho") {
										return true;
									}
								}
								
								return false;
								})();`)
					Trace.Printf("processo %s [ 5] chegou no final... ?", processo)
					api.waitNotBusySidaWindow(SIDA_WINDOW)
					Trace.Printf("processo %s [ 6] chegou no final... ?", processo)
					api.evalOnSidaWindow([]byte(jsPolyfills))
					Trace.Printf("processo %s [ 7] chegou no final... ?", processo)
					api.evalOnSidaWindow([]byte(jsPolyfillShimInjectScript))
					Trace.Printf("\n\nprocesso %s [ 8] VERIFICAR ?", processo)
					api.waitNotBusySidaWindow(SIDA_WINDOW)
					api.waitForCondition(SIDA_WINDOW, `(function () {
						try {
						  var arr__ = Array.from(document.querySelectorAll('td')).filter(function (a) {
							return a === 1;
						  });
						} catch (e) {
						  return false;
						}
					  
						return true;
					  })();`)
					Trace.Printf("\n\nprocesso %s [ 9] ???? ?\n\n", processo)
					api.evalOnSidaWindow([]byte(jsSidaGetInscInfo))
					api.evalOnSidaWindow([]byte(JSUnicodeHandle))
					Trace.Printf("\n\nprocesso %s [10] injetou... \nchamar stringify()\n", processo)
					jsonStr := api.getInscricoesFromSidaMulti()
					Trace.Printf("\n\nprocesso %s [11] JSON: %s\n", processo, jsonStr)
					resp = resp + "_MULTI_||>" + jsonStr + "\n"
					Trace.Printf("processo %s [20] chegou no final... ?", processo)
				}

				respostaFinal = respostaFinal + "###" + processo + "$$$" + resp

			}

			w.Write([]byte(respostaFinal))
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
	http.HandleFunc("/", handlerWithInitialTemplate(api, pdfPath, port))  // set router
	http.HandleFunc("/json", handlerWithInitialJson(api, pdfPath))        // set router
	http.HandleFunc("/deletefiles", deleteFilesHandler(pdfPath))          // set router
	http.HandleFunc("/initSida", initSidaHandler(api))                    // set router
	http.HandleFunc("/pesquisa_sida_processo", pesquisaSidaProcesso(api)) // set router
	http.HandleFunc("/eval_js", handleInject(api))
	http.HandleFunc("/eval_sida_window_js", handleInjectSidaWindow(api))
	http.HandleFunc("/abre_sida_window", abreSidaWindow(api))
	http.HandleFunc("/pesquisa_sida_varios_processos", pesquisaSidaVariosProcessos(api)) // set router
	http.HandleFunc("/ws", WebSocketHandle(api))                                         // set router
	Info.Printf("\nServir na porta " + port + "... Visite http://localhost:" + port + " no Chrome (ou Firefox se tiver atualizado)")
	Trace.Printf("\n-x")
	err := http.ListenAndServe(":"+port, nil) // set listen port
	Trace.Printf("\n-x")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func websocketGoroutine(api *apiConn, c *websocket.Conn) {
	defer c.Close()
	for {

		msgType, msg, err := c.ReadMessage()
		if err != nil {
			Trace.Println("err read:", err)
			break
		}
		Trace.Printf("raw msgType: %v", msgType)
		Trace.Printf("raw msg: %v", msg)
		Trace.Printf("pra string msg : %s   ", msg)

		var obj WebSocketMessage
		er := json.Unmarshal(msg, &obj)
		if er != nil {
			Trace.Printf("ERRO UNMARSHALL : %s   ", er)
		}

		Trace.Printf("pra string msg : %v   ", obj)
		switch obj.Tipo {
		case "sida_pesquisa":
			Trace.Print("sida pesquisa")
			processosArr := strings.Split(obj.Payload, ",")
			api.SIDAInit()
			for _, processo := range processosArr {
				resp := api.SIDAConsultaProcesso(processo)
				resp = resp + "processo||>" + processo + "\n"
				jj, err := json.Marshal(WebSocketMessage{Tipo: "sida_resp", Payload: resp})
				if err != nil {
					Trace.Printf("ERRO MARSHALL : %s   ", err)
				}
				e := c.WriteMessage(msgType, jj)
				if e != nil {
					Trace.Println("err write:", e)
					break
				}
			}
		}
	}
}
