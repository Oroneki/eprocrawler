package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"runtime"
	"sync"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// structs ---------------------------------------------------

type mensagem struct {
	tipo    string
	payload interface{}
}

type apiConn struct {
	linksMap    map[int]*ole.IDispatch
	mutex       *sync.Mutex
	perguntaCh  chan mensagem
	respostaCh  chan interface{}
	window      *ole.IDispatch
	windowJsObj *ole.IDispatch
	windowObj   *ole.IDispatch
}

type Processo struct {
	oleDlinkref  int
	numStrImpuro string
}

type SendDC struct {
	ol *ole.IDispatch
	ch chan *Processo
}

type SendJanProc struct {
	janid    string
	processo *Processo
}

// methods -----------------------------------------------------

func (api *apiConn) createObject(str string) *ole.IUnknown {
	// Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "CREATEOBJECT",
		payload: str,
	}
	resposta := <-api.respostaCh
	return resposta.(*ole.IUnknown)
}

func (api *apiConn) queryInterface(iun *ole.IUnknown) *ole.IDispatch {
	// Trace.Println("x")

	api.perguntaCh <- mensagem{
		tipo:    "QUERYINTERFACE",
		payload: iun,
	}
	resposta := <-api.respostaCh
	return resposta.(*ole.IDispatch)
}

func (api *apiConn) janelaEprocesso() bool {
	Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "JANELAEPROCESSO",
		payload: nil,
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) sendProcessosDaJanelaToChannel(ch chan *Processo) int {
	// Trace.Println("x")
	// Trace.Printf("\nsendProcessosDaJanelaToChannel -> %#v --- %#v\n", ol, ch)
	api.perguntaCh <- mensagem{
		tipo:    "SENDPROCESSOSTOCHANNEL",
		payload: &SendDC{api.window, ch},
	}
	resposta := <-api.respostaCh
	return resposta.(int)
}

func (api *apiConn) patchWinPrincipal() bool {
	Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "PATCHWINDOWPRINCIPAL",
		payload: nil,
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) abreProcesso(janID string, processo *Processo) bool {
	Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "ABREPROCESSO_0",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) paginaProcessoCarregou(janID string, processo *Processo) bool {
	Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "TESTA_PAGINAPROCESSOCARREGOU",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) paginaProcessoPatcheVaiProDownload(janID string, processo *Processo) bool {
	Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "TESTA_PAGINAPROCESSOPATCH_VAI_PRO_DOWNLOAD",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) clicaParaGerarPDF(janID string, processo *Processo) bool {
	Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "CLICA_PRA_GERAR_PDF_0",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) paginaDocumentosCarregou(janID string, processo *Processo) bool {
	// Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "PAGINA_DOCUMENTOS_CARREGOU_0",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) getHrefStringOrNot(janID string, processo *Processo) string {
	// Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "PEGA_HREF_STRING_OR_NOT_0",
		payload: &SendJanProc{janID, processo},
	}

	resposta := <-api.respostaCh

	// Trace.Printf("%T -> href: %v", resposta, resposta)
	return resposta.(string)
}

func (api *apiConn) getCookies() string {
	// Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "GET_COOKIES_0",
		payload: nil,
	}

	resposta := <-api.respostaCh
	// Info.Printf("%T -> cookies: %v", resposta, resposta)
	return resposta.(string)
}

func (api *apiConn) getInitialJSONData() string {
	// Trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "GET_JSON_DATA_0",
		payload: nil,
	}

	resposta := <-api.respostaCh
	// Info.Printf("%T -> cookies: %v", resposta, resposta)
	return resposta.(string)
}

// methods bootstrap ---------------------------------------------

