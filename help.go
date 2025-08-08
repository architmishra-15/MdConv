package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	RESET     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	White     = "\033[97m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Strike    = "\033[9m"

	BrightRed    = "\033[91m"
	BrightGreen  = "\033[92m"
	BrightYellow = "\033[93m"
	BrightBlue   = "\033[94m"
	BrightPurple = "\033[95m"
	BrightCyan   = "\033[96m"
	BrightWhite  = "\033[97m"
)

func VersionInfo(version string) {

	msg := fmt.Sprintf("%s%sVersion:%s  %s", Bold, Green, RESET, version)
	width := len(msg)

	top := "╭" + strings.Repeat("─", width-2) + "╮"
	middle := fmt.Sprintf("│  %s         │", msg)
	bottom := "╰" + strings.Repeat("─", width-2) + "╯"

	fmt.Println(top)
	fmt.Println(middle)
	fmt.Println(bottom)

}

func Help() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Println()
	fmt.Println(Bold + Green + Underline + "Usage:\033[0m")
	fmt.Printf("%7s", "")
	fmt.Println(Bold + "md2html " + RESET + "\033[4m<input md>\033[0m" + " \033[4m<output.html>\033[0m" + Dim + " [OPTIONAL]\033[0m")

	fmt.Println()

	fmt.Printf("%s%s%sOptions:\n%s", Bold, Underline, BrightGreen, RESET)
	fmt.Printf("\n")
	fmt.Printf("%s%sGlobal Options:\n%s", Bold, BrightCyan, RESET)
	fmt.Fprintf(w, "%s%s\thelp, --help, -h%s\t:\tShow this help message\n", Bold, Blue, RESET)
	fmt.Fprintf(w, "%s%s\tversion, -v, --version%s\t:\tShow version\n", Bold, Blue, RESET)

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s%s%sFlags:%s\n", Bold, Underline, BrightGreen, RESET)
	fmt.Fprintf(w, "\t%s%s--style=<name>%s\t:\tProvide name of custom syntax styling\n", Bold, Blue, RESET)
	fmt.Fprintf(w, "\t%s%s--title=<name>%s\t:\tGives the title to the HTML generated\n", Bold, Blue, RESET)
	fmt.Fprint(w, "\t\033[1m\033[34m--bg-black\033[0m\t:\tPass this flag to use generate webpage with dark background\n")

	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "%s%s%sExample:%s\n", Bold, Underline, BrightGreen, RESET)
	fmt.Fprintf(w, "%s\tmd2html input.md output.html%s\n", Yellow, RESET)
	fmt.Fprintf(w, "%s\tmd2html input.md --style=onedark%s\n", Yellow, RESET)
	fmt.Fprintf(w, "%s\tmd2html input.md output.html --style=monokai --title='Test Document'%s\n", Yellow, RESET)
	fmt.Fprintf(w, "%s\tmd2html input.md --title='Test Documen' -bg-black%s\n", Yellow, RESET)

	fmt.Fprintf(w, "%s%sAuthor:%s\t%s%s%sArchit Mishra\n%s", BrightWhite, Bold, RESET, Bold, Italic, BrightCyan, RESET)
	fmt.Fprintf(w, "%s%sVersion:%s\t%s%s%s%s%s\n", BrightWhite, Bold, RESET, Bold, Italic, BrightCyan, CurrentVersion, RESET)

}
