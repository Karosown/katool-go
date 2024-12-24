package rss

const SourceCodeGetFunc = `()=>{
return fetch(window.location.href)  
  .then(response => response.text())  
  .then(text => text);  
}`
