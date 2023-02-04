# Interchain Borsasını kullanın

Bu bölümde, borsa ve hayata geçirildiğinde nasıl işleyeceği hakkında bilgi edineceksiniz. Bu, ilerleyen bölümlerde ne inşa edeceğinizi daha iyi anlamanızı sağlayacaktır.

Bunu başarmak için aşağıdaki görevleri yerine getireceğiz:

* İki yerel blok zinciri başlatın
* İki zincir arasında bir IBC aktarıcısı kurun
* İki zincir üzerinde bir token çifti için bir borsa emir defteri oluşturun
* Mars zincirinde satış emirleri gönderin
* Venüs zincirinde satın alma emirleri gönderin
* Satış veya alış emirlerini iptal etme

İki yerel blok zincirini başlatmak ve IBC aktarıcısını kurmak, iki zincir arasında bir değişim emri defteri oluşturmamızı sağlayacaktır. Bu emir defteri, satış ve alış emirleri göndermemize ve artık sürdürmek istemediğimiz emirleri iptal etmemize olanak tanıyacaktır.

Bu bölümdeki komutların yalnızca bu eğitimdeki sonraki tüm bölümleri tamamladıysanız düzgün çalışacağını unutmamak önemlidir. Bu bölümün sonunda, borsanın nasıl çalışacağı konusunda iyi bir anlayışa sahip olmalısınız.

### Blok zinciri node'larını başlatma

Zincirler arası değişimi kullanmaya başlamak için iki ayrı blok zinciri başlatmanız gerekecektir. Bu, `ignite chain serve` komutu ve ardından `-c` bayrağı ve her bir blok zinciri için yapılandırma dosyasının yolu çalıştırılarak yapılabilir. Örneğin, `mars` blok zincirini başlatmak için şu komutu çalıştırırsınız:

```
ignite chain serve -c mars.yml
```

`Venus` blockchain'ini başlatmak için benzer bir komut çalıştırırsınız, ancak `venus.yml` yapılandırma dosyasının yolunu girersiniz:

```
ignite chain serve -c venus.yml
```

Her iki blockchain de çalıştıktan sonra, iki chain arasında zincirler arası alışverişi etkinleştirmek için relayer'ı yapılandırmaya devam edebilirsiniz.

### Aktarıcı

Ardından, iki chain arasında bir IBC relayer kuralım. Eğer geçmişte bir relayer kullandıysanız, relayer yapılandırma dizinini sıfırlayın:

```
rm -rf ~/.ignite/relayer
```

Şimdi `ignite relayer configure` komutunu kullanabilirsiniz. Bu komut, kaynak ve hedef chain'lerin yanı sıra ilgili RPC uç noktalarını, musluk URL'lerini, bağlantı noktası numaralarını, sürümleri, gas fiyatlarını ve gas limitlerini belirtmenize olanak tanır.

```
ignite relayer configure -a --source-rpc "http://0.0.0.0:26657" --source-faucet "http://0.0.0.0:4500" --source-port "dex" --source-version "dex-1" --source-gasprice "0.0000025stake" --source-prefix "cosmos" --source-gaslimit 300000 --target-rpc "http://0.0.0.0:26659" --target-faucet "http://0.0.0.0:4501" --target-port "dex" --target-version "dex-1" --target-gasprice "0.0000025stake" --target-prefix "cosmos" --target-gaslimit 300000
```

İki chain arasında bir bağlantı oluşturmak için ignite relayer connect komutunu kullanabilirsiniz. Bu komut, kaynak ve hedef zincirler arasında bir bağlantı kurarak aralarında veri ve varlık aktarımı yapmanızı sağlar.

```
ignite relayer connect
```

Artık iki ayrı blockchain ağımız ve bunlar arasındaki iletişimi kolaylaştırmak için kurulmuş bir aktarıcı bağlantımız olduğuna göre, bu ağlarla etkileşim kurmak için zincirler arası değişim ikilisini kullanmaya başlamaya hazırız. Bu, emir defterleri ve alım/satım emirleri oluşturmamızı ve iki chain arasında varlık ticareti yapmamızı sağlayacaktır.

### Sipariş Defteri

