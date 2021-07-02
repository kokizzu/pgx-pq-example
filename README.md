

# PGX Example

```

alias time='/usr/bin/time -f "\nCPU: %Us\tReal: %es\tRAM: %MKB"'

docker-compose up

go mod tidy

time go run main.go
1000

CPU: 0.64s      Real: 3.51s     RAM: 94016KB
```
