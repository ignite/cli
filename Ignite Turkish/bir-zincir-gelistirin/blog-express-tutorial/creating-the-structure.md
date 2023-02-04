# Creating the structure

Aşağıdaki komutla yeni bir blok zinciri oluşturun:

```
ignite scaffold chain blog
```

Bu, [blok zinciri uygulamanız](https://docs.cosmos.network/main/basics/app-anatomy) için gerekli dosya ve dizinleri içeren `blog/` adında yeni bir dizin oluşturacaktır. Ardından, çalıştırarak yeni oluşturulan dizine gidin:

```
cd blog
```

Uygulamanız blog gönderilerini depolayacağından ve bunlarla çalışacağından, bu gönderileri temsil etmek için bir `Post` türü oluşturmanız gerekecektir. Bunu aşağıdaki Ignite CLI komutunu kullanarak yapabilirsiniz:

```
ignite scaffold type post title body creator id:uint
```

Bu, dört alana sahip bir Yazı türü oluşturacaktır: `title`, `body`, `creator`, hepsi `string` ve `id uint` türünde.

Ignite'ın kod iskelesi komutlarını kullandıktan sonra değişikliklerinizi Git gibi bir sürüm kontrol sistemine işlemek iyi bir uygulamadır. Bu, Ignite tarafından otomatik olarak yapılan değişiklikler ile geliştiriciler tarafından manuel olarak yapılan değişiklikler arasında ayrım yapmanıza ve gerekirse değişiklikleri geri almanıza olanak tanır. Aşağıdaki komutları kullanarak değişikliklerinizi Git'e işleyebilirsiniz:

```
git add .
git commit -am "ignite scaffold type post title body"
```

### Mesaj oluşturma

Daha sonra, blog gönderileriniz için CRUD (oluşturma, okuma, güncelleme ve silme) işlemlerini uygulayacaksınız. Oluşturma, güncelleme ve silme işlemleri uygulamanın durumunu değiştirdiğinden, yazma işlemleri olarak kabul edilirler. Cosmos SDK blok zincirlerinde durum, durum geçişlerini tetikleyen mesajlar içeren işlemlerin [yayınlanmasıyla](https://docs.cosmos.network/main/basics/tx-lifecycle) değiştirilir. "Create post" mesajı içeren işlemleri yayınlama ve işleme mantığını oluşturmak için aşağıdaki Ignite CLI komutunu kullanabilirsiniz:

```
ignite scaffold message create-post title body --response id:uint
```

Bu, her ikisi de `string` türünde olan `title` ve `body` olmak üzere iki alana sahip bir "gönderi oluştur" mesajı oluşturacaktır. Gönderiler anahtar-değer deposunda liste benzeri bir veri yapısında saklanacak ve burada artan bir tamsayı ID ile indekslenecektir. Yeni bir gönderi oluşturulduğunda, ona bir ID tamsayısı atanacaktır. `--response` bayrağı, "create post" mesajına yanıt olarak `uint` türünde bir kimlik döndürmek için kullanılır.

Uygulamanızda belirli bir blog gönderisini güncellemek için, üç argüman kabul eden "update post" adlı bir mesaj oluşturmanız gerekecektir: başlık, gövde ve id. Uint türündeki id bağımsız değişkeni, hangi blog gönderisini güncellemek istediğinizi belirtmek için gereklidir. Bu mesajı Ignite CLI komutunu kullanarak oluşturabilirsiniz:

```
ignite scaffold message update-post title body id:uint
```

Uygulamanızda belirli bir blog gönderisini silmek için, yalnızca silinecek gönderinin id'sini kabul eden "gönderiyi sil" adlı bir mesaj oluşturmanız gerekecektir. Bu mesajı Ignite CLI komutunu kullanarak oluşturabilirsiniz:

```
ignite scaffold message delete-post id:uint
```

### Sorgu oluşturma

[Sorgular](https://docs.cosmos.network/main/basics/query-lifecycle), kullanıcıların blok zinciri durumundan bilgi almalarını sağlar. Uygulamanızda iki sorguya sahip olacaksınız: "show post" ve "list post". "show post" sorgusu, kullanıcıların belirli bir gönderiyi kimliğine göre almasına izin verirken, "list post" sorgusu depolanan tüm gönderilerin sayfalandırılmış bir listesini döndürür.

"show post" sorgusunu oluşturmak için aşağıdaki Ignite CLI komutunu kullanabilirsiniz:

```
ignite scaffold query show-post id:uint --response post:Post
```

Bu sorgu, argüman olarak uint türünde bir id kabul edecek ve yanıt olarak Post türünde bir post döndürecektir.

"list post" sorgusunu oluşturmak için aşağıdaki Ignite CLI komutunu kullanabilirsiniz:

```
ignite scaffold query list-post --response post:Post --paginated
```

Bu sorgu, Post türündeki bir gönderiyi sayfalandırılmış bir çıktı olarak döndürür. `--paginated` bayrağı, sorgunun sonuçlarını sayfalandırılmış bir biçimde döndürmesi gerektiğini belirtir ve kullanıcıların bir seferde belirli bir sonuç sayfasını almasına olanak tanır.

### Özet

Blok zinciri uygulamanızın ilk kurulumunu tamamladığınız için tebrikler! Bir "post" veri türünü başarıyla oluşturdunuz ve üç tür mesaj (oluşturma, güncelleme ve silme) ile iki tür sorguyu (mesajları listeleme ve gösterme) işlemek için gerekli kodu oluşturdunuz.

Ancak bu noktada, oluşturduğunuz mesajlar herhangi bir durum geçişini tetiklemeyecek ve oluşturduğunuz sorgular herhangi bir sonuç döndürmeyecektir. Bunun nedeni, Ignite'ın bu özellikler için yalnızca şablon kodu oluşturması ve bunları işlevsel hale getirmek için gerekli mantığı uygulamanın size kalmış olmasıdır.

Eğitimin sonraki bölümlerinde, blok zinciri uygulamanızı tamamlamak için mesaj işleme ve sorgu mantığını nasıl uygulayacağınızı öğreneceksiniz. Bu, oluşturduğunuz mesajları ve sorguları işlemek için kod yazmayı ve bunları blok zincirinin durumundan veri değiştirmek veya almak için kullanmayı içerecektir. Bu sürecin sonunda, Cosmos SDK blok zinciri üzerinde tamamen işlevsel bir blog uygulamasına sahip olacaksınız.
