# Proyecto2-DB1
Este repositorio contiene todo el código fuente para el segundo proyecto de la clase de Bases de Datos 1: Un sistema de control de negocio
para una empresa ficticia.

---
# Sección 1: Diseño de Base de Datos
Para esta sección, se diseñó un diagrama de entidad-relación que se puede ver en `docx/DER Proyecto #2.pdf`, un modelo relacional en `docx/Modelo Relacional Proyecto #2.pdf` y también se llevó a cabo un proceso de normalización hasta 3FN que es visible en `docx/Normalización Proyecto #2.pdf`. Dentro de este último archivo, también se indican las dependencias funcionales de todo el sistema, junto con una captura del modelo relacional final luego de dicho proceso.

Con respecto a el código SQL, se realizaron tres cosas principales: un DDL basado en el modelo relacional anterior, una población de la base de datos inicial de 25 entradas por tabla y un proceso de indexado básico para las columnas. El DDL mencionado se puede encontrar dentro de `sql/ddl.sql`, y el script utilizado para crear una población base de la base de datos se encuentra también en el mismo directorio bajo `sql/pruebas.sql`. Con respecto a los índices, estos fueron definidos dentro del mismo archivo que el DDL, justo hasta el final del archivo junto con la justificación de la implementación de cada uno.


---
# Sección #2: Docker
En este proyecto se puede también encontrar un archivo `.env.example`, el cual tiene las credenciales y variables de entorno necesarias para correr el código (para no tener que subir el mismo `.env`).

Para levantar el proyecto desde cero:

```bash
cp .env.example .env
docker compose up --build
```

Servicios disponibles:

- Frontend React/Vite: `http://localhost:5173`
- Backend Go HTTP API: `http://localhost:8080`
- Endpoint de prueba del backend: `http://localhost:8080/api/health`
- PostgreSQL local: `localhost:5433`

Las credenciales obligatorias para la base de datos son:

```env
DB_USER=proy2
DB_PASSWORD=secret
```
