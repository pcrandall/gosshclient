package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/pcrandall/figlet4go"
)

var clear map[string]func()

func callClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func printHeader(hostname string) {
	callClear()
	padding := ""
	signature := "pcrandall '21"
	paddingLen := 0

	if hostname == "" {
		hostname = "SSH CLIENT"
	} else {
		hostname = "SSH " + hostname
	}

	ascii := figlet4go.NewAsciiRender()
	// change the font color
	// uncomment to activate colors
	colors := [...]color.Attribute{
		color.FgWhite,
		// color.FgMagenta,
		// color.FgYellow,
		// color.FgCyan,
		// color.FgRed,
		// color.FgBlue,
		// color.FgHiGreen,
		// color.FgGreen,
	}
	options := figlet4go.NewRenderOptions()
	options.FontColor = make([]color.Attribute, len(hostname))
	for i := range options.FontColor {
		options.FontColor[i] = colors[i%len(colors)]
	}

	// you can add more fonts like this if you want
	// ascii.LoadFont("./fonts/bigMoneyNE.flf")
	renderStr, _ := ascii.RenderOpts(hostname, options)

	// calculate the correct padding for the signature 11 is the font height
	var last, longestRow int
	for i := 0; i < len(renderStr)-1; i++ {
		if renderStr[i] == 10 {
			curlongest := i - last
			last = i
			if curlongest > longestRow {
				longestRow = curlongest
			}
		}
	}

	//check if even or odd, add some more padding
	if longestRow%2 == 1 {
		longestRow /= 2
	} else {
		longestRow = (longestRow / 2) + 4
	}

	// TODO fix the padding calc
	// paddingLen = len(renderStr)/(11*2) - len(signature)
	paddingLen = longestRow - len(signature)
	for i := 0; i <= paddingLen; i++ {
		padding += " "
	}
	// remove the last three blank rows, all uppercase chars have a height of 8, the font height for default font is 11
	fmt.Println(renderStr[:len(renderStr)-len(renderStr)/11*3-1])
	// print the signature
	fmt.Printf("%s%s\n", padding, signature)
}
