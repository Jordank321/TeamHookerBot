FROM lambci/lambda:build-nodejs8.10

ENV AWS_DEFAULT_REGION us-east-1

COPY . .

RUN npm install

# Assumes you have a .lambdaignore file with a list of files you don't want in your zip
RUN cat .lambdaignore | xargs zip -9qyr lambda.zip . -x

CMD aws lambda update-function-code --function-name mylambda --zip-file fileb://lambda.zip

# docker build -t mylambda .
# docker run --rm -e AWS_ACCESS_KEY_ID -e AWS_SECRET_ACCESS_KEY mylambda