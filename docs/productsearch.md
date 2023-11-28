# Product search

1. create bucket for product data
2. upload product data csv to bucket
3. create store for relevant product fields with vectors, meta data product name, product id
4. add pipeline and map fields
5. create query with user input, sim search on fields with weights, configure treshold and top 5 products, add chatcomp with fields from top results
6. expose query as http api