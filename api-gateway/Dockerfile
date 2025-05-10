FROM traefik:v2.10

# Configuración para Traefik
COPY ./traefik.yml /etc/traefik/traefik.yml
COPY ./dynamic_conf /etc/traefik/dynamic_conf

# Exponer puertos
EXPOSE 80
EXPOSE 8080

# Establecer etiquetas de metadatos
LABEL maintainer="GrowDesk Team" \
      description="API Gateway para GrowDesk V2" \
      version="1.0"

# Comando de inicio predeterminado (usar configuración predefinida)
CMD ["--configFile=/etc/traefik/traefik.yml"] 