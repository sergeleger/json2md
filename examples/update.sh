curl -s "https://api.nobelprize.org/v1/prize.json" | ../bin/json2md nobel.md.template - > nobel.md
curl -s "https://swapi.dev/api/films/" | ../bin/json2md starwars.md.template - > starwars.md
curl -s "http://api.worldbank.org/v2/countries?per_page=10&format=json" | ../bin/json2md worldbank.md.template - > worldbank.md