Bir çift token için sipariş defteri oluşturmak için aşağıdaki komutu kullanabilirsiniz:

```
interchanged tx dex send-create-pair dex channel-0 marscoin venuscoin --from alice --chain-id mars --home ~/.mars
```

Bu komut, `marscoin` ve `venuscoin` token çifti için bir sipariş defteri oluşturacaktır. Komut, Mars blok zincirinde `alice` kullanıcısı tarafından yürütülecektir. `--home` parametresi Mars blok zinciri için yapılandırma dizininin konumunu belirtir.

Bir emir defteri oluşturmak, işlemin yayınlandığı Mars blok zincirindeki ve Venüs blok zincirindeki durumu etkiler.

Mars blockchain'inde `send-create-pair` komutu boş bir satış emri defteri oluşturur.

```
interchanged q dex list-sell-order-book
```

```
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 0
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Venus blok zincirinde, aynı `send-createPair` komutu bir alış emri defteri oluşturur:

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 0
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Mars blockchain'indeki `create-pair` komutunda, Venüs chain'ine bir IBC paketi gönderilir. Bu paket, Venüs chain üzerinde bir alım emri defteri oluşturmak için kullanılan bilgileri içerir.

Venüs chain'i IBC paketini aldığında, pakette yer alan bilgileri işler ve bir alım emri defteri oluşturur. Venüs chain'i daha sonra alım emri defterinin başarıyla oluşturulduğunu teyit etmek için Mars zincirine bir onay gönderir.

Venüs chain'inden onayı alan Mars chain'i bir satış emri defteri oluşturur. Bu satış emri defteri Venüs zincirindeki alış emri defteri ile ilişkilendirilir ve kullanıcıların iki chain arasında varlık ticareti yapmasına olanak tanır.

### Satış Emri

Bir emir defteri oluşturduktan sonra, bir sonraki adım bir satış emri oluşturmaktır. Bu, belirli miktarda tokenı kilitleyen ve Mars blockchain'inde bir satış emri oluşturan bir mesajla bir işlem yayınlamak için kullanılan `send-sell-order` komutu kullanılarak yapılabilir.

```
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 15  --from alice --chain-id mars --home ~/.mars
```

Verilen örnekte, `send-sell-order` komutu 10 `marscoin` token ve 15 `venuscoin` token için bir satış emri oluşturmak için kullanılır. Bu satış emri Mars blockchaini'ndeki emir defterine eklenecektir.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```
balances:
- amount: "990"  # decreased from 1000
  denom: marscoin
- amount: "1000"
  denom: token
```

```
interchanged q dex list-sell-order-book
```

```
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders: # a new sell order is created
    - amount: 10
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

### Sipariş al

Bir satış emri oluşturduktan sonra, alım satım sürecindeki bir sonraki adım genellikle bir alış emri oluşturmaktır. Bu, belirli miktarda tokenı kilitlemek ve Venus blok zincirinde bir satın alma emri oluşturmak için kullanılan `send-buy-order` komutu kullanılarak yapılabilir

```
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --chain-id venus --home ~/.venus --node tcp://localhost:26659
```

Verilen örnekte, `send-buy-order` komutu 10 `marscoin` token ve 5 `venuscoin` token için bir satın alma emri oluşturmak için kullanılır. Bu satın alma emri Venus blok zincirindeki emir defterine eklenecektir.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```
balances:
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "950" # decreased from 1000
  denom: venuscoin
```

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders: # a new buy order is created
    - amount: 10
      creator: cosmos1mrrttwtdcp47pl4hq6sar3mwqpmtc7pcl9e6ss
      id: 0
      price: 5
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

### Satış Emri ile Değişim Gerçekleştirin

Şu anda `marscoin` için iki açık emriniz var:

* Mars zincirinde 15 `venuscoin` karşılığında 10 `marscoin` satmayı teklif ettiğiniz bir satış emri.
* Venüs zincirinde, 5 `venuscoin` karşılığında 5 `marscoin` satın almak istediğiniz bir alış emri.

Bir takas gerçekleştirmek için, aşağıdaki komutu kullanarak Mars zincirine bir satış emri gönderebilirsiniz:

