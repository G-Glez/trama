# Start

TRAMA (Tournament Records and Metrics Assistant) es un servicio que se utiliza para guardar datos sobre partidas, crear pairings, guardar datos de tornos... de warhammer 40k, igual en el futuro de más juegos.

La idea es que esto sea un Back-end en golang, que guardará los datos en una base de datos SQL

## Setup

Debes crear el setup del proyecto siguiendo estos pasos. Pregunta uno a uno por cada paso

- Crear un git init para usar git como control de versiones
- Empezar un BE en golang. Usa Gin como framework. Crea un fichero de config de variables de entorno
- Añadir una bbdd SQLite al repo. Añadela al repo pero que no pueda cambiar en la config de git, que esté siempre vacía.
- Haz un dockerfile para poder correr el servicio y que exponga mediante el puerto 8080 de localhost la api
- Genera un endpoint de ejemplo con un "Hola mundo"
- Dame alternativas para testing de esta api, con ejemplos. No se si se peude ejecutar openapi desde go en funcion del entorno etc

Entre cada tarea harás un commit usando conventional commits donde:

new: setup - para el primer commit
feat: descripcion de la feature - para las features
fix: descripcion del fix - para los fixes