package rss

// SourceCodeGetFunc 用于获取网页源代码的JavaScript函数
// SourceCodeGetFunc is a JavaScript function for getting web page source code
const SourceCodeGetFunc = `()=>{
return fetch(window.location.href)  
  .then(response => response.text())  
  .then(text => text);  
}`
