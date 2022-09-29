# traer la imagen de postgres (hacer pull, e instalar contenedor)
FROM postgres:10.3 

# copiar el archivo up.sql dentro del directorio /docker-entrypoint-initdb.d y será lo que se ejecutará como 1.sql
COPY up.sql /docker-entrypoint-initdb.d/1.sql

# Ejecutar el comando postgres para inicializar la db
CMD ["postgres"]