# Levantar entorno de desarrollo
## correr proyecto proyecto en  local

1. agregar .env al archivo raiz ([buscar en drive del grupo)](https://drive.google.com/drive/u/2/folders/1W0yLxniB3MZg-YSfSXnIQl65ZBLRV1ms) )
2. correr docker-compose del proyecto
```bash
sudo docker system prune -a --volumes // elimina persistencia de ejecuciones anteriores (opcional)      
sudo docker-compose up -d         # backend
```
## ver pruebas pruebas en local :

3. correr el docker-compose de sonarqube
```bash
sudo docker-compose -f docker-compose.sonarqube.yml up -d
```

4. optener token de sonarqube como muestra [el laboratorio](https://docs.google.com/document/d/1M4jfM4QFLrdFof22xvSUrf-4OPilO0gvXoNjAB5STLA/edit?tab=t.0) y ponerlo en .env "SONAR_TOKEN" 
5. correr sonar scanner en la raiz del proyecto 
### en linux
```
docker run --rm \
  --env-file .env \
  -e SONAR_HOST_URL="http://172.17.0.1:9000" \
  -v "$(pwd):/usr/src" \
  sonarsource/sonar-scanner-cli

```

### en windows docker-desktop
```
docker run --rm \
  --env-file .env \
  -e SONAR_HOST_URL="http://host.docker.internal:9000" \
  -v "$(pwd):/usr/src" \
  sonarsource/sonar-scanner-cli
```

5. visualizar proyecto en [sonarqube GUI](http://localhost:9000/dashboard?id=recommender&codeScope=overall) para ver reportes



## ver covertura (en CLI)
```
go tool cover -func=coverage.out
```

## cerrar  royecto 
```
sudo docker-compose down --rmi all -v --remove-orphans

```

