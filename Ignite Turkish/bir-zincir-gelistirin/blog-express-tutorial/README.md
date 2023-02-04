# Blog, Express Tutorial

### 5 dakikada "bir blog oluşturun"

Bu eğitimde, blok zincirine veri yazmamızı ve blok zincirinden veri okumamızı sağlayan bir modül ile bir blok zinciri oluşturacağız. Bu modül, bir blog uygulamasına benzer şekilde blog gönderileri oluşturma ve okuma yeteneğini uygulayacaktır. Son kullanıcı yeni blog gönderileri gönderebilecek ve blok zincirindeki mevcut gönderilerin bir listesini görüntüleyebilecek. Bu eğitim, blok zinciri ile etkileşim kurmak için bu modülü oluşturma ve kullanma sürecinde size rehberlik edecektir.

Bu eğitimin amacı, blok zincirine veri göndermenize ve bu verileri blok zincirinden geri okumanıza olanak tanıyan bir geri bildirim döngüsü oluşturmak için adım adım talimatlar sağlamaktır. Bu eğitimin sonunda, eksiksiz bir geri bildirim döngüsü oluşturmuş olacak ve bunu blok zinciri ile etkileşim kurmak için kullanabileceksiniz.

İlk olarak, Ignite CLI ile yeni bir blog blok zinciri oluşturun:

```
ignite scaffold chain blog
```

Blok zinciri kullanan bir blog uygulaması oluşturmak için uygulamamızın gereksinimlerini tanımlamamız gerekir. Uygulamanın `Post` tipindeki nesneleri blok zincirinde saklamasını istiyoruz. Bu nesneler iki özelliğe sahip olmalıdır: bir `title` ve bir `body`.

Gönderileri blok zincirinde depolamanın yanı sıra, kullanıcılara bu gönderiler üzerinde CRUD (oluşturma, okuma, güncelleme ve silme) işlemleri gerçekleştirme olanağı da sağlamak istiyoruz. Bu, kullanıcıların yeni gönderiler oluşturmasına, mevcut gönderileri okumasına, mevcut gönderilerin içeriğini güncellemesine ve artık ihtiyaç duyulmayan gönderileri silmesine olanak tanıyacaktır.

Ignite CLI'nin özelliklerinden biri, temel CRUD işlevselliğini uygulayan kod üretme yeteneğidir. Bu, uygulamanızda veri oluşturmak, okumak, güncellemek ve silmek için gerekli kodu hızlı bir şekilde oluşturmak için kullanılabilen iskele komutlarının kullanılmasıyla gerçekleştirilir.

Ignite CLI, farklı veri yapılarında depolanan veriler için kod üretme yeteneğine sahiptir. Bunlar arasında, artan bir tamsayı tarafından indekslenen veri koleksiyonları olan listeler, özel bir anahtar tarafından indekslenen koleksiyonlar olan eşlemeler ve tek veri örnekleri olan tekler yer alır. Bu farklı veri yapılarını kullanarak uygulamanızı özel ihtiyaçlarınıza uyacak şekilde özelleştirebilirsiniz. Örneğin, bir blog uygulaması oluşturuyorsanız, tüm gönderileri depolamak için her gönderinin bir tamsayı tarafından indekslendiği bir liste kullanmak isteyebilirsiniz. Alternatif olarak, her gönderiyi benzersiz başlığına göre dizinlemek için bir harita veya tek bir gönderiyi depolamak için bir single kullanabilirsiniz. Veri yapısı seçimi, uygulamanızın özel gereksinimlerine bağlı olacaktır.

Seçtiğiniz veri yapısına ek olarak Ignite CLI, kod oluşturacağı veri türünün adını ve veri türünü tanımlayan alanları da sağlamanızı gerektirir. Örneğin, bir blog uygulaması oluşturuyorsanız, yazının "başlığı" ve "gövdesi" için alanları olan "Yazı" adlı bir tür oluşturmak isteyebilirsiniz. Ignite CLI bu bilgileri kullanarak uygulamanızda bu türden verilerin oluşturulması, okunması, güncellenmesi ve silinmesi için gerekli kodu oluşturacaktır.

Blog dizinine geçin ve ignite scaffold list komutunu çalıştırın:

```
cd blog
ignite scaffold list post title body
```

Uygulamanız için kod oluşturmak üzere Ignite CLI'yı kullandığınıza göre, şimdi CLI'nın ne oluşturduğunu gözden geçirelim. Ignite CLI, belirttiğiniz veri yapısı ve veri türü için kodun yanı sıra bu verileri işlemek için gereken temel CRUD işlemleri için de kod oluşturmuş olacaktır. Bu kod, uygulamanız için sağlam bir temel oluşturacaktır ve özel ihtiyaçlarınıza uyacak şekilde daha da özelleştirebilirsiniz. Ignite CLI tarafından oluşturulan kodu inceleyerek gereksinimlerinizi karşıladığından emin olabilir ve bu aracı kullanarak uygulamanızı nasıl oluşturacağınızı daha iyi anlayabilirsiniz.

