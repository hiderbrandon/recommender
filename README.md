# Levantar entorno de desarrollo
## correr proyecto proyecto en  local

1. agregar .env al archivo raiz ([buscar en drive del grupo)](https://drive.google.com/drive/u/2/folders/1W0yLxniB3MZg-YSfSXnIQl65ZBLRV1ms) )
2. correr docker-compose del proyecto
```bash
sudo docker system prune -a --volumes // elimina persistencia de ejecuciones anteriores (opcional)      
sudo docker-compose up -d         # backend
```

## Análisis de calidad con SonarCloud

El proyecto está configurado para enviar análisis de calidad de código automáticamente a **SonarCloud** usando GitHub Actions.

### ¿Qué se analiza?
- Complejidad ciclomática
- Duplicación de código
- Cobertura de pruebas (debe ser >60%)
- Code smells y deuda técnica

### ¿Cómo se ejecuta?

Cada vez que haces `push` o un `pull request`, GitHub Actions ejecuta automáticamente el análisis y lo envía a SonarCloud.

### Ver resultados

Puedes consultar los resultados en:

🔗 [https://sonarcloud.io/project/overview?id=hiderbrandon_recommender](https://sonarcloud.io/project/overview?id=hiderbrandon_recommender)

---

**Nota**: No es necesario correr nada manualmente para esto.
---

## Ejecutar con Docker Compose en entorno local 

1. Clona este repositorio.
2. Crea un archivo `.env` con las variables necesarias (ejemplo abajo).
3. Ejecuta:
   ```bash
   docker-compose up -d
   ```

## Ejemplo de Variables de Entorno
```
DB_PORT=26257
COCKROACH_UI_PORT=8082
APP_PORT=8080
DB_HOST=cockroachdb
DB_USER=root
DB_PASSWORD=
DB_NAME=defaultdb
SSL_MODE=disable 
```
## cerrar  Proyecto 
```
sudo docker-compose down --rmi all -v --remove-orphans

```

