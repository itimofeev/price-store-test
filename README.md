# Тестовое задание Lead Golang Engineer
Copyright © Tensigma 
LTD Contact: hello@atlant.io


Требуется написать gRPC-сервер на языке GoLang (1.13+), с постоянным хранилищем MongoDB, реализующий 2 метода:
- Fetch(URL) - запросить внешний CSV-файл со списком продуктов по внешнему адресу. 
  CSV-файл имеет вид PRODUCT NAME;PRICE. 
  Последняя цена каждого продукта должна быть сохранена в базе с датой запроса. 
  Также нужно сохранять количество изменений цены продукта.
- List(<paging params>, <sorting params>) - получить постраничный список продуктов с их ценами, количеством изменений цены и датами их последнего обновления. 
  Предусмотреть все варианты сортировки для реализации интерфейса в виде бесконечного скролла.

Сервер должен быть запущен в 2+ экземплярах (каждый в своем Docker-контейнере) и закрыт балансировщиком, соответствующие конфигурации также должны быть предоставлены для тестовой среды.


# Getting started
You need to have installed make, docker and docker-compose.
Having that, just run `make deploy` that will
* deploy mongo and postgres
* build and deploy in two instances price-store app
* deploy nginx that will round robin requests to one of price-store instance
