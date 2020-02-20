## Usage
```
1. bash# go build
2. bash# for i in {1..10};do { curl -i -XPOST  -d '{"cmd":"sleep 3"}' localhost:11000/launch} &;done; wait
```

