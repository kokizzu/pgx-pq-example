

# PGX and PQ Example

```

alias time='/usr/bin/time -f "\nCPU: %Us\tReal: %es\tRAM: %MKB"'

docker-compose up

go mod tidy

time go run main.go

1000
PQ 3.162459197s
1000
PGX 832.63568ms

CPU: 0.73s      Real: 4.34s     RAM: 97524KB
```
