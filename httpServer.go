package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type eprocessoData struct {
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

func ignoreThisShit(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     ignoreThisShit,
}

func WebSocketHandler(chWrite chan WebSocketMessage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		trace.Printf("WS ENDPOINT")

		conn, err := upgrader.Upgrade(w, r, nil)
		// defer conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		for msg := range chWrite {			
			jsonMsg, errr := json.Marshal(msg)
			if errr != nil {
				info.Printf("msg ERROR WS")
				trace.Printf("msg ERROR WS \n%#v \n%#v", errr, msg)
				continue
			}
			conn.WriteMessage(websocket.TextMessage, jsonMsg)
		}

	}
}

func withData(data *eprocessoData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		t := template.Must(template.ParseFiles("front_build/index.html"))
		// Trace.Printf("\n\n\n\nTEMPLATE:\n%#v\n\n<- do template\n", data)
		e := t.Execute(w, data)
		if e != nil {
			info.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error())
			trace.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error())
		}
	}
}

func handlerWithInitialTemplate(api *apiConn, pasta_down string, port string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trace.Printf("-")
		var data map[string]map[string]string
		trace.Printf("\n-x pegar dados do api pra o servidor")
		api.patchWinPrincipal()
		trace.Printf("-")
		strData := api.getInitialJSONData()
		// Trace.Printf("\nString Retornada do api\n%#v\n\n\n", strData)
		// Trace.Printf("\nmap antes\n%#v\n\n\n", data)
		dataBytes := []byte(strData)
		// Trace.Printf("\ndata bytes\n%#v\n\n\n", dataBytes)
		err := json.Unmarshal(dataBytes, &data)
		if err != nil {
			trace.Println("Erro no Unmarshall", err)
			trace.Println(data)
		}
		trace.Printf("-")

		data["__META__"]["pasta_download"] = pasta_down

		trace.Printf("\n-x")
		t := template.Must(template.ParseFiles("front_build/index.html"))
		trace.Printf("-")

		// Trace.Printf("\n\n\n\nTEMPLATE:\n%#v\n\n<- do template\n", data)

		e := t.Execute(w, &eprocessoData{data, port})
		trace.Printf("-")

		if e != nil {
			info.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error())
			trace.Printf("\n\n\nERRO no parse: %#v\n%#v\n", e, e.Error())
		}
		trace.Printf("-")

	}
}

