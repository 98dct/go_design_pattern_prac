package main

import (
	"fmt"
)

/**

 */

type Computer interface {
	Print()
	SetPrinter(Printer)
}

type Mac struct {
	Printer Printer
}

func (m *Mac) Print() {
	fmt.Println("Print request for Mac")
	m.Printer.PrintFile()
}

func (m *Mac) SetPrinter(p Printer) {
	m.Printer = p
}

type Windows struct {
	Printer Printer
}

func (w *Windows) Print() {
	fmt.Println("Print request for Windows")
	w.Printer.PrintFile()
}

func (w *Windows) SetPrinter(p Printer) {
	w.Printer = p
}

type Printer interface {
	PrintFile()
}

type Canon struct {
}

func (p *Canon) PrintFile() {
	fmt.Println("Printing by a canon Printer")
}

type Lenovo struct {
}

func (l *Lenovo) PrintFile() {
	fmt.Println("Printing by a Lenovo Printer")
}

func main() {
	lenovoPrinter := &Lenovo{}
	canonPrinter := &Canon{}
	macComputer := &Mac{}
	macComputer.SetPrinter(lenovoPrinter)
	macComputer.Print()
	fmt.Println()

	macComputer.SetPrinter(canonPrinter)
	macComputer.Print()
	fmt.Println()

	windowsComputer := Windows{}
	windowsComputer.SetPrinter(lenovoPrinter)
	windowsComputer.Print()
	fmt.Println()

	windowsComputer.SetPrinter(canonPrinter)
	windowsComputer.Print()
	fmt.Println()
}