```
interchanged tx dex send-sell-order dex channel-0 marscoin 5 venuscoin 3 --from alice --home ~/.mars
```

Bu satış emri, 3 `venuscoin` karşılığında 5 `marscoin` satmayı teklif ediyor, Venüs zincirinde mevcut satın alma emri tarafından doldurulacak. Bu, Venüs zincirindeki alış emrinin miktarının 5 `marscoin` azalmasıyla sonuçlanacaktır.

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders:
    - amount: 5 # decreased from 10
      creator: cosmos1mrrttwtdcp47pl4hq6sar3mwqpmtc7pcl9e6ss
      id: 0
      price: 5
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Doldurulan satış emrinin göndericisi, 25 `venuscoin` token karşılığında 5 `marscoin` takas etmiştir. Bu, takasın değerini belirlemek için satış emrinin miktarının (5 `marscoin`) alış emrinin fiyatı (5 `venuscoin`) ile çarpıldığı anlamına gelir. Bu durumda, takasın değeri 25 `venuscoin` kuponu olmuştur.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```
balances:
- amount: "25" # increased from 0
  denom: ibc/BB38C24E9877
- amount: "985" # decreased from 990
  denom: marscoin
- amount: "1000"
  denom: token
```

Karşı taraf veya `marscoin` satın alma emrini gönderen, takas sonucunda 5 `marscoin` alacaktır.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```
balances:
- amount: "5" # increased from 0
  denom: ibc/745B473BFE24 # marscoin voucher
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "950"
  denom: venuscoin
```

`Venuscoin` bakiyesi değişmeden kaldı çünkü uygun miktarda `venuscoin` (`50`) bir önceki adımda satın alma emri oluşturulduğunda zaten kilitlenmişti.

### Alış Emri ile Değişim Gerçekleştirin

Satın alma emri ile bir değişim gerçekleştirmek için, merkezi olmayan borsaya `15 venuscoin` karşılığında `5` `marscoin` satın almak için bir işlem gönderin. Bu, aşağıdaki komut çalıştırılarak yapılır:

```
interchanged tx dex send-buy-order dex channel-0 marscoin 5 venuscoin 15 --from alice --home ~/.venus --node tcp://localhost:26659
```

Bu satın alma emri Mars zincirinde hemen doldurulacak ve satış emrini oluşturan kişi ödeme olarak `75` `venuscoin` kuponu alacaktır.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```
balances:
- amount: "100" # increased from 25
  denom: ibc/BB38C24E9877 # venuscoin voucher
- amount: "985"
  denom: marscoin
- amount: "1000"
  denom: token
```

Satış emrinin miktarı, doldurulan alış emrinin miktarı kadar azaltılacaktır, yani bu durumda 5 `marscoin` azaltılacaktır.

```
interchanged q dex list-sell-order-book
```

```
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders:
    - amount: 5 # decreased from 10
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Satın alma emrini oluşturan kişi 75 venuscoin karşılığında 5 marscoin kuponu alır (5marscoin \* 15venuscoin):

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```
balances:
- amount: "10" # increased from 5
  denom: ibc/745B473BFE24 # marscoin vouchers
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "875" # decreased from 950
  denom: venuscoin
```

### Kısmen Doldurulmuş Satış Emri ile Tam Değişim

Kısmen doldurulmuş bir satış emriyle değişimi tamamlamak için, merkezi olmayan borsaya 3 `venuscoin` karşılığında 10 `marscoin` satmak üzere bir işlem gönderin. Bu, aşağıdaki komut çalıştırılarak yapılır:

```
interchanged tx dex send-sell-order dex channel-0 marscoin 10 venuscoin 3 --from alice --home ~/.mars
```

Bu senaryoda, satış tutarı 10 `marscoin'dir`, ancak yalnızca 5 `marscoin` için mevcut bir satın alma emri vardır. Alış emri tamamen doldurulacak ve emir defterinden kaldırılacaktır. Önceden oluşturulmuş satın alma emrinin yazarı borsadan 10 `marscoin` kuponu alacaktır.

