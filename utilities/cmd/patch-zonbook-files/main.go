package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config defines the configuration values from config.json
type Config struct {
	Path        string `json:"Path"`
	Language    string `json:"Language"`
	DebugString string `json:"Debug"`
	Debug       bool
}

var configFileName = "config.json"

var globalConfig Config

func populateConfiguration() error {
	content, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	text := string(content)

	err = json.Unmarshal([]byte(text), &globalConfig)
	if err != nil {
		return err
	}

	globalConfig.Debug = globalConfig.DebugString == "true"

	return nil
}

func debugPrint(debug bool, s string) {
	if debug {
		fmt.Println(s)
	}
}

// Utility functions to get the first and last N chars from a string
func getFirstN(s string, n int) string {
	if len(s) < n {
		return ""

	}

	return s[0:n]
}

func getLastN(s string, n int) string {
	if len(s) < n {
		return ""

	}

	return s[len(s)-n:]
}

func stripFirstN(s string, n int) string {
	if len(s) < n {
		return ""

	}

	return getLastN(s, len(s)-n)
}

func stripLastN(s string, n int) string {
	if len(s) < n {
		return ""

	}

	return getFirstN(s, len(s)-n)
}

/*
	Returns an empty string if text does not match:
	  <!ENTITY % phrases-go SYSTEM "../../shared/go.ent">
      <!ENTITY % phrases-go SYSTEM "../../shared/go.ent">
*/
func patchEntRef(debug bool, text string, lang string, depth int) string {
	entString := "<!ENTITY % phrases-" + lang + " SYSTEM \"../../shared/" + lang + ".ent\">"

	text = strings.Trim(text, " \t")

	if text != entString {
		return ""
	}

	// Split string by "
	// The string we want to change has the following parts:
	//   "<!ENTITY % phrases-go SYSTEM "
	//   "../../shared/go.ent"
	//   ">"
	// so make sure there are three parts
	parts := strings.Split(text, "\"")

	if len(parts) != 3 {
		return ""
	}

	// Now build up the tag
	tag := "<!ENTITY % phrases-" + lang + " SYSTEM \""

	i := 0
	debugPrint(debug, "Adding "+strconv.Itoa(depth)+" redirects to entity file")
	for i < depth {
		tag += "../"
		i++
	}

	return tag + "shared/" + lang + ".ent\">"
}

// Determines whether text ends with:
//   </title></info>
func isFullTitle(debug bool, text string) bool {
	// Trim off any leading or trailing white space
	text = strings.Trim(text, " \t")
	endChars := "</title></info>"

	// Is the string long enough?
	minLen := len("<info><title>." + endChars)

	if len(text) < minLen {
		debugPrint(debug, "Text is not long enough to be a full title tag")
		return false
	}

	lastCharsLength := len(endChars)

	// Does it end with </info</title>
	lastChars := getLastN(text, lastCharsLength)

	if lastChars != endChars {
		debugPrint(debug, "Text does not end with "+endChars)
		return false
	}

	return true
}

/*
   Returns "" if it's not a title tag;
   i.e, does NOT start with:
     <info><title>

   Returns string, where string is:
     <info><title id="SECTION-ID.title">...
   if we have part of a title, such as:
	 <info><title>...
   OR
	 <info><title id="SECTION-ID.title"...</info></title>
   if we have an entire title.
*/
func getTitle(debug bool, text string, sectionID string) string {
	// Trim off any leading white space
	text = strings.TrimLeft(text, " \t")

	// It's not a title tag if the line isn't at least as long as:
	preamble := "<info><title>"
	minLen := len(preamble)

	if len(text) < minLen {
		return ""
	}

	firstTags := getFirstN(text, minLen)

	if firstTags != preamble {
		// It's not at least <info><title>
		return ""
	}

	debugPrint(debug, "Original title tag:")
	debugPrint(debug, "  "+text)

	/* text that has a title tag is something like:
		 <info><title>AssumeRolev2.go</title></info>
		 <info><title>TEXT
	  so check for the first case by calling isFullTitle().
	*/

	fullTitle := isFullTitle(debug, text)
	if fullTitle {
		debugPrint(debug, "Full title")
	} else {
		debugPrint(debug, "NOT a full title")
	}

	// Strip off leading tags
	titleText := stripFirstN(text, len(preamble))

	/* Should now be something like:
		 TEXT</info></title>
	   OR
		 TEXT
	   so strip off traling "</info></title>"
	   (we'll restore them when we return if fullTitle)
	*/
	if fullTitle {
		titleText = stripLastN(titleText, len("</info></title>"))
	}

	debugPrint(debug, "Title text before removing suffix:")
	debugPrint(debug, titleText)

	// If we have a full title
	if fullTitle {
		return "<info><title id=\"" + sectionID + ".title\">" + titleText + "</title></info>"
	}

	// titleText did NOT end with "</title></title>"
	return "<info><title id=\"" + sectionID + ".title\">" + titleText
}

