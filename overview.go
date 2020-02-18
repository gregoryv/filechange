/*
Package filechange provides a sensor of file modifications.

  s := new(filechange.Sensor)
  s.Visit = func(modified ...string) {
    for _, f := range modified {
      fmt.Println(f)
    }
  }
  s.Recursive = true
  s.Run(context.Background())
*/
package filechange
