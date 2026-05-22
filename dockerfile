# Use the official lightweight Caddy image
FROM caddy:latest

# Copy your configuration file
COPY Caddyfile /etc/caddy/Caddyfile

# Copy your static site content to the default Caddy server root
COPY ./dist /srv