// If text is a section tag, such as:
//   <section id="anythingv2.go">
// and lang == "go"
// return "anythingv2"
func getSectionID(debug bool, text string, lang string) string {
	// Trim off any potential leading or trailing white space
	text = strings.Trim(text, " \t")

	// Split by quotes, so we should get something like:
	//   <section id=
	//   anythingv2.go
	//   >
	parts := strings.Split(text, "\"")

	if len(parts) != 3 {
		return ""
	}

	if parts[0] != "<section id=" {
		return ""
	}

	if lang != "" {
		// Lop off .LANG if it exists
		lastPart := "." + lang
		lastLen := len(lastPart)
		tail := getLastN(parts[1], lastLen)

		if tail == lastPart {
			text = stripLastN(parts[1], lastLen)
		} else {
			text = parts[1]
		}
	}

	// Do something similar if the section ID ends with a period
	lastPart := "."
	lastLen := len(lastPart)
	tail := getLastN(parts[1], lastLen)

	if tail == lastPart {
		text = stripLastN(parts[1], lastLen)
	}

	return text
}

func patchFiles(debug bool, sourcePath string, lang string) error {
	srcParts := strings.Split(sourcePath, "\\")
	srcLength := len(srcParts) // Should be 3 for C:\test\test
	debugPrint(debug, "Depth is "+strconv.Itoa(srcLength)+" for "+sourcePath)

	// srcParts := strings.Split(sourcePath, "/")
	// srcLength := len(srcParts)
	// Get list of files and directories
	// and call patchFile for every file.
	err := filepath.Walk(sourcePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Only parse XML files
			// First get filename from path
			fileParts := strings.Split(path, "\\")
			// we get: {D:, src, aws-cdk-examples, file.xml}
			// Should be 5 for C:\test\test\Dynamodb\README.xml
			// and 6 for C:\test\test\dynamodb\DescribeTable\DescribeTablev2.xml
			length := len(fileParts)

			debugPrint(debug, "Depth is "+strconv.Itoa(length)+" for "+path)
			depth := length - srcLength - 3

			fileName := fileParts[length-1] // Now we have file.ext, so does ext == xml?

			parts := strings.Split(fileName, ".") // gives us file and xml

			if len(parts) < 2 || parts[1] != "xml" {
				return nil
			}

			// If it's an xml file, call patchFile
			err = patchFile(debug, path, depth, lang)

			return err
		})
	if err != nil {
		return err
	}

	return nil
}

func myWrite(debug bool, w bufio.Writer, s string) {
	_, err := w.WriteString(s + "\n")
	if err != nil {
		panic(err)
	}

	w.Flush()
}

// Patches the section ID and creates the associated title ID in an XML (Zonbook) file
// Name is going to be something like:
//    CopyObject.xml

// Patches the section ID and creates the associated title ID in an XML (Zonbook) file
// Name is going to be something like:
//    CopyObject.xml

// Returns true if we created patch file name.lmx
func patchFile(debug bool, name string, depth int, lang string) error {
	fmt.Println("Parsing " + name)

	parts := strings.Split(name, ".") // gives us filename and xml
	// make sure it has two parts
	if len(parts) != 2 {
		msg := "The filename " + name + " did not match NAME.EXTENSION"
		return errors.New(msg)
	}

	// Get contents of file
	inFile, err := os.Open(name)
	if err != nil {
		msg := "Got error opening " + name
		return errors.New(msg)
	}

	defer inFile.Close()

	// Create topic file in that directory
	fmt.Println("Creating new file " + parts[0] + ".lmx")
	f, err := os.Create(parts[0] + ".lmx")
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	lineNumber := 0

	// If not empty, we have found a section tag
	// (from call to getSectionID)
	sectionID := ""
	rootSectionID := ""

	// If not empty, we have at least part of a title
	// (from call to getTitle)
	patchedTitle := ""

	scanner := bufio.NewScanner(inFile)

	for scanner.Scan() {
		text := scanner.Text()
		lineNumber++

		debugPrint(debug, strconv.Itoa(lineNumber)+": "+text)

		// Skip empty lines
		if len(text) == 0 {
			myWrite(debug, *w, text)
			continue
		}

		// Do we already have part of a title?
		if patchedTitle != "" {
			// Just append the next line
			patchedTitle = strings.TrimRight(patchedTitle, " ")
			text = strings.TrimLeft(text, " ")
			debugPrint(debug, "")
			debugPrint(debug, "Created new title from two lines: ")
			debugPrint(debug, "  "+patchedTitle+" "+text)
			debugPrint(debug, "")

			myWrite(debug, *w, patchedTitle+" "+text)

			sectionID = ""
			patchedTitle = ""
			continue
		}

		// If we aren't in a section, try to find one so we can get its ID and save it for the title ID
		if sectionID == "" {
			// If input line contains a reference to the enitity file,
			// fix that reference.
			fixedEntRef := patchEntRef(debug, text, lang, depth)

			if fixedEntRef != "" {
				myWrite(debug, *w, fixedEntRef)
				continue
			}

			sectionID = getSectionID(debug, text, lang)

			// If section ID == "", it's not a section, so just pass it through
			if sectionID == "" {
				myWrite(debug, *w, text)
				continue
			}

			// Stash first section ID as we need it for later section IDs
			if rootSectionID == "" {
				rootSectionID = sectionID
			} else {
				// Do we want to do this for all subsequent section IDs?
				sectionID = rootSectionID + "-" + sectionID
			}

			sectionTag := "<section id=\"" + sectionID + "\">"

			debugPrint(debug, "Original section tag:")
			debugPrint(debug, "  "+text)

			// We have a section ID, so barf out the section tag with that ID
			debugPrint(debug, "New section tag: ")
			debugPrint(debug, "  "+sectionTag)

			myWrite(debug, *w, sectionTag)

			continue
		}

		// We are in a section, so look for a title tag
		patchedTitle = getTitle(debug, text, sectionID)

		// We just pass it through as it's not a title if patchedTitle == ""
		if patchedTitle == "" {
			myWrite(debug, *w, text)
			continue
		}

		fullTitle := isFullTitle(debug, patchedTitle)

		if fullTitle {
			// We have a complete title, so no need to look for another until we hit a new section tag
			debugPrint(debug, "New title with ID: ")
			debugPrint(debug, "  "+patchedTitle)

			myWrite(debug, *w, patchedTitle)

			sectionID = ""
			patchedTitle = ""
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		msg := "Got error scanning " + name
		return errors.New(msg)
	}

	w.Flush()
	return nil
}

// Updates a Zonbook XML file so the title has an ID (section ID + ".title")
// If a -l flag is supplied, lops off the .lang from the section ID and title value
func main() {
	err := populateConfiguration()
	if err != nil {
		globalConfig.Path = ""
		globalConfig.Language = ""
		globalConfig.Debug = false
	}
	// Process command line args
	pathPtr := flag.String("p", globalConfig.Path, "Fully-qualified path to XML files")
	langPtr := flag.String("l", globalConfig.Language, "The language for the entity filename.")
	debugPtr := flag.Bool("d", globalConfig.Debug, "Whether to show debug output")

	flag.Parse()

	// Validate args
	path := *pathPtr
	lang := *langPtr
	debug := *debugPtr

	if path == "" {
		fmt.Println("You must specify where to start looking for XML/Zonbook files (-p PATH)")
		return
	}

	err = patchFiles(debug, path, lang)
	if err != nil {
		fmt.Println("Got an error patching files in " + path + ":")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Patched XML files in " + path)
	}
}
