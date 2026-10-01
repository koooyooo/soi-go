package execute

import "fmt"

func (e *Executor) help(in string) error {
	fmt.Println(`
General:
      Soi is a url management CLI system. which could add url, find url and rename url.

Commands:

  [add]:
      Desc:  add url to soi
      Usage: add (dir) (name) (URL) (#tags...)
      Option:
        -n: logical name of the url    (default: <title> of the URL)
        -d: directory to store the url (default: YYYY-MM)

  [dig]:
      Desc:  dig url directory with [Tab] key completion and [→] key listing next suggestions
      Usage: dig (path)

  [list]:
      Desc:  list all urls with filtering
      Usage: list (free words)
      Option:
        -c/-f/-s/-e: browser
        -n/-a/-v: sort
        -p: private mode

  [tag]:
      Desc:  replace tags of a soi
      Usage: tag (hash) (#tags...)

  [mv]:
      Desc:  move file or dir within the current bucket
      Usage: mv (from) (to)

  [rm]:
      Desc:  remove file or dir within the current bucket
      Usage: rm (path)

  [cb]:
      Desc:  show or change bucket
      Usage: cb [bucket]

  [quit]:
      Desc:  quit soi> and go back to console. Ctrl+D works too.
      Usage: quit`)
	return nil
}
