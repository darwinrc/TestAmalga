# Parte 1

### Calcular resumen de compras

Ejecutar el programa con el siguiente comando (debe tener instalado Go 1.20):
```
go run main.go
```

Llamar el endpoint:
```
curl --location --request GET 'localhost:5000/resumen/2019-12-01?dias=5'
```

Ejemplo de la respuesta:
```
{
    "total": 29815625.82,
    "comprasPorTDC": {
        "amex": 2861283.64,
        "amex corp": 2923605.31,
        "maestro": 2901466.24,
        "master classic": 2930035.72,
        "master gold": 2988878.78,
        "master plat": 2395721.04,
        "privada": 2538281.95,
        "visa classic": 2697760.23,
        "visa debit": 2318518.49,
        "visa gold": 2750625.22,
        "visa plat": 2509449.2
    },
    "nocompraron": 1888,
    "compraMasAlta": 30477.77
}
```

#### Ejecutar la suite de pruebas

```
go test ./...
```