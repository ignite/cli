# Yapılandırma dosyası

### Yapılandırma dosyası referansı

Blockchain klasörünüzde oluşturulan `config.yml` dosyası, blockchaininiz için geliştirme ortamını tanımlamak için anahtar-değer çiftleri kullanır.

Yalnızca varsayılan bir parametre kümesi sağlanmıştır. Daha ayrıntılı yapılandırma gerekiyorsa, bu parametreleri `config.yml` dosyasına ekleyebilirsiniz.

Blockchain'in oluşumu sırasında oluşturulan kullanıcı hesaplarının bir listesi.

```
accounts:
  - name: alice
    coins: ['20000token', '200000000stake']
  - name: bob
    coins: ['10000token', '100000000stake']
```

Ignite, zinciri `ignite chain init` ve `ignite chain start` ile başlatırken `accounts`'dan gelen bilgileri kullanır. Yukarıdaki örnekte Ignite, zincirin `genesis.json` dosyasına iki hesap ekleyecektir.

`name`, bir hesapla ilişkili bir anahtar çiftinin yerel adıdır. Zincir başlatıldıktan ve başlatıldıktan sonra, işlemleri imzalarken `name`'i kullanabileceksiniz. Yukarıdaki yapılandırmayla, hem Alice'in hem de Bob'un hesaplarıyla işlemleri imzalayabilirsiniz, `exampled tx bank send ... --from alice`.

`coins`, hesap için token bakiyelerinin bir listesidir. Eğer bir token değeri bu listede yer alıyorsa, genesis bakiyesinde mevcut olacak ve geçerli bir token olacaktır. Yukarıdaki yapılandırma dosyası ile başlatıldığında, bir zincirin genesis'te yalnızca iki hesabı (Alice ve Bob) ve iki yerel jetonu (`token` ve `stake` değerleriyle) olacaktır.

Varsayılan olarak, bir zincir her yeniden başlatıldığında, Ignite her hesap için yeni bir anahtar çifti oluşturacaktır. Dolayısıyla, hesap adı aynı kalsa da (`bob`), her zincir yeniden başlatıldığında farklı bir anımsatıcıya ve adrese sahip olacaktır.

Bir hesabın belirli bir adrese sahip olmasını istiyorsanız, `address` alanına geçerli bir bech32 adresi girin. Önek (varsayılan olarak `cosmos`) zinciriniz tarafından beklenen önekle eşleşmelidir. Bir hesap bir `address`'le sağlandığında bir anahtar çifti oluşturulmayacaktır, çünkü bir adresten bir anahtar türetmek imkansızdır. Belirli bir adrese sahip bir hesap genesis dosyasına eklenecektir (ilişkili bir token bakiyesi ile), ancak anahtar çifti olmadığından, bu adresten işlem yayınlayamazsınız. Bu, Ignite dışında bir anahtar çifti oluşturduğunuzda (örneğin, zincirinizin CLI'sini kullanarak veya bir uzantı cüzdanında) ve bu anahtar çiftinin adresiyle ilişkili bir token bakiyesine sahip olmak istediğinizde kullanışlıdır.

```
accounts:
  - name: bob
    coins: ['20000token', '200000000stake']
    address: cosmos1s39200s6v4c96ml2xzuh389yxpd0guk2mzn3mz
```

Bir hesabın belirli bir mnemonic'ten başlatılmasını istiyorsanız, `mnemonic` alanına geçerli bir mnemonic girin. Bir özel anahtar, bir açık anahtar ve bir adres bir anımsatıcıdan türetilecektir.

```
accounts:
  - name: bob
    coins: ['20000token', '200000000stake']
    mnemonic: cargo ramp supreme review change various throw air figure humble soft steel slam pole betray inhale already dentist enough away office apple sample glue
```

Tek bir hesap için hem `address` hem de `mnemonic` tanımlayamazsınız.

Bazı hesaplar validatör hesaplar olarak kullanılır (`validators` bölümüne bakın). Validatör hesapların bir `adress` alanı olamaz, çünkü Ignite'ın özel bir anahtar türetebilmesi gerekir (rastgele bir mneomic'den (anımsatıcı) veya anımsatıcı alanında sağlanan belirli bir `mnemonic`'den). Validatör hesaplar, kendi kendine delegasyon için yeterli miktarda stake değerine sahip olmalıdır.

Varsayılan olarak, `alice` hesabı validatör hesabı olarak kullanılır, anahtarı genesis'te rastgele oluşturulan bir anımsatıcıdan türetilir, stake değeri `stake`'dir ve bu hesap kendi kendine delegasyon için yeterli `stake`'e sahiptir.

