version: '3.8'
services:
  mysql:
    image: mysql:8.0
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: 'root_password'
      MYSQL_DATABASE: 'mt_order'
    ports:
      - "3306:3306"
    volumes:
      - ./server/migrations:/docker-entrypoint-initdb.d
  redis:
    image: redis:7-alpine
    restart: always
    ports:
      - "6379:6379"