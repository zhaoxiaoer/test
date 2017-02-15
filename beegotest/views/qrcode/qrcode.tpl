<p>{{date .time .format}}</p>
<form action={{.url}} method="post">
请输入内容:<br>
<input type="text" name="content" value={{.content}} style="width:600px;height:30px"><br>
<input type="submit" value="生成二维码">
</form><br>
{{if .errinfo}}
<p>{{.errinfo}}</p>
{{end}}
{{if .image}}
<!--img src={{.image}} alt="主页二维码" /-->
<img src="data:;base64,{{.image}}" alt="二维码" /><br>
{{end}}