Ignite CLI, `proto/blog/blog` dizininde çeşitli dosyalar ve değişiklikler oluşturmuştur. Bunlar şunları içerir:

* `post.proto`: Bu, `Post` türünü tanımlayan ve `title`, `body`, `id` ve `creator` alanlarını içeren bir protokol buffer dosyasıdır.
* `tx.proto`: Bu dosya üç RPC (uzaktan yordam çağrısı) içerecek şekilde değiştirilmiştir: CreatePost, UpdatePost ve DeletePost. Bu RPC'lerin her biri, bir gönderi üzerinde ilgili CRUD işlemini gerçekleştirmek için kullanılabilecek bir Cosmos SDK mesajına karşılık gelir.
* `query.proto`: Bu dosya iki sorgu içerecek şekilde değiştirilmiştir: `Post` ve `PostAll`. Post sorgusu, `ID`'sine göre tek bir gönderiyi almak için kullanılabilirken, `PostAll` sorgusu sayfalandırılmış bir gönderi listesini almak için kullanılabilir.
* `genesis.proto`: Bu dosya, blok zincirinin ilk başlatıldığında başlangıç durumunu tanımlayan modülün genesis durumuna gönderileri dahil edecek şekilde değiştirilmiştir.

Ignite CLI ayrıca `x/blog/keeper` dizininde uygulamanız için CRUD'a özgü mantığı uygulayan birkaç yeni dosya oluşturdu. Bunlar şunları içerir:

* msg\_server\_post.go: Bu dosya `CreatePost`, `UpdatePost` ve DeletePost mesajları için keeper yöntemlerini uygular. Bu yöntemler, ilgili bir mesaj modül tarafından işlendiğinde çağrılır ve CRUD işlemlerinin her biri için özel mantığı ele alır.
* query\_post.go: Bu dosya, sırasıyla ID'ye göre tek tek gönderileri veya sayfalandırılmış bir gönderi listesini almak için kullanılan `Post` ve `PostAll` sorgularını uygular.
* post.go: Bu dosya, keeper yöntemlerinin bağlı olduğu temel işlevleri uygular. Bu işlevler, gönderileri depoya ekleme, tek tek gönderileri alma, gönderi sayısını alma ve uygulamadaki gönderileri yönetmek için gereken diğer işlemleri içerir.

Genel olarak, bu dosyalar blog uygulamanızın CRUD işlevselliği için gerekli uygulamayı sağlar. CRUD işlemlerinin her biri için özel mantığın yanı sıra bu işlemlerin bağlı olduğu temel işlevleri de ele alırlar.

Dosyalar `x/blog/types` dizininde oluşturulmuş ve değiştirilmiştir.

* messages\_post.go: Bu yeni dosya Cosmos SDK mesaj kurucularını ve `Route()`, `Type()`, `GetSigners()`, `GetSignBytes()` ve `ValidateBasic()` gibi ilişkili yöntemleri içerir.
* keys.go: Bu dosya, blog gönderilerini depolamak için anahtar önekleri içerecek şekilde değiştirildi. Anahtar önekleri kullanarak, blog gönderilerimizin verilerinin veritabanındaki diğer veri türlerinden ayrı tutulmasını ve gerektiğinde bunlara kolayca erişilebilmesini sağlayabiliriz.
* genesis.go: Bu dosya, blog modülünün başlangıç (genesis) durumunu ve bu başlangıç durumunu doğrulamak için `Validate()` işlevini tanımlamak üzere değiştirildi. Bu, başlangıç verilerini tanımladığı ve uygulamamızın kurallarına göre geçerli olmasını sağladığı için blok zincirimizin kurulumunda önemli bir adımdır.
* codec.go: Bu dosya, mesaj türlerimizi kodlayıcıya kaydetmek için değiştirildi, böylece ağ üzerinden iletildiklerinde düzgün bir şekilde serileştirilmeleri ve serileştirilmeleri sağlandı.

Ayrıca, `*.proto` dosyalarından `*.pb.go` dosyaları oluşturuldu ve bunlar uygulamamız tarafından kullanılan mesajlar, RPC'ler ve sorgular için tip tanımları içeriyor. Bu dosyalar, verilerimizin yapısını dilden bağımsız bir şekilde tanımlamamızı sağlayan Protokol Buffers (protobuf) aracı kullanılarak `*.proto` dosyalarından otomatik olarak oluşturulur.

Ignite CLI, birkaç dosya oluşturarak ve değiştirerek `x/blog/client/cli` dizinine işlevsellik eklemiştir.

