#local
docker build -t my-lambda --target local .
docker run -d -p 9000:8080 my-lambda

#prod
docker build -t my-lambda-prod --target prod .   
docker run -d  --env-file .env -p 9000:8080 kavalife-erp-backend