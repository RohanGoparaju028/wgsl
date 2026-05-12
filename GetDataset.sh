#!/bin/bash
curl -L  -o leukemia-gene-expression-cumida.zip\
  https://www.kaggle.com/api/v1/datasets/download/brunogrisci/leukemia-gene-expression-cumida

chmod +x

unzip leukemia-gene-expression-cumida.zip

mv Leukemia_GSE9476.csv Leukemia.csv

rm leukemia-gene-expression-cumida.zip

pip3 install -r requirments.txt