* tx\_post.go: Bu dosya, blog modülü için mesaj içeren işlemlerin yayınlanması için CLI komutlarını uygulamak üzere oluşturulmuştur. Bu komutlar, kullanıcıların Ignite CLI kullanarak blok zincirine kolayca mesaj göndermelerini sağlar.
* query\_post.go: Bu dosya, blog modülünü sorgulamak için CLI komutlarını uygulamak üzere oluşturulmuştur. Bu komutlar, kullanıcıların blok zincirinden blog gönderilerinin listesi gibi bilgileri almasına olanak tanır.
* tx.go: Bu dosya, işlemlerin zincirin ikili yapısına yayınlanması için CLI komutlarını eklemek üzere değiştirildi.
* query.go: Bu dosya da zinciri sorgulamaya yönelik CLI komutlarını zincirin binary'sine eklemek için değiştirildi.

Gördüğünüz gibi, `ignite scaffold list` komutu bir dizi kaynak kod dosyası oluşturmuş ve değiştirmiştir. Bu dosyalar mesaj türlerini, bir mesaj işlendiğinde yürütülecek mantığı ve her şeyi birbirine bağlayan kabloları tanımlar. Bu, blog gönderilerini oluşturma, güncelleme ve silme mantığının yanı sıra bu bilgileri almak için gereken sorguları da içerir.

Oluşturulan kodu çalışırken görmek için blok zincirini başlatmamız gerekecek. Bunu, bizim için blok zincirini oluşturacak, başlatacak ve başlatacak olan `ignite chain serve` komutunu kullanarak yapabiliriz:

```
ignite chain serve
```

Blok zinciri çalıştığında, onunla etkileşim kurmak ve kodun blog gönderilerini oluşturma, güncelleme ve silme işlemlerini nasıl gerçekleştirdiğini görmek için ikiliyi kullanabiliriz. Ayrıca sorguları nasıl işlediğini ve yanıtladığını da görebiliriz. Bu bize uygulamamızın nasıl çalıştığını daha iyi anlamamızı ve işlevselliğini test etmemizi sağlayacaktır.

Bir terminal penceresinde ignite chain serve çalışırken, başka bir terminal açın ve blok zincirinde yeni bir blog gönderisi oluşturmak için zincirin binary'sini kullanın:

```
blogd tx blog create-post 'Hello, World!' 'This is a blog post' --from alice
```

Bir işlemi imzalamak için kullanılacak hesabı belirtmek üzere `--from` bayrağını kullanırken, belirtilen hesabın kullanıma hazır olduğundan emin olmak önemlidir. Bir geliştirme ortamında, `ignite chain serve` komutunun çıktısında veya `config.yml` dosyasında kullanılabilir hesapların bir listesini görebilirsiniz.

İşlemleri yayınlarken `--from` bayrağının gerekli olduğunu da belirtmek gerekir. Bu bayrak, işlem sürecinde çok önemli bir adım olan işlemi imzalamak için kullanılacak hesabı belirtir. Geçerli bir imza olmadan, işlem blok zinciri tarafından kabul edilmeyecektir. Bu nedenle, `--from` bayrağı ile belirtilen hesabın kullanılabilir olduğundan emin olmak önemlidir.

İşlem başarıyla yayınlandıktan sonra, blog gönderilerinin listesi için blok zincirini sorgulayabilirsiniz. Bunu yapmak için, blok zincirine eklenen tüm blog gönderilerinin sayfalandırılmış bir listesini döndüren `blogd q blog list-post` komutunu kullanabilirsiniz.

```
blogd q blog list-post

Post:
- body: This is a blog post
  creator: cosmos1xz770h6g55rrj8vc9ll9krv6mr964tzhqmsu2v
  id: "0"
  title: Hello, World!
pagination:
  next_key: null
  total: "0"
```

Blok zincirini sorgulayarak işleminizin başarıyla işlendiğini ve blog gönderisinin zincire eklendiğini doğrulayabilirsiniz. Ayrıca, hesaplar, bakiyeler ve yönetim teklifleri gibi blok zincirindeki diğer veriler hakkında bilgi almak için diğer sorgu komutlarını kullanabilirsiniz.

Az önce oluşturduğumuz blog gönderisini gövde içeriğini değiştirerek modifiye edelim. Bunu yapmak için, blok zincirindeki mevcut bir blog gönderisini güncellememizi sağlayan `blogd tx blog update-post` komutunu kullanabiliriz. Bu komutu çalıştırırken, değiştirmek istediğimiz blog gönderisinin ID'sini ve kullanmak istediğimiz yeni gövde içeriğini belirtmemiz gerekecektir. Bu komutu çalıştırdıktan sonra, işlem blok zincirine yayınlanacak ve blog yazısı yeni gövde içeriğiyle güncellenecektir.

