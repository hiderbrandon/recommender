## Levantar entorno de desarrollo

pasos para correr el proyecto :

1. agregar .env al archivo raiz (buscar en drive del grupo)
2. correr docker-compose del proyecto
```bash
sudo docker-compose up -d         # backend
```
3. correr docker-compose de sonarqube
```bash
sudo docker-compose -f docker-compose.sonarqube.yml up -d
```
4. visualizar proyecto en sonarqube GUI 


## ver covertura ```
go tool cover -func=coverage.out
```