Bakiyeleri kontrol etmek için aşağıdaki komutu çalıştırabilir:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```
balances:
- amount: "15" # increased from 5
  denom: ibc/745B473BFE24 # marscoin voucher
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "875"
  denom: venuscoin
```

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 1
    orders: [] # buy order with amount 5marscoin has been closed
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

Satış emrinin yazarı başarılı bir şekilde 5 marscoin bozdurdu ve 25 venuscoin kuponu aldı. Diğer 5marscoin bir satış emri oluşturdu:

```
balances:
- amount: "125" # increased from 100
  denom: ibc/BB38C24E9877 # venuscoin vouchers
- amount: "975" # decreased from 985
  denom: marscoin
- amount: "1000"
  denom: token
```

```
interchanged q dex list-sell-order-book
```

```
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders:
    - amount: 5 # hasn't changed
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
    - amount: 5 # new order is created
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 1
      price: 3
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

### Kısmen Doldurulmuş Alış Emri ile Tam Değişim

Borsayı kısmen doldurulmuş bir satın alma emriyle tamamlamak için, 5 `venuscoin` karşılığında 10 `marscoin` satın almak üzere merkezi olmayan borsaya bir işlem gönderin. Bu, aşağıdaki komut çalıştırılarak yapılır:

```
interchanged tx dex send-buy-order dex channel-0 marscoin 10 venuscoin 5 --from alice --home ~/.venus --node tcp://localhost:26659
```

Bu senaryoda, alış emri 5 `marscoin` için yalnızca kısmen karşılanır. Mars zincirinde 5 `marscoin` için (fiyatı 3 venuscoin olan) mevcut bir satış emri vardır ve bu emir tamamen karşılanır ve emir defterinden kaldırılır. Kapatılan satış emrinin sahibi, ödeme olarak 5 `marscoin` ve 3 `venuscoin'in` çarpımı olan 15 venuscoin kuponu alacaktır.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars)
```

```
balances:
- amount: "140" # increased from 125
  denom: ibc/BB38C24E9877 # venuscoin vouchers
- amount: "975"
  denom: marscoin
- amount: "1000"
  denom: token
```

```
interchanged q dex list-sell-order-book
```

```
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders:
    - amount: 5 # order hasn't changed
      creator: cosmos14ntyzr6d2dx4ppds9tvenx53fn0xl5jcakrtm4
      id: 0
      price: 15
    # a sell order for 5 marscoin has been closed
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Bu senaryoda, satın alma emrinin sahibi ödeme olarak 5 `marscoin` kuponu alacak ve bu da tokenlerinin 50 `venuscoin`'ini kilitleyecektir. Satış emri tarafından doldurulmayan kalan 5 `marscoin`, Venüs zincirinde yeni bir satın alma emri oluşturacaktır. Bu, satın alma emrinin yazarının hala 5 `marscoin` satın almakla ilgilendiği ve bunun için belirtilen fiyatı ödemeye istekli olduğu anlamına gelir. Yeni satın alma emri, başka bir satış emri tarafından doldurulana veya alıcı tarafından iptal edilene kadar emir defterinde kalacaktır.

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

```
balances:
- amount: "20" # increased from 15
  denom: ibc/745B473BFE24 # marscoin vouchers
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "825" # decreased from 875
  denom: venuscoin
```

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders:
    - amount: 5 # new buy order is created
      creator: cosmos1mrrttwtdcp47pl4hq6sar3mwqpmtc7pcl9e6ss
      id: 1
      price: 5
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

### Sipariş İptal Etme

Açıklanan değişimlerden sonra, hala iki açık emir vardır: Mars zincirinde bir satış emri (15 `venuscoin` karşılığında 5 `marscoin`) ve Venüs zincirinde bir alış emri (5 `venuscoin` karşılığında 5 `marscoin`).

Bir blok zincirindeki bir emri iptal etmek için, iptal etmek istediğiniz emrin türüne bağlı olarak `cancel-sell-order` veya `cancel-buy-order` komutunu kullanabilirsiniz. Komut, IBC bağlantısının `channel-id`'si, emrin miktar-denom'u ve `amount-denom`'u ve iptal etmek istediğiniz emrin `order-id`'si dahil olmak üzere birkaç argüman alır.

