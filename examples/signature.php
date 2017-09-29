<?php

$key = '943b421c9eb07c830af81030552c86009268de4e532ba2ee2eab8247c6da0881';
$salt = '520f986b998545b4785e0defbc4f3c1203f22de2374a3d53cb7a7fe9fea309c5';

$keyBin = pack("H*" , $key);
if(empty($keyBin)) {
	die('Key expected to be hex-encoded string');
}

$saltBin = pack("H*" , $salt);
if(empty($saltBin)) {
	die('Salt expected to be hex-encoded string');
}

$resize = 'fill';
$width = 300;
$height = 300;
$gravity = 'no';
$enlarge = 1;
$extension = 'png';

$url = 'http://img.example.com/pretty/image.jpg';
$encodedUrl = rtrim(strtr(base64_encode($url), '+/', '-_'), '=');

$path = sprintf("/%s/%d/%d/%s/%d/%s.%s", $resize, $width, $height, $gravity, $enlarge, $encodedUrl, $extension);

$signature = rtrim(strtr(base64_encode(hash_hmac('sha256', $saltBin.$path, $keyBin, true)), '+/', '-_'), '=');

print(sprintf("/%s%s", $signature, $path));