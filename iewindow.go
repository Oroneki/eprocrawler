package main

import (
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func waitForConditionOnIEWindow(ie *ole.IDispatch, condition string) {
	init := time.Now()
	trace.Println("    waitForConditionOnIEWindow")
	for {
		ok, err := oleutil.CallMethod(ie, "eval", condition)
		if err != nil {
			time.Sleep(250 * time.Millisecond)
			// Trace.Printf("     erro no waiting for condition > \n-----------------\n%s\n-----------------\n", condition)
			// Trace.Printf("     erro                          > \n\n %s \n    ++++ \n %v\n\n", err, err)
			continue
		}
		okk := ok.Value().(bool)
		// Trace.Println("       -")
		// Trace.Printf("     waiting for condition > \n( %s ) === %s\n", condition, ok)
		// Trace.Println("       okk: ", okk)
		if okk {
			break
		}
		// time.Sleep(100 * time.Millisecond)
		if time.Since(init) > 90*time.Second {
			panic("Timeout! waitForConditionOnIEWindow")
		}
	}
	final := time.Since(init)
	time.Sleep(100 * time.Millisecond)
	trace.Printf("    esperou condição por [%s]", final)
}
