#!/usr/bin/env bash

echo "Building Application for AWS Linux Platform..."
./build.sh

if [[ ! -f ./bin/mortgage-calculator-api ]] ; then
    echo 'Application binary not found!'
    exit
fi

echo "Creating application artifact..."
zip -r AppBundle.zip Procfile ./bin/mortgage-calculator-api

echo "Deploying to AWS Elastic Beanstalk..."
eb deploy mortgage-calculator-api

echo "Cleaning deployment files..."
rm -rv AppBundle.zip
rm -rv ./bin/mortgage-calculator-api