Mars zincirindeki bir satış emrini iptal etmek için aşağıdaki komutu çalıştırırsınız:

```
interchanged tx dex cancel-sell-order dex channel-0 marscoin venuscoin 0 --from alice --home ~/.mars
```

Bu, satış emrini iptal edecek ve emir defterinden kaldıracaktır. Alice'in `marscoin` bakiyesi, iptal edilen satış emrinin tutarı kadar artacaktır.

Güncellenmiş `marscoin` bakiyesi de dahil olmak üzere Alice'in bakiyelerini kontrol etmek için aşağıdaki komutu çalıştırın:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.mars) 
```

Bu, güncellenmiş `marscoin` bakiyesi de dahil olmak üzere Alice'in bakiyelerinin bir listesini döndürecektir.

```
balances:
- amount: "140"
  denom: ibc/BB38C24E9877
- amount: "980" # increased from 975
  denom: marscoin
- amount: "1000"
  denom: token
```

Mars zincirindeki satış emri iptal edildikten sonra, bu blockchain üzerindeki satış emri defteri boş olacaktır. Bu, Mars zincirinde artık aktif bir satış emri olmadığı ve `marscoin` satın almak isteyen herkesin yeni bir satın alma emri oluşturması gerektiği anlamına gelir. Satış emri defteri, yeni bir satış emri oluşturulup eklenene kadar boş kalacaktır.

```
interchanged q dex list-sell-order-book
```

```
sellOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Venüs zincirindeki bir satın alma emrini iptal etmek için aşağıdaki komutu çalıştırabilirsiniz:

```
interchanged tx dex cancel-buy-order dex channel-0 marscoin venuscoin 1 --from alice --home ~/.venus --node tcp://localhost:26659
```

Bu, satın alma emrini iptal edecek ve emir defterinden kaldıracaktır. Alice'in `venuscoin` bakiyesi, iptal edilen satın alma emrinin tutarı kadar artacaktır.

Güncellenmiş `venuscoin` bakiyesi de dahil olmak üzere Alice'in bakiyelerini kontrol etmek için aşağıdaki komutu çalıştırabilirsiniz:

```
interchanged q bank balances $(interchanged keys show -a alice --home ~/.venus) --node tcp://localhost:26659
```

`Venuscoin` miktarı artırıldı:

```
balances:
- amount: "20"
  denom: ibc/745B473BFE24
- amount: "900000000"
  denom: stake
- amount: "1000"
  denom: token
- amount: "850" # increased from 825
  denom: venuscoin
```

Bu, güncellenmiş `venuscoin` bakiyesi de dahil olmak üzere Alice'in bakiyelerinin bir listesini döndürecektir.

Bir satın alma emrini iptal ettikten sonra, Venus blockchaini'ndeki satın alma emri defteri boş olacaktır. Bu, chain üzerinde artık aktif bir alım emri olmadığı ve `marscoin` satmak isteyen herkesin yeni bir satış emri oluşturması gerektiği anlamına gelir. Satın alma emri defteri, yeni bir satın alma emri oluşturulup eklenene kadar boş kalacaktır.

```
interchanged q dex list-buy-order-book --node tcp://localhost:26659
```

```
buyOrderBook:
- amountDenom: marscoin
  book:
    idCount: 2
    orders: []
  index: dex-channel-0-marscoin-venuscoin
  priceDenom: venuscoin
```

Bu kılavuzda, iki farklı blockchain ağı arasında token ticareti yapmak için zincirler arası bir borsanın nasıl kurulacağını gösterdik. Bu, belirli bir token çifti için bir borsa emir defteri oluşturmayı ve ikisi arasında sabit bir döviz kuru belirlemeyi içeriyordu.

Borsa kurulduktan sonra, kullanıcılar Mars zincirinde satış emirleri gönderebiliyor ve Venüs zincirinde alım emirleri verebiliyordu. Bu sayede tokenlarını satışa sunabiliyor ya da borsadan token satın alabiliyorlardı. Ayrıca kullanıcılar gerektiğinde emirlerini iptal de edebiliyordu.
