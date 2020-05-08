// The following is very simple Markdown processor based upon the goldmark 
// Readme example. (See: [yuin/goldmark](https://github.com/yuin/goldmark))
//
// It walks through all *.md files in the `docs/content` directory and 
// renders them as HTML in the `docs/html` directory. 
//
// Since these *.md files are under our control, we explicilty allow 
// (unsafe) embedded html. We can use this to provide titles for all 
// webpages. 
//
package main

import (
   "bytes"
   "fmt"
   "github.com/yuin/goldmark"
   "github.com/yuin/goldmark/extension"
   "github.com/yuin/goldmark/parser"
   "github.com/yuin/goldmark/renderer/html"
   "io/ioutil"
   "os"
   "path/filepath"
   "strings"
)
func main() {

  md := goldmark.New(
    goldmark.WithExtensions(extension.GFM),
    goldmark.WithParserOptions(
      parser.WithAutoHeadingID(),
    ),
    goldmark.WithRendererOptions(
      //html.WithHardWraps(),
      html.WithXHTML(),
      html.WithUnsafe(),
    ),
  )

  os.MkdirAll("docs/content", 0755)
  os.MkdirAll("docs/html",    0755)
  
  err := filepath.Walk("docs/content", func(path string, info os.FileInfo, err error) error {
    if err != nil { return err }
    if strings.HasSuffix(path, ".md") {
    
      htmlPath := strings.Replace(path, "docs/content", "docs/html", 1)
      htmlPath =  strings.Replace(htmlPath, ".md", ".html", 1)
      fmt.Printf("converting [%s] to [%s]\n", path, htmlPath)
      
      fileBytes, err := ioutil.ReadFile(path)
      if err != nil { return err }
      
      var buf bytes.Buffer
      err = md.Convert(fileBytes, &buf)
      if err != nil { return err }
      
      err = ioutil.WriteFile(htmlPath, buf.Bytes(), 0644)
      if err != nil { return err }
    }
    return nil
  })
  if err != nil {
    fmt.Printf("docTool-goIPy error: %v", err)
  }
  
}