Zinciriniz kendi [cointype'ını ](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)kullanıyorsa, tamsayı değerini sağlamak için cointype alanını kullanabilirsiniz

```
accounts:
  - name: bob
    coins: ['20000token', '200000000stake']
    cointype: 7777777
```

`ignite chain init` ve `ignite chain serve` gibi komutlar, geliştirme amacıyla bir validatör node'u başlatır ve başlatır.

```
validators:
  - name: alice
    bonded: '100000000stake'
```

`name`, `accounts` listesindeki anahtar adını ifade eder.

`bonded`, bir validatörün kendi delegasyon tutarıdır. `bonded` tutar `1000000`'den düşük veya `account` listesindeki hesap bakiyesinden yüksek olmamalıdır.

Validatörler node yapılandırma dosyalarını veri dizininde saklar. Varsayılan olarak Ignite, veri dizininin adı olarak projenin adını kullanır, örneğin `$HOME/.example/`. Veri dizini için farklı bir yol kullanmak için home özelliğini özelleştirebilirsiniz.

Veri dizinindeki yapılandırma Ignite tarafından sık sık sıfırlanır. Yapılandırma dosyalarındaki bazı değişiklikleri kalıcı hale getirmek için `$HOME/.example/config/app.toml`, `$HOME/.example/config/config.toml` ve `$HOME/.example/config/client.toml` adreslerine karşılık gelen `app`, `config` ve `client` özelliklerini kullanabilirsiniz.

```
validators:
  - name: alice
    bonded: '100000000stake'
    home: "~/.mychain"
    app:
      pruning: "nothing"
    config:
      moniker: "mychain"
    client:
      output: "json"
```

`config.toml`, `app.toml` ve `client.toml` için hangi özelliklerin mevcut olduğunu görmek için `ignite chain init` ile bir zincir başlatın ve hakkında daha fazla bilgi edinmek istediğiniz dosyayı açın.

Şu anda Ignite yalnızca bir validatör node'ubaşlatır, bu nedenle `validators` listesindeki ilk öğe kullanılır (gerisi yok sayılır). Birden fazla validatör için destek devam etmektedir.

`build` özelliği, Ignite'ın zincirinizin ikili dosyasını nasıl oluşturacağını özelleştirmenizi sağlar.

Ignite varsayılan olarak `main` paketini `cmd/PROJECT_NAME/main.go` dizininden oluşturur. Projenizde birden fazla `main` paketi varsa veya dizini yeniden adlandırdıysanız, `main` özelliğini kullanarak `main` Go paketinin yolunu belirtin:

```
build:
  main: cmd/hello/cmd
```

Ignite projenizi bir binary olarak derler ve binary adı olarak projenin adını `d` son ekiyle birlikte kullanır. Binary adını özelleştirmek için `binary` özelliğini kullanın:

```
build:
  binary: "helloworldd"
```

Derleme işleminde kullanılan bağlayıcı bayraklarını özelleştirmek için:

```
build:
  ldflags: [ "-X main.Version=development", "-X main.Date=01/05/2022T19:54" ]
```

Varsayılan olarak, özel protokol tampon (proto) dosyaları `proto` dizininde bulunur. Projeniz proto dosyalarını farklı bir dizinde tutuyorsa, Ignite'a bunu bildirmelisiniz:

```
build:
  proto:
    path: "myproto"
```

Ignite, kutudan gerekli üçüncü taraf proto ile birlikte çıkar. Ignite ayrıca ekstra proto dosyaları için `third_party/proto` ve `proto_vendor` dizinlerine bakar. Projeniz üçüncü taraf proto dosyalarını farklı bir dizinde tutuyorsa, Ignite'a bunu bildirmelisiniz:

```
build:
  proto:
    third_party_paths: ["my_third_party/proto"]
```

Musluk hizmeti adreslere jeton gönderir.

```
faucet:
  name: bob
  coins: ["5token", "100000stake"]
```

`name`, `accounts` listesindeki bir anahtar adını ifade eder. Bu gerekli bir özelliktir.

`coins`, musluk tarafından bir kullanıcıya gönderilecek token miktarıdır. Bu gerekli bir özelliktir.

`coins_max`, tek bir adrese gönderilebilecek maksimum token miktarıdır. Jeton limitini sıfırlamak için `rate_limit_window` özelliğini kullanın (saniye cinsinden).

Varsayılan olarak musluk `4500` numaralı bağlantı noktasında çalışır. Farklı bir `port` numarası kullanmak için port özelliğini kullanın.

```
faucet:
  name: faucet
  coins: [ "100token", "5foo" ]
  coins_max: [ "2000token", "1000foo" ]
  port: 4500
  rate_limit_window: 3600
```

Genesis dosyası blockchain'deki ilk bloktur. Token bakiyeleri ve modüllerin durumu gibi önemli bilgiler içerdiğinden bir blockchain başlatmak için gereklidir. Genesis `$DATA_DIR/config/genesis.json` içinde saklanır.

Genesis dosyası geliştirme sırasında sık sık yeniden başlatıldığından, `genesis` özelliğinde kalıcı seçenekler ayarlayabilirsiniz:

```
genesis:
  app_state:
    staking:
      params:
        bond_denom: "denom"
```

Bir genesis dosyasının hangi özellikleri desteklediğini öğrenmek için, bir zinciri başlatın ve veri dizininde genesis dosyasına bakın.

Ignite, `ignite generate` komut kümesi ile zincirinizle etkileşim için istemci tarafı kodu oluşturabilir. İstemci tarafı kodunun oluşturulduğu yolları özelleştirmek için aşağıdaki özellikleri kullanın.

```
client:
  openapi:
    path: "docs/static/openapi.yml"
  typescript:
    path: "ts-client"
  composables:
    path: "vue/src/composables"
  hooks:
    path: "react/src/hooks"
```
