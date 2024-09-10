FROM node:20.11.1-alpine3.19

WORKDIR /usr/src/frontend

COPY ./frontend/package.json .

RUN npm install

COPY ./frontend .

RUN npm run build
CMD [ "npm", "run" , "start" ]
