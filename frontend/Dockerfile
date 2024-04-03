# Node runtime
FROM node:20-alpine
# Frontend docker deployment responsible
LABEL maintainer="andreas.malde@lyse.net"
LABEL deployedBy="Andreas Follevaag Malde"
# Change working directory
WORKDIR /app
# Copy relevant elements of the project to the docker container
COPY package*.json ./
COPY ./src ./src
COPY ./public ./public
COPY .env .env
# Install node dependencies, make production build and
# delete node_modules afterward 
RUN npm install \
    && npm install -g serve\
    && npm run build \
    && rm -fr node_modules
# Expose port 3000 internally
EXPOSE 3000
# Start the production server
CMD [ "serve","-s","build" ]