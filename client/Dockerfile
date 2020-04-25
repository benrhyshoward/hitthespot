FROM node:14.0.0-alpine3.10

WORKDIR /usr/src/app

COPY package.json ./
COPY yarn.lock ./

RUN yarn

EXPOSE 3000

CMD [ "yarn", "start" ]