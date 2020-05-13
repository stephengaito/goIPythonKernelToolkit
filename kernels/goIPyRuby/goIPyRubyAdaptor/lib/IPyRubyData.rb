# This ruby script implements the MakeData method

require 'json'
require 'pp'
require 'base64'

  # The following are the "standard" "MIMETypes" for IPython Data
  #
	MIMETypeHTML       = "text/html"
	MIMETypeJavaScript = "application/javascript"
	MIMETypeJPEG       = "image/jpeg"
	MIMETypeJSON       = "application/json"
	MIMETypeLatex      = "text/latex"
	MIMETypeMarkdown   = "text/markdown"
	MIMETypePNG        = "image/png"
	MIMETypePDF        = "application/pdf"
	MIMETypeSVG        = "image/svg+xml"
	MIMETypeText       = "text/plain"
  #
  # The following are allowed "mimetypes" for error reporting
  #
  MIMETypeEName      = "ename"
  MIMETypeEValue     = "evalue"
  MIMETypeETraceback = "traceback"
  MIMETypeEStatus    = "status"

# The following are the "standard" "MIMETypes" for IPython Data
#
IPyRubyMIMEMapKeys = [
	MIMETypeHTML,
	MIMETypeJavaScript,
	MIMETypeJPEG,
	MIMETypeJSON,
	MIMETypeLatex,
	MIMETypeMarkdown,
	MIMETypePNG,
	MIMETypePDF,
	MIMETypeSVG,
	MIMETypeText,
]

# The following are allowed "mimetypes" for error reporting
#
IPyRubyMIMEMapErrorKeys = [
  MIMETypeEName,
  MIMETypeEValue,
  MIMETypeETraceback,
  MIMETypeEStatus,
]

IPyRubyMIMEMapAllKeys = [IPyRubyMIMEMapKeys, IPyRubyMIMEMapErrorKeys].flatten

def IsIPyRubyMIMEMap(aValue)
  return false unless aValue.is_a?(Hash)
  return false unless aValue.keys.difference(IPyRubyMIMEMapAllKeys).length < 1
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

def MakeFileData(mimeType, filePath)
end

def MakeHTMLData(someHtml)
  return MakeData(MIMETypeHTML, someHtml.to_s)
end

def MakeJavaScriptData(someJavaScript)
  return MakeData(MIMETYpeJavaScript, someJavaScript.to_s)
end

# Note that JPEG image bytes are stored as bytes inside a Ruby String
# (which *can* include null bytes).
#
def MakeJPEGData(someJPEGImageBytes)
  return MakeDataAndText(
    MIMETypeJPEG,
    someJPEGImageBytes.to_s,
    Base64.encode64(someJPEGImageBytes.to_s)
  )
end

def MakeJSONData(aValue)
  return MakeData(MIMETypeJSON, JSON.generate(aValue))
end

def MakeLatexData(someLatex)
  return MakeDataAndText(
    MIMETypeLatex,
    "$"+someLatex.to_s.strip+"$",
    someLatex.to_s
  )
end

def MakeMarkdownData(someMarkdown)
  return MakeData(MIMETypeMarkdown, someMarkdown.to_s)
end

def MakeMathData(someLatex)
  return MakeDataAndText(
    MIMETypeLatex,
    "$$"+someLatex.to_s.strip+"$$",
    someLatex.to_s
  )
end

# Note that PDF bytes are stored as bytes inside a Ruby String
# (which *can* include null bytes).
#
def MakePDFData(somePDFBytes)
  return MakeData(MIMETypePDF, somePDFBytes.to_s)
end

# Note that PNG image bytes are stored as bytes inside a Ruby String
# (which *can* include null bytes).
#
def MakePNGData(somePNGImageBytes)
  return MakeDataAndText(
    MIMETypePNG,
    somePNGImageBytes.to_s,
    Base64.encode(somePNGImageBytes.to_s)
  )
end

def MakeSVGData(someSVG)
  return MakeData(MIMETypeSVG, someSVG.to_s)
end

def MakeData(mimeType, data)
  textData = data
  textData = data.pretty_inspect unless data.is_a?(String)
  return MakeDataAndText(
    mimeType,
    data,
    textData
  )
end

def MakeDataAndText(mimeType, data, textData)
  return data if IsIPyRubyData(data)
  dataValue         = Hash.new
  dataValue['Data'] = Hash.new
  mimeType = MIMETypeText unless IPyRubyMIMEMapKeys.include?(mimeType)
  jsonData = JSON.generate(data) 
  data     = data.pretty_inspect unless data.is_a?(Stirng)
  dataValue['Data'][MIMETypeJSON] = jsonData
  dataValue['Data'][mimeType]     = data
  dataValue['Metadata']  = Hash.new
  dataValue['Transient'] = Hash.new
  return dataValue
end

def Convert2Data(origValue)

  # ensure origValue IS an IPyRubyData object
  #
  origValue = MakeData("", origValue)

  # Now work with the goIPyRuby callbacks to convert this IPyRubyData object 
  # into a goIPyKernel Data object
  #
  dataObj = IPyKernelData.new
  origValue['Data'].each_pair do | aMIMEKey, aValue |
    dataObj.addData(aMIMEKey, aValue)
  end
  origValue['Metadata'].each_pair do | aMIMEKey, aValue |
    #
    # Metadata is a collection of hashed of key-value pairs corresponding to 
    # each IPyRubyMIMEMapKeys.
    #
    # We simply ignore any metadata which is not an IPyRubyMIMEMapKey or not
    # a hash of key-value pairs.
    #
    if aValue.is_a?(Hash) && IPyRubyMIMEMapKeys.include?(aMIMEKey) then
      aValue.each_pair do | aMetaKey, aMetaValue |
        dataObj.addMetadata(aMIMEKey, aMetaKey.to_s, aMetaValue.to_s)
      end
    end
  end
  #
  # At the moment there is NO example of the use of the 'Transient' data-key 
  # in gophernotes and/or the IPython documentation I can find... 
  #
  # SO we quitely ignore 'Transient' data.
  #
  return dataObj
end

void MakeLastErrorData() {
// could use: 
// $! 	latest error message
// $@ 	location of error
// $_ 	string last read by gets
// $. 	line number last read by interpreter 

  return {
    Data: {
      ename: "ERROR",
      evalue: $!,
      traceback: nil,
      status: "error"
    }
  }
}

def IPyRubyEval(aString)
  evalResult = begin
    TOPLEVEL_BINDING.eval(aString)
  rescue
    MakeLastErrorData()
  end
  
  begin
    return Convert2Data(evalResult)
  rescue
    # not much we can do... should not get here!
  end
  # SHOULD NOT end up here!
  return nil
end