func (api *apiConn) olePoolInicio() {
	Trace.Println("x Inicio Pool")
	runtime.LockOSThread()
	Trace.Println("x")
	err := ole.CoInitialize(0)
	Trace.Println("x")
	if err != nil {
		Trace.Println("x Erro")
		oleerr := err.(*ole.OleError)
		// S_FALSE           = 0x00000001 // CoInitializeEx was already called on this thread
		if oleerr.Code() != ole.S_OK && oleerr.Code() != 0x00000001 {
			fmt.Println(err)
		}
	} else {
		// Only invoke CoUninitialize if the thread was not initizlied before.
		// This will allow other go packages based on go-ole play along
		// with this library.
		Trace.Println("x Tranquilo")
		defer ole.CoUninitialize()
	}

	var regexHrefLinkProcesso = regexp.MustCompile(`\'.*?\'`)
	for mensagem := range api.perguntaCh {
		Trace.Printf(`
============================================================================
	Mensagem:
	| TIPO: 	%s
	| PAYLOAD:	%#v
============================================================================
			`, mensagem.tipo, mensagem.payload)
		switch tipo := mensagem.tipo; tipo {
		case "CREATEOBJECT":
			api.mutex.Lock()
			unknown, err := oleutil.CreateObject(mensagem.payload.(string))
			api.mutex.Unlock()
			if err != nil {
				panic(err)
			}
			api.respostaCh <- unknown

		case "QUERYINTERFACE":
			api.mutex.Lock()
			iun := mensagem.payload.(*ole.IUnknown)
			id, err := iun.QueryInterface(ole.IID_IDispatch)
			api.mutex.Unlock()
			if err != nil {
				panic(err)
			}
			api.respostaCh <- id

		case "JANELAEPROCESSO":
			Trace.Println("x")
			api.mutex.Lock()
			Trace.Println("x")
			unknown, _ := oleutil.CreateObject("shell.Application")
			Trace.Println("x")
			shell, _ := unknown.QueryInterface(ole.IID_IDispatch)
			Trace.Println("x")
			windows, _ := shell.CallMethod("Windows")
			Trace.Println("x")
			wins := windows.ToIDispatch()
			Trace.Println("x")
			nois, _ := wins.GetProperty("Count")
			Trace.Println("x")
			valConta := int(nois.Val)
			Trace.Printf("\n %d janelas identificadas.", valConta)
			var re = regexp.MustCompile(`eprocesso\.suiterfb\.receita\.fazenda\/ControleAcessarCaixaTrabalho\.asp\?psAcao=apresentarPagina&psLimpaEquipe=`)
			Trace.Println("x")
			var itemjanela *ole.IDispatch
			for i := 0; i < valConta; i++ {
				Trace.Println(" o")
				item, _ := wins.CallMethod("Item", i)
				Trace.Println(" o")
				itemd := item.ToIDispatch()
				Trace.Printf(" \n            o    %#v", itemd)
				locationURLV, _ := itemd.GetProperty("LocationURL")
				Trace.Println(" o")
				urlV := locationURLV.Value()
				Trace.Println(" o")
				url := urlV.(string)
				Trace.Println(" o")
				Trace.Printf("\nJanela Identificada: (id: %d) %s", i, url)
				Trace.Println(" o")

				testeRegex := re.MatchString(url)
				Trace.Println(" o")
				if testeRegex {
					Trace.Println(" o!")
					Trace.Printf(`
						[!] E-PROCESSO : (id: %d) %s
						`, i, locationURLV)
					itemjanela = itemd
					Trace.Println(" o!")
					break
				}
			}
			busy := oleutil.MustGetProperty(itemjanela, "Busy")
			container := oleutil.MustGetProperty(itemjanela, "Container")
			application := oleutil.MustGetProperty(itemjanela, "Application")
			fmt.Printf(`
			Janela Internet Explorer:
			+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
			Busy              %v
			Container:        %v
			Application       %v
			HWND              %v
			name              %v
			*****************************************************************************************	
				`,
				busy.Value(),
				container.Value(),
				application.Value(),
				oleutil.MustGetProperty(itemjanela, "HWND").Value(),
				oleutil.MustGetProperty(itemjanela, "Name").Value(),
			)

			api.mutex.Unlock()
			Trace.Println("x")
			api.window = itemjanela
			Trace.Println("x")
			api.respostaCh <- true

		case "PATCHWINDOWPRINCIPAL":
			api.mutex.Lock()
			Trace.Println("x")
			iePrincipalD := api.window
			Trace.Println("x")
			iePrincipalDocumentV, _ := iePrincipalD.GetProperty("Document")
			Trace.Println("x")
			iePrincipalDocumentD := iePrincipalDocumentV.ToIDispatch()
			Trace.Println("x")
			windowPrincipal := oleutil.MustGetProperty(iePrincipalDocumentD, "parentWindow").ToIDispatch()
			Trace.Println("x")
			api.windowJsObj = windowPrincipal
			Trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", `window.oro_obj = {};`)
			Trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", `
				jQuery.extend({
					stringify  : function stringify(obj) {         
						if ("JSON" in window) {
							return JSON.stringify(obj);
						}
						
						var t = typeof (obj);
						if (t != "object" || obj === null) {
							// simple data type
							if (t == "string") obj = '"' + obj + '"';
							
							return String(obj);
							} else {
								// recurse array or object
								var n, v, json = [], arr = (obj && obj.constructor == Array);
								
								for (n in obj) {
									v = obj[n];
									t = typeof(v);
									if (obj.hasOwnProperty(n)) {
										if (t == "string") {
											v = '"' + v + '"';
											} else if (t == "object" && v !== null){
												v = jQuery.stringify(v);
											}
											
											json.push((arr ? "" : '"' + n + '":') + String(v));
										}
									}
									
									return (arr ? "[" : "{") + String(json) + (arr ? "]" : "}");
								}
							}
							});
							
							window.getJsonData = function() {
								var td_regex = /ddrivetip\(\'(.*?)\'\,/m;
								var tableCol = document.getElementsByTagName("table");
								var num_campos;
								var headers;
								var map_headers = {};
								var map_final = {};
								console.log(tableCol.length);
								for (i = 0; i < tableCol.length; i++) {
									if (tableCol[i].id !== "tblProcessos") {
										continue;
									}
									headers = tableCol[i].getElementsByTagName("th");
									num_campos = headers.length;
									console.log(num_campos, " campos");
									for (n = 0; n < num_campos; n++) {
										//console.log(n, " -> ", headers[n].innerText);
										map_headers[n] = headers[n].innerText.replace(/[^a-zA-Z ]/g, "_").toLowerCase().replace(/[^a-zA-Z0-9\_]+$/gm, "");
										//console.log(n, " -> ", map_headers[n]);
									}
									var trs = tableCol[i].getElementsByTagName("tr");
									var trlen = trs.length;
									//console.log();
									//console.log("trlen", trlen);
									//console.log(trs);
									for (t = 0; t < trlen; t++) {
										var tds = trs[t].getElementsByTagName("td");
										var mapinha = {};
										if (tds.length !== num_campos) { continue };
										//console.log(t, tds);
										for (k = 0; k < tds.length; k++) {
											//console.log("---------------------------------------");
											var value
											if (!tds[k].onmouseover) {
												value = tds[k].innerText
												} else {
													var match = td_regex.exec(tds[k].outerHTML);
													value = match[1];
												}
												//console.log(k, " - ", map_headers[k], " --> ", value);
												value = value;

												//if (value === "" || value === "-") {
												//	console.log("+Descartando: ", value);
												//	continue
												//}
												
												mapinha[map_headers[k]] = value;
											}
											console.log(Object.keys(mapinha));
											key_p = mapinha["n_mero processo"].replace(/\(\d+\)/g, "");
											key_p = key_p.replace(/\D/g, "");
											//console.log("key_p ", key_p);
											map_final[key_p] = mapinha;
										}
										
									}
									map_final["__META__"] = {codEquipe: document.getElementById("hidEquipeSelecionadaCaixaTrabalho").value};
									return jQuery.stringify(map_final);
									};
									`)
			Trace.Println("x")

			oleutil.MustCallMethod(windowPrincipal, "eval", `window.hacked_visualizarProcesso = function(TARGET_JANELA, psNumeroProcesso, psNumeroEquipeAtividade, psNomeEquipeAtual, psNomeAtividadeAtual) {
										
				abrirPopupProcesso(TARGET_JANELA, { 'psNumeroProcesso': psNumeroProcesso, 'psNumeroEquipeAtividade': psNumeroEquipeAtividade, 'psNomeEquipeAtual': psNomeEquipeAtual, 'psNomeAtividadeAtual': psNomeAtividadeAtual });
			};`)

			Trace.Println("x")

			oleutil.MustCallMethod(windowPrincipal, "eval", `window.abrirPopupProcesso = function (TARGET, params) {
				var parametros = "psAcao=exibir&" ;
				var lnHeight = screen.height;
				var lnWidth = screen.width;
				parametros += "psNumeroProcesso=" + params["psNumeroProcesso"] ;
				parametros += "&psNumeroEquipeAtividade="+params["psNumeroEquipeAtividade"] ;
				parametros += "&psNomeEquipeAtual=" + params["psNomeEquipeAtual"] ;
				parametros += "&psNomeAtividadeAtual=" + params["psNomeAtividadeAtual"] ;
				var posJan = Object.keys(window.oro_obj).length * 20 + 5;
				if (window.oro_obj[TARGET]) {
					window.oro_obj[TARGET] = window.open("about:blank", TARGET, "width="+(lnWidth-300)+",height="+(lnHeight-270)+",scrollbars=no,resizable=yes,left=100,top=100");
				}
				window.oro_obj[TARGET] = window.open("/ControleVisualizacaoProcesso.asp?" + parametros, TARGET, "width="+(lnWidth-360)+",height="+(lnHeight-298)+",scrollbars=no,resizable=yes,left="+(posJan)+",top="+(posJan));
				console.log('Abrindo processo ', params["psNumeroProcesso"], ' na janela ', TARGET);
			};`)
			Trace.Println("x")

			oleutil.MustCallMethod(windowPrincipal, "eval", `window.test_load_processo = function(TARGET, processo_str) {
				var t_1_b_ =  window.oro_obj[TARGET] && window.oro_obj[TARGET].document && window.oro_obj[TARGET].document.readyState === "complete";
				console.log(t_1_b_, "  readystate")
				if (t_1_b_) {
					console.log('testa_readystate');
					return window.oro_obj[TARGET].document.title.replace(/\D/g, "").indexOf(processo_str.replace(/\D/g, "").slice(0, 10)) > -1 &&
					window.oro_obj[TARGET].document.getElementsByTagName("area") && 
					window.oro_obj[TARGET].document.getElementsByTagName("area").length > 1;
				} else {
					return false;
				};    
			};

			window.clica_pra_gerar_pdf = function(TARGET) {
				window.oro_obj[TARGET].document.getElementById("chkMetaDados").checked = false;
				window.oro_obj[TARGET].document.getElementById("chkMetaDados").checked = false;
				var naoPag = window.oro_obj[TARGET].document.getElementById("chkNaoPaginavel");
				if (naoPag) {
					naoPag.checked = false;
					naoPag.checked = false;
				};				
				window.oro_obj[TARGET].document.getElementById("imgPdf").click();
				window.oro_obj[TARGET].document.getElementById("imgPdf").click();
			};

			window.get_download_href_or_false = function(TARGET) {
				var string_get_donload_retiurn = window.oro_obj[TARGET].document.getElementById("linkDownloadParte1") || {title: ""};
				if (string_get_donload_retiurn.title === "") {
					if (!window.oro_obj[TARGET].document.getElementById("imgPdf").disabled) {
						window.oro_obj[TARGET].document.getElementById("imgPdf").click()
					};
				};
				return string_get_donload_retiurn.title;
				};		
				
				`)
			Trace.Println("x")

			oleutil.MustCallMethod(windowPrincipal, "eval", `window.test_load_pagina_download = function(TARGET) {
				var t_2_a_ = window.oro_obj[TARGET].document && window.oro_obj[TARGET].document.readyState === "complete";
				if (t_2_a_) {
					return window.oro_obj[TARGET].document.getElementsByTagName("img") && window.oro_obj[TARGET].document.getElementsByTagName("img").length > 1;
				} else {
					return false
				}
			};`)
			api.mutex.Unlock()
			Trace.Println("x")

			api.respostaCh <- true

		case "SENDPROCESSOSTOCHANNEL":
			payload := mensagem.payload.(*SendDC)
			pD := payload.ol
			pC := payload.ch
			api.mutex.Lock()
			iePrincipalD := pD
			iePrincipalDocumentV, _ := iePrincipalD.GetProperty("Document")
			iePrincipalDocumentD := iePrincipalDocumentV.ToIDispatch()
			iePrincipalHTMLDocumentQ, _ := iePrincipalDocumentD.QueryInterface(ole.IID_IDispatch)
			iePrincipalHTMLDocumentgetElementsByTagNameV, _ := iePrincipalHTMLDocumentQ.CallMethod("getElementsByTagName", "a")
			iePrincipalHTMLDocumentgetElementsByTagNameD := iePrincipalHTMLDocumentgetElementsByTagNameV.ToIDispatch()
			nodeListLengthV, _ := iePrincipalHTMLDocumentgetElementsByTagNameD.GetProperty("length")
			len := int(nodeListLengthV.Val)
			var procNuminnerHTMLregex = regexp.MustCompile(`D&nbsp;&nbsp;(\d+[\.\d\/\-]+)`)

			var resp int
			canalDeProcessos := pC
			for loc := 0; loc < len; loc++ {

				stringNode := fmt.Sprint(loc)

				linkaV, _ := iePrincipalHTMLDocumentgetElementsByTagNameD.GetProperty(stringNode)

				linkaD := linkaV.ToIDispatch()
				innerHTML, _ := linkaD.GetProperty("innerHTML")

				innerHTMLVal := fmt.Sprint(innerHTML.Value())

				if procNuminnerHTMLregex.MatchString(innerHTMLVal) {
					processo := procNuminnerHTMLregex.FindStringSubmatch(innerHTMLVal)[1]

					api.linksMap[resp] = linkaD
					go func(p *Processo) {
						canalDeProcessos <- p
						fmt.Printf("\n + Microrotina encaminhou processo %v pro canal.\n", p)
					}(&Processo{resp, processo})
					resp++
				}
			}
			api.mutex.Unlock()
			api.respostaCh <- resp

		case "ABREPROCESSO_0":
			api.mutex.Lock()
			pld := mensagem.payload.(*SendJanProc)
			janid := pld.janid
			processo := pld.processo

			href := oleutil.MustGetProperty(api.linksMap[processo.oleDlinkref], "href").Value().(string)

			args := regexHrefLinkProcesso.FindAllString(href, 4)

			cmdAbreJs := fmt.Sprintf(
				"javascript:hacked_visualizarProcesso('%s', %s, %s, %s, %s)",
				janid,
				args[0],
				args[1],
				args[2],
				args[3],
			)

			oleutil.PutProperty(api.linksMap[processo.oleDlinkref], "href", cmdAbreJs)

			oleutil.MustCallMethod(api.linksMap[processo.oleDlinkref], "click")

			api.mutex.Unlock()
			api.respostaCh <- true

		case "TESTA_PAGINAPROCESSOCARREGOU":
			pld := mensagem.payload.(*SendJanProc)
			janid := pld.janid
			processo := pld.processo
			// Trace.Println("x")
			api.mutex.Lock()
			areaNodeList, e := oleutil.CallMethod(api.windowJsObj, "eval", fmt.Sprintf(`window.test_load_processo("%s", "%s")`, janid, processo.numStrImpuro))
			api.mutex.Unlock()
			if e != nil {
				api.respostaCh <- false
				Trace.Printf("\n%s - ERRO na chamada da pagina de carregamento do processo\n", janid)
			} else {
				resp := areaNodeList.Value().(bool)
				api.respostaCh <- resp
				Trace.Printf("\n\n%s - Resposta do carregamento da pagina do processo: %v\n\n", janid, resp)
			}

		case "TESTA_PAGINAPROCESSOPATCH_VAI_PRO_DOWNLOAD":
			pld := mensagem.payload.(*SendJanProc)
			janid := pld.janid
			// processo := pld.processo
			// Trace.Println("x")
			api.mutex.Lock()
			oleutil.MustCallMethod(
				api.windowJsObj,
				"eval",
				fmt.Sprintf(`					
					window.oro_obj.%s.obterMultiplosDocumentos = function () {
					var psOperacao = "C";
					var psEscopo = "I";
					var winbind = window.oro_obj.%s;
					var gsNumeroProcesso = winbind.document.getElementById('hidNumeroProcesso').value;
					var gsNumeroEquipeAtividade = winbind.document.getElementById('hidNumeroEquipeAtividade').value;
					var lsSituacaoProcesso = winbind.document.getElementById('hidSituacaoProcesso').value;					
					var gsResponsavelProcesso = winbind.document.getElementById('hidResponsavelProcesso').value;
					var gsNomeEquipe = winbind.document.getElementById('hidNomeEquipeAtual').value;
					var gsNomeAtividade = winbind.document.getElementById('hidNomeAtividadeAtual').value;
					var gsNumeroProcessoFormatado = winbind.document.getElementById('hidNumeroProcessoFormatado').value;					
					//Pega só as chaves dos documentos
					var gaDocumentosSelecionados = winbind.gsDocSelecionados.split("@");
					//valida se selecionou no máximo 1000 documentos
					if (gaDocumentosSelecionados.length>1000){
						alert("Selecione no máximo 1000 documentos.");
						return;
					}			
				
					var laDocumento = new Array();
					var laNumeroDocumentos = new Array();
					for (var lnIndice = 0; lnIndice < gaDocumentosSelecionados.length; lnIndice++){
				
						laDocumento = gaDocumentosSelecionados[lnIndice].split("|");
						laNumeroDocumentos[lnIndice] = laDocumento[0];
					}
					var lsURL = "ControleMultiplosDocumentos.asp?psAcao=apresentarPagina&psNumeroProcesso=" + gsNumeroProcesso + "&paDocSelecionados=" +
							laNumeroDocumentos + "&psNumeroEquipeAtividade=" + gsNumeroEquipeAtividade + "&psOperacao=" + psOperacao + "&psEscopo=" + psEscopo +
							"&psSituacaoProcesso=" + lsSituacaoProcesso + "&psResponsavelProcesso=" + gsResponsavelProcesso + "&psNomeEquipe=" + gsNomeEquipe + 
							"&psNomeAtividade=" + gsNomeAtividade + "&psNumeroProcessoFormatado=" + gsNumeroProcessoFormatado;
				
					console.log(winbind.name)
					window.open(lsURL, winbind.name);
				}();				
								
				`, janid, janid))
			api.mutex.Unlock()
			// Trace.Println("x")
			api.respostaCh <- true

		case "CLICA_PRA_GERAR_PDF_0":
			pld := mensagem.payload.(*SendJanProc)
			janid := pld.janid
			// processo := pld.processo
			// Trace.Println("x")
			api.mutex.Lock()
			// Trace.Println("x --")
			oleutil.MustCallMethod(
				api.windowJsObj,
				"eval",
				fmt.Sprintf(`window.clica_pra_gerar_pdf("%s");`,
					janid),
			)
			// Trace.Println("x --------------------")
			api.mutex.Unlock()
			api.respostaCh <- true

		case "PAGINA_DOCUMENTOS_CARREGOU_0":
			pld := mensagem.payload.(*SendJanProc)
			janid := pld.janid
			// processo := pld.processo
			api.mutex.Lock()
			// Trace.Println("x")
			res := oleutil.MustCallMethod(
				api.windowJsObj,
				"eval",
				fmt.Sprintf(`window.test_load_pagina_download("%s")`,
					janid),
			)
			api.mutex.Unlock()
			api.respostaCh <- res.Value().(bool)

		case "PEGA_HREF_STRING_OR_NOT_0":
			pld := mensagem.payload.(*SendJanProc)
			janid := pld.janid
			// processo := pld.processo
			api.mutex.Lock()
			// Trace.Println("x")
			res := oleutil.MustCallMethod(
				api.windowJsObj,
				"eval",
				fmt.Sprintf(`window.get_download_href_or_false("%s")`,
					janid),
			)
			api.mutex.Unlock()
			resposta := res.Value()
			// Trace.Printf("\n%T -> %v\n", resposta, resposta)
			api.respostaCh <- resposta

		case "GET_COOKIES_0":
			// pld := mensagem.payload.(*SendJanProc)
			// janid := pld.janid
			// processo := pld.processo
			api.mutex.Lock()
			// Trace.Println("x")
			res := oleutil.MustCallMethod(
				api.windowJsObj,
				"eval",
				`window.document.cookie`,
			)
			api.mutex.Unlock()
			resposta := res.Value()
			// Info.Printf("\n%T -> %v\n", resposta, resposta)
			api.respostaCh <- resposta.(string)

		case "GET_JSON_DATA_0":
			// pld := mensagem.payload.(*SendJanProc)
			// janid := pld.janid
			// processo := pld.processo
			Trace.Println("x getJsonData")
			api.mutex.Lock()
			Trace.Println("x")
			Trace.Printf(`
				api.windowJsObj -> %#v;
				`, api.windowJsObj)
			var res1 *ole.VARIANT
			for {
				Trace.Println(" -- loop -- ")
				res1, err = api.windowJsObj.CallMethod("eval", `window.getJsonData();`)
				if err != nil {
					fmt.Printf("%#v", err)
					time.Sleep(350 * time.Millisecond)
					continue
				}
				break
			}

			// Trace.Printf("%#v \n %T \n", res1, res1)
			api.mutex.Unlock()
			Trace.Println("x")
			resposta := res1.Value()
			Trace.Printf("\n %#v", resposta)
			respSalvar := resposta.(string)
			d1 := []byte(respSalvar)
			e := ioutil.WriteFile("jsonstr.json", d1, 0666)
			if e != nil {
				panic(err)
			}
			// Info.Printf("\n%T -> %v\n", resposta, resposta)
			api.respostaCh <- respSalvar
			Trace.Println("x")

		}
	}
}

func instantiateNewAPIConn() *apiConn {
	Trace.Println("x")
	apInst := apiConn{
		make(map[int]*ole.IDispatch),
		&sync.Mutex{},
		make(chan mensagem),
		make(chan interface{}),
		nil,
		nil,
		nil,
	}
	Trace.Println("x")
	go apInst.olePoolInicio()
	Trace.Println("x")
	return &apInst
}