```
blogd tx blog update-post 0 'Hello, World!' 'This is a blog post from Alice' --from alice
```

Artık blog gönderisini yeni içerikle güncellediğimize göre, değişiklikleri görmek için blok zincirini tekrar sorgulayalım. Bunu yapmak için, blok zincirindeki tüm blog gönderilerinin bir listesini döndürecek olan `blogd q blog list-post` komutunu kullanabiliriz. Bu komutu tekrar çalıştırarak, güncellenmiş blog gönderisini listede görebilir ve yaptığımız değişikliklerin blok zincirine başarıyla uygulandığını doğrulayabiliriz.

```
blogd q blog list-post

Post:
- body: This is a blog post from Alice
  creator: cosmos1xz770h6g55rrj8vc9ll9krv6mr964tzhqmsu2v
  id: "0"
  title: Hello, World!
pagination:
  next_key: null
  total: "0"
```

Bob'un hesabını kullanarak blog gönderilerinden birini silmeye çalışalım. Ancak, blog gönderisi Alice'in hesabı kullanılarak oluşturulduğundan, blok zincirinin kullanıcının gönderiyi silme yetkisine sahip olup olmadığını kontrol etmesini bekleyebiliriz. Bu durumda, Bob yazının yazarı olmadığından, işlemi blok zinciri tarafından reddedilmelidir.

Bir blog gönderisini silmek için, blok zincirindeki mevcut bir blog gönderisini silmemizi sağlayan `blogd tx blog delete-post` komutunu kullanabiliriz. Bu komutu çalıştırırken, silmek istediğimiz blog gönderisinin kimliğinin yanı sıra işlemi imzalamak için kullanmak istediğimiz hesabı da belirtmemiz gerekecektir. Bu durumda, işlemi imzalamak için Bob'un hesabını kullanacağız.

Bu komutu çalıştırdıktan sonra, işlem blok zincirinde yayınlanacaktır. Ancak, Bob yazının yazarı olmadığından, blok zinciri onun işlemini reddetmelidir ve blog yazısı silinmeyecektir. Bu, blok zincirinin kuralları ve izinleri nasıl uygulayabileceğine dair bir örnektir ve yalnızca yetkili kullanıcıların blok zincirinde değişiklik yapabileceğini gösterir.

```
blogd tx blog delete-post 0 --from bob

raw_log: 'failed to execute message; message index: 0: incorrect owner: unauthorized'
```

Şimdi blog gönderisini tekrar silmeyi deneyelim, ancak bu kez Alice'in hesabını kullanarak. Alice blog gönderisinin yazarı olduğu için, gönderiyi silme yetkisine sahip olmalıdır.

```
blogd tx blog delete-post 0 --from alice
```

Blog gönderisinin Alice tarafından başarılı bir şekilde silinip silinmediğini kontrol etmek için, gönderilerin bir listesi için blok zincirini tekrar sorgulayabiliriz.

```
blogd q blog list-post

Post: []
pagination:
  next_key: null
  total: "0"
```

Ignite CLI ile bir blog oluşturma eğitimini başarıyla tamamladığınız için tebrikler! Talimatları izleyerek yeni bir blok zinciri oluşturmayı, CRUD işlevselliğine sahip bir "post" türü için kod oluşturmayı, yerel bir blok zinciri başlatmayı ve blogunuzun işlevselliğini test etmeyi öğrendiniz.

Artık basit bir uygulamanın çalışan bir örneğine sahip olduğunuza göre, Ignite tarafından oluşturulan kodu deneyebilir ve değişikliklerin uygulamanın davranışını nasıl etkilediğini görebilirsiniz. Bu, uygulamanızı özel ihtiyaçlarınıza uyacak şekilde özelleştirmenize ve uygulamanızın işlevselliğini geliştirmenize olanak tanıyacağı için sahip olunması gereken değerli bir beceridir. Veri yapısında veya veri türünde değişiklikler yapmayı deneyebilir ya da koda ek alanlar veya işlevler ekleyebilirsiniz.

Aşağıdaki eğitimlerde, blok zincirlerinin nasıl oluşturulacağını daha iyi anlamak için Ignite'ın ürettiği koda daha yakından bakacağız. Kodun bir kısmını kendimiz yazarak Ignite'ın nasıl çalıştığını ve bir blok zinciri üzerinde uygulama oluşturmak için nasıl kullanılabileceğini daha iyi anlayabiliriz. Bu, Ignite CLI'nin yetenekleri ve sağlam ve güçlü uygulamalar oluşturmak için nasıl kullanılabileceği hakkında daha fazla bilgi edinmemize yardımcı olacaktır. Bu eğitimleri kaçırmayın ve Ignite ile blok zinciri dünyasının derinliklerine dalmaya hazır olun!
