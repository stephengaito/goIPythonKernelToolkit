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
  return false unless (aValue.keys - IPyRubyMIMEMapAllKeys).length < 1
  return true
end

def IsIPyRubyData(aValue)
  return false unless aValue.is_a?(Hash)
  return false unless aValue.has_key?('Data')
  return false unless aValue.has_key?('Metadata')
  return false unless aValue.has_key?('Transient')
  return false unless IsIPyRubyMIMEMap(aValue['Data'])
  return false unless IsIPyRubyMIMEMap(aValue['Metadata'])
  return false unless IsIPyRubyMIMEMap(aValue['Transient'])
  return true
end

def MakeFileData(mimeType, filePath)
  fileContents = IO.read(filePath)
  case mimeType
  when MIMETypeJPEG
    MakeJPEGData(fileContents)
  when MIMETypePDF
    MakePDFData(fileContents)
  when MIMETypePNG
    MakePNGData(fileContents)
  else
    MakeData(mimeType, fileContents)
  end
end

def MakeHTMLData(someHtml)
  return MakeData(MIMETypeHTML, someHtml.to_s)
end

def MakeJavaScriptData(someJavaScript)
  return MakeData(MIMETypeJavaScript, someJavaScript.to_s)
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
  return MakeDataAndText(
    MIMETypePDF,
    somePDFBytes.to_s,
    Base64.encode64(somePDFBytes.to_s)
  )
end

# Note that PNG image bytes are stored as bytes inside a Ruby String
# (which *can* include null bytes).
#
def MakePNGData(somePNGImageBytes)
  return MakeDataAndText(
    MIMETypePNG,
    somePNGImageBytes.to_s,
    Base64.encode64(somePNGImageBytes.to_s)
  )
end

def MakeSVGData(someSVG)
  return MakeData(MIMETypeSVG, someSVG.to_s)
end

def MakeTextData(someText)
  return MakeData(MIMETypeText, someText.to_s)
end

def MakeData(mimeType, data)
  textData = data
  textData = data.pretty_inspect.chomp unless data.is_a?(String)
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
  data     = data.pretty_inspect.chomp unless data.is_a?(String)
  textData = "" if textData.nil?
  textData = textData.pretty_inspect.chomp unless textData.is_a?(String)
  dataValue['Data'][MIMETypeText] = textData
  dataValue['Data'][mimeType]     = data
  dataValue['Metadata']  = Hash.new
  dataValue['Transient'] = Hash.new
  return dataValue
end

def Convert2Data(origValue)
 
  # ensure origValue IS an IPyRubyData object
  #
  origValue = 
    if origValue.nil? then MakeDataAndText(MIMETypeText, "", "")
    else                   MakeData("", origValue)
    end

  # Now work with the goIPyRuby callbacks to convert this IPyRubyData object 
  # into a goIPyRuby Data object
  #
  dataObj = IPyRubyData_New(origValue)
  
  origValue['Data'].each_pair do | aMIMEKey, aValue |
    IPyRubyData_AddData(dataObj, aMIMEKey, aValue)
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
        IPyRubyData_AddMetadata(dataObj, aMIMEKey, aMetaKey.to_s, aMetaValue.to_s)
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

def MakeLastErrorData(err, errMsg) 
# could use: 
# $! 	latest error message
# $@ 	location of error
# $_ 	string last read by gets
# $. 	line number last read by interpreter 
  
  dataObj = IPyRubyData_New([err, errMsg, "lastErrorData"])
  IPyRubyData_AddData(dataObj, "ename", "ERROR")
  IPyRubyData_AddData(dataObj, "evalue", err.to_s)
  IPyRubyData_AppendTraceback(dataObj, errMsg)
  IPyRubyData_AddData(dataObj, "status", "error")
    
  return dataObj
end

def IPyRubyEval(aString)

  # prepare to capture any errors from the TOPLEVEL_BINDING.eval
  #
  evalTrace = TracePoint.new(:raise) do | tp |
    errReport = [ tp.raised_exception.to_s ]
    strLines  = aString.lines
    if tp.lineno < strLines.length then 
      errReport.push "in line: (#{tp.lineno})"
      errReport.push aString.lines[tp.lineno]
    end
    errorData = MakeLastErrorData(
      errReport.join("\n"),
      "IPyRuby kernel evaluation of Ruby code FAILED"
    )
    evalTrace.disable
    return errorData
  end
  evalTrace.enable
  evalResult = TOPLEVEL_BINDING.eval(aString)
  evalTrace.disable

  # prepare to capture any errors from the Convert2Data
  #
  convertTrace = TracePoint.new(:raise) do | tp |
    errReport = [
      tp.raised_exception.to_s,
      caller.shift
    ]
    errorData = MakeLastErrorData(
      errReport.join("\n"),
      "Convert2Data FAILED"
    )
    convertTrace.disable
    return errorData
  end
  convertTrace.enable
  convertedResult = Convert2Data(evalResult)
  convertTrace.disable
  
  return convertedResult
end
