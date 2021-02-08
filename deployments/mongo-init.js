
let products = db.createCollection("products")

db.products.createIndex(
  {
      "name": 1
  },
  {
      unique: true,
  }
)


db.products.createIndex(
  {
      "lastUpdate": 1,
      "lastUpdate": -1,
  }
)

