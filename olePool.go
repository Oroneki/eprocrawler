package main

import (
	"fmt"
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

type WindowIdentification int32

const (
	EPROCESSO_WINDOW WindowIdentification = 0
	SIDA_WINDOW      WindowIdentification = 1
)

type apiConn struct {
	linksMap     map[int]*ole.IDispatch
	mutex        *sync.Mutex
	perguntaCh   chan mensagem
	respostaCh   chan interface{}
	window       *ole.IDispatch
	windowJsObj  *ole.IDispatch
	windowObj    *ole.IDispatch
	sidaIE       *ole.IDispatch
	sidaIEWindow *ole.IDispatch
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
	trace.Println("x")
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
	trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "PATCHWINDOWPRINCIPAL",
		payload: nil,
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) abreProcesso(janID string, processo *Processo) bool {
	trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "ABREPROCESSO_0",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) paginaProcessoCarregou(janID string, processo *Processo) bool {
	trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "TESTA_PAGINAPROCESSOCARREGOU",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) paginaProcessoPatcheVaiProDownload(janID string, processo *Processo) bool {
	trace.Println("x")
	api.perguntaCh <- mensagem{
		tipo:    "TESTA_PAGINAPROCESSOPATCH_VAI_PRO_DOWNLOAD",
		payload: &SendJanProc{janID, processo},
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) SIDAInit() bool {
	trace.Println("x SIDAInit")
	api.perguntaCh <- mensagem{
		tipo:    "SIDA_INIT_0",
		payload: nil,
	}
	resposta := <-api.respostaCh
	return resposta.(bool)
}

func (api *apiConn) SIDAConsultaProcesso(processo string) *grabConsultaProcessoSidaResult {
	trace.Println("x SIDAInit")
	api.perguntaCh <- mensagem{
		tipo:    "SIDA_CONSULTA_PROCESSO_0",
		payload: processo,
	}
	resposta := <-api.respostaCh
	return resposta.(*grabConsultaProcessoSidaResult)

}

func (api *apiConn) clicaParaGerarPDF(janID string, processo *Processo) bool {
	trace.Println("x")
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

func (api *apiConn) grabSidaWindow() bool {
	api.perguntaCh <- mensagem{
		tipo:    "SIDA_DEZAJUIZA_GRAB_WINDOW",
		payload: nil,
	}
	<-api.respostaCh
	return true
}

func (api *apiConn) waitForCondition(window WindowIdentification, condition string) bool {
	trace.Println("x api method...")
	var msg mensagem
	if window == SIDA_WINDOW {
		trace.Println("x tipo de mensagem eh sida")
		msg = mensagem{
			tipo:    "WAIT_FOR_CONDITION_ON_SIDA_WINDOW",
			payload: condition,
		}
	} else if window == EPROCESSO_WINDOW {
		msg = mensagem{
			tipo:    "WAIT_FOR_CONDITION_ON_EPROCESSO_WINDOW",
			payload: condition,
		}
	} else {
		trace.Println("x erro -- nao existe essa window")
		return false
	}
	trace.Println("x enviar mensagem")
	api.perguntaCh <- msg
	trace.Println("x menagem voltou, mandar resposta...")
	<-api.respostaCh
	trace.Println("x foi resposta...")
	time.Sleep(100 * time.Millisecond)
	return true
}

func (api *apiConn) evalOnWindow(codeStr []byte) []byte {
	api.perguntaCh <- mensagem{
		tipo:    "EVAL_CODE",
		payload: string(codeStr),
	}
	resposta := <-api.respostaCh
	info.Printf("evalOnWindow resp %T ||| %v", resposta, resposta)
	return []byte(resposta.(string))
}

func (api *apiConn) evalOnSidaWindow(codeStr []byte) []byte {
	api.perguntaCh <- mensagem{
		tipo:    "EVAL_CODE_SIDA_WINDOW",
		payload: string(codeStr),
	}
	resposta := <-api.respostaCh
	info.Printf("evalOnSIDAWindow resp %T ||| %v", resposta, resposta)
	return []byte(resposta.(string))
}

func (api *apiConn) waitNotBusySidaWindow(window WindowIdentification) bool {
	trace.Printf("api.waitNotBusySidaWindow: args: %v", window)
	api.perguntaCh <- mensagem{
		tipo:    "WAIT_NOT_BUSY",
		payload: window,
	}
	resposta := <-api.respostaCh
	info.Printf("evalOnSIDAWindow resp %T ||| %v", resposta, resposta)
	return true
}

func (api *apiConn) getInscricoesFromSidaMulti() string {
	trace.Println("     getInscrições")
	api.perguntaCh <- mensagem{
		tipo:    "GET_INSC_FROM_SIDA_MULTI",
		payload: nil,
	}
	resposta := <-api.respostaCh
	info.Printf("resposta json stringifado %T --> %v", resposta, resposta)
	return resposta.(string)
}

// methods bootstrap ---------------------------------------------

func (api *apiConn) olePoolInicio() {
	trace.Println("x Inicio Pool")
	runtime.LockOSThread()
	trace.Println("x")
	err := ole.CoInitialize(0)
	trace.Println("x")
	if err != nil {
		trace.Println("x Erro")
		oleerr := err.(*ole.OleError)
		// S_FALSE           = 0x00000001 // CoInitializeEx was already called on this thread
		if oleerr.Code() != ole.S_OK && oleerr.Code() != 0x00000001 {
			info.Println(err)
			trace.Println(err)
		}
	} else {
		// Only invoke CoUninitialize if the thread was not initizlied before.
		// This will allow other go packages based on go-ole play along
		// with this library.
		trace.Println("x Tranquilo")
		defer ole.CoUninitialize()
	}

	var regexHrefLinkProcesso = regexp.MustCompile(`\'.*?\'`)
	for mensagem := range api.perguntaCh {
		trace.Printf(`
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
			trace.Println("x")
			api.mutex.Lock()
			trace.Println("x")
			unknown, e := oleutil.CreateObject("shell.Application")
			if e != nil {
				trace.Println("shell.Application não criada")
			}
			trace.Println("x")
			shell, e := unknown.QueryInterface(ole.IID_IDispatch)
			if e != nil {
				trace.Println("query interface falha")
			}
			trace.Println("x")
			windows, e := shell.CallMethod("Windows")
			if e != nil {
				trace.Println("metodo windows nao possivel de ser chamado...")
			}
			trace.Println("x")
			wins := windows.ToIDispatch()
			trace.Println("x")
			nois, _ := wins.GetProperty("Count")
			trace.Println("x")
			valConta := int(nois.Val)
			trace.Printf("\n %d janelas identificadas.", valConta)
			var re = regexp.MustCompile(`eprocesso\.suiterfb\.receita\.fazenda\/ControleAcessarCaixaTrabalho.*?apresentarPagina`)
			trace.Println("x")
			var itemjanela *ole.IDispatch

			for i := 0; i < valConta; i++ {
				trace.Println("\n----\nitem", i)
				item, e := wins.CallMethod("Item", i)
				if e != nil {
					trace.Printf("\nitem %d miow\n--- continue ---", i)
					continue
				}
				trace.Println(" o")
				itemd := item.ToIDispatch()
				trace.Printf(" \n            o    %#v", itemd)
				locationURLV, e := itemd.GetProperty("LocationURL")
				trace.Printf(" \n----------------o    %#v", itemd)
				if e != nil {
					trace.Println("janela sem LocationURL")
					i++
					continue
				}
				trace.Println(" item ", i, " URL ->", locationURLV)
				urlV := locationURLV.Value()
				trace.Println(" o")
				url := urlV.(string)
				trace.Println(" o")
				trace.Printf("\nJanela Identificada: (id: %d) %s", i, url)
				trace.Println(" o")

				testeRegex := re.MatchString(url)
				trace.Println(" o")
				if testeRegex {
					trace.Println(" o!")
					trace.Printf(`


	+++++++++++++++++++++++++++++
	++ IDENTIFICADA PELO REGEX ++
	+++++++++++++++++++++++++++++

	E-PROCESSO : (i: %d) 
		URL: %s
		

		`, i, url)
					itemjanela = itemd
					trace.Println(" o!")
					// break
				}
			}
			busy, e := oleutil.GetProperty(itemjanela, "Busy")
			if e != nil {
				trace.Println("busy nao deu certo")

			}
			container, e := oleutil.GetProperty(itemjanela, "Container")
			if e != nil {
				trace.Println("busy nao deu certo")

			}
			application, e := oleutil.GetProperty(itemjanela, "Application")
			if e != nil {
				trace.Println("busy nao deu certo")

			}
			info.Printf(`Janela Internet Explorer identificada: HWND %v Busy: %v`,
				oleutil.MustGetProperty(itemjanela, "HWND").Value(),
				busy.Value(),
			)
			trace.Printf(`
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

			trace.Println("x")
			api.window = itemjanela
			trace.Println("x")
			api.mutex.Unlock()
			api.respostaCh <- true

		case "PATCHWINDOWPRINCIPAL":
			api.mutex.Lock()
			trace.Println("x")
			iePrincipalD := api.window
			trace.Println("x iePrincipalD -->", iePrincipalD)
			iePrincipalDocumentV, _ := iePrincipalD.GetProperty("Document")
			trace.Println("x")
			iePrincipalDocumentD := iePrincipalDocumentV.ToIDispatch()
			trace.Println("x")
			title := oleutil.MustGetProperty(iePrincipalDocumentD, "title").ToIDispatch()
			trace.Println("  title: ", title)
			windowPrincipal := oleutil.MustGetProperty(iePrincipalDocumentD, "parentWindow").ToIDispatch()
			trace.Println("x")
			api.windowJsObj = windowPrincipal
			variant, err := oleutil.CallMethod(windowPrincipal, "eval", `_____OWNED____`)
			if err != nil {
				trace.Printf("[X] nao tem owned. PATCH WINDOW!!! ")
			} else {
				trace.Printf("Janela já ta OWNED, seguir e liberar o mytex e o canal")
				api.mutex.Unlock()
				trace.Println("x")
				api.respostaCh <- true
				trace.Println("x")
				break

			}
			trace.Println("x")
			variantVal := variant.Value()
			trace.Println("variantVal", variantVal)
			oleutil.MustCallMethod(windowPrincipal, "eval", jJSconsole)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", `console.log('owned->', window._____OWNED____);`)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", `window.oro_obj = {};`)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", JShackedobj)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", JSabrirJanela)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", JSpatchInicial)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", JStestLoadPagina)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", JSUnicodeHandle)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", JSgetJsonData)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", JSjqueryStringify)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", `window._____OWNED____ = true;`)
			trace.Println("x")
			oleutil.MustCallMethod(windowPrincipal, "eval", `console.log('owned->', window._____OWNED____);`)
			// POLYFILLS --------------------------
			// querySelector e QuerySelectorAll
			oleutil.MustCallMethod(windowPrincipal, "eval", jsPolyfills)
			trace.Println("x")
			api.mutex.Unlock()
			trace.Println("x")

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
						trace.Printf("\n + Microrotina encaminhou processo %v pro canal.\n", p)
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
				trace.Printf("\n%s - ERRO na chamada da pagina de carregamento do processo\n", janid)
			} else {
				resp := areaNodeList.Value().(bool)
				api.respostaCh <- resp
				trace.Printf("\n\n%s - Resposta do carregamento da pagina do processo: %v\n\n", janid, resp)
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
				fmt.Sprintf(JSpaginaProcessoPatch, janid, janid))
			api.mutex.Unlock()
			// Trace.Println("x")
			api.respostaCh <- true

		case "EVAL_CODE":
			pld := mensagem.payload.(string)
			api.mutex.Lock()
			res, err := oleutil.CallMethod(
				api.windowJsObj,
				"eval",
				pld,
			)
			if err != nil {
				trace.Println("opa... erro no eval.")

			}
			trace.Printf("res -> %v", res)

			api.mutex.Unlock()

			api.respostaCh <- "respos"

		case "EVAL_CODE_SIDA_WINDOW":
			pld := mensagem.payload.(string)
			api.mutex.Lock()
			res, err := oleutil.CallMethod(
				api.sidaIEWindow,
				"eval",
				pld,
			)
			if err != nil {
				trace.Println("sida opa... erro no eval.")

			}
			trace.Printf("sida res -> %v", res)

			api.mutex.Unlock()

			api.respostaCh <- "sida_respos"

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
			trace.Printf("-")
			pld := mensagem.payload.(*SendJanProc)
			trace.Printf("-")
			janid := pld.janid
			trace.Printf("-")
			// processo := pld.processo
			trace.Printf("-")

			api.mutex.Lock()
			// Trace.Println("x")
			trace.Printf("-")

			res := oleutil.MustCallMethod(
				api.windowJsObj,
				"eval",
				fmt.Sprintf(`window.get_download_href_or_false("%s")`,
					janid),
			)
			trace.Printf("-")

			api.mutex.Unlock()
			trace.Printf("-")

			resposta := res.Value()
			trace.Printf("-")

			// Trace.Printf("\n%T -> %v\n", resposta, resposta)
			trace.Printf("-")

			api.respostaCh <- resposta
			trace.Printf("-")

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
			trace.Println("x getJsonData")
			api.mutex.Lock()
			trace.Println("x")
			trace.Printf(`
				api.windowJsObj -> %#v;
				`, api.windowJsObj)
			var res1 *ole.VARIANT
			for {
				trace.Println(" -- loop -- ")
				res1, err = api.windowJsObj.CallMethod("eval", `window.getJsonData();`)
				if err != nil {
					trace.Printf("%#v", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}
				break
			}

			// Trace.Printf("%#v \n %T \n", res1, res1)
			api.mutex.Unlock()
			trace.Println("x")
			resposta := res1.Value()
			trace.Printf("\n %#v", resposta)
			respSalvar := resposta.(string)
			// d1 := []byte(respSalvar)
			// e := ioutil.WriteFile("jsonstr.json", d1, 0666)
			// if e != nil {
			// 	panic(err)
			// }
			// Info.Printf("\n%T -> %v\n", resposta, resposta)
			api.respostaCh <- respSalvar
			trace.Println("x")

		case "SIDA_INIT_0":
			trace.Println("x Iniciando SIDA...")
			api.mutex.Lock()
			api.respostaCh <- sSIDAjanelaInit(api)
			api.mutex.Unlock()

		case "SIDA_CONSULTA_PROCESSO_0":
			trace.Println("x Iniciando Consulta por Processo no SIDA...")
			api.mutex.Lock()
			sidaVaiPraPaginaResultadoConsulta(api, mensagem.payload.(string))
			// api.respostaCh <- "grabConsultaInfo(api)"
			respostaGrab := grabConsultaInfo(api, mensagem.payload.(string))
			trace.Printf("respostaGrab : %s", respostaGrab)
			api.respostaCh <- respostaGrab
			api.mutex.Unlock()

		case "WAIT_FOR_CONDITION_ON_SIDA_WINDOW":
			trace.Println("x waitf for sida condition")
			api.mutex.Lock()
			trace.Println("x lock mutex")
			waitForConditionOnIEWindow(api.sidaIEWindow, mensagem.payload.(string))
			api.mutex.Unlock()
			api.respostaCh <- true

		case "WAIT_FOR_CONDITION_ON_EPROCESSO_WINDOW":
			trace.Println("x waitf for eproc condition")
			api.mutex.Lock()
			waitForConditionOnIEWindow(api.windowObj, mensagem.payload.(string))
			api.mutex.Unlock()
			api.respostaCh <- true

		case "WAIT_NOT_BUSY":
			trace.Println("x wait not busy")
			api.mutex.Lock()
			//esperar todas :)
			WaitIEWindow(api.sidaIE)
			api.mutex.Unlock()
			api.respostaCh <- true

		case "GET_INSC_FROM_SIDA_MULTI":
			trace.Println("x GET_INSC_FROM_SIDA_MULTI")
			api.mutex.Lock()
			res, err := api.sidaIEWindow.CallMethod(
				"eval",
				`(function () {
					return stringify_INJECTED();
				  })();`,
			)
			if err != nil {
				trace.Println("x ERRO NA CHAMADA - GET_INSC_FROM_SIDA_MULTI")
				trace.Printf("x ERRO %v", err)
				api.respostaCh <- "erro"
				panic(err)
			} else {
				trace.Printf("res --> %v | %v", res, res.Val)
				api.respostaCh <- res.Value().(string) // string!
			}
			api.mutex.Unlock()

		case "SIDA_DEZAJUIZA_GRAB_WINDOW":
			trace.Println("x")
			api.mutex.Lock()
			trace.Println("x")
			unknown, _ := oleutil.CreateObject("shell.Application")
			trace.Println("x")
			shell, _ := unknown.QueryInterface(ole.IID_IDispatch)
			trace.Println("x")
			windows, _ := shell.CallMethod("Windows")
			trace.Println("x")
			wins := windows.ToIDispatch()
			trace.Println("x")
			nois, _ := wins.GetProperty("Count")
			trace.Println("x")
			valConta := int(nois.Val)
			trace.Printf("\n %d janelas identificadas.", valConta)
			var re = regexp.MustCompile(`www\d?\.pgfn\.fazenda\/PGFN\/Milenio\/PrincipalFrames\.asp`)
			trace.Println("x")
			var itemjanela *ole.IDispatch

			for i := 0; i < valConta; i++ {
				trace.Println("\n-------------------------\n\nitem ", i, "\n\n")
				item, e := wins.CallMethod("Item", i)
				if e != nil {
					trace.Printf("\nitem %d miow\n--- continue ---", i)
					continue
				}
				trace.Println(" o")
				itemd := item.ToIDispatch()
				trace.Printf(" \n            o    %#v", itemd)
				locationURLV, err := itemd.GetProperty("LocationURL")
				if err != nil {
					trace.Printf(" \n erro ao pegar url... continuar")
					continue
				}
				trace.Println(" item ", i, " URL ->", locationURLV)
				urlV := locationURLV.Value()
				trace.Println(" o")
				url := urlV.(string)
				trace.Println(" o")
				trace.Printf("\nJanela Identificada: (id: %d) %s", i, url)
				trace.Println(" o")

				testeRegex := re.MatchString(url)
				trace.Println(" o")
				if testeRegex {
					trace.Println(" o!")
					trace.Printf(`


	+++++++++++++++++++++++++++++
	++ IDENTIFICADA PELO REGEX ++
	+++++++++++++++++++++++++++++

	E-PROCESSO : (i: %d) 
		URL: %s
		

		`, i, url)
					itemjanela = itemd
					trace.Println(" o!")
					// break
				}
			}
			busy := oleutil.MustGetProperty(itemjanela, "Busy")
			container := oleutil.MustGetProperty(itemjanela, "Container")
			application := oleutil.MustGetProperty(itemjanela, "Application")
			info.Printf(`Janela Internet Explorer identificada: HWND %v Busy: %v`,
				oleutil.MustGetProperty(itemjanela, "HWND").Value(),
				busy.Value(),
			)
			trace.Printf(`
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

			trace.Println("x")
			api.window = itemjanela
			trace.Println("x")
			api.mutex.Unlock()
			api.respostaCh <- true

		}

	}
}

func instantiateNewAPIConn() *apiConn {
	trace.Println("x")
	apInst := apiConn{
		make(map[int]*ole.IDispatch),
		&sync.Mutex{},
		make(chan mensagem),
		make(chan interface{}),
		nil,
		nil,
		nil,
		nil,
		nil,
	}
	trace.Println("x")
	go apInst.olePoolInicio()
	trace.Println("x")
	return &apInst
}
