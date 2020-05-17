#!/usr/bin/env ruby 

# This ruby script uses test/unit to unit test the IPyRubyData ruby code 

require 'test/unit'
require './lib/IPyRubyData'

###############################################################
# Start by mocking the IPyKernelData code used by Convert2Data

class IPyKernelData

  def initialize(aDataValue)
    puts "--------------IPyKernelData::new---------------------"
    pp aDataValue
    puts "-----------------------------------------------------"
  end

  def addData(aMIMEKey, aValue)
    puts "--------------IPyKernelData::addData-----------------"
    pp aMIMEKey
    pp aValue
    puts "-----------------------------------------------------"
  end
  
  def appendTraceback(aValue)
    puts "--------------IPyKernelData::appendTraceback---------"
    pp aValue
    puts "-----------------------------------------------------"
  end

  def addMetadata(aMIMEKey, aMetaKey, aMetaValue)
    puts "--------------IPyKernelData::addMetadata-------------"
    pp aMIMEKey
    pp aMetaKey
    pp aMetaValue
    puts "-----------------------------------------------------"
  end
end

class TestIPyRubyData < Test::Unit::TestCase

  def test_IsIPyRubyMIMEMap
    assert(!IsIPyRubyMIMEMap("hello"),
      "hello is not a MIMEMap-hash")
    assert(!IsIPyRubyMIMEMap({"hello" => "world"} ),
      "hello is not a MIMEMap-key")
    assert(IsIPyRubyMIMEMap({"text/plain" => "Hello world!"} ),
      "text/plain is a MIMEMap-key")
  end

  def test_IsIPyRubyData
    assert(!IsIPyRubyData("hello"),
      "hello is not a Data-hash")
    assert(!IsIPyRubyData({"hello" => "world"}),
      "hello is not a Data-key")
     assert(!IsIPyRubyData({
      "Data" => "hello world"
     }),
      "missing Metadata and Transient")
     assert(!IsIPyRubyData({
      "Metadata" => {},
      "Transient" => {}
     }),
      "missing Data")
     assert(!IsIPyRubyData({
      "Data" => "hello world",
      "Metadata" => {},
      "Transient" => {}
     }),
      "Data is not a MIMEMap")
     assert(IsIPyRubyData({
      "Data" => {"text/plain" => "Hello world!"},
      "Metadata" => {},
      "Transient" => {}
     }),
      "We have a IPyRubyData")
  end
 
#  def test_MakeLastErrorData
#    Can not test MakeLastErrorData directly...
#    ... since it (now) manipulates Go data in TheObjectStore
#    MOVED TO ../rubyEval_IPyKernelData_test.go
#  end
  
  def test_MakeDataAndText
    lastErrData = {
    "Data" => {
      "ename" => "ERROR",
      "evalue" => "testError",
      "traceback" => ["test_MakeDataAndText"],
      "status" => "error"
    },
    "Metadata" => {},
    "Transient" => {}
    }
    assert(MakeDataAndText("silly", lastErrData, "sillier") == lastErrData,
      "should have returned data")
    someData = MakeDataAndText("sillyMimeType", "someData", "someText")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "someData",
      "should have correct MIMETypeText value")
    someData = MakeDataAndText(MIMETypeJSON, "someData", "someText")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "someText",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeJSON),
      "should have MIMETypeJSON")
    assert(someData['Data'][MIMETypeJSON] == "someData",
      "should have correct MIMETypeJSON value")
  end
  
  def test_MakeData
    someData = MakeData(MIMETypeJSON, [10, 20])
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "[10, 20]",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeJSON),
      "should have MIMETypeJSON")
    assert(someData['Data'][MIMETypeJSON] == "[10, 20]",
      "should have correct MIMETypeJSON value")
  end
  
  def test_MakeSVGData
    someData = MakeSVGData("some SVG data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "some SVG data",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeSVG),
      "should have MIMETypeSVG")
    assert(someData['Data'][MIMETypeSVG] == "some SVG data",
      "should have correct MIMETypeSVG value")
  end

  def test_MakePNGData
    someData = MakePNGData("some PNG data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "c29tZSBQTkcgZGF0YQ==\n",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypePNG),
      "should have MIMETypePNG")
    assert(someData['Data'][MIMETypePNG] == "some PNG data",
      "should have correct MIMETypePNG value")
  end

  def test_MakePDFData
    someData = MakePDFData("some PDF data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "c29tZSBQREYgZGF0YQ==\n",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypePDF),
      "should have MIMETypePDF")
    assert(someData['Data'][MIMETypePDF] == "some PDF data",
      "should have correct MIMETypePDF value")
  end
  
  def test_MakeMathData
    someData = MakeMathData("some Math data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "some Math data",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeLatex),
      "should have MIMETypeLatex")
    assert(someData['Data'][MIMETypeLatex] == "$$some Math data$$",
      "should have correct MIMETypeLatex value")
  end

  def test_MakeMarkdownData
    someData = MakeMarkdownData("some Markdown data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "some Markdown data",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeMarkdown),
      "should have MIMETypeMarkdown")
    assert(someData['Data'][MIMETypeMarkdown] == "some Markdown data",
      "should have correct MIMETypeMarkdown value")
  end

  def test_MakeLatexData
    someData = MakeLatexData("some Latex data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "some Latex data",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeLatex),
      "should have MIMETypeLatex")
    assert(someData['Data'][MIMETypeLatex] == "$some Latex data$",
      "should have correct MIMETypeLatex value")
  end

  def test_MakeJSONData
    someData = MakeJSONData("some JSON data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "\"some JSON data\"",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeJSON),
      "should have MIMETypeJSON")
    assert(someData['Data'][MIMETypeJSON] == "\"some JSON data\"",
      "should have correct MIMETypeJSON value")
  end

  def test_MakeJPEGData
    someData = MakeJPEGData("some JPEG data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "c29tZSBKUEVHIGRhdGE=\n",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeJPEG),
      "should have MIMETypeJPEG")
    assert(someData['Data'][MIMETypeJPEG] == "some JPEG data",
      "should have correct MIMETypeJPEG value")
  end

  def test_MakeJavaScriptData
    someData = MakeJavaScriptData("some JavaScript data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "some JavaScript data",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeJavaScript),
      "should have MIMETypeJavaScript")
    assert(someData['Data'][MIMETypeJavaScript] == "some JavaScript data",
      "should have correct MIMETypeJavaScript value")
  end

  def test_MakeHTMLData
    someData = MakeHTMLData("some HTML data")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText] == "some HTML data",
      "should have correct MIMETypeText value")
    assert(someData['Data'].has_key?(MIMETypeHTML),
      "should have MIMETypeHTML")
    assert(someData['Data'][MIMETypeHTML] == "some HTML data",
      "should have correct MIMETypeHTML value")
  end
  
  def test_MakeFileData
    someData = MakeFileData(MIMETypeText, "lib/IPyRubyData.rb")
    assert(IsIPyRubyData(someData), "should have been Data")
    assert(someData['Data'].has_key?(MIMETypeText),
      "should have MIMETypeText")
    assert(someData['Data'][MIMETypeText].include?("def IPyRubyEval"),
      "should have correct MIMETypeText value")
  end
end
