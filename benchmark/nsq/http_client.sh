#! /bin/zsh
for ((i=1; i<=100000; i++)); do
   curl -d "hello world $i" 'http://192.168.31.49:4151/pub?topic=topic'
done