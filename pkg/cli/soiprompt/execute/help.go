package execute

import "fmt"

func (e *Executor) help(in string) error {
	fmt.Println(`
General:
      Soi is a url management CLI system. which could add url, find url and rename url.

Commands:

  [add]: 
      Desc:  add url to model
      Usage: add (URL)
      Option:
        -n: logical name of the url    (default: <title> of the URL)
        -d: directory to store the url (default: "new")

  [dig]: 
      Desc:  dig url directory with [Tab] key completion and [→] key listing next suggestions
      Usage: dig (URL Completion with [Tab] and [→] key listing next suggestions)

  [list]:
      Desc:  list all urls with filtering
      Usage: list (free words)
		
  [tag]:
      Desc:  not implemented now

  [mv]:
      Desc:  move file to dir 
      Usage: mv (current path) to (dir)

  [quit]:
      Desc:  quit model> and go back to console. Ctrl+D works too.
      Usage: quit`)
	return nil
}
