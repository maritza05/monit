package main

import (
	"flag"
	"io"
)

// type FileLog struct {
// 	name       string
// 	lineOffset int64
// }

func main() {
	// var notifier Notifier
	// f := &flags{
	// 	offset:       1,
	// 	pattern:      "ERROR|WARN",
	// 	verbose:      false,
	// 	slack:        false,
	// 	hide_finding: true,
	// }

	// if err := f.parse(); err != nil {
	// 	os.Exit(1)
	// }

	// ERROR_REGEX, err := regexp.Compile(f.pattern)
	// if err != nil {
	// 	panic("Invalid regex")
	// }

	// switch {
	// case f.slack:
	// 	notifier = NewSlackNotifier(
	// 		os.Getenv("SLACK_BOT_TOKEN"),
	// 		os.Getenv("SLACK_CHANNEL_ID"))
	// default:
	// 	notifier = NewFileNotifier(f.output_file)
	// }

	// results := make(chan Result, len(f.filepaths))

	// // Launching go routines for files
	// for _, n := range f.filepaths {
	// 	// Doing this because of range variable f captured by func literal warning
	// 	ff := n
	// 	go func() {
	// 		monitor := &FileMonitor{offset: f.offset,
	// 			regex:     ERROR_REGEX,
	// 			storeFile: f.output_file,
	// 			filepath:  ff}
	// 		monitor.Start(results)
	// 	}()
	// }

	// // Wait for goroutines
	// var notification string
	// for i := 0; i < len(f.filepaths); i++ {
	// 	foundError := <-results

	// 	// If we find an error show notification to show what happened
	// 	if foundError.error != nil {
	// 		notification = fmt.Sprintf("Something failed when trying to read file: %s\n%s", foundError.filepath, foundError.error)
	// 		notifier.Notify(notification)
	// 		return
	// 	}

	// 	// If we dont got an error and found is false
	// 	if foundError.error == nil && !foundError.found {
	// 		return
	// 	}

	// 	// If we found something
	// 	if foundError.found {
	// 		if f.hide_finding && f.slack {
	// 			notification = fmt.Sprintf("Error found in file: %s", foundError.filepath)
	// 			notifier.Notify(notification)
	// 		} else {
	// 			notification = fmt.Sprintf("This is the error found in file: %s\n%s", foundError.filepath, foundError.content)
	// 			notifier.Notify(notification)
	// 		}
	// 	}

	// }

}

func run(s *flag.FlagSet, args []string, out io.Writer) error {

}
