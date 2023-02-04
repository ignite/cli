# App Init

### Blockchain'i başlatma

Bu bölümde, chain'ler arası değişim uygulaması için temel blockchain modülünü oluşturacaksınız. Blockchain, modül, işlem, IBC paketleri ve mesajların iskeletini oluşturuyorsunuz. Daha sonraki bölümlerde, işlem işleyicilerinin her birine daha fazla kod entegre edeceksiniz.

### Blockchain'i Oluşturun

Scaffold `interchange` adında yeni bir blockchain:

```
ignite scaffold chain interchange --no-module
```

`Interchange` adında yeni bir dizin oluşturulur.

Modülleri, türleri ve eşlemeleri scaffold haline getirebileceğiniz bu dizine geçin:

```
cd interchange
```

`interchange` dizini çalışan bir blok zinciri uygulaması içerir.

İlk iskele ile sizin için yerel bir GitHub deposu oluşturuldu.

Ardından, yeni bir IBC modülü oluşturun.

### dex Modülünü Oluşturun

Blok zincirinizin içinde IBC özelliklerine sahip `dex` adında bir modül oluşturun.

Dex modülü, sipariş defterlerini oluşturma ve sürdürme ve bunları IBC aracılığıyla ikinci blok zincirine yönlendirme mantığını içerir.

```
ignite scaffold module dex --ibc --ordering unordered --dep bank
```

### Alış ve Satış Emir Defterleri için CRUD mantığı oluşturma

İki türü oluşturma, okuma, güncelleme ve silme (CRUD) eylemleriyle `scaffold` edin.

`sellOrderBook` ve `buyOrderBook` `type`'larını oluşturmak için aşağıdaki Ignite CLI türü komutlarını çalıştırın:

```
ignite scaffold map sell-order-book amountDenom priceDenom --no-message --module dex
ignite scaffold map buy-order-book amountDenom priceDenom --no-message --module dex
```

Değerler şunlardır:

* `amountDenom`: satılacak token ve hangi miktarda satılacağı&#x20;
* `priceDenom`: token satış fiyatı

`no-message` bayrağı mesaj oluşturma işleminin atlanacağını belirtir. Özel mesajlar sonraki adımlarda oluşturulacaktır.

`--module dex` bayrağı, `dex` modülündeki türün iskeletlenmesini belirtir.

### IBC Paketlerini Oluşturun

IBC için üç paket oluşturun:

* Bir sipariş defteri çifti `createPair`
* Bir satış emri `sellOrder`
* Bir satın alma emri `buyOrder`

```
ignite scaffold packet create-pair sourceDenom targetDenom --module dex
ignite scaffold packet sell-order amountDenom amount:int priceDenom price:int --ack remainingAmount:int,gain:int --module dex
ignite scaffold packet buy-order amountDenom amount:int priceDenom price:int --ack remainingAmount:int,purchase:int --module dex
```

İsteğe bağlı `--ack` bayrağı, paket hedef zincir tarafından alındıktan sonra döndürülen onaylamanın alan adlarını ve türlerini tanımlar. `--ack` bayrağının değeri virgülle ayrılmış bir ad listesidir (boşluk yok). İsteğe bağlı türleri iki nokta üst üste (`:`) işaretinden sonra ekleyin.

### İptal iletileri

Emirlerin iptali ağda yerel olarak yapılır, gönderilecek bir paket yoktur.

Bir satış veya alış emrini iptal etmek üzere bir mesaj oluşturmak için `message` komutunu kullanın:

```
ignite scaffold message cancel-sell-order port channel amountDenom priceDenom orderID:int --desc "Cancel a sell order" --module dex
ignite scaffold message cancel-buy-order port channel amountDenom priceDenom orderID:int --desc "Cancel a buy order" --module dex
```

İletiyle birlikte bir işlemi yayınlamak için kullanılan CLI komutunun açıklamasını tanımlamak için isteğe bağlı `--desc` bayrağını kullanın.

### Denom'u İzleme

Token demonları `ibc-transfer` modülünde açıklandığı gibi aynı davranışa sahip olmalıdır:

* Bir chain'den alınan harici bir token, `voucher` olarak adlandırılan benzersiz bir `denom`a sahiptir.
* Bir token bir blockchain'e gönderildiğinde ve daha sonra geri gönderildiğinde ve alındığında, chain kuponu çözebilir ve orijinal token değerine geri dönüştürebilir.

`Voucher` tokenlar hash olarak temsil edilir, bu nedenle hangi orijinal kuponun bir voucher ile ilişkili olduğunu saklamanız gerekir. Bunu indeksli bir tip ile yapabilirsiniz.

Sakladığınız bir `voucher` için kaynak port ID'sini, kaynak kanal ID'sini ve orijinal kupon değerini tanımlayın:

```
ignite scaffold map denom-trace port channel origin --no-message --module dex
```

### İki Blok Zinciri için Yapılandırma Oluşturun

Her biri için özel token içeren iki blockchain ağını test etmek için iki yapılandırma dosyası `mars.yml` ve `venus.yml` ekleyin.

Yapılandırma dosyalarını `interchange` klasörüne ekleyin.

Mars için yerel denomlar `marscoin` ve Venüs için `venuscoin`'dir.

İçeriğinizle birlikte `mars.yml` dosyasını oluşturun:

mars.yml

```
version: 1
build:
  proto:
    path: proto
    third_party_paths:
    - third_party/proto
    - proto_vendor
accounts:
- name: alice
  coins:
  - 1000token
  - 100000000stake
  - 1000marscoin
- name: bob
  coins:
  - 500token
  - 1000marscoin
  - 100000000stake
faucet:
  name: bob
  coins:
  - 5token
  - 100000stake
  host: 0.0.0.0:4500
genesis:
  chain_id: mars
validators:
- name: alice
  bonded: 100000000stake
  home: $HOME/.mars
```

İçeriğinizle `venus.yml` dosyasını oluşturun:

venus.yml

```
version: 1
build:
  proto:
    path: proto
    third_party_paths:
    - third_party/proto
    - proto_vendor
accounts:
- name: alice
  coins:
  - 1000token
  - 1000000000stake
  - 1000venuscoin
- name: bob
  coins:
  - 500token
  - 1000venuscoin
  - 100000000stake
faucet:
  name: bob
  coins:
  - 5token
  - 100000stake
  host: :4501
genesis:
  chain_id: venus
validators:
- name: alice
  bonded: 100000000stake
  app:
    api:
      address: :1318
    grpc:
      address: :9092
    grpc-web:
      address: :9093
  config:
    p2p:
      laddr: :26658
    rpc:
      laddr: :26659
      pprof_laddr: :6061
  home: $HOME/.venus
```

İki blockchain'i tek bir makinede yan yana çalıştırmak için bunları farklı portlarda başlatmanız gerekir. `venus.yml`, HTTP API, gRPC, P2P ve RPC servislerini özel portlarda başlatan bir validators yapılandırmasına sahiptir.

İskele kurduktan sonra, şimdi sizin için oluşturulan yerel GitHub deposuna bir taahhütte bulunmak için iyi bir zaman.

```
git add .
git commit -m "Scaffold module, maps, packages and messages for the dex"
```

Sipariş defteri için kodu bir sonraki bölümde uygulayın.
