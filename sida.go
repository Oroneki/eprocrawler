package main

import (
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func SIDAjanelaInit(api *apiConn) bool {
	if api.sidaIE != nil {
		return true
	}
	unknown, _ := oleutil.CreateObject("InternetExplorer.Application")
	ie, _ := unknown.QueryInterface(ole.IID_IDispatch)
	oleutil.CallMethod(ie, "Navigate", "http://www3.pgfn.fazenda/pgfn/milenio/aplicativos2.asp")
	oleutil.PutProperty(ie, "Visible", true)
	for {
		time.Sleep(250 * time.Millisecond)
		if oleutil.MustGetProperty(ie, "Busy").Val == 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	document := oleutil.MustGetProperty(ie, "document").ToIDispatch()
	window := oleutil.MustGetProperty(document, "parentWindow").ToIDispatch()
	time.Sleep(3000 * time.Millisecond)
	for {
		logado := SIDAcheckLogado(window)
		Trace.Println("logado: ", logado)
		if logado {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	api.sidaIE = ie
	api.sidaIEWindow = window
	return true
}

func SIDAcheckLogado(ieWindow *ole.IDispatch) bool {

	Trace.Println("--------")
	variant := oleutil.MustCallMethod(ieWindow, "eval", `window.location.pathname === "/PGFN/Milenio/PrincipalFrames.asp";`)
	Trace.Println("--------")
	return variant.Value().(bool)

}

func SIDAVaiPraConsulta(api *apiConn, processo string) string {
	Trace.Println("--------")
	oleutil.CallMethod(api.sidaIEWindow, "Navigate", "http://www3.pgfn.fazenda/PGFN/Divida/Consulta/Inscricao/Cons11.asp")
	Trace.Println("--------")
	time.Sleep(50 * time.Millisecond)
	Trace.Println("--------")
	WaitIEWindow(api.sidaIE)
	Trace.Println("--------")
	for {
		okk_, err := oleutil.CallMethod(api.sidaIEWindow, "eval", `window.location.pathname === "/PGFN/Divida/Consulta/Inscricao/Cons11.asp";`)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			Trace.Println("- :( -")
			continue
		}
		okk := okk_.Value().(bool)
		Trace.Println("-")
		Trace.Println("okk: ", okk)
		if okk {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	Trace.Println("\nCarregou consulta pro processo: ", processo, "\n")
	Trace.Println("--------")

	oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementById("op_numProcAntigo").click();`)
	Trace.Println("--------")
	time.Sleep(100 * time.Millisecond)
	Trace.Println("--------")
	oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByTagName("input")[12].value = "`+processo+`";`)
	Trace.Println("--------")

	valAntes := oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByTagName("input")[20].value;`).Value().(string)
	Trace.Println("--------")

	Trace.Println("\nVal antes: ", valAntes, "\n")
	Trace.Println("--------")

	oleutil.MustCallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByName("ok")[0].click();`)
	Trace.Println("--------")
	time.Sleep(50 * time.Millisecond)
	Trace.Println("--------")
	WaitIEWindow(api.sidaIE)
	Trace.Println("--------")

	keyValues := oleutil.MustCallMethod(api.sidaIEWindow, "eval", SidaKeyValuesConsulta).Value().(string)
	Trace.Println("\n\nResultado:\n", keyValues)

	valDepois, err := oleutil.CallMethod(api.sidaIEWindow, "eval", `window.document.getElementsByTagName("input")[20].value;`)
	if err != nil {
		Trace.Println("\n\nERRO NO 'VAL DEPOIS'... ISSO EH BOM ???:\n")
		return keyValues
	}
	valDepois_ := valDepois.Value().(string)
	return keyValues + "_CAMPO_20_||>" + valDepois_ + "\n"
}

func WaitIEWindow(ie *ole.IDispatch) {
	for {
		time.Sleep(250 * time.Millisecond)
		Trace.Println("-")
		if oleutil.MustGetProperty(ie, "Busy").Val == 0 {
			break
		}
	}
	time.Sleep(100 * time.Millisecond)
}
