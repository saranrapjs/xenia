
# sfix
    import "github.com/coralproject/xenia/pkg/script/sfix"






## func Add
``` go
func Add(db *db.DB, scr script.Script) error
```
Add inserts a script for testing.


## func Get
``` go
func Get(fileName string) (script.Script, error)
```
Get retrieves a set document from the filesystem for testing.


## func Remove
``` go
func Remove(db *db.DB, pattern string) error
```
Remove is used to clear out all the test sets from the collection.
All test documents must start with STEST in their name.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)