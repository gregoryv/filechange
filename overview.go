/*
Package filechange provides a sensor for sensing file modifications.

  s := new(filechange.Sensor)
  s.UseDefaults()
  s.React = func(modified ...string) {
    for _, f := range modified {
      fmt.Println(f)
    }
  }
  s.Recursive = true
  s.Run(context.Background())
*/
package filechange
