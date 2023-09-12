package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lukasz-lobocki/tabby"
)

/*
parseBool returns thisString if thisBool

	'thisBool' value to be checked
	'thisString' value to be returned if true
*/
func parseBool(thisBool bool, thisString string) string {

	if thisBool {
		return thisString
	}

	return ""
}

/*
escapeMarkdown returns same string but safeguarderd against markdown interpretation

	'text' text to be safeguarded
*/
func escapeMarkdown(text string) string {

	// These characters need to be escaped in Markdown in order to appear as literal characters instead of performing some markdown functions
	needEscape := []string{
		`\`, "`", "*", "_",
		"{", "}",
		"[", "]",
		"(", ")",
		"#", ".", "!",
		"+", "-",
	}

	for _, thisNeed := range needEscape {
		text = strings.Replace(text, thisNeed, `\`+thisNeed, -1)
	}

	return text
}

/*
emitJson prints result in the form of a json

	'repos' slice of structures describing the repos
*/
func emitJson(repos []tRepo) error {

	jsonInfo, err := json.MarshalIndent(repos, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling json failed. %w", err)
	}
	fmt.Println(string(jsonInfo))

	if loggingLevel >= 2 {
		logInfo.Printf("%d records marshalled.\n", len(repos))
	}

	return nil
}

/*
emitTable prints result in the form of a table

	'repos' slice of structures describing the repos
*/
func emitTable(repos []tRepo) error {

	table := new(tabby.Table)

	thisColumns := getColumns()

	var thisHeader []string

	/* Building slice of titles */

	for _, thisColumn := range thisColumns {
		if thisColumn.isShown(config) {

			thisHeader = append(thisHeader,
				color.New(thisColumn.titleColor).SprintFunc()(
					thisColumn.title(config),
				),
			)

		}
	}

	/* Set the header */

	if err := table.SetHeader(thisHeader); err != nil {
		return fmt.Errorf("emitTable: setting header failed. %w", err)
	}

	if loggingLevel >= 1 {
		logInfo.Println("header set.")
	}

	/* Populate the table */

	for _, thisRepo := range repos {

		var thisRow []string

		/* Building slice of columns within a single row*/

		for _, thisColumn := range thisColumns {

			if thisColumn.isShown(config) {
				thisRow = append(thisRow,
					color.New(thisColumn.contentColor(thisRepo)).SprintFunc()(
						thisColumn.contentSource(config, thisRepo),
					),
				)
			}
		}

		if err := table.AppendRow(thisRow); err != nil {
			return fmt.Errorf("emitTable: appending row failed. %w", err)
		}
		if loggingLevel >= 3 {
			logInfo.Printf("row [%s] appended.", thisRepo.ShortName)
		}

	}

	if loggingLevel >= 2 {
		logInfo.Printf("%d rows appended.\n", len(repos))
	}

	/* Emit the table */

	if loggingLevel >= 3 {
		table.Print(&tabby.Config{Spacing: "|", Padding: "."})
	} else {
		table.Print(nil)
	}
	return nil
}

/*
emitMarkdown prints result in the form of markdown table

	'repos' slice of structures describing the repos
*/
func emitMarkdown(repos []tRepo) {
	thisColumns := getColumns()

	var thisHeader []string

	/* Building slice of titles */

	for _, thisColumn := range thisColumns {
		if thisColumn.isShown(config) {
			thisHeader = append(thisHeader, thisColumn.title(config))
		}
	}

	/* Emiting titles */

	fmt.Println("| " + strings.Join(thisHeader, " | ") + " |")

	if loggingLevel >= 1 {
		logInfo.Println("header printed.")
	}

	/* Emit markdown line that separates header from body table */

	var thisSeparator []string

	for _, thisColumn := range thisColumns {
		if thisColumn.isShown(config) {
			thisSeparator = append(thisSeparator, getThisAlignChar()[thisColumn.contentAlignMD])
		}
	}
	fmt.Println("| " + strings.Join(thisSeparator, " | ") + " |")

	if loggingLevel >= 1 {
		logInfo.Println("separator printed.")
	}

	/* Iterating through repos */

	for _, thisRepo := range repos {

		var thisRow []string

		/* Building slice of columns within a single row*/

		for _, thisColumn := range thisColumns {
			if thisColumn.isShown(config) {
				if thisColumn.contentEscapeMD {
					thisRow = append(thisRow, escapeMarkdown(thisColumn.contentSource(config, thisRepo)))
				} else {
					thisRow = append(thisRow, thisColumn.contentSource(config, thisRepo))
				}
			}
		}

		/* Emitting row */

		fmt.Println("| " + strings.Join(thisRow, " | ") + " |")
	}

	if loggingLevel >= 2 {
		logInfo.Printf("%d rows printed.\n", len(repos))
	}
}
