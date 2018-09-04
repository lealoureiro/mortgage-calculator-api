#!/usr/bin/env bash

echo "Building Application for AWS Linux Platform..."
./build.sh

echo "Creating application artifact..."
zip -r AppBundle.zip Procfile ./bin/mortgage-calculator-api

echo "Deploying to AWS Elastic Beanstalk..."
eb deploy mortgage-calculator-api

echo "Cleaning deploying files..."
rm -rv AppBundle.zip