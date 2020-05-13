# This ruby script implements the MakeData method

require 'json'
require 'pp'

#	MIMETypeHTML       = "text/html"
#	MIMETypeJavaScript = "application/javascript"
#	MIMETypeJPEG       = "image/jpeg"
#	MIMETypeJSON       = "application/json"
#	MIMETypeLatex      = "text/latex"
#	MIMETypeMarkdown   = "text/markdown"
#	MIMETypePNG        = "image/png"
#	MIMETypePDF        = "application/pdf"
#	MIMETypeSVG        = "image/svg+xml"
#	MIMETypeText       = "text/plain"

IPyRubyMIMEMapKeys = [
  "text/html",
	"application/javascript",
	"image/jpeg",
	"application/json",
	"text/latex",
	"text/markdown",
	"image/png",
	"application/pdf",
	"image/svg+xml",
	"text/plain",
]
  
def IsIPyRubyMIMEMap(aValue)
  return false unless aValue.is_a?(Hash)
  return false unless aValue.keys.difference(IPyRubyMIMEMapKeys).length < 1
  return true
end

def IsIPyRubyData(aValue)
  return false unless aValue.is_a?(Hash)
  return false unless aValue.has_key?('Data')
  return false unless aValue.has_key?('Metadata')
  return false unless aValue.has_key?('Transient')
  return false unless IsIPyRubyMIMEMap(aValue['Data'])
  return false unless IsIPyRubyMIMEMap(aValue['MetaData'])
  return false unless IsIPyRubyMIMEMap(aValue['Transient'])
  return true
end

def MakeData(mimeType, data)
  dataValue = Hash.new
  

end

def Convert2Data(origValue) 
  return origValue if IsIPyRubyData(origValue)
  newValue = Hash.new
  newValue['Data'] = Hash.new
  

end