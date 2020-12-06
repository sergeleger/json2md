# World Bank Data

See [The World Bank](https://www.worldbank.org/) -
[Country API Queries](https://datahelpdesk.worldbank.org/knowledgebase/articles/898590-country-api-queries).


Created with the following execution of `json2md`

```shell
$ curl -s "http://api.worldbank.org/v2/countries?per_page=10&format=json" | \
    json2md worldbank.md.template - > worldbank.md
```



## First 10 Entries

| ISO Code | Country | Capital |
| -------- | ------- | ------- |
| AW | Aruba | Oranjestad |
| AF | Afghanistan | Kabul |
| A9 | Africa |  |
| AO | Angola | Luanda |
| AL | Albania | Tirane |
| AD | Andorra | Andorra la Vella |
| L5 | Andean Region |  |
| 1A | Arab World |  |
| AE | United Arab Emirates | Abu Dhabi |
| AR | Argentina | Buenos Aires |