func handlerWithInitialJson(api *apiConn, pasta_down string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trace.Printf("-")
		var data map[string]map[string]string
		trace.Printf("-")

		trace.Printf("\n-x pegar dados do api pra o servidor")
		api.patchWinPrincipal()
		trace.Printf("-")
		strData := api.getInitialJSONData()
		trace.Printf("-")
		json.Unmarshal([]byte(strData), &data)
		trace.Printf("-")

		data["__META__"]["pasta_download"] = pasta_down
		trace.Printf("-")

		bytejson, _ := json.Marshal(data)
		trace.Printf("-")

		// Trace.Printf("\n%#v", strData)
		w.Header().Set("Content-Type", "application/json")
		trace.Printf("-")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		trace.Printf("-")

		w.Write(bytejson)
		trace.Printf("-")
		trace.Printf("\n-x")
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
			trace.Println("Apagar arquivos handler... ok?", ok)
			var mapa = make(map[string][]string)
			mapa["certo"] = ok
			mapa["errado"] = errado
			trace.Println("MAPA\n", mapa)
			arrBytes, err := json.Marshal(mapa)
			if err != nil {
				trace.Println("nao fez o []bytes :(")
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
		trace.Println("init Sida apenas")
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

// func pesquisaSidaProcesso(api *apiConn) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == "POST" {
// 			body, err := ioutil.ReadAll(r.Body)
// 			if err != nil {
// 				http.Error(w, "Error reading request body",
// 					http.StatusInternalServerError)
// 			}
// 			processo := string(body)
// 			trace.Println("processo no handler (body)\n", processo)
// 			resp := api.SIDAConsultaProcesso(processo)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.Header().Set("Access-Control-Allow-Origin", "*")
// 			w.Write([]byte(resp))
// 		} else {
// 			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		}
// 	}
// }

func handleInject(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			trace.Printf("handleInject request info: \n%v", &r)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			trace.Println("handleInject")
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
			trace.Printf("handleInject request info: \n%v", &r)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			trace.Println("handleInject")

			res := api.evalOnSidaWindow(body)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(res)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func handleToggleScrapDownloads(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// logic
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			isLocked := api.toggleMutexLock()
			if isLocked {
				w.Write([]byte("{\"isLocked\": 1}"))
			} else {
				w.Write([]byte("{\"isLocked\": 0}"))
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func pesquisaSidaVariosProcessos(api *apiConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trace.Println("endpoint reached: pesquisaSidaVariosProcessos")

		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "POST" {
			trace.Println("      METHODO: post ...")
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				trace.Printf("body ... %o", err)
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}
			processos := string(body)
			trace.Printf("   processos no handler (body) : %s \n", processos)
			arrProcs := strings.Split(processos, "|")
			trace.Printf("    processos : %s \n", arrProcs)

			api.SIDAInit()
			s := make([]map[string]interface{}, 0)
			for y, processo := range arrProcs {
				info.Printf("Baixando processo %v (%d de %d)", processo, y, len(arrProcs))

				m := make(map[string]interface{})
				trace.Printf(" Pesquisar processo : %s  ", processo)
				resp := api.SIDAConsultaProcesso(processo)
				trace.Printf(" Resposta: \n\n%s\n\n", resp)
				m["processo_key"] = processo
				switch resp.quantidadeIdentificador {

				case NENHUMA_INSCRICAO:
					trace.Println("nehuma inscrição... faz json aki mesmo")
					m["qtd_inscricoes"] = 0
					m["lista_inscrições"] = make([]map[string]interface{}, 0)
					s = append(s, m)

				case UMA_INSCRICAO:
					trace.Println("uma inscrição... ja vem json")
					var i map[string]interface{}
					err := json.Unmarshal([]byte(resp.Json), &i)
					if err != nil {
						trace.Printf("erro ao parsear... %v", err)
					}
					trace.Printf("uma inscrição mapa i %v", i)
					m["qtd_inscricoes"] = 1
					m["lista_inscrições"] = make([]interface{}, 0)
					m["lista_inscrições"] = append(m["lista_inscrições"].([]interface{}), i)
					s = append(s, m)
					trace.Printf("slice s %v", s)

				case VARIAS_INSCRICOES:
					trace.Printf("processo %s tem mais de uma inscrição.", processo)
					api.waitForCondition(SIDA_WINDOW, `(function () {
						return document.getElementsByTagName('img').length > 0;
						})();`)
					trace.Printf("processo %s [ 0]", processo)
					api.evalOnSidaWindow([]byte("function abrirJanela(href) {window.navigate(href);}"))
					trace.Printf("processo %s [ 1]", processo)
					api.evalOnSidaWindow([]byte(jsPolyfills))
					trace.Printf("processo %s [ 2] injetou polyfillss", processo)
					api.evalOnSidaWindow([]byte(`var arraYImages = document.querySelectorAll('img');
					arraYImages[arraYImages.length - 2].click();`))
					trace.Printf("processo %s [ 3] clickou ?", processo)
					api.waitForCondition(SIDA_WINDOW, "document.getElementById('formatoHtml').checked === true")
					trace.Printf("processo %s [ 4] avaliou true a condição ?", processo)
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
					trace.Printf("processo %s [ 5] chegou no final... ?", processo)
					api.waitNotBusySidaWindow(SIDA_WINDOW)
					trace.Printf("processo %s [ 6] chegou no final... ?", processo)
					api.evalOnSidaWindow([]byte(jsPolyfills))
					trace.Printf("processo %s [ 7] chegou no final... ?", processo)
					api.evalOnSidaWindow([]byte(jsPolyfillShimInjectScript))
					trace.Printf("\n\nprocesso %s [ 8] VERIFICAR ?", processo)
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
					trace.Printf("\n\nprocesso %s [ 9] ???? ?\n\n", processo)
					api.evalOnSidaWindow([]byte(jsSidaGetInscInfo))
					api.evalOnSidaWindow([]byte(JSUnicodeHandle))
					trace.Printf("\n\nprocesso %s [10] injetou... \nchamar stringify()\n", processo)
					jsonStr := api.getInscricoesFromSidaMulti()
					trace.Printf("\n\nprocesso %s [11] JSON: %s\n", processo, jsonStr)
					// resp = resp + "_MULTI_||>" + jsonStr + "@@\n"
					var arr []interface{}
					err := json.Unmarshal([]byte(jsonStr), &arr)
					if err != nil {
						trace.Printf("erro %v", err)
					}
					m["lista_inscrições"] = arr
					m["qtd_inscricoes"] = len(arr)

					s = append(s, m)

					trace.Printf("processo %s [20] chegou no final... ?", processo)

				}

			}
			jsonFinal, err := json.Marshal(s)
			trace.Printf("slice final %v", s)
			if err != nil {
				trace.Printf("erro %v", err)
			}
			trace.Println("final do for para todos os processos... mandar resposta final...")
			w.Write(jsonFinal)
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

func serveHttp(api *apiConn, pdfPath string, port string, wsWriteChan chan WebSocketMessage) {
	trace.Printf("\n-x")
	fs := http.FileServer(http.Dir("front_build/static/"))
	sw := http.FileServer(http.Dir("front_build/sw/"))
	trace.Printf("\n-x")
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/sw/", http.StripPrefix("/sw/", sw))
	trace.Printf("\n-File Server")
	fsPdf := http.FileServer(http.Dir(pdfPath))
	trace.Printf("\n-x")
	http.Handle("/pdf/", corsMiddleware(http.StripPrefix("/pdf/", fsPdf)))
	trace.Printf("\n-x")
	// Trace.Printf("\n\nDATA (serve api)%#v\n---\n-< do serveApi", data)
	http.HandleFunc("/", handlerWithInitialTemplate(api, pdfPath, port)) // set router
	http.HandleFunc("/json", handlerWithInitialJson(api, pdfPath))       // set router
	http.HandleFunc("/deletefiles", deleteFilesHandler(pdfPath))         // set router
	http.HandleFunc("/initSida", initSidaHandler(api))                   // set router
	// http.HandleFunc("/pesquisa_sida_processo", pesquisaSidaProcesso(api)) // set router
	http.HandleFunc("/eval_js", handleInject(api))
	http.HandleFunc("/eval_sida_window_js", handleInjectSidaWindow(api))
	http.HandleFunc("/abre_sida_window", abreSidaWindow(api))
	http.HandleFunc("/pesquisa_sida_varios_processos", pesquisaSidaVariosProcessos(api)) // set router
	http.HandleFunc("/toggle_lock", handleToggleScrapDownloads(api))                     // set router
	http.HandleFunc("/ws", WebSocketHandler(wsWriteChan))                                //  set router
	info.Printf("\nServir na porta " + port + "... Visite http://localhost:" + port + " no Chrome (ou Firefox se tiver atualizado)")
	trace.Printf("\n-x")
	err := http.ListenAndServe(":"+port, nil) // set listen port
	trace.Printf("\n-x")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
