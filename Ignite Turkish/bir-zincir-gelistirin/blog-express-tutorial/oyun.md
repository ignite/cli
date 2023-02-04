# Oyun

Versiyon: v0.26.1

### Alice tarafından bir blog yazısı oluşturun

```
blogd tx blog create-post hello world --from alice
```

### Bir blog gönderisi gösterin

```
post:
  body: world
  creator: cosmos1x33ummgkjdd6h2frlugt3tft7vnc0nxyfxnx9h
  id: "0"
  title: hello
```

### Bob tarafından bir blog yazısı oluşturun

```
blogd tx blog create-post foo bar --from bob
```

### Tüm blog gönderilerini sayfalandırma ile listeleme[​](https://docs.ignite.com/guide/blog/play#list-all-blog-posts-with-pagination) <a href="#list-all-blog-posts-with-pagination" id="list-all-blog-posts-with-pagination"></a>

```
blogd q blog list-post       
```

```
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

### Update a blog post[​](broken-reference) <a href="#update-a-blog-post" id="update-a-blog-post"></a>

```
blogd tx blog update-post hello cosmos 0 --from alice
```

```
post:
  body: cosmos
  creator: cosmos1x33ummgkjdd6h2frlugt3tft7vnc0nxyfxnx9h
  id: "0"
  title: hello
```

### Delete a blog post[​](broken-reference) <a href="#delete-a-blog-post" id="delete-a-blog-post"></a>

```
blogd tx blog update-post hello cosmos 0 --from alice
```

```
blogd q blog show-post 0
```

```
post:
  body: cosmos
  creator: cosmos1x33ummgkjdd6h2frlugt3tft7vnc0nxyfxnx9h
  id: "0"
  title: hello
```

### Blog gönderisini silme

```
blogd tx blog delete-post 0 --from alice
```

```
blogd q blog list-post
```

```
pagination:
  next_key: null
  total: "1"
post:
- body: bar
  creator: cosmos1ysl9ws3fdamrrj4fs9ytzrrzw6ul3veddk7gz3
  id: "1"
  title: foo
```

### Başarısız bir blog gönderisini silme

```
blogd tx blog delete-post 1 --from alice
```

```
raw_log: 'failed to execute message; message index: 0: incorrect owner: unauthorized'
```
