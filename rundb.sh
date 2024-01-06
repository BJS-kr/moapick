#!/bin/bash
docker run --rm -d --name test_pg -e POSTGRES_PASSWORD=test -v dbdata:/var/lib/postgresql/data -p 5432:5432 postgres