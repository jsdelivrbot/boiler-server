package controllers

import (
	"github.com/AzureTech/goazure"
)

type ErrorController struct {
	goazure.Controller
}
//var E402 string = "<!DOCTYPE html PUBLIC\"-//W3C//DTD XHTML 1.0 Strict//EN\"\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><meta http-equiv=\"Content-Type\"content=\"text/html; charset=gb2312\"/><title>402-禁止访问:访问被拒绝。</title><style type=\"text/css\"><!--body{margin:0;font-size:.7em;font-family:Verdana,Arial,Helvetica,sans-serif;background:#EEEEEE}fieldset{padding:0 15px 10px 15px}h1{font-size:2.4em;margin:0;color:#FFF}h2{font-size:1.7em;margin:0;color:#CC0000}h3{font-size:1.2em;margin:10px 0 0 0;color:#000000}#header{width:96%;margin:0 0 0 0;padding:6px 2%6px 2%;font-family:\"trebuchet MS\",Verdana,sans-serif;color:#FFF;background-color:#555555}#content{margin:0 0 0 2%;position:relative}.content-container{background:#FFF;width:96%;margin-top:8px;padding:10px;position:relative}--></style></head><body><div id=\"header\"><h1>服务器错误</h1></div><div id=\"content\"><div class=\"content-container\"><fieldset><h2>403-禁止访问:访问被拒绝。</h2><h3>您无权使用所提供的凭据查看此目录或页面。</h3></fieldset></div></div></body></html>"
var E403 string = "<!DOCTYPE html PUBLIC\"-//W3C//DTD XHTML 1.0 Strict//EN\"\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><meta http-equiv=\"Content-Type\"content=\"text/html; charset=gb2312\"/><title>403-禁止访问:访问被拒绝。</title><style type=\"text/css\"><!--body{margin:0;font-size:.7em;font-family:Verdana,Arial,Helvetica,sans-serif;background:#EEEEEE}fieldset{padding:0 15px 10px 15px}h1{font-size:2.4em;margin:0;color:#FFF}h2{font-size:1.7em;margin:0;color:#CC0000}h3{font-size:1.2em;margin:10px 0 0 0;color:#000000}#header{width:96%;margin:0 0 0 0;padding:6px 2%6px 2%;font-family:\"trebuchet MS\",Verdana,sans-serif;color:#FFF;background-color:#555555}#content{margin:0 0 0 2%;position:relative}.content-container{background:#FFF;width:96%;margin-top:8px;padding:10px;position:relative}--></style></head><body><div id=\"header\"><h1>服务器错误</h1></div><div id=\"content\"><div class=\"content-container\"><fieldset><h2>403-禁止访问:访问被拒绝。</h2><h3>您无权使用所提供的凭据查看此目录或页面。</h3></fieldset></div></div></body></html>"
//var E404 string = "<!DOCTYPE html PUBLIC\"-//W3C//DTD XHTML 1.0 Strict//EN\"\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><meta http-equiv=\"Content-Type\"content=\"text/html; charset=gb2312\"/><title>404-找不到文件或目录。</title><style type=\"text/css\"><!--body{margin:0;font-size:.7em;font-family:Verdana,Arial,Helvetica,sans-serif;background:#EEEEEE}fieldset{padding:0 15px 10px 15px}h1{font-size:2.4em;margin:0;color:#FFF}h2{font-size:1.7em;margin:0;color:#CC0000}h3{font-size:1.2em;margin:10px 0 0 0;color:#000000}#header{width:96%;margin:0 0 0 0;padding:6px 2%6px 2%;font-family:\"trebuchet MS\",Verdana,sans-serif;color:#FFF;background-color:#555555}#content{margin:0 0 0 2%;position:relative}.content-container{background:#FFF;width:96%;margin-top:8px;padding:10px;position:relative}--></style></head><body><div id=\"header\"><h1>服务器错误</h1></div><div id=\"content\"><div class=\"content-container\"><fieldset><h2>404-找不到文件或目录。</h2><h3>您要查找的资源可能已被删除，已更改名称或者暂时不可用。</h3></fieldset></div></div></body></html>"
//var E500 string = "<!DOCTYPE html PUBLIC\"-//W3C//DTD XHTML 1.0 Strict//EN\"\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\"><head><meta http-equiv=\"Content-Type\"content=\"text/html; charset=gb2312\"/><title>500-内部服务器错误。</title><style type=\"text/css\"><!--body{margin:0;font-size:.7em;font-family:Verdana,Arial,Helvetica,sans-serif;background:#EEEEEE}fieldset{padding:0 15px 10px 15px}h1{font-size:2.4em;margin:0;color:#FFF}h2{font-size:1.7em;margin:0;color:#CC0000}h3{font-size:1.2em;margin:10px 0 0 0;color:#000000}#header{width:96%;margin:0 0 0 0;padding:6px 2%6px 2%;font-family:\"trebuchet MS\",Verdana,sans-serif;color:#FFF;background-color:#555555}#content{margin:0 0 0 2%;position:relative}.content-container{background:#FFF;width:96%;margin-top:8px;padding:10px;position:relative}--></style></head><body><div id=\"header\"><h1>服务器错误</h1></div><div id=\"content\"><div class=\"content-container\"><fieldset><h2>500-内部服务器错误。</h2><h3>您查找的资源存在问题，因而无法显示。</h3></fieldset></div></div></body></html>"

func (ctl *ErrorController) Error402() {
	tplName := "custerrs/402.html"
	ctl.TplName = tplName
	ctl.Render()
}

func (ctl *ErrorController) Error403() {
	tplName := "custerrs/402.html"
	ctl.TplName = tplName
	ctl.Render()
}

func (ctl *ErrorController) Error404() {
	tplName := "custerrs/404.html"
	ctl.TplName = tplName
	ctl.Render()
}

func (ctl *ErrorController) Error501() {
	tplName := "custerrs/501.html"
	ctl.TplName = tplName
	ctl.Render()
}

func (ctl *ErrorController) Error500() {
	tplName := "custerrs/500.html"
	ctl.TplName = tplName
	ctl.Render()
}

func (ctl *ErrorController) Error502() {
	tplName := "custerrs/502.html"
	ctl.TplName = tplName
	ctl.Render()
}