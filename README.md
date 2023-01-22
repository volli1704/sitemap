# sitemap
Parses site url's and return list of all site url in parallel

```
  go build -o ./parser cmd/parser/main.go
```

```
  ./parser -url https://gobyexample.com/ -depth 3 -workers 15 -o result.out
```
