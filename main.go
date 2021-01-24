package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/theckman/yacspin"
	"golang.org/x/crypto/ssh"
)

var (
	GlobalPassWord string
	b              bytes.Buffer
	client         *ssh.Client
	config         Config
	session        *ssh.Session
)

func main() {

	remoteHost, userInput, bannerInput := "", "", ""

	hostIdx, commandIdx := -1, -1

	cfg := yacspin.Config{
		Frequency:       200 * time.Millisecond,
		CharSet:         yacspin.CharSets[54],
		Suffix:          "Connecting to host, if unable to connect verify network connection",
		SuffixAutoColon: false,
		Message:         "",
		StopCharacter:   "√",
		StopMessage:     "Completed!",
		StopColors:      []string{"fgGreen"},
		Colors:          []string{"fgYellow"},
	}

	GetConfig()

	printHeader(remoteHost)

	if len(config.Hosts) > 0 {
		// print out the hosts config.yml
		for idx, val := range config.Hosts {
			fmt.Printf("[%d] %s\t", idx+1, val.HostName)
			if idx == len(config.Hosts)-1 {
				fmt.Println("")
			}
		}

		for hostIdx <= 0 || hostIdx > len(config.Hosts) {
			hostIdx, _ = strconv.Atoi(getUserInput("Which host would you like to connect to?", true))
		}
		hostIdx--
	} else {
		hostIdx = 0
	}

	remoteHost = string(config.Hosts[hostIdx].HostName)

	spinner, err := yacspin.New(cfg) // handle the error
	if err != nil {
		panic(err)
	}
	spinner.Start() // Start the spinner

	// connect to remote machine
	client, session = connectViaSsh(config.Hosts[hostIdx].UserName, config.Hosts[hostIdx].Connection, config.Hosts[hostIdx].PassWord)
	spinner.Stop() // connected stop spinner

	printHeader(remoteHost)

	if len(config.Hosts[hostIdx].Commands) > 1 {
		for idx, val := range config.Hosts[hostIdx].Commands {
			fmt.Printf("[%d] %s\n", idx+1, val.Name)
		}
		for commandIdx < 0 || commandIdx > len(config.Hosts[hostIdx].Commands) {
			commandIdx, _ = strconv.Atoi(getUserInput("Which command would you like to run?", true))
		}
		commandIdx--
	} else {
		commandIdx = 0
	}

	printHeader(remoteHost)

	stamp := time.Now().Format(time.RFC3339)
	stamp = strings.ReplaceAll(stamp[0:19], "T", "_") // replace the T in timestamp and slice the string
	stamp = strings.ReplaceAll(stamp, ":", "")        // don't put colons in windows filenames please
	fileName := remoteHost + "_" + stamp + ".txt"     // append the timestamp to filename

	if config.Hosts[hostIdx].Commands[commandIdx].UserInput == true {
		userInput = getUserInput("Enter search term", false)
		bannerInput = userInput
		if config.Hosts[hostIdx].Commands[commandIdx].WhiteSpace == true {
			userInput = string(config.Hosts[hostIdx].Commands[commandIdx].String) + " " + userInput
		} else {
			userInput = string(config.Hosts[hostIdx].Commands[commandIdx].String) + userInput
		}
	} else {
		userInput = string(config.Hosts[hostIdx].Commands[commandIdx].String)
		bannerInput = ""
	}

	banner := color.New(color.FgBlack, color.BgYellow).SprintFunc()
	if bannerInput != "" {
		spinner.Suffix("Searching logs for " + banner(bannerInput) + " may take a minute or two...")
	} else {
		spinner.Suffix(banner("\tSearches may take a minute or two..."))
	}

	spinner.Start() // Start the spinner
	session.Stdout = &b
	session.Run(userInput) // run the commands on the remote host
	spinner.Stop()         // Finished running stop the spinner

	fmt.Println(b.String())

	userInput = ""

	for userInput == "" {
		userInput = getUserInput("Would you like to save output to file? [Y/n]", true)
		if strings.ToUpper(userInput) == "Y" {
			output, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()
			io.WriteString(output, b.String())
			output.Sync()
			fmt.Printf("Results of search saved as %s in current directory\n", banner(fileName))
		}
	}
	fmt.Printf("[Ctrl+c to close]")
	client.Wait()
	client.Close()
}

func getUserInput(prompt string, selection bool) string {

	userInput := ""

	if !selection {
		banner := color.New(color.FgBlack, color.BgYellow).SprintFunc()
		fmt.Printf("%s\n", banner("Default screen buffer size may be too small to view all results. Right click top"))
		fmt.Printf("%s\n", banner("left of window > Properties > Layout > Change Screen Buffer size height to 9999 "))

		banner = color.New(color.FgBlack, color.BgWhite).SprintFunc()
		if prompt != "" {
			fmt.Printf("\nSearches are %s case sensitive\n%s", banner("NOT"), prompt)
		} else {
			fmt.Printf("\nSearches are %s case sensitive: ", banner("NOT"))
		}

	} else {
		fmt.Printf("%s: ", prompt)
	}

	for userInput == "" {
		fmt.Scanf("%s", &userInput)
	}

	return userInput
}
