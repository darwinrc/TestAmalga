# Parte 2

### Obtener organizaciones

Ejecutar el programa con el siguiente comando (debe tener instalado Go 1.20):
```
go run main.go
```

Llamar el endpoint:
```
curl --location --request GET 'localhost:5000/organization?file=organization.csv'
```

Ejemplo de la respuesta:
```
[
    {
        "organization": "org1",
        "users": [
            {
                "username": "jperez",
                "roles": [
                    "admin",
                    "superadmin"
                ]
            },
            {
                "username": "asosa",
                "roles": [
                    "writer"
                ]
            }
        ]
    },
    {
        "organization": "org2",
        "users": [
            {
                "username": "jperez",
                "roles": [
                    "admin"
                ]
            },
            {
                "username": "rrodriguez",
                "roles": [
                    "writer",
                    "editor"
                ]
            }
        ]
    }
]
```

#### Ejecutar la suite de pruebas

```
go test ./...
```