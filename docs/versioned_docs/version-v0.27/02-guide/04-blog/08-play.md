# Play

## Create a blog post by Alice

```
blogd tx blog create-post hello world --from alice
```

## Show a blog post

```
blogd q blog show-post 0
```

```yml
post:
  body: world
  creator: cosmos1x33ummgkjdd6h2frlugt3tft7vnc0nxyfxnx9h
  id: "0"
  title: hello
```

## Create a blog post by Bob

```
blogd tx blog create-post foo bar --from bob
```

## List all blog posts with pagination

```
blogd q blog list-post       
```

```yml               
pagination:
  next_key: null
  total: "2"
post:
- body: world
  creator: cosmos1x33ummgkjdd6h2frlugt3tft7vnc0nxyfxnx9h
  id: "0"
  title: hello
- body: bar
  creator: cosmos1ysl9ws3fdamrrj4fs9ytzrrzw6ul3veddk7gz3
  id: "1"
  title: foo
```

## Update a blog post

```
blogd tx blog update-post hello cosmos 0 --from alice
```

```
blogd q blog show-post 0
```

```yml
post:
  body: cosmos
  creator: cosmos1x33ummgkjdd6h2frlugt3tft7vnc0nxyfxnx9h
  id: "0"
  title: hello
```

## Delete a blog post

```
blogd tx blog delete-post 0 --from alice
```

```
blogd q blog list-post
```

```yml
pagination:
  next_key: null
  total: "1"
post:
- body: bar
  creator: cosmos1ysl9ws3fdamrrj4fs9ytzrrzw6ul3veddk7gz3
  id: "1"
  title: foo
```

## Delete a blog post unsuccessfully

```
blogd tx blog delete-post 1 --from alice
```

```yml
raw_log: 'failed to execute message; message index: 0: incorrect owner: unauthorized'
```