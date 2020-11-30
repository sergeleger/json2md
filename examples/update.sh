#!/usr/bin/env bash

curl -s "https://api.nobelprize.org/v1/prize.json" | ../bin/md2json nobel.md.template - > nobel.md
curl -s "https://swapi.dev/api/films/" | ../bin/md2json starwars.md.template - > starwars.md
curl -s "http://api.worldbank.org/v2/countries?per_page=10&format=json" | ../bin/md2json worldbank.md.template - > worldbank.md