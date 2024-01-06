#!/bin/bash
docker run --rm -d --name test_pg -e POSTGRES_PASSWORD=test -v dbdata:/var/lib/postgresql/data -p 5433:5432 postgres