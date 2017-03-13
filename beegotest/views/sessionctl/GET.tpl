<!DOCTYPE html>
<html lang="zh-cn">
<!--<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />-->
  <head>
    <meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>小功能，方便测试</title>
	
    <!-- Bootstrap -->
	<link href="/public/bootstrap/3.3.0/css/bootstrap.min.css" rel="stylesheet">
	
	<!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="http://cdn.bootcss.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="http://cdn.bootcss.com/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
{{if .username}}
<p>欢迎您，{{.username}}</p>
{{else}}    
<form action="/sessionctl" method="post">
用户:<input type="text" name="username"><br>
密码:<input type="password" name="password"><br>
<input type="submit" value="登录">
</form><br>
{{end}}

  <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
  <script src="/public/jquery/3.1.1/jquery.min.js"></script>
  <!-- Include all compiled plugins (below), or include individual files as needed -->
  <script src="/public/bootstrap/3.3.0/js/bootstrap.min.js"></script>
  </body>
</html>