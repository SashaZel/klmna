# Specifies a parent image
FROM golang:1.23.3
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY . .

ENV POSTGRES_DB $POSTGRES_DB
ENV POSTGRES_USER $POSTGRES_USER
ENV POSTGRES_PASSWORD $POSTGRES_PASSWORD
ENV POSTGRES_HOST $POSTGRES_HOST
ENV POSTGRES_PORT $POSTGRES_PORT
 
# # Installs Go dependencies
RUN go mod download
 
# # Builds your app with optional configuration
RUN go build -o klmna

# # Tells Docker which network port your container listens on
EXPOSE 80
 
# Specifies the executable command that runs when the container starts
CMD [ "./klmna" ]