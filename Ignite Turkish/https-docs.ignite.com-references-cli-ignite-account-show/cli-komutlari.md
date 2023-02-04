# CLI Komutları

Ignite CLI için dokümantasyon.&#x20;

### ignite <a href="#ignite" id="ignite"></a>

Ignite CLI, blockchain'inizi iskeletlemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar

**Özet**

Ignite CLI, dünyanın en popüler modüler blockchain framework'ü olan Cosmos SDK ile inşa edilmiş egemen blockchain'ler oluşturmaya yönelik bir araçtır. Ignite CLI, blockchaininizi iskeletlemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar.

Başlamak için bir blockchain oluşturun:

```
ignite scaffold chain example
```

**Seçenekler**

```
  -h, --help   help for ignite
```

**Ayrıca Bakınız**

* [ignite account](broken-reference) - Ignite hesapları oluşturma, silme ve gösterme
* [ignite chain](broken-reference) - Bir blockchain node'u oluşturun, çalıştırın ve başlatın
* [ignite completion](broken-reference) - Belirtilen shell için otomatik tamamlama script'i oluşturur
* [ignite docs](broken-reference) - Ignite CLI dokümanlarını gösterir
* [ignite generate](broken-reference) - Kaynak koddan istemciler, API dokümanları oluşturun
* [ignite network](broken-reference) - Üretimde bir blockchain başlatın
* [ignite node](broken-reference) - Canlı bir blockchain düğümüne istekte bulunun
* [ignite plugin](broken-reference) - Eklentilerle başa çıkma
* [ignite relayer](broken-reference) - Blockchainleri bir IBC aktarıcı ile bağlayın
* [ignite scaffold](broken-reference) - Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun
* [ignite tools](broken-reference) - İleri düzey kullanıcılar için araçlar
* [ignite version](broken-reference) - Geçerli yapı bilgilerini yazdırma

### ignite account <a href="#ignite-account" id="ignite-account"></a>

Ignite hesapları oluşturma, silme ve gösterme

**Özet**

Ignite hesaplarını yönetmek için komutlar. Bir Ignite hesabı, bir anahtarlıkta saklanan özel/genel bir anahtar çiftidir. Şu anda Ignite hesapları, Ignite aktarıcı komutları ile etkileşim kurarken ve "ignite network" komutlarını kullanırken kullanılmaktadır.

Not: Ignite hesap komutları zincirinizin anahtarlarını ve hesaplarını yönetmek için değildir. Hesapları "config.yml "den yönetmek için zincirinizin ikili dosyasını kullanın. Örneğin, blockchaininizin adı "mychain" ise, zincirin anahtarlarını yönetmek için "mychaind keys" komutunu kullanın.

**Seçenekler**

```
  -h, --help                     help for account
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Ayrıca Bakınız**

* [ignite](broken-reference) - Ignite CLI, blockchaininizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite account create](broken-reference) - Bir hesap oluşturma
* [ignite account delete](broken-reference) - Bir hesabı isme göre silme
* [ignite account export](broken-reference) - Bir hesabı özel anahtar olarak dışa aktarma
* [ignite account import](broken-reference) - Bir anımsatıcı veya özel anahtar kullanarak bir hesabı içe aktarma
* [ignite account list](broken-reference) - Tüm hesapların bir listesini göster
* [ignite account show](broken-reference) - Show detailed information about a particular account

### ignite account create <a href="#ignite-account" id="ignite-account"></a>

Yeni bir hesap oluşturun

```
ignite account create [name] [flags]
```

**Seçenekler**

```
  -h, --help   help for create
```

**Üst komutlardan devralınan seçenekler**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Ayrıca Bakınız**

* [ignite account](broken-reference) - Ignite hesapları oluşturma, silme ve gösterme

### ignite account delete <a href="#ignite-account-delete" id="ignite-account-delete"></a>

Bir hesabı isme göre silme

```
ignite account delete [name] [flags]
```

**Seçenekler**

```
  -h, --help   help for delete
```

**Üst komutlardan devralınan seçenekler**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Ayrıca Bakınız**

* [ignite account](broken-reference) - Ignite hesapları oluşturma, silme ve gösterme

### ignite account export <a href="#ignite-account-export" id="ignite-account-export"></a>

Bir hesabı özel anahtar olarak dışa aktarma

```
ignite account export [name] [flags]
```

**Seçenekler**

```
  -h, --help                help for export
      --non-interactive     do not enter into interactive mode
      --passphrase string   passphrase to encrypt the exported key
      --path string         path to export private key. default: ./key_[name]
```

**Üst komutlardan devralınan seçenekler**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Ayrıca Bakınız**

* [ignite account](broken-reference) - Ignite hesapları oluşturma, silme ve gösterme

### ignite account import <a href="#ignite-account-import" id="ignite-account-import"></a>

Bir anımsatıcı veya özel anahtar kullanarak bir hesabı içe aktarma

```
ignite account import [name] [flags]
```

**Seçenekler**

```
  -h, --help                help for import
      --non-interactive     do not enter into interactive mode
      --passphrase string   passphrase to decrypt the imported key (ignored when secret is a mnemonic)
      --secret string       Your mnemonic or path to your private key (use interactive mode instead to securely pass your mnemonic)
```

**Üst komutlardan devralınan seçenekler**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Ayrıca Bakınız**

* [ignite account](broken-reference) - Ignite hesapları oluşturma, silme ve gösterme

### ignite account list <a href="#ignite-account-list" id="ignite-account-list"></a>

Tüm hesapların bir listesini göster

```
ignite account list [flags]
```

**Seçenekler**

```
      --address-prefix string   account address prefix (default "cosmos")
  -h, --help                    help for list
```

**Üst komutlardan devralınan seçenekler**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Ayrıca Bakınız**

* [ignite account](broken-reference) - Ignite hesapları oluşturma, silme ve gösterme

### ignite account show <a href="#ignite-account-show" id="ignite-account-show"></a>

Belirli bir hesap hakkında ayrıntılı bilgi gösterme

```
ignite account show [name] [flags]
```

**Seçenekler**

```
      --address-prefix string   account address prefix (default "cosmos")
  -h, --help                    help for show
```

**Üst komutlardan devralınan seçenekler**

```
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Ayrıca Bakınız**

* [ignite account](broken-reference) - Ignite hesapları oluşturma, silme ve gösterme

### ignite chain <a href="#ignite-chain" id="ignite-chain"></a>

Bir blockchain node'u oluşturun, çalıştırın ve başlatın

**Özet**

Bu ad alanındaki komutlar, geliştirme amacıyla blockchain node'unuzu yerel olarak oluşturmanıza, başlatmanıza ve başlatmanıza olanak tanır.

Bu komutları çalıştırmak için, Ignite'ın kaynak kodunu bulabilmesi için proje dizininin içinde olmalısınız. Bulunduğunuzdan emin olmak için "ls" komutunu çalıştırın, çıktıda aşağıdaki dosyaları görmelisiniz: "go.mod", "x", "proto", "app", vb.

Varsayılan olarak "build" komutu projenin "main" paketini belirleyecek, gerekirse bağımlılıkları yükleyecek, derleme bayraklarını ayarlayacak, projeyi bir binary'e derleyecek ve binary'i yükleyecektir. "build" komutu, örneğin zinciri manuel olarak başlatmak ve başlatmak için sadece derlenmiş binary'i istiyorsanız kullanışlıdır. Sürekli entegrasyon iş akışının bir parçası olarak zincirinizin binary dosyalarını otomatik olarak yayınlamak için de kullanılabilir.

"init" komutu zincirin binary'sini oluşturacak ve bunu yerel bir validatör node'unu başlatmak için kullanacaktır. Varsayılan olarak validatör node'u $HOME dizininizde projenizin adıyla eşleşen gizli bir dizinde başlatılacaktır. Bu dizine veri dizini denir ve bir zincirin oluşum dosyasını ve bir validatör anahtarını içerir. Bu komut, veri dizinini hızlı bir şekilde oluşturmak ve başlatmak ve blockchain'i manuel olarak başlatmak için zincirin binary'sini kullanmak istiyorsanız kullanışlıdır. "init" komutu yalnızca geliştirme amaçlıdır, üretim için değildir.

"serve" komutu, geliştirme amacıyla blockchaininizi tek bir validatör node ile yerel olarak oluşturur, başlatır ve başlatır. "serve" ayrıca dosya değişiklikleri için kaynak kod dizinini izler ve zinciri akıllı bir şekilde yeniden oluşturur/başlatır/başlatır, esasen "kod-yeniden yükleme" sağlar. "serve" komutu yalnızca geliştirme amaçlıdır, üretim için değildir.

Üretim ve geliştirme arasında ayrım yapmak için aşağıdakileri göz önünde bulundurun.

Üretimde, blockchainler genellikle aynı yazılımı farklı kişi ve kuruluşlar tarafından işletilen birçok validatör node üzerinde çalıştırır. Bir blockchain'i üretimde başlatmak için validatör kuruluşlar, düğümlerini eş zamanlı olarak başlatmak üzere başlatma sürecini koordine eder.

Geliştirme sırasında, bir blockchain tek bir validatör node üzerinde yerel olarak başlatılabilir. Bu kullanışlı süreç, bir zinciri hızlı bir şekilde yeniden başlatmanıza ve daha hızlı yinelemenize olanak tanır. Geliştirme sırasında tek bir node üzerinde bir zincir başlatmak, yerel bir sunucu üzerinde geleneksel bir web uygulaması başlatmaya benzer.

"Faucet" komutu, "config.yml" içinde tanımlanan "faucet" hesabından bir adrese token göndermenizi sağlar. Alternatif olarak, zincirde bulunan herhangi bir hesaptan token göndermek için zincirin ikilisini kullanabilirsiniz.

"simulate" komutu, zinciriniz için bir simülasyon test süreci başlatmanıza yardımcı olur.

**Seçenekler**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -h, --help            help for chain
  -y, --yes             answers interactive yes/no questions with yes
```

**Ayrıca Bakınız**

* [ignite](broken-reference) - Ignite CLI, blockchaininizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite chain build](broken-reference) - Bir node binary'si oluşturun
* [ignite chain debug](broken-reference) - Bir blockchain uygulaması için hata ayıklayıcı başlatma
* [ignite chain faucet](broken-reference) - Bir hesaba para gönderme
* [ignite chain init](broken-reference) - Zincirinizi başlatın
* [ignite chain serve](broken-reference) - Geliştirme aşamasında bir blockchain node'u başlatın
* [ignite chain simulate](broken-reference) - Blockchain için simülasyon testi çalıştırın

### ignite chain build <a href="#ignite-chain-build" id="ignite-chain-build"></a>

Bir node binary'si oluşturun

**Özet**

Build komutu projenin kaynak kodunu bir binary olarak derler ve binary'yi $(go env GOPATH)/bin dizinine yükler.

Bir bayrak kullanarak binary için çıktı dizinini özelleştirebilirsiniz:

```
ignite chain build --output dist
```

Binary'yi derlemek için Ignite önce protokol tampon (proto) dosyalarını Go kaynak koduna derler. Proto dosyaları gerekli tip ve servis tanımlarını içerir. Proto dosyalarını derlemek için başka bir program kullanıyorsanız, Ignite'a proto derleme adımını atlamasını söylemek için bir bayrak kullanabilirsiniz:

```
ignite chain build --skip-proto
```

Daha sonra, Ignite go.mod dosyasında belirtilen bağımlılıkları yükler. Varsayılan olarak Ignite, modül önbelleğinde depolanan ana modülün bağımlılıklarının indirildiklerinden beri değiştirilmediğini kontrol etmez. Bağımlılık kontrolünü zorlamak için (aslında, "go mod verify" çalıştırmak) bir bayrak kullanın:

```
ignite chain build --check-dependencies
```

Daha sonra, Ignite projenin "ana" paketini tanımlar. Varsayılan olarak "ana" paket "cmd/{app}d" dizininde bulunur; burada "{app}" iskele projesinin adıdır ve "d" daemon anlamına gelir. Projeniz birden fazla "ana" paket içeriyorsa, Ignite'ın config.yml dosyasında derlemesi gereken paketin yolunu belirtin:

```
build:
  main: custom/path/to/main
```

Varsayılan olarak binary adı, "d" son ekiyle birlikte üst düzey modül adıyla (go.mod'da belirtilen) eşleşecektir. Bu config.yml dosyasında özelleştirilebilir:

Ayrıca özel bağlayıcı bayrakları da belirtebilirsiniz:

```
build:
  ldflags:
    - "-X main.Version=development"
    - "-X main.Date=01/05/2022T19:54"
```

Bir sürüm için binary oluşturmak için --release bayrağını kullanın. Belirtilen bir veya daha fazla sürüm hedefi için binaryd, projenin kaynak dizinindeki bir "release/" dizininde oluşturulur. Sürüm hedeflerini GOOS:GOARCH derleme etiketleri ile belirtin. İsteğe bağlı --release.targets belirtilmezse, mevcut ortamınız için bir binary oluşturulur.

```
ignite chain build --release -t linux:amd64 -t darwin:amd64 -t darwin:arm64
```

```
ignite chain build [flags]
```

**Seçenekler**

```
      --check-dependencies        verify that cached dependencies have not been modified since they were downloaded
      --clear-cache               clear the build cache (advanced)
      --debug                     build a debug binary
  -h, --help                      help for build
  -o, --output string             binary output path
  -p, --path string               path of the app (default ".")
      --release                   build for a release
      --release.prefix string     tarball prefix for each release target. Available only with --release flag
  -t, --release.targets strings   release targets. Available only with --release flag
      --skip-proto                skip file generation from proto
  -v, --verbose                   verbose output
```

**Üst komutlardan devralınan seçenekler**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**Ayrıca Bakınız**

* [ignite chain](broken-reference) - Bir blockchain node'u oluşturun, çalıştırın ve başlatın

### ignite chain debug <a href="#ignite-chain-debug" id="ignite-chain-debug"></a>

Bir blockchain uygulaması için hata ayıklayıcı başlatma

**Özet**

Debug komutu bir hata ayıklama sunucusu başlatır ve bir hata ayıklayıcı başlatır.

Ignite varsayılan olarak Delve hata ayıklayıcısını kullanır. Delve, sürecin yürütülmesini kontrol ederek, değişkenleri değerlendirerek ve iş parçacığı / goroutine durumu, CPU kayıt durumu ve daha fazlası hakkında bilgi sağlayarak programınızla etkileşime girmenizi sağlar.

Varsayılan terminal istemcisinin istenmediği durumlarda isteğe bağlı olarak bir hata ayıklama sunucusu başlatılabilir. Sunucu başladığında önce blok zinciri uygulamasını çalıştırır, ona bağlanır ve son olarak bir istemci bağlantısı bekler. Hem JSON-RPC hem de DAP istemci bağlantılarını kabul eder.

Bir hata ayıklama sunucusu başlatmak için aşağıdaki bayrağı kullanın:

```
ignite chain debug --server
```

Bir hata ayıklama sunucusunu özel bir adresle başlatmak için aşağıdaki bayrakları kullanın:

```
ignite chain debug --server --server-address 127.0.0.1:30500
```

İstemci bağlantısı kapatıldığında hata ayıklama sunucusu otomatik olarak durur.

```
ignite chain debug [flags]
```

**Seçenekler**

```
  -h, --help                    help for debug
  -p, --path string             path of the app (default ".")
      --server                  start a debug server
      --server-address string   debug server address (default "127.0.0.1:30500")
```

**Üst komutlardan devralınan seçenekler**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**Ayrıca Bakınız**

* [ignite chain](broken-reference) - Bir blockchain node'u oluşturun, çalıştırın ve başlatın

### ignite chain faucet <a href="#ignite-chain-faucet" id="ignite-chain-faucet"></a>

Bir hesaba para gönderme

```
ignite chain faucet [address] [coin<,...>] [flags]
```

**Seçenekler**

```
  -h, --help          help for faucet
      --home string   directory where the blockchain node is initialized
  -p, --path string   path of the app (default ".")
  -v, --verbose       verbose output
```

**Üst komutlardan devralınan seçenekler**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**Ayrıca Bakınız**

* [ignite chain](broken-reference) - Bir blockchain node'u oluşturun, çalıştırın ve başlatın

### ignite chain init <a href="#ignite-chain-init" id="ignite-chain-init"></a>

Zincirinizi başlatın

**Özet**

init komutu binary'yi derler ve yükler ("ignite chain build" gibi) ve bu binary'yi blockchain'in veri dizinini bir validatör için başlatmak için kullanır. Derleme işleminin nasıl çalıştığını öğrenmek için "ignite chain build --help" bölümüne bakın.

Varsayılan olarak, veri dizini $HOME/.mychain olarak başlatılır; burada "mychain" projenin adıdır. Özel bir veri dizini ayarlamak için --home bayrağını kullanın veya config.yml dosyasında değeri ayarlayın:

```
validators:
  - name: alice
    bonded: '100000000stake'
    home: "~/.customdir"
```

Veri dizini "config" dizininde üç dosya içerir: app.toml, config.toml, client.toml. Bu dosyalar blockchain node'unuzun ve istemci çalıştırılabilir dosyasının davranışını özelleştirmenizi sağlar. Bir zincir yeniden başlatıldığında veri dizini sıfırlanabilir. Bu dosyalardaki bazı değerleri kalıcı hale getirmek için config.yml dosyasında ayarlayın:

```
validators:
  - name: alice
    bonded: '100000000stake'
    app:
      minimum-gas-prices: "0.025stake"
    config:
      consensus:
        timeout_commit: "5s"
        timeout_propose: "5s"
    client:
      output: "json"
```

Yukarıdaki yapılandırma doğrulayıcının minimum gas fiyatını değiştirir (varsayılan olarak gas fiyatı "ücretsiz" işlemlere izin vermek için 0 olarak ayarlanmıştır), blok süresini 5s olarak ayarlar ve çıktı biçimini JSON olarak değiştirir. Bu yapılandırmanın ne tür değerleri kabul ettiğini görmek için veri dizininde oluşturulan TOML dosyalarına bakın.

Ignite, başlatma sürecinin bir parçası olarak token bakiyeleri ile zincir üzerinde hesaplar oluşturur. Varsayılan olarak config.yml dosyasının üst düzey "accounts" özelliğinde iki hesap bulunur. Daha fazla hesap ekleyebilir ve token bakiyelerini değiştirebilirsiniz. Hangi değerleri ayarlayabileceğinizi görmek için config.yml kılavuzuna bakın.

Bu hesaplardan biri bir validatör hesabıdır ve kendi kendine devredilen token miktarı üst düzey "validator" özelliğinde ayarlanabilir.

Başlatılmış bir zincirin en önemli bileşenlerinden biri, zincirin 0. bloğu olan genesis dosyasıdır. Genesis dosyası veri dizini "config" alt dizininde saklanır ve konsensüs ve modül parametreleri de dahil olmak üzere zincirin başlangıç durumunu içerir. Genesis'in değerlerini config.yml dosyasında özelleştirebilirsiniz:

```
genesis:
  app_state:
    staking:
      params:
        bond_denom: "foo"
```

Yukarıdaki örnek staking token'ını "foo" olarak değiştirir. Staking denom'u değiştirirseniz, validator hesabının doğru token'lara sahip olduğundan emin olun.

init komutu SADECE GELİŞTİRME AMAÇLARI İÇİN kullanılmak üzere tasarlanmıştır. Kaputun altında "appd init", "appd add-genesis-account", "appd gentx" ve "appd collect-gentx" gibi komutları çalıştırır. Üretim için, üretim düzeyinde bir node başlatma sağlamak için bu komutları manuel olarak çalıştırmak isteyebilirsiniz.

```
ignite chain init [flags]
```

**Seçenekler**

```
      --check-dependencies   verify that cached dependencies have not been modified since they were downloaded
      --clear-cache          clear the build cache (advanced)
      --debug                build a debug binary
  -h, --help                 help for init
      --home string          directory where the blockchain node is initialized
  -p, --path string          path of the app (default ".")
      --skip-proto           skip file generation from proto
```

**Üst komutlardan devralınan seçenekler**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**Ayrıca Bakınız**

* [ignite chain](broken-reference) - Bir blockchain node'u oluşturun, çalıştırın ve başlatın

### ignite chain serve <a href="#ignite-chain-serve" id="ignite-chain-serve"></a>

Geliştirme aşamasında bir blockchain node'u başlatın

**Özet**

Serv komutu binary'yi derler ve yükler ("ignite chain build" gibi), blockchain'in veri dizinini bir validatör için başlatmak için bu binary'yi kullanır ("ignite chain init" gibi) ve otomatik kod yeniden yükleme ile geliştirme amacıyla node'u yerel olarak başlatır.

Otomatik kod yeniden yükleme, Ignite'ın proje dizinini izlemeye başladığı anlamına gelir. Bir dosya değişikliği tespit edildiğinde, Ignite otomatik olarak node'u yeniden oluşturur, yeniden başlatır ve yeniden başlatır.

Mümkün olduğunda Ignite, genesis dosyasını dışa ve içe aktararak zincirin mevcut durumunu korumaya çalışacaktır.

Bir genesis dosyası mevcut olsa bile Ignite'ı temiz bir başlangıçtan başlamaya zorlamak için aşağıdaki bayrağı kullanın:

```
ignite chain serve --reset-once
```

Ignite'ı kaynak kod her değiştirildiğinde durumu sıfırlamaya zorlamak için aşağıdaki bayrağı kullanın:

```
ignite chain serve --force-reset
```

Ignite ile farklı yapılandırma dosyaları kullanarak aynı kaynak kodundan birden fazla blockchain başlatmak mümkündür. Bu, blok zincirleri arası işlevsellik oluşturuyorsanız ve örneğin bir blockchainden diğerine paket göndermeyi denemek istiyorsanız kullanışlıdır. Belirli bir yapılandırma dosyası kullanarak bir node başlatmak için:

```
ignite chain serve --config mars.yml
```

serve komutu SADECE GELİŞTİRME AMAÇLI kullanılmak üzere tasarlanmıştır. Kaputun altında, "appd start" çalıştırır, burada "appd" zincirinizin ikili dosyasının adıdır. Üretim için "appd start" komutunu manuel olarak çalıştırmak isteyebilirsiniz.

```
ignite chain serve [flags]
```

**Seçenekler**

```
      --check-dependencies   verify that cached dependencies have not been modified since they were downloaded
      --clear-cache          clear the build cache (advanced)
  -f, --force-reset          force reset of the app state on start and every source change
      --generate-clients     generate code for the configured clients on reset or source code change
  -h, --help                 help for serve
      --home string          directory where the blockchain node is initialized
  -p, --path string          path of the app (default ".")
      --quit-on-fail         quit program if the app fails to start
  -r, --reset-once           reset the app state once on init
      --skip-proto           skip file generation from proto
  -v, --verbose              verbose output
```

**Üst komutlardan devralınan seçenekler**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**Ayrıca Bakınız**

* [ignite chain](broken-reference) - Bir blockchain node'u oluşturun, çalıştırın ve başlatın

### ignite chain simulate <a href="#ignite-chain-simulate" id="ignite-chain-simulate"></a>

Blockchain için simülasyon testi çalıştırın

**Synopsis**

Blockchain için simülasyon testi çalıştırın. Simüle edilmiş bir node'a her modülden çok sayıda rastgele giriş mesajı gönderir ve değişmezlerin bozulup bozulmadığını kontrol eder

```
ignite chain simulate [flags]
```

**Seçenekler**

```
      --blockSize int             operations per block (default 30)
      --exportParamsHeight int    height to which export the randomly generated params
      --exportParamsPath string   custom file path to save the exported params JSON
      --exportStatePath string    custom file path to save the exported app state JSON
      --exportStatsPath string    custom file path to save the exported simulation statistics JSON
      --genesis string            custom simulation genesis file; cannot be used with params file
      --genesisTime int           override genesis UNIX time instead of using a random UNIX time
  -h, --help                      help for simulate
      --initialBlockHeight int    initial block to start the simulation (default 1)
      --lean                      lean simulation log output
      --numBlocks int             number of new blocks to simulate from the initial block height (default 200)
      --params string             custom simulation params file which overrides any random params; cannot be used with genesis
      --period uint               run slow invariants only once every period assertions
      --printAllInvariants        print all invariants if a broken invariant is found
      --seed int                  simulation random seed (default 42)
      --simulateEveryOperation    run slow invariants every operation
  -v, --verbose                   verbose log output
```

**Üst komutlardan devralınan seçenekler**

```
  -c, --config string   path to Ignite config file (default: ./config.yml)
  -y, --yes             answers interactive yes/no questions with yes
```

**Ayrıca Bakınız**

* [ignite chain](broken-reference) - Bir blockchain node'u oluşturun, çalıştırın ve başlatın

### ignite completion <a href="#ignite-completion" id="ignite-completion"></a>

Belirtilen shell için otomatik tamamlama script'i oluşturur

**Özet**

Belirtilen shell için ignite otomatik tamamlama script'i oluşturur. Oluşturulan script'in nasıl kullanılacağına ilişkin ayrıntılar için her bir alt komutun yardımına bakın.

**Seçenekler**

```
  -h, --help   help for completion
```

**Ayrıca Bakınız**

* [ignite](broken-reference) -Ignite CLI, blockchaininizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite completion bash](broken-reference) - Bash için otomatik tamamlama script'i oluşturun
* [ignite completion fish](broken-reference) - Fish için otomatik tamamlama script'i oluşturun
* [ignite completion powershell](broken-reference) - Powershell için otomatik tamamlama script'i oluşturun
* [ignite completion zsh](broken-reference) - Zsh için otomatik tamamlama script'i oluşturun

### ignite completion bash <a href="#ignite-completion-bash" id="ignite-completion-bash"></a>

Bash için otomatik tamamlama script'i oluşturun

**Özet**

Bash shell için otomatik tamamlama script'i oluşturun.

Bu script 'bash-completion' paketine bağlıdır. Zaten yüklü değilse, işletim sisteminizin paket yöneticisi aracılığıyla yükleyebilirsiniz.

Mevcut shell oturumunuzdaki tamamlamaları yüklemek için:

```
source <(ignite completion bash)
```

Her yeni oturumun tamamlamalarını yüklemek için bir kez çalıştırın:

**#### Linux:**

```
ignite completion bash > /etc/bash_completion.d/ignite
```

**#### macOS:**

```
ignite completion bash > $(brew --prefix)/etc/bash_completion.d/ignite
```

Bu kurulumun etkili olması için yeni bir shell başlatmanız gerekecektir.

**Seçenekler**

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

**Ayrıca Bakınız**

* [ignite completion](broken-reference) - Belirtilen shell için otomatik tamamlama script'i oluşturur

### ignite completion fish <a href="#ignite-completion-fish" id="ignite-completion-fish"></a>

Fish için otomatik tamamlama script'i oluşturun

**Özet**

Fish shell'i için otomatik tamamlama script'i oluşturun

Geçerli shell oturumunuzdaki tamamlamaları yüklemek için:

```
ignite completion fish | source
```

Her yeni oturumun tamamlamalarını yüklemek için bir kez çalıştırın:

```
ignite completion fish > ~/.config/fish/completions/ignite.fish
```

Bu kurulumun etkili olması için yeni bir shell başlatmanız gerekecektir.

```
ignite completion fish [flags]
```

**Seçenekler**

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

**Ayrıca Bakınız**

* [ignite completion](broken-reference) - Belirtilen shell için otomatik tamamlama script'i oluşturur

### ignite completion powershell <a href="#ignite-completion-powershell" id="ignite-completion-powershell"></a>

Powershell için otomatik tamamlama script'i oluşturun

**Özet**

Powershell için otomatik tamamlama script'i oluşturun

Geçerli shell oturumunuzdaki tamamlamaları yüklemek için:

```
ignite completion powershell | Out-String | Invoke-Expression
```

Her yeni oturum için tamamlananları yüklemek için yukarıdaki komutun çıktısını powershell profilinize ekleyin.

```
ignite completion powershell [flags]
```

**Seçenekler**

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

**Ayrıca Bakınız**

* [ignite completion](broken-reference) - Belirtilen shell için otomatik tamamlama script'i oluşturur

### ignite completion zsh <a href="#ignite-completion-zsh" id="ignite-completion-zsh"></a>

Zsh için otomatik tamamlama script'i oluşturun

**Özet**

Zsh shell'i için otomatik tamamlama script'i oluşturun

Ortamınızda shell completion zaten etkin değilse, etkinleştirmeniz gerekecektir. Aşağıdakileri bir kez çalıştırabilirsiniz:

```
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

Geçerli shell oturumunuzdaki tamamlamaları yüklemek için:

```
source <(ignite completion zsh); compdef _ignite ignite
```

Her yeni oturumun tamamlamalarını yüklemek için bir kez çalıştırın:

**#### Linux:**

```
ignite completion zsh > "${fpath[1]}/_ignite"
```

**#### macOS:**

```
ignite completion zsh > $(brew --prefix)/share/zsh/site-functions/_ignite
```

Bu kurulumun etkili olması için yeni bir shell başlatmanız gerekecektir.

```
ignite completion zsh [flags]
```

**Seçenekler**

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

**Ayrıca Bakınız**

* [ignite completion](broken-reference) - Belirtilen shell için otomatik tamamlama script'i oluşturur

### ignite docs <a href="#ignite-docs" id="ignite-docs"></a>

Ignite CLI dokümanlarını göster

```
ignite docs [flags]
```

**Seçenekler**

```
  -h, --help   help for docs

```

**Ayrıca Bakınız**

* [ignite](broken-reference) - Ignite CLI, blockchaininizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar

### ignite generate <a href="#ignite-generate" id="ignite-generate"></a>

Kaynak koddan istemciler, API dokümanları oluşturun

**Özet**

Kaynak koddan istemciler, API dokümanları oluşturun.

Protokol tampon dosyalarını Go'ya derlemek veya belirli işlevleri uygulamak, örneğin bir OpenAPI spesifikasyonu oluşturmak gibi.

Üretilen kaynak kod, bir komut tekrar çalıştırılarak yeniden oluşturulabilir ve elle düzenlenmesi amaçlanmamıştır.

**Seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -h, --help          help for generate
  -p, --path string   path of the app (default ".")
```

**Ayrıca Bakınız**

* [ignite](broken-reference) - Ignite CLI, blockchaininizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite generate composables](broken-reference) - TypeScript ön uç istemcisi ve Vue 3 bileşikleri
* [ignite generate hooks](broken-reference) - TypeScript ön uç istemcisi ve React kancaları
* [ignite generate openapi](broken-reference) - Zinciriniz için OpenAPI özellikleri
* [ignite generate proto-go](broken-reference) - Protokol tampon dosyalarını Cosmos SDK için gerekli Go kaynak koduna derleyin
* [ignite generate ts-client](broken-reference) - TypeScript frontend istemcisi
* [ignite generate vuex](broken-reference) - DEPRECATED TypeScript ön uç istemcisi ve Vuex mağazaları

### ignite generate composables <a href="#ignite-generate-composables" id="ignite-generate-composables"></a>

TypeScript ön uç istemcisi ve Vue 3 bileşikleri

```
ignite generate composables [flags]
```

**Seçenekler**

```
  -h, --help            help for composables
  -o, --output string   Vue 3 composables output path
  -y, --yes             answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**Ayrıca Bakınız**

* [ignite generate](broken-reference) - Kaynak koddan istemciler, API dokümanları oluşturun

### ignite generate hooks[​](https://docs.ignite.com/references/cli#ignite-generate-hooks) <a href="#ignite-generate-hooks" id="ignite-generate-hooks"></a>

TypeScript ön uç istemcisi ve React kancaları

```
ignite generate hooks [flags]
```

**Seçenekler**

```
  -h, --help            help for hooks
  -o, --output string   React hooks output path
  -y, --yes             answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**Ayrıca Bakınız**

* [ignite generate](broken-reference) - Kaynak koddan istemciler, API dokümanları oluşturun

### ignite generate openapi <a href="#ignite-generate-openapi" id="ignite-generate-openapi"></a>

Zinciriniz için OpenAPI özellikleri

```
ignite generate openapi [flags]
```

**Seçenekler**

```
  -h, --help   help for openapi
  -y, --yes    answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**Ayrıca Bakınız**

* [ignite generate](broken-reference) - Kaynak koddan istemciler, API dokümanları oluşturun

### ignite generate proto-go <a href="#ignite-generate-proto-go" id="ignite-generate-proto-go"></a>

Protokol tampon dosyalarını Cosmos SDK için gerekli Go kaynak koduna derleyin

```
ignite generate proto-go [flags]
```

**Seçenekler**

```
  -h, --help   help for proto-go
  -y, --yes    answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**Ayrıca Bakınız**

* [ignite generate](broken-reference) - Kaynak koddan istemciler, API dokümanları oluşturun

### ignite generate ts-client <a href="#ignite-generate-ts-client" id="ignite-generate-ts-client"></a>

TypeScript frontend istemcisi

**Özet**

Blockchain projeniz için framework agnostik bir TypeScript istemcisi oluşturun.

Varsayılan olarak TypeScript istemcisi "ts-client/" dizininde oluşturulur. Çıktı dizinini config.yml dosyasında özelleştirebilirsiniz:

```
client:
  typescript:
    path: new-path
```

Çıktı, bir bayrak kullanılarak da özelleştirilebilir:

```
ignite generate ts-client --output new-path
```

TypeScript istemci kodu, blockchain bir bayrakla başlatıldığında sıfırlama veya kaynak kodu değişikliklerinde otomatik olarak yeniden oluşturulabilir:

```
ignite chain serve --generate-clients
```

```
ignite generate ts-client [flags]
```

**Seçenekler**

```
  -h, --help            help for ts-client
  -o, --output string   TypeScript client output path
      --use-cache       use build cache to speed-up generation
  -y, --yes             answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**Ayrıca Bakınız**

* [ignite generate](broken-reference) - Kaynak koddan istemciler, API dokümanları oluşturun

### ignite generate vuex <a href="#ignite-generate-vuex" id="ignite-generate-vuex"></a>

DEPRECATED TypeScript ön uç istemcisi ve Vuex mağazaları

```
ignite generate vuex [flags]
```

**Seçeneakler**

```
  -h, --help            help for vuex
  -o, --output string   Vuex store output path
  -y, --yes             answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -p, --path string   path of the app (default ".")
```

**Ayrıca Bakınız**

* [ignite generate](broken-reference) - Kaynak koddan istemciler, API dokümanları oluşturun

### ignite network <a href="#ignite-network" id="ignite-network"></a>

Üretimde bir blockchain başlatın

**Özet**

Ignite Network komutları, egemen Cosmos blockchain'lerinin başlatılmasını koordine etmeyi sağlar.

Bir Cosmos blockchain'i başlatmak için birinin koordinatör ve diğerlerinin de validatör olması gerekir. Bunlar sadece rollerdir, herkes koordinatör veya validatör olabilir. Bir koordinatör, Ignite blockchain'inde başlatılacak bir zincir hakkında bilgi yayınlar, validatör taleplerini onaylar ve başlatmayı koordine eder. Validatörler bir zincire katılmak için talepler gönderir ve bir blockchain başlatılmaya hazır olduğunda düğümlerini başlatır.

Zinciriniz hakkındaki bilgileri bir koordinatör olarak yayınlamak için aşağıdaki komutu çalıştırın (URL, Cosmos SDK zinciri içeren bir depoyu işaret etmelidir):

```
ignite network chain publish github.com/ignite/example
```

Bu komut, aşağıdaki komutlarda kullanacağınız bir başlatma tanımlayıcısı döndürecektir. Bu tanımlayıcının 42 olduğunu varsayalım.

Ardından, validatörlerden node'larını başlatmalarını ve validatör olarak ağa katılma talebinde bulunmalarını isteyin. Bir test ağı için CLI tarafından önerilen varsayılan değerleri kullanabilirsiniz.

```
ignite network chain init 42

ignite network chain join 42 --amount 95000000stake
```

Koordinatör olarak tüm validatör taleplerini listeleyin:

```
ignite network request list 42
```

Validatör taleplerini onaylayın:

```
ignite network request approve 42 1,2
```

Validatör setinde ihtiyacınız olan tüm validatörleri onayladıktan sonra zincirin başlatılmaya hazır olduğunu duyurun:

```
ignite network chain launch 42
```

Validatörler artık node'larını launch için hazırlayabilirler:

```
ignite network chain prepare 42
```

Bu komutun çıktısı, bir validatörün node'unu başlatmak için kullanacağı bir komut gösterecektir, örneğin "exampled --home \~/.example". Yeterli sayıda validatör node'larını başlattıktan sonra, bir blockchain yayınlanmış olacaktır.

**Seçenekler**

```
  -h, --help                        help for network
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite](broken-reference) - Ignite CLI, blockchaininizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın
* [ignite network coordinator](broken-reference) - Koordinatör profilini gösterme ve güncelleme
* [ignite network profile](broken-reference) - Adres profil bilgilerini göster
* [ignite network project](broken-reference) - Projeleri ele alın
* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama
* [ignite network reward](broken-reference) - Ağ ödüllerini yönetin
* [ignite network tool](broken-reference) - Yardımcı araçları çalıştırmak için komutlar
* [ignite network validator](broken-reference) - Validatör profilini gösterme ve güncelleme
* [ignite network version](broken-reference) - Eklentinin sürümü

### ignite network chain <a href="#ignite-network-chain" id="ignite-network-chain"></a>

Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

**Özet**

"chain" ad alanı, Ignite ile blockchainleri başlatmak için en sık kullanılan komutları içerir.

Bir koordinatör olarak blockchaininizi Ignite'ta "yayınlarsınız". Genesis için yeterli sayıda validatör onaylandığında ve genesis üzerinde herhangi bir değişiklik yapılması gerekmediğinde, bir koordinatör "launch" komutuyla zincirin başlatılmaya hazır olduğunu duyurur. Başarısız bir başlatma durumunda, koordinatör "revert-launch" komutunu kullanarak bunu geri alabilir.

Bir validatör olarak, node'unuzu "init" eder ve "join" komutu ile bir blockchain için validatör olmak üzere başvuruda bulunursunuz. Zincirin başlatıldığı duyurulduktan sonra, validatörler nihai genesis'i oluşturabilir ve "prepare" komutu ile eşlerin listesini indirebilirler.

"install" komutu kaynak kodunu indirmek, derlemek ve zincirin binary'sini yerel olarak kurmak için kullanılabilir. Binary, örneğin bir validatör node'unu başlatmak veya başlatıldıktan sonra zincirle etkileşime geçmek için kullanılabilir.

Ignite'ta yayınlanan tüm zincirler "list" komutu kullanılarak listelenebilir.

**Seçenekler**

```
  -h, --help   help for chain
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network](broken-reference) - Üretimde bir blockchain başlatın
* [ignite network chain init](broken-reference) - Yayınlanmış bir zincir kimliğinden bir zincir başlatma
* [ignite network chain install](broken-reference) - Bir başlatma (launch) için zincir binary'sini yükleme
* [ignite network chain join](broken-reference) - Bir ağa validatör olarak katılma talebi
* [ignite network chain launch](broken-reference) - Bir zincirin başlatılmasını tetikleyin
* [ignite network chain list](broken-reference) - Yayınlanan zincirleri listeleyin
* [ignite network chain prepare](broken-reference) - Zinciri başlatmak için hazırlayın
* [ignite network chain publish](broken-reference) - Yeni bir ağ başlatmak için yeni bir zincir yayınlayın
* [ignite network chain revert-launch](broken-reference) - Bir ağın koordinatör olarak başlatılmasını geri alma
* [ignite network chain show](broken-reference) - Bir zincirin ayrıntılarını göster

### ignite network chain init <a href="#ignite-network-chain-init" id="ignite-network-chain-init"></a>

Yayınlanmış bir zincir kimliğinden bir zincir başlatma

**Özet**

Ignite network chain init, Ignite zincirinde depolanan bilgilerden bir blockchain için bir validatör node'u başlatmak için validatörler tarafından kullanılan bir komuttur.

```
ignite network chain init 42
```

Bu komut, başlatma kimliği 42 olan bir zincir hakkındaki bilgileri getirir. Zincirin kaynak kodu geçici bir dizine klonlanır ve node'un binary'si kaynaktan derlenir. İkili kod daha sonra node'u başlatmak için kullanılır. Ignite varsayılan olarak blockchain için ana dizin olarak "\~/spn/\[launch-id]/" kullanır.

Bir validatör node'unun başlatılmasının önemli bir parçası gentx'in (zincirin başlangıcına bir validatör ekleyen bir işlem) oluşturulmasıdır.

"init" komutu, kendi kendini yetkilendirme ve komisyon gibi değerler isteyecektir. Bu değerler validatörün gentx'inde kullanılacaktır. Değerleri interaktif olmayan modda sağlamak için bayrakları kullanabilirsiniz.

Blockchain'in ana dizini için farklı bir yol seçmek için "--home" bayrağını kullanın:

```
ignite network chain init 42 --home ~/mychain
```

"init" komutunun sonucu, bir genesis validator transaction (gentx) dosyası içeren bir validatör ana dizinidir.

```
ignite network chain init [launch-id] [flags]
```

**Seçenekler**

```
      --check-dependencies                  verify that cached dependencies have not been modified since they were downloaded
      --clear-cache                         clear the build cache (advanced)
      --from string                         account name to use for sending transactions to SPN (default "default")
  -h, --help                                help for init
      --home string                         home directory used for blockchains
      --keyring-backend string              keyring backend to store your account keys (default "test")
      --keyring-dir string                  accounts keyring directory (default "/home/runner/.ignite/accounts")
      --validator-account string            account for the chain validator (default "default")
      --validator-details string            details about the validator
      --validator-gas-price string          validator gas price
      --validator-identity string           validator identity signature (ex. UPort or Keybase)
      --validator-moniker string            custom validator moniker
      --validator-security-contact string   validator security contact email
      --validator-self-delegation string    validator minimum self delegation
      --validator-website string            associate a website with the validator
  -y, --yes                                 answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain install <a href="#ignite-network-chain-install" id="ignite-network-chain-install"></a>

Bir başlatma (launch) için zincir binary'sini yükleme

```
ignite network chain install [launch-id] [flags]
```

**Seçenekler**

```
      --check-dependencies   verify that cached dependencies have not been modified since they were downloaded
      --clear-cache          clear the build cache (advanced)
      --from string          account name to use for sending transactions to SPN (default "default")
  -h, --help                 help for install
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain join <a href="#ignite-network-chain-join" id="ignite-network-chain-join"></a>

Bir ağa validatör olarak katılma talebi

**Özet**

"join" komutu validatörler tarafından bir blockchain'e katılma isteği göndermek için kullanılır. Gerekli argüman bir blockchainin başlatma kimliğidir.

"join" komutu, validatörün blockchain için bir ana dizin oluşturmuş olmasını ve "ignite network chain init" komutunu çalıştırarak ya da veri dizinini zincirin ikili dosyası ile manuel olarak başlatarak bir gentx'e sahip olmasını bekler.

Varsayılan olarak "join" komutu yalnızca bir validatör olarak katılma isteği gönderir. Bununla birlikte, genellikle bir validatörün kendi kendine delegasyon sağlayabilmesi için token bakiyesi olan bir genesis hesabı da talep etmesi gerekir.

Aşağıdaki komut, 42 numaralı başlatma kimliği ile blockchain'e validatör olarak katılma isteği gönderecek ve 95000000 STAKE token bakiyesine sahip bir hesap olarak eklenme talebinde bulunacaktır.

```
ignite network chain join 42 --amount 95000000stake
```

Validatör olarak katılma isteği bir gentx dosyası içerir. Ignite gentx dosyasını varsayılan olarak "ignite network chain init" tarafından kullanılan bir ev dizininde arar. Farklı bir dizin kullanmak için "--home" bayrağını kullanın veya "--gentx" bayrağı ile doğrudan bir gentx dosyası iletin.

Bir zincire validatör olarak katılmak için, diğer validatörlerin bağlanabilmesi için düğümünüzün IP adresini sağlamalısınız. Join komutu sizden IP adresini isteyecek ve değeri otomatik olarak tespit edip doldurmaya çalışacaktır. IP adresini manuel olarak belirtmek istiyorsanız "--peer-address" bayrağını kullanabilirsiniz:

```
ignite network chain join 42 --peer-address 0.0.0.0
```

"Join" Ignite blockchain'ine bir işlem yayınladığından, Ignite blockchain'inde bir hesaba ihtiyacınız olacaktır. Ancak testnet aşamasında Ignite otomatik olarak bir musluktan token talep eder.

```
ignite network chain join [launch-id] [flags]
```

**Seçenekler**

```
      --amount string            amount of coins for account request (ignored if coordinator has fixed the account balances or if --no-acount flag is set)
      --check-dependencies       verify that cached dependencies have not been modified since they were downloaded
      --from string              account name to use for sending transactions to SPN (default "default")
      --gentx string             path to a gentx json file
  -h, --help                     help for join
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --no-account               prevent sending a request for a genesis account
      --peer-address string      peer's address
  -y, --yes                      answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain launch <a href="#ignite-network-chain-launch" id="ignite-network-chain-launch"></a>

Bir zincirin başlatılmasını tetikleyin

**Özet**

Başlatma komutu dünyaya zincirin başlatılmaya hazır olduğunu bildirir.

Yalnızca zincirin koordinatörü başlatma komutunu uygulayabilir.

```
ignite network chain launch 42
```

Başlat komutu yürütüldükten sonra genesis üzerinde hiçbir değişiklik kabul edilmez. Örneğin, validatörler artık validatör olarak başvurmak için "ignite network chain join" komutunu başarıyla yürütemeyecektir.

Başlatma komutu zincirin başlayacağı tarih ve saati belirler. Varsayılan olarak geçerli saat ayarlanır. Validatörlere başlatmaya hazırlanmaları için daha fazla zaman vermek için zamanı "--launch-time" bayrağı ile ayarlayın:

```
ignite network chain launch 42 --launch-time 2023-01-01T00:00:00Z
```

Başlatma komutu yürütüldükten sonra, validatörler nihai genesis'i oluşturabilir ve node'larını başlatma için hazırlayabilir. Örneğin, doğrulayıcılar genesis oluşturmak ve eş listesini doldurmak için "ignite network chain prepare" komutunu çalıştırabilir.

Başlatma zamanını değiştirmek veya genesis dosyasını değişikliklere açmak istiyorsanız, örneğin yeni validatörleri kabul etmeyi ve hesap eklemeyi mümkün kılmak için "ignite network chain revert-launch" komutunu kullanabilirsiniz.

```
ignite network chain launch [launch-id] [flags]
```

**Seçenekler**

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for launch
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --launch-time string       timestamp the chain is effectively launched (example "2022-01-01T00:00:00Z")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain list <a href="#ignite-network-chain-list" id="ignite-network-chain-list"></a>

Yayınlanan zincirleri listeleyin

```
ignite network chain list [flags]
```

**Seçenekler**

```
      --advanced     show advanced information about the chains
  -h, --help         help for list
      --limit uint   limit of results per page (default 100)
      --page uint    page for chain list result (default 1)
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain prepare <a href="#ignite-network-chain-prepare" id="ignite-network-chain-prepare"></a>

Zinciri başlatmak (launch) için hazırlayın

**Özet**

prepare komutu, son oluşumu oluşturarak ve eşlerin IP adreslerini validatörün yapılandırma dosyasına ekleyerek bir validatör node'unu zinciri başlatma için hazırlar.

```
ignite network chain prepare 42
```

Ignite varsayılan olarak veri dizini olarak "$HOME/spn/LAUNCH\_ID" kullanır. Node'u başlatırken farklı bir veri dizini kullandıysanız, "--home" bayrağını kullanın ve veri dizinine giden doğru yolu ayarlayın.

Ignite "config/genesis.json" içinde genesis dosyasını oluşturur ve "config/config.toml" dosyasını değiştirerek eş IP'leri ekler.

Prepar komutu, koordinatör zincir başlatmayı tetikledikten ve "ignite network chain launch" ile genesis'i sonlandırdıktan sonra çalıştırılmalıdır. Ignite'ı "--force" bayrağı ile başlatmanın tetiklenip tetiklenmediğini kontrol etmeden prepare komutunu çalıştırmaya zorlayabilirsiniz (bu önerilmez).

Prepar komutu çalıştırıldıktan sonra node başlatılmaya hazırdır.

```
ignite network chain prepare [launch-id] [flags]
```

**Seçenekler**

```
      --check-dependencies       verify that cached dependencies have not been modified since they were downloaded
      --clear-cache              clear the build cache (advanced)
  -f, --force                    force the prepare command to run even if the chain is not launched
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for prepare
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain publish <a href="#ignite-network-chain-publish" id="ignite-network-chain-publish"></a>

Yeni bir ağ başlatmak için yeni bir zincir yayınlayın

**Özet**

Ignite ile bir blockchain başlatma sürecine başlamak için, bir koordinatörün bir blockchain hakkındaki bilgileri yayınlaması gerekir. Gerekli olan tek bilgi blockchainin kaynak kodunun URL'sidir.

Aşağıdaki komut örnek bir blockchain hakkındaki bilgileri yayınlar:

```
ignite network chain publish github.com/ignite/example
```

Bu komut blockchainin kaynak kodunu getirir, ikiliyi derler, bir blockchainin ikiliyle başlatılabileceğini doğrular ve blockchain hakkındaki bilgileri Ignite'ta yayınlar. Şu anda yalnızca halka açık depolar desteklenmektedir. Komut, Ignite'ta zincirin tanımlayıcısı olarak işlev gören bir tamsayı numarası döndürür.

Ignite'ta bir blockchain yayınlayarak bu blockchain'in "koordinatörü" olursunuz. Koordinatör, validatör taleplerini onaylama ve reddetme, blockchain parametrelerini ayarlama ve zincirin başlatılmasını tetikleme yetkisine sahip bir hesaptır.

Bir zincir yayınlanırken varsayılan Git dalı kullanılır. Belirli bir dal, etiket veya commit hash kullanmak istiyorsanız, sırasıyla "--branch", "--tag" veya "--hash" bayraklarını kullanın.

Depo adı varsayılan zincir kimliği olarak kullanılır. Ignite zincir kimliklerinin benzersiz olmasını sağlamaz, ancak geçerli bir biçime sahip olmaları gerekir: \[string]-\[integer]. Özel bir zincir kimliği ayarlamak için "--chain-id" bayrağını kullanın.

```
ignite network chain publish github.com/ignite/example --chain-id foo-1
```

Zincir yayınlandıktan sonra kullanıcılar, zincirin oluşumuna eklenecek jeton bakiyelerine sahip hesaplar talep edebilir. Varsayılan olarak, kullanıcılar istedikleri sayıda token talep etmekte serbesttir. Token talep eden tüm kullanıcıların aynı miktarda token almasını istiyorsanız, bir coin listesiyle birlikte "--account-balance" bayrağını kullanın.

```
ignite network chain publish github.com/ignite/example --account-balance 2000foocoin
```

```
ignite network chain publish [source-url] [flags]
```

**Seçenekler**

```
      --account-balance string   balance for each approved genesis account for the chain
      --amount string            amount of coins for account request
      --branch string            Git branch to use for the repo
      --chain-id string          chain ID to use for this network
      --check-dependencies       verify that cached dependencies have not been modified since they were downloaded
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
      --genesis-config string    name of an Ignite config file in the repo for custom Genesis
      --genesis-url string       URL to a custom Genesis
      --hash string              Git hash to use for the repo
  -h, --help                     help for publish
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --mainnet                  initialize a mainnet project
      --metadata string          add chain metadata
      --no-check                 skip verifying chain's integrity
      --project uint             project ID to use for this network
      --reward.coins string      reward coins
      --reward.height int        last reward height
      --shares string            add shares for the project
      --tag string               Git tag to use for the repo
      --total-supply string      add a total of the mainnet of a project
  -y, --yes                      answers interactive yes/no questions with yes
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain revert-launch <a href="#ignite-network-chain-revert-launch" id="ignite-network-chain-revert-launch"></a>

Bir ağın koordinatör olarak başlatılmasını geri alma

**Özet**

Revert launch komutu, bir zincirin önceden planlanmış başlatılmasını geri alır.

Yalnızca zincirin koordinatörü başlatma komutunu uygulayabilir.

```
ignite network chain revert-launch 42
```

Revert launch komutu çalıştırıldıktan sonra, zincirin oluşumunda değişiklik yapılmasına tekrar izin verilir. Örneğin, validatörler zincire katılma talebinde bulunabilecektir. Revert launch ayrıca başlatma zamanını da sıfırlar.

```
ignite network chain revert-launch [launch-id] [flags]
```

**Seçenekler**

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for revert-launch
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u launch için hazırlayın

### ignite network chain show <a href="#ignite-network-chain-show" id="ignite-network-chain-show"></a>

Bir zincirin ayrıntılarını göster

**Seçenekler**

```
  -h, --help   help for showssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain](broken-reference) - Bir zincir yayınlayın, validatör olarak katılın ve node'u başlatmak için hazırlayın
* [ignite network chain show accounts](broken-reference) - Zincirin tüm hakediş ve oluşum hesaplarını gösterin
* [ignite network chain show genesis](broken-reference) - Zincir genesis dosyasını göster
* [ignite network chain show info](broken-reference) - Zincirin bilgi ayrıntılarını göster
* [ignite network chain show peers](broken-reference) - Zincirin eşler listesini göster
* [ignite network chain show validators](broken-reference) - Zincirin tüm validatörlerini göster

### ignite network chain show accounts <a href="#ignite-network-chain-show-accounts" id="ignite-network-chain-show-accounts"></a>

Zincirin tüm hakediş ve oluşum hesaplarını gösterin

```
ignite network chain show accounts [launch-id] [flags]
```

**Seçenekler**

```
  -h, --help            help for accounts
      --prefix string   account address prefix (default "spn")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain show](broken-reference) - Bir zincirin ayrıntılarını göster

### ignite network chain show genesis <a href="#ignite-network-chain-show-genesis" id="ignite-network-chain-show-genesis"></a>

Zincir genesis dosyasını göster

```
ignite network chain show genesis [launch-id] [flags]
```

**Seçenekler**

```
      --clear-cache   clear the build cache (advanced)
  -h, --help          help for genesis
      --out string    path to output Genesis file (default "./genesis.json")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain show](broken-reference) - Bir zincirin ayrıntılarını göster

### ignite network chain show info <a href="#ignite-network-chain-show-info" id="ignite-network-chain-show-info"></a>

Zincirin bilgi ayrıntılarını göster

```
ignite network chain show info [launch-id] [flags]
```

**Seçenekler**

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain show](broken-reference) - Bir zincirin ayrıntılarını göster

### ignite network chain show peers <a href="#ignite-network-chain-show-peers" id="ignite-network-chain-show-peers"></a>

Zincirin eşler listesini göster

```
ignite network chain show peers [launch-id] [flags]
```

**Seçenekler**

```
  -h, --help         help for peers
      --out string   path to output peers list (default "./peers.txt")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain show](broken-reference) - Bir zincirin ayrıntılarını göster

### ignite network chain show validators <a href="#ignite-network-chain-show-validators" id="ignite-network-chain-show-validators"></a>

Zincirin tüm validatörlerini göster

```
ignite network chain show validators [launch-id] [flags]
```

**Seçenekler**

```
  -h, --help            help for validators
      --prefix string   account address prefix (default "spn")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network chain show](broken-reference) - Bir zincirin ayrıntılarını göster

### ignite network coordinator <a href="#ignite-network-coordinator" id="ignite-network-coordinator"></a>

Koordinatör profilini gösterme ve güncelleme

**Seçenekler**

```
  -h, --help   help for coordinator
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network](broken-reference) - Üretimde bir blockchain başlatın
* [ignite network coordinator set](broken-reference) - Koordinatör profilinde bir bilgi ayarlama
* [ignite network coordinator show](broken-reference) - Koordinatör profilini göster

### ignite network coordinator set <a href="#ignite-network-coordinator-set" id="ignite-network-coordinator-set"></a>

Koordinatör profilinde bir bilgi ayarlama

**Özet**

Ignite'taki koordinatörler, koordinatör için bir açıklama içeren bir profil ayarlayabilirler. Koordinatör seti komutu, koordinatör için bilgi ayarlamaya izin verir. Aşağıdaki bilgiler ayarlanabilir:

* details: koordinatör hakkında genel bilgi.
* identity: Keybase veya Veramo gibi bir sistemle koordinatörün kimliğini doğrulamak için bir bilgi parçası.
* website: koordinatörün web sitesi.

```
ignite network coordinator set details|identity|website [value] [flags]
```

**Seçenekler**

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for set
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network coordinator](broken-reference) - Koordinatör profilini gösterme ve güncelleme

### ignite network coordinator show <a href="#ignite-network-coordinator-show" id="ignite-network-coordinator-show"></a>

Koordinatör profilini göster

```
ignite network coordinator show [address] [flags]
```

**Seçenekler**

```
 -h, --help   help for show
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network coordinator](broken-reference) - Koordinatör profilini gösterme ve güncelleme

### ignite network profile <a href="#ignite-network-profile" id="ignite-network-profile"></a>

Adres profil bilgilerini göster

```
ignite network profile [project-id] [flags]
```

**Seçenekler**

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for profile
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network](broken-reference) - Üretimde bir blockchain başlatın

### ignite network project <a href="#ignite-network-project" id="ignite-network-project"></a>

Projelerle ilgilenin

**Seçenekler**

```
  -h, --help   help for project
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network](broken-reference) - Üretimde bir blockchain başlatın
* [ignite network project account](broken-reference) - Proje hesaplarını yönetin
* [ignite network project create](broken-reference) - Bir proje oluşturun
* [ignite network project list](broken-reference) - Yayınlanan projeleri listeleyin
* [ignite network project show](broken-reference) - Yayınlanan projeleri gösterin
* [ignite network project update](broken-reference) - Proje detaylarının güncellenmesi

### ignite network project account <a href="#ignite-network-project-account" id="ignite-network-project-account"></a>

Proje hesaplarını yönetin

**Seçenekler**

```
  -h, --help   help for account
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network project](broken-reference) - Projeleri yönetin
* [ignite network project account list](broken-reference) - Projenin tüm mainnet ve mainnet hakedişlerini göster

### ignite network project account list <a href="#ignite-network-project-account-list" id="ignite-network-project-account-list"></a>

Projenin tüm mainnet ve mainnet hakedişlerini göster

```
ignite network project account list [project-id] [flags]
```

**Seçenekler**

```
 -h, --help   help for list
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network project account](broken-reference) - Proje hesaplarını idare edin

### ignite network project create <a href="#ignite-network-project-create" id="ignite-network-project-create"></a>

Bir proje oluşturun

```
ignite network project create [name] [total-supply] [flags]
```

**Seçenekler**

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for create
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --metadata string          Add a metadata to the chain
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network project](broken-reference) - Projeleri idare edin

### ignite network project list <a href="#ignite-network-project-list" id="ignite-network-project-list"></a>

Yayınlanan projeleri listeleyin

```
ignite network project list [flags]
```

**Seçenekler**

```
  -h, --help   help for list
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network project](broken-reference) - Projeleri idare edin

### ignite network project show <a href="#ignite-network-project-show" id="ignite-network-project-show"></a>

Yayınlanmış projeyi göster

```
ignite network project show [project-id] [flags]
```

**Seçenekler**

```
  -h, --help   help for show
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network project](broken-reference) - Projeleri idare edin

### ignite network project update <a href="#ignite-network-project-update" id="ignite-network-project-update"></a>

Proje detaylarının güncellenmesi

```
ignite network project update [project-id] [flags]
```

**Seçenekler**

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for update
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --metadata string          update the project metadata
      --name string              update the project name
      --total-supply string      update the total of the mainnet of a project
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**Ayrıca Bakınız**

* [ignite network project](broken-reference) - Projeleri idare edin

### ignite network request <a href="#ignite-network-request" id="ignite-network-request"></a>

Talep oluşturma, gösterme, reddetme ve onaylama

**Özet**

"request" ad alanı, istek oluşturma, gösterme, onaylama ve reddetme komutlarını içerir.

Bir istek, Ignite'ta token bakiyeleri ve validatörleri olan hesaplar eklemek gibi genesis dosyasında değişiklik yapılmasına izin veren bir mekanizmadır. Herkes bir istek gönderebilir, ancak yalnızca bir zincirin koordinatörü bir isteği onaylayabilir veya reddedebilir.

Her talebin bir durumu vardır:

* Pending: koordinatörün onayını bekliyor
* Approved: koordinatör tarafından onaylanmış, içeriği lansman bilgilerine uygulanmıştır.
* Rejected: koordinatör veya talep oluşturucu tarafından reddedildi

**Seçenekler**

```
  -h, --help   help for request
```

**Üst komutlardan devralınan seçenekler**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

Ayrıca Bakınız

* [ignite network](broken-reference) - Üretimde bir blok zinciri başlatın
* [ignite network request add-account](broken-reference) - Hesap eklemek için istek gönder
* [ignite network request approve](broken-reference) - Talepleri onaylayın
* [ignite network request change-param](broken-reference) - Bir modül parametresini değiştirmek için istek gönderme
* [ignite network request list](broken-reference) - Bir zincir için tüm talepleri listeleyin
* [ignite network request reject](broken-reference) - Talepleri reddetme
* [ignite network request remove-account](broken-reference) - Bir genesis hesabını kaldırmak için istek gönderin
* [ignite network request remove-validator](broken-reference) - Validatörü kaldırmak için istek gönderme
* [ignite network request show](broken-reference) - Bir istek hakkında ayrıntılı bilgi gösterme
* [ignite network request verify](broken-reference) - Talebi doğrulayın ve bunlardan zincir oluşumunu simüle edin

### ignite network request add-account <a href="#ignite-network-request-add-account" id="ignite-network-request-add-account"></a>

Hesap eklemek için istek gönder

#### Özet

"add account" komutu, belirli bir adrese ve belirli bir coin bakiyesine sahip bir hesabı zincirin oluşumuna eklemek için yeni bir istek oluşturur.

Eğer başlatma bilgilerinde zaten aynı adrese sahip bir genesis hesabı ya da bir vesting hesabı belirtilmişse, talep otomatik olarak uygulanamaz.

Bir koordinatör bir zincirdeki tüm genesis hesaplarının aynı bakiyeye sahip olması gerektiğini belirtmişse (örneğin test ağları için kullanışlıdır), "hesap ekle" argüman olarak yalnızca bir adres bekler. Bir token bakiyesi sağlamaya çalışmak hatayla sonuçlanacaktır.

```
ignite network request add-account [launch-id] [address] [coins] [flags]
```

**Seçenekler**

```
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for add-account
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Create, show, reject and approve requests

### ignite network request approve <a href="#ignite-network-request-approve" id="ignite-network-request-approve"></a>

Talepleri onaylayın

#### Özet

"approve" komutu bir zincirin koordinatörü tarafından talepleri onaylamak için kullanılır. Virgülle ayrılmış bir liste ve/veya tire sözdizimi kullanılarak birden fazla talep onaylanabilir.

```
ignite network request approve 42 1,2,3-6,7,8
```

Yukarıdaki komut, başlatma kimliği 42 olan bir zincire dahil edilen 1'den 8'e kadar kimliklere sahip istekleri onaylar.

İstekler onaylandığında Ignite istenen değişiklikleri uygular ve zinciri yerel olarak başlatma ve başlatma simülasyonu yapar. Zincir başarılı bir şekilde başlarsa, talepler "doğrulanmış" olarak kabul edilir ve onaylanır. İstenen bir veya daha fazla değişiklik zincirin yerel olarak başlatılmasını engellerse, doğrulama işlemi başarısız olur ve tüm isteklerin onayı iptal edilir. Doğrulama işlemini atlamak için "--no-verification" bayrağını kullanın.

Ignite'ın istekleri, istek kimliklerinin "approve" komutuna gönderildiği sırayla onaylamaya çalışacağını unutmayın.

```
ignite network request approve [launch-id] [number<,...>] [flags]
```

**Seçenekler**

```
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for approve
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --no-verification          approve the requests without verifying them
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

**SEE ALSO**

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite network request change-param <a href="#ignite-network-request-change-param" id="ignite-network-request-change-param"></a>

Bir modül parametresini değiştirmek için istek gönderme

```
ignite network request change-param [launch-id] [module-name] [param-name] [value (json, string, number)] [flags]
```

**Seçenekler**

```
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for change-param
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite network request list <a href="#ignite-network-request-list" id="ignite-network-request-list"></a>

Bir zincir için tüm talepleri listeleyin

```
ignite network request list [launch-id] [flags]
```

**Seçenekler**

```
  -h, --help            help for list
      --prefix string   account address prefix (default "spn")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite network request reject

Talepleri reddetme

**Özet**

"reject" komutu bir zincirin koordinatörü tarafından istekleri reddetmek için kullanılır.

```
ignite network request reject 42 1,2,3-6,7,8
```

"reject" komutunun sözdizimi "approve" komutunun sözdizimine benzer.

```
ignite network request reject [launch-id] [number<,...>] [flags]
```

**Seçenekler**

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for reject
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite ağ isteği remove-validator

Validatörü kaldırmak için istek gönderme

```
ignite network request remove-account [launch-id] [address] [flags]
```

#### Seçenekler

```
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for remove-account
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite network istek gösterisi

Bir istek hakkında ayrıntılı bilgi gösterme

```
ignite network request remove-validator [launch-id] [address] [flags]
```

**Seçenekler**

```
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for remove-validator
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite network talep doğrulama

Talebi doğrulayın ve bunlardan zincir oluşumunu simüle edin

#### Özet

"verify" komutu, bu istekleri onaylamanın zincirin sorunsuz bir şekilde başlatılmasına izin veren geçerli bir oluşumla sonuçlanacağını doğrulamak için seçilen istekleri yerel olarak bir zincirin oluşumuna uygular. Bu komut istekleri onaylamaz, sadece kontrol eder.

```
ignite network request show [launch-id] [request-id] [flags]
```

**Seçenekler**

```
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for verify
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite network talep doğrulama

Talebi doğrulayın ve bunlardan zincir oluşumunu simüle edin

#### Özet

"verify" komutu, bu istekleri onaylamanın zincirin sorunsuz bir şekilde başlatılmasına izin veren geçerli bir oluşumla sonuçlanacağını doğrulamak için seçilen istekleri yerel olarak bir zincirin oluşumuna uygular. Bu komut istekleri onaylamaz, sadece kontrol eder.

```
ignite network request verify [launch-id] [number<,...>] [flags]
```

**Seçenekler**

```
      --clear-cache              clear the build cache (advanced)
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for verify
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network request](broken-reference) - Talep oluşturma, gösterme, reddetme ve onaylama

### ignite network ödülü

Ağ ödüllerini yönetin

**Seçenekler**

```
  -h, --help   help for reward
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network](broken-reference) - Üretimde bir blockchain başlatın
* [ignite network reward release](broken-reference) - Başlatılan zincirlerin izleme modüllerini SPN ile bağlayın
* [ignite network reward set](broken-reference) - bir ağ zinciri ödülü belirleyin

### ignite network ödül açiklamasi

Başlatılan zincirlerin izleme modüllerini SPN ile bağlayın

```
ignite network reward release [launch-id] [chain-rpc] [flags]
```

**Seçenekler**

```
      --create-client-only        only create the network client id
      --from string               account name to use for sending transactions to SPN (default "default")
  -h, --help                      help for release
      --keyring-backend string    keyring backend to store your account keys (default "test")
      --spn-gaslimit int          gas limit used for transactions on SPN (default 400000)
      --spn-gasprice string       gas price used for transactions on SPN (default "0.0000025uspn")
      --testnet-account string    testnet chain account (default "default")
      --testnet-faucet string     faucet address of the testnet chain
      --testnet-gaslimit int      gas limit used for transactions on testnet chain (default 400000)
      --testnet-gasprice string   gas price used for transactions on testnet chain (default "0.0000025stake")
      --testnet-prefix string     address prefix of the testnet chain (default "cosmos")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network reward](broken-reference) - Ağ ödüllerini yönetin

### ignite network ödül seti̇

bir ağ zinciri ödülü belirleyin

```
ignite network reward set [launch-id] [last-reward-height] [coins] [flags]
```

#### Seçenekler

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for set
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network reward](broken-reference) - Ağ ödüllerini yönetin

### ignite ağ aracı

Yardımcı araçları çalıştırmak için komutlar

#### Seçenekler

```
  -h, --help   help for tool
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network](broken-reference) - Üretimde bir blockchain başlatın
* [ignite network tool proxy-tunnel](broken-reference) - HTTP üzerinden bir proxy tüneli kurma

### ignite ağ aracı proxy-tunnel

HTTP üzerinden bir proxy tüneli kurma

#### Özet

HTTP tünellemesine ihtiyaç duyan her düğüm için bir HTTP proxy sunucusu ve HTTP proxy istemcileri başlatır.

HTTP tünelleme SADECE SPN\_CONFIG\_FILE içinde tünellenmiş eşlerin/düğümlerin listesini içeren "tunneled\_peers" alanı varsa etkinleştirilir.

SPN'yi koordinatör olarak kullanıyorsanız ve HTTP tünelleme özelliğine hiç izin vermek istemiyorsanız, düz TCP bağlantıları yerine HTTP tünellemenin etkin olduğu doğrulayıcı isteklerini onaylamayarak "spn.yml" dosyasının oluşturulmasını önleyebilirsiniz.

```
ignite network tool proxy-tunnel SPN_CONFIG_FILE [flags]
```

**Seçenekler**

```
  -h, --help   help for proxy-tunnel
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network tool](broken-reference) - Yardımcı araçları çalıştırmak için komutlar

### ignite ağ validatörü

Validatör profilini gösterme ve güncelleme

#### Seçenekler

```
  -h, --help   help for validator
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network](broken-reference) - Üretimde bir blok zinciri başlatın
* [ignite network validator set](broken-reference) - Validatör profilinde bir bilgi ayarlama
* [ignite network validator show](broken-reference) - Validatör profilini göster

### ignite ağ doğrulayıcı seti

Validatör profilinde bir bilgi ayarlama

#### Özet

Ignite üzerindeki validatörler, validatör için bir açıklama içeren bir profil ayarlayabilir. Validator set komutu validator için bilgi ayarlamaya izin verir. Aşağıdaki bilgiler ayarlanabilir:

* details: validatör hakkında genel bilgi.
* identity: validatörün kimliğini aşağıdaki gibi bir sistemle doğrulamak için bilgi parçası Veramo'nun anahtar tabanı.
* &#x20;web sitesi: validatörün web sitesi.
* güvenlik: validatör için güvenlik irtibatı.

```
ignite network validator set details|identity|website|security [value] [flags]
```

#### Seçenekler

```
      --from string              account name to use for sending transactions to SPN (default "default")
  -h, --help                     help for set
      --home string              home directory used for blockchains
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network validator](broken-reference) - Doğrulayıcı profilini gösterme ve güncelleme

### ignite ağ doğrulayici gösteri̇si̇

Doğrulayıcı profilini göster

```
ignite network validator show [address] [flags]
```

**Options**

```
  -h, --help   help for show
```

**Options inherited from parent commands**

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network validator](broken-reference) - Doğrulayıcı profilini gösterme ve güncelleme

### ignite ağ sürümü

Eklentinin sürümü

#### Özet

Bir zincirle etkileşim için kullanılacak eklentinin sürümü koordinatör tarafından belirtilebilir.

```
ignite network version [flags]
```

**Seçenekler**

```
  -h, --help   help for version
```

#### Üst komutlardan devralınan seçenekler

```
      --local                       Use local SPN network
      --nightly                     Use nightly SPN network
      --spn-faucet-address string   SPN faucet address (default "https://faucet.devnet.ignite.com:443")
      --spn-node-address string     SPN node address (default "https://rpc.devnet.ignite.com:443")
```

#### Ayrıca Bakınız

* [ignite network](broken-reference) - Üretimde bir blockchain başlatın

### ignite node <a href="#ignite-node" id="ignite-node"></a>

Canlı bir blockchain node'una istekte bulunun

**Seçenekler**

```
  -h, --help          help for node
      --node string   <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite](broken-reference) - Ignite CLI, blockchain'inizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite node query](broken-reference) - Alt komutları sorgulama
* [ignite node tx](broken-reference) - İşlemler alt komutları

### ignite node query <a href="#ignite-node-query" id="ignite-node-query"></a>

Querying subcommands

#### Seçenekler

```
  -h, --help   help for query
```

#### Üst komutlardan devralınan seçenekler

```
      --node string   <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite node](broken-reference) - Canlı bir blockchain node'una istekte bulunun
* [ignite node query bank](broken-reference) - Banka modülü için sorgulama komutları
* [ignite node query tx](broken-reference) - Hash'e göre işlem sorgusu

### ignite node sorgu bankası

Banka modülü için sorgulama komutları

**Seçenekler**

```
  -h, --help   help for bank
```

#### Üst komutlardan devralınan seçenekler

```
      --node string   <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite node query](broken-reference) - Alt komutları sorgulama
* [ignite node query bank balances](broken-reference) - Hesap bakiyelerini hesap adına veya adrese göre sorgulama

### ignite node banka bakiyelerini sorgula

Query for account balances by account name or address

```
ignite node query bank balances [from_account_or_address] [flags]
```

#### Seçenekler

```
      --address-prefix string    account address prefix (default "cosmos")
      --count-total              count total number of records in all balances to query for
  -h, --help                     help for balances
      --home string              directory where the blockchain node is initialized
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --limit uint               pagination limit of all balances to query for (default 100)
      --offset uint              pagination offset of all balances to query for
      --page uint                pagination page of all balances to query for. This sets offset to a multiple of limit (default 1)
      --page-key string          pagination page-key of all balances to query for
      --reverse                  results are sorted in descending order
```

#### Üst komutlardan devralınan seçenekler

```
      --node string   <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite node query bank](broken-reference) - Banka modülü için sorgulama komutları

### ignite node query tx <a href="#ignite-node-query-tx" id="ignite-node-query-tx"></a>

Hash'e göre işlem sorgusu

```
ignite node query tx [hash] [flags]
```

#### Seçenekler

```
  -h, --help   help for tx
```

#### Üst komutlardan devralınan seçenekler

```
      --node string   <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite node query](broken-reference) - Alt komutları sorgulama

### ignite node tx <a href="#ignite-node-tx" id="ignite-node-tx"></a>

İşlemler alt komutları

**Options**

```
      --address-prefix string    account address prefix (default "cosmos")
      --fees string              fees to pay along with transaction; eg: 10uatom
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default "auto")
      --gas-prices string        gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            build an unsigned transaction and write it to STDOUT
  -h, --help                     help for tx
      --home string              directory where the blockchain node is initialized
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Üst komutlardan devralınan seçenekler

```
      --node string   <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite node](broken-reference) - Canlı bir blockchain node'una istekte bulunun
* [ignite node tx bank](broken-reference) - Banka işlemi alt komutları

### ignite node tx bankasi

Banka işlemi alt komutları

**A**Seçenekler**yarlar**

```
  -h, --help   help for bank
```

#### Üst komutlardan devralınan seçenekler

```
      --address-prefix string    account address prefix (default "cosmos")
      --fees string              fees to pay along with transaction; eg: 10uatom
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default "auto")
      --gas-prices string        gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            build an unsigned transaction and write it to STDOUT
      --home string              directory where the blockchain node is initialized
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --node string              <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite node tx](broken-reference) - İşlemler alt komutları
* [ignite node tx bank send](broken-reference) - Bir hesaptan diğerine para gönderin.

Bir hesaptan diğerine para gönderin.

```
ignite node tx bank send [from_account_or_address] [to_account_or_address] [amount] [flags]
```

#### Seçenekler

```
  -h, --help   help for send
```

Üst komutlardan devralınan seçenekler

```
      --address-prefix string    account address prefix (default "cosmos")
      --fees string              fees to pay along with transaction; eg: 10uatom
      --gas string               gas limit to set per-transaction; set to "auto" to calculate sufficient gas automatically (default "auto")
      --gas-prices string        gas prices in decimal format to determine the transaction fee (e.g. 0.1uatom)
      --generate-only            build an unsigned transaction and write it to STDOUT
      --home string              directory where the blockchain node is initialized
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
      --node string              <host>:<port> to tendermint rpc interface for this chain (default "https://rpc.cosmos.network:443")
```

#### Ayrıca Bakınız

* [ignite node tx bank](broken-reference) - Banka işlemi alt komutları

### ignite eklentisi

Eklentileri idare edin

#### Seçenekler

```
  -h, --help   help for plugin
```

#### Ayrıca Bakınız

* [ignite](broken-reference) - Ignite CLI, blockchain'inizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite plugin add](broken-reference) - Eklenti yapılandırmasına bir eklenti bildirimi ekler
* [ignite plugin describe](broken-reference) - Kayıtlı bir eklenti hakkında çıktı bilgisi
* [ignite plugin list](broken-reference) - Bildirilen eklentileri ve durumlarını listeleme
* [ignite plugin remove](broken-reference) - Zincirin eklenti yapılandırmasından bir eklenti bildirimini kaldırır
* [ignite plugin scaffold](broken-reference) - Scaffold yeni bir eklenti
* [ignite plugin update](broken-reference) - Eklentileri güncelleyin

### ignite eklentisi ekle

Eklenti yapılandırmasına bir eklenti bildirimi ekler

#### Özet

Bir eklenti yapılandırmasına bir eklenti bildirimi ekler. Oluşturulan yapılandırma tanımına eklenecek eklenti yolundan sonra bildirilen anahtar değer çiftlerine uyar. Örnek: ignite plugin add github.com/org/my-plugin/ foo=bar baz=qux

```
ignite plugin add [path] [key=value]... [flags]
```

#### Seçenekler

```
  -g, --global   use global plugins configuration ($HOME/.ignite/plugins/plugins.yml)
  -h, --help     help for add
```

#### Ayrıca Bakınız

* [ignite plugin](https://docs.ignite.com/references/cli#ignite-plugin) - Eklentileri idare et

### ignite eklentisi açıklaması

Kayıtlı bir eklenti hakkında çıktı bilgisi

#### Özet

Kayıtlı bir eklenti komutları ve kancaları hakkında çıktı bilgisi.

```
ignite plugin describe [path] [flags]
```

#### Seçenekler

```
  -h, --help   help for describe
```

#### Ayrıca Bakınız

* [ignite plugin](broken-reference) - Eklentileri idare edin

### ignite eklenti listesi

Bildirilen eklentileri ve durumlarını listeleme

#### Özet

Bildirilen eklentilerin durumunu ve bilgilerini yazdırır

```
ignite plugin list [flags]
```

**Options**

```
  -h, --help   help for list
```

**SEE ALSO**

* [ignite plugin](broken-reference) - Handle plugins

### ignite eklentisi kaldırma

Removes a plugin declaration from a chain's plugin configuration

```
ignite plugin remove [path] [flags]
```

#### Seçenekler

```
  -g, --global   use global plugins configuration ($HOME/.ignite/plugins/plugins.yml)
  -h, --help     help for remove
```

#### Ayrıca Bakınız

[ignite plugin](broken-reference) - Eklentileri idare et

### ignite eklenti scaffold'u

Scaffold yeni bir eklenti

#### Özet

Geçerli dizinde, verilen depo yolu yapılandırılmış yeni bir eklenti iskeleler. Geçerli dizin zaten bir git deposu değilse, verilen modül adıyla bir git deposu oluşturulacaktır.

```
ignite plugin scaffold [github.com/org/repo] [flags]
```

#### Seçenekler

```
  -h, --help   help for scaffold
```

#### Ayrıca Bakınız

ignite eklentisi - Eklentileri idare et

### ignite eklentisi güncellemesi

Eklentileri güncelleyin

#### Özet

Yol ile belirtilen bir eklentiyi günceller. Yol belirtilmezse, bildirilen tüm eklentiler güncellenir

```
ignite plugin update [path] [flags]
```

#### Seçenekler

```
  -h, --help   help for update
```

#### Ayrıca Bakınız

* [ignite plugin](broken-reference) - Eklentileri idare edin

### ignite relayer <a href="#ignite-relayer" id="ignite-relayer"></a>

Blockchain'leri bir IBC aktarıcı ile bağlayın

#### Seçenekler

```
  -h, --help   help for relayer
```

#### Ayrıca Bakınız

* [ignite](broken-reference) - Ignite CLI, blok zincirinizi scaffold etmek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite relayer configure](broken-reference) - Aktarma için kaynak ve hedef zincirlerini yapılandırma
* [ignite relayer connect](broken-reference) - Aktarma için kaynak ve hedef zincirlerini yapılandırma

Aktarma için kaynak ve hedef zincirlerini yapılandırma

```
ignite relayer configure [flags]
```

#### Seçenekler

```
  -a, --advanced                  advanced configuration options for custom IBC modules
  -h, --help                      help for configure
      --keyring-backend string    keyring backend to store your account keys (default "test")
      --keyring-dir string        accounts keyring directory (default "/home/runner/.ignite/accounts")
      --ordered                   set the channel as ordered
  -r, --reset                     reset the relayer config
      --source-account string     source Account
      --source-client-id string   use a custom client id for source
      --source-faucet string      faucet address of the source chain
      --source-gaslimit int       gas limit used for transactions on source chain
      --source-gasprice string    gas price used for transactions on source chain
      --source-port string        IBC port ID on the source chain
      --source-prefix string      address prefix of the source chain
      --source-rpc string         RPC address of the source chain
      --source-version string     module version on the source chain
      --target-account string     target Account
      --target-client-id string   use a custom client id for target
      --target-faucet string      faucet address of the target chain
      --target-gaslimit int       gas limit used for transactions on target chain
      --target-gasprice string    gas price used for transactions on target chain
      --target-port string        IBC port ID on the target chain
      --target-prefix string      address prefix of the target chain
      --target-rpc string         RPC address of the target chain
      --target-version string     module version on the target chain
```

#### Ayrıca Bakınız

[ignite relayer](https://docs.ignite.com/references/cli#ignite-relayer) - Blockchain'leri bir IBC relayer ile bağlayın

### ignite relayer connect

Yollarla ilişkili zincirleri bağlayın ve aradaki tx paketlerini aktarmaya başlayın

```
ignite relayer connect [<path>,...] [flags]
```

#### Seçenekler

```
  -h, --help                     help for connect
      --keyring-backend string   keyring backend to store your account keys (default "test")
      --keyring-dir string       accounts keyring directory (default "/home/runner/.ignite/accounts")
```

#### Ayrıca Bakınız

ignite relayer - Blockchain'leri bir IBC relayer ile bağlayın

### ignite scaffold

Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun

#### Özet

İskele, uygulamanızın önemli parçaları için kod oluşturmanın hızlı bir yoludur.

Her bir skaffolding hedefi (zincir, modül, mesaj, vb.) hakkında ayrıntılar için ilgili komutu "--help" bayrağı ile çalıştırın, örneğin, "ignite scaffold chain --help".

Ignite ekibi, scaffolding komutlarını çalıştırmadan önce kodun bir sürüm kontrol sistemine işlenmesini şiddetle tavsiye eder. Bu, kaynak koddaki değişiklikleri görmeyi ve değişiklikleri geri almaya karar verdiyseniz komutu geri almayı kolaylaştıracaktır.

Chain scaffolding komutu ile oluşturduğunuz bu blockchain, modüler Cosmos SDK çerçevesini kullanır ve proof of stake, token transferi, blockchainler arası bağlantı, yönetişim ve daha fazlası gibi işlevler için birçok standart modülü içe aktarır. Özel işlevsellik, geleneksel olarak "x/" dizininde bulunan modüllerde uygulanır. Varsayılan olarak, blockchain'iniz boş bir özel modülle birlikte gelir. Ek bir modül oluşturmak için module scaffolding komutunu kullanın.

Boş bir özel modül fazla bir şey yapmaz, temelde işlemlerin işlenmesinden ve uygulama durumunun değiştirilmesinden sorumlu olan mantık için bir kapsayıcıdır. Cosmos SDK blockchain'leri, bir veya daha fazla mesaj içeren, kullanıcı tarafından gönderilen imzalı işlemleri işleyerek çalışır. Bir mesaj, bir durum geçişini tanımlayan veriler içerir. Bir modül, herhangi bir sayıda mesajın işlenmesinden sorumlu olabilir.

Bir mesaj scaffolding komutu, yeni bir Cosmos SDK mesajı türünü işlemek için kod oluşturacaktır. Mesaj alanları, mesajın hatasız işlenmesi durumunda üretmesi amaçlanan durum geçişini tanımlar.

Mesajları scaffolding etmek, modülünüzün gerçekleştirebileceği ayrı "eylemler" oluşturmak için kullanışlıdır. Ancak bazen, blockchain'inizin belirli bir türün örneklerini oluşturma, okuma, güncelleme ve silme (CRUD) işlevselliğine sahip olmasını istersiniz. Verileri nasıl saklamak istediğinize bağlı olarak, bir tür için CRUD işlevselliğinin scaffold'unu oluşturan üç komut vardır: list, map ve single. Bu komutlar dört mesaj (her CRUD eylemi için bir tane) ve verileri depoya ekleme, silme ve depodan alma mantığını oluşturur. Yalnızca mantığı scaffold etmek istiyorsanız, örneğin mesajları ayrı ayrı scaffold etmeye karar verdiyseniz, bunu "--no-message" bayrağıyla da yapabilirsiniz.

Bir blockchain'den veri okumak sorgular yardımıyla gerçekleşir. Veri yazmak için mesajları nasıl scaffold haline getirebildiğinize benzer şekilde, blockchain uygulamanızdan verileri geri okumak için sorguları scaffold haline getirebilirsiniz.

Ayrıca, sadece proto mesaj açıklamasına sahip yeni bir protokol tampon dosyası üreten bir türü de scaffold edebilirsiniz. Proto mesajlarının Go türleri ürettiğini (ve bunlara karşılık geldiğini), Cosmos SDK mesajlarının ise "Msg" hizmetindeki proto "rpc" ye karşılık geldiğini unutmayın.

Özel IBC mantığına sahip bir uygulama oluşturuyorsanız, IBC paketlerini scaffold etmeniz gerekebilir. Bir IBC paketi, bir blockchain'den diğerine gönderilen verileri temsil eder. IBC paketlerini yalnızca "--ibc" bayrağı ile scaffold edilmiş IBC özellikli modüllerde scaffold edebilirsiniz. Varsayılan modülün IBC özellikli olmadığını unutmayın.

#### Seçenekler

```
  -h, --help   help for scaffold
```

#### Ayrıca Bakınız

* [ignite](broken-reference) - Ignite CLI, blockchain'inizi scaffolding etmek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar
* [ignite scaffold chain](broken-reference) - Yeni Cosmos SDK blockchain
* [ignite scaffold list](broken-reference) - Dizi olarak saklanan veriler için CRUD
* [ignite scaffold map](broken-reference) - Anahtar-değer çiftleri olarak saklanan veriler için CRUD
* [ignite scaffold message](broken-reference) - Blockchain üzerinde durum geçişi gerçekleştirmek için mesaj
* [ignite scaffold module](broken-reference) - Özel Cosmos SDK modülü
* [ignite scaffold packet](broken-reference) - Bir IBC paketi göndermek için mesaj
* [ignite scaffold query](broken-reference) - Bir blok zincirinden veri almak için sorgu
* [ignite scaffold react](broken-reference) - React web uygulaması şablonu
* [ignite scaffold single](broken-reference) - Tek bir konumda depolanan veriler için CRUD
* [ignite scaffold type](broken-reference) - Tip tanımı
* [ignite scaffold vue](broken-reference) - Vue 3 web uygulaması şablonuignite scaffold zinciri

Yeni Cosmos SDK blockchain

#### Özet

Uygulamaya özel yeni bir Cosmos SDK blok zinciri oluşturun.

Örneğin, aşağıdaki komut "hello/" dizininde "hello" adında bir blok zinciri oluşturacaktır:

```
ignite scaffold chain hello
```

Proje adı basit bir ad veya bir URL olabilir. İsim, proje için Go modül yolu olarak kullanılacaktır. Proje adı örnekleri:

```
ignite scaffold chain foo
ignite scaffold chain foo/bar
ignite scaffold chain example.org/foo
ignite scaffold chain github.com/username/foo
    
```

Geçerli dizinde kaynak kod dosyalarını içeren yeni bir dizin oluşturulacaktır. Farklı bir yol kullanmak için "--path" bayrağını kullanın.

Blockchain'inizin mantığının çoğu özel modüllerde yazılır. Her modül, bağımsız bir işlevsellik parçasını etkili bir şekilde kapsüller. Cosmos SDK kurallarına göre, özel modüller "x/" dizini içinde saklanır. Varsayılan olarak Ignite, projenin adıyla eşleşen bir ada sahip bir modül oluşturur. Varsayılan modül olmadan bir blok zinciri oluşturmak için "--no-module" bayrağını kullanın. Bir proje oluşturulduktan sonra "ignite scaffold module" komutu ile ek modüller eklenebilir.

Cosmos SDK tabanlı blok zincirlerindeki hesap adreslerinin dize önekleri vardır. Örneğin, Cosmos Hub blok zinciri varsayılan "cosmos" önekini kullanır, böylece adresler aşağıdaki gibi görünür: "cosmos12fjzdtqfrrve7zyg9sv8j25azw2ua6tvu07ypf". Özel bir adres öneki kullanmak için "--address-prefix" bayrağını kullanın. Örneğin:

```
ignite scaffold chain foo --address-prefix bar
```

Ignite, bir blockchain'in kaynak kodunu derlerken varsayılan olarak derleme sürecini hızlandırmak için bir önbellek oluşturur. Bir blockchain oluştururken önbelleği temizlemek için "--clear-cache" bayrağını kullanın. Bu bayrağı kullanmanız pek gerekmeyecektir.

Blok zinciri Cosmos SDK modüler blockchain Cosmos SDK hakkında daha fazla bilgi için [https://docs.cosmos.network](https://docs.cosmos.network)

```
ignite scaffold chain [name] [flags]
```

#### Seçenekler

```
      --address-prefix string   account address prefix (default "cosmos")
      --clear-cache             clear the build cache (advanced)
  -h, --help                    help for chain
      --no-module               create a project without a default module
  -p, --path string             create a project in a specific path (default ".")
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite scaffold listesi

Dizi olarak saklanan veriler için CRUD

#### Özet

"list" scaffolding komutu, blockchain durumunda bir liste olarak saklanan verilerin saklanması ve bunlarla etkileşim kurulması mantığını uygulayan dosyalar oluşturmak için kullanılır.

Komut, yeni bir veri türünün adı olarak kullanılacak bir NAME argümanını kabul eder. Ayrıca türü tanımlayan bir ALAN listesi de kabul eder.

Verilerle etkileşim, oluşturma, okuma, güncelleme ve silme (CRUD) modelini takip eder. Her tür için blockchain'e veri yazmak üzere üç Cosmos SDK mesajı tanımlanmıştır: MsgCreate{Name}, MsgUpdate{Name}, MsgDelete{Name}. Veri okumak için iki sorgu tanımlanmıştır: {Name} ve {Name}All. Tip, mesajlar ve sorgular "proto/" dizininde protokol tampon mesajları olarak tanımlanır. Mesajlar ve sorgular sırasıyla "Msg" ve "Query" servislerine bağlanır.

Mesajlar işlendiğinde, uygun keeper yöntemleri çağrılır. Geleneksel olarak, yöntemler "x/{moduleName}/keeper/msgserver{name}.go" içinde tanımlanır. Alma, ayarlama, kaldırma ve ekleme için yararlı yöntemler "{name}.go" içinde aynı " keeper" paketinde tanımlanmıştır.

"list" komutu temelde yeni bir veri türü tanımlamanıza olanak tanır ve türün örneklerini oluşturma, okuma, güncelleme ve silme mantığını sağlar. Örneğin, bir gönderi listesini işlemek için kod üreten ve her gönderinin "başlık" ve "gövde" alanlarına sahip olduğu bir komutu inceleyelim:

```
ignite scaffold list post title body
```

Bu size bir "Post" tipi, MsgCreatePost, MsgUpdatePost, MsgDeletePost ve iki sorgu sağlar: Post ve PostAll. Derlenmiş CLI, diyelim ki ikilisi "blogd" ve modülü "blog", zinciri sorgulamak (bkz. "blogd q blog") ve yukarıdaki mesajlarla işlemleri yayınlamak (bkz. "blogd tx blog") için komutlara sahiptir.

List komutu ile oluşturulan kodun düzenlenmesi ve uygulama ihtiyaçlarınıza göre uyarlanması amaçlanmıştır. Bu kodu, daha sonra uygulayacağınız gerçek iş mantığı için bir "iskelet" olarak düşünün.

Varsayılan olarak, tüm alanların string olduğu varsayılır. Farklı türde bir alan istiyorsanız, bunu iki nokta üst üste ":" işaretinden sonra belirtebilirsiniz. Şu türler desteklenir: string, bool, int, uint, coin, array.string, array.int, array.uint, array.coin. Alan türlerinin kullanımına bir örnek:

```
ignite scaffold list pool amount:coin tags:array.string height:int
```

Desteklenen tipler:

| Type         | Alias   | Index | Code Type | Description                          |
| ------------ | ------- | ----- | --------- | ------------------------------------ |
| string       | -       | yes   | string    | Metin türü                           |
| array.string | strings | no    | \[]string | Metin türü listesi                   |
| bool         | -       | yes   | bool      | Boolean tipi                         |
| int          | -       | yes   | int32     | Tamsayı tipi                         |
| array.int    | ints    | no    | \[]int32  | Tamsayı türlerinin listesi           |
| uint         | -       | yes   | uint64    | İşaretsiz tamsayı türü               |
| array.uint   | uints   | no    | \[]uint64 | İşaretsiz tamsayı türlerinin listesi |
| coin         | -       | no    | sdk.Coin  | Cosmos SDK coin türü                 |
| array.coin   | coins   | no    | sdk.Coins | Cosmos SDK coin türlerinin listesi   |

"Dizin", türün "ignite iskele haritasında" bir dizin olarak kullanılıp kullanılamayacağını gösterir.

Ignite ayrıca özel türleri de destekler:

```
ignite scaffold list product-details name desc
ignite scaffold list product price:coin details:ProductDetails
```

Yukarıdaki örnekte, önce "ProductDetails" türü tanımlanmış ve ardından "details" alanı için özel tür olarak kullanılmıştır. Ignite henüz özel tür dizilerini desteklememektedir.

Zinciriniz JSON notasyonundaki özel türleri kabul edecektir:

```
exampled tx example create-product 100coin '{"name": "x", "desc": "y"}' --from alice
```

Varsayılan olarak kod, projenizin adıyla eşleşen modülde scaffold haline getirilecektir. Projenizde birden fazla modül varsa, farklı bir modül belirtmek isteyebilirsiniz:

```
ignite scaffold list post title body --module blog
```

Varsayılan olarak her mesaj, işlemi imzalayanın adresini temsil eden bir "yaratıcı" alanıyla birlikte gelir. Bu alanın adını bir bayrakla özelleştirebilirsiniz:

```
ignite scaffold list post title body --signer author
```

CRUD mesajları olmadan sadece getter/setter mantığını iskelelemek mümkündür. Bu, yöntemlerin bir türü işlemesini istediğinizde, ancak mesajları manuel olarak iskele etmek istediğinizde kullanışlıdır. Mesaj iskelesini atlamak için bir bayrak kullanın:

```
ignite scaffold list post title body --no-message
```

Bir liste "--no-message" bayrağı ile scaffold haline getirilirse "creator" alanı oluşturulmaz.

```
ignite scaffold list NAME [field]... [flags]
```

#### Seçenekler

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for list
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite scaffold haritası

Anahtar-değer çiftleri olarak saklanan veriler için CRUD

#### Özet

"map" iskele komutu, blok zinciri durumunda anahtar-değer çiftleri (veya bir sözlük) olarak depolanan verilerin depolanması ve bunlarla etkileşim kurulması mantığını uygulayan dosyalar oluşturmak için kullanılır.

"map" komutu "ignite scaffold list" komutuna çok benzer, tek fark değerlerin nasıl indekslendiğidir. "List" ile değerler artan bir tamsayı ile indekslenirken, "list" değerleri kullanıcı tarafından sağlanan bir değerle (ya da birden fazla değerle) indekslenir.

Aynı blog yazısı örneğini kullanalım:

```
ignite scaffold map post title body
```

Bu komut, gönderi oluşturmak, okumak, güncellemek ve silmek için bir "Gönderi" türü ve CRUD işlevselliği scaffold'u oluşturur. Ancak, zincirinizin binary'si ile yeni bir gönderi oluştururken (veya zincirin API'si aracılığıyla bir işlem gönderirken) bir "index" sağlamanız gerekecektir:

```
blogd tx blog create-post [index] [title] [body]
blogd tx blog create-post hello "My first post" "This is the body"
```

Bu komut bir gönderi oluşturacak ve bunu blockchain durumunda "hello" dizini altında saklayacaktır. "hello" anahtarı için sorgulama yaparak gönderinin değerini geri getirebileceksiniz.

```
blogd q blog show-post hello
```

İndeksi özelleştirmek için "--index" bayrağını kullanın. Değerleri sorgulamayı basitleştiren birden fazla indeks sağlanabilir. Örneğin:

```
ignite scaffold map product price desc --index category,guid
```

Bu komutla, hem bir kategori hem de bir GUID (global olarak benzersiz kimlik) tarafından indekslenen bir "Ürün" değeri elde edersiniz. Bu, aynı kategoriye sahip ancak farklı GUID'ler kullanan ürün değerlerini programlı olarak getirmenize olanak tanır.

"list" ve "map" scaffolding'in davranışı çok benzer olduğundan, "--no-message", "--module", "--signer" bayraklarının yanı sıra özel tipler için iki nokta üst üste sözdizimini de kullanabilirsiniz.

```
ignite scaffold map NAME [field]... [flags]
```

#### Seçenekler

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for map
      --index strings   fields that index the value (default [index])
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blok zinciri, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite iskele mesajı

Blockchain üzerinde durum geçişi gerçekleştirmek için mesaj

#### Özet

Mesaj scaffold'u, belirli Cosmos SDK mesajlarını işlemek için blockchain'inize hızlı bir şekilde işlevsellik eklemek için kullanışlıdır.

Mesajlar, nihai amacı blockchain üzerinde durum geçişlerini tetiklemek olan nesnelerdir. Bir mesaj, blockchain'in durumunun nasıl değişeceğini etkileyen veri alanları için bir kaptır. Mesajları bir kullanıcının gerçekleştirebileceği "eylemler" olarak düşünebilirsiniz.

Örneğin, banka modülünde hesaplar arasında token transferleri için bir "Gönder" mesajı bulunur. Gönder mesajının üç alanı vardır: adresten (gönderici), adrese (alıcı) ve token miktarı. Bu mesaj başarıyla işlendiğinde, token tutarı gönderenin hesabından düşülecek ve alıcının hesabına eklenecektir.

Ignite'ın mesaj scaffold, yeni mesaj türleri oluşturmanıza ve bunları zincirinize eklemenize olanak tanır. Örneğin:

```
ignite scaffold message add-pool amount:coins denom active:bool --module dex
```

Yukarıdaki komut, üç alana sahip yeni bir MsgAddPool mesajı oluşturacaktır: miktar (jeton cinsinden), denom (bir dize) ve aktif (bir boolean). Mesaj "dex" modülüne eklenecektir.

Varsayılan olarak, mesaj "proto/{app}/{module}/tx.proto" içinde bir proto mesajı olarak tanımlanır ve "Msg" servisine kaydedilir. MsgAddPool ile bir işlem oluşturmak ve yayınlamak için bir CLI komutu modülün "cli" paketinde oluşturulur. Ek olarak Ignite, bir mesaj kurucusunu ve sdk.Msg arayüzünü karşılayacak ve mesajı modüle kaydedecek kodu scaffold eder.

En önemlisi "keeper" paketinde Ignite bir "AddPool" fonksiyonunu iskeleye yerleştirir. Bu fonksiyonun içinde mesaj işleme mantığını uygulayabilirsiniz.

Bir mesaj başarıyla işlendiğinde veri döndürebilir. Yanıt alanlarını ve türlerini belirtmek için -response bayrağını kullanın. Örneğin:

```
ignite scaffold message create-post title body --response id:int,title
```

Yukarıdaki komut, hem bir ID (bir tamsayı) hem de bir başlık (bir dize) döndüren MsgCreatePost'u scaffold edecektir

Mesaj scaffold'u, "ignite scaffold list/map/single" kurallarını izler ve standart ve özel türlere sahip alanları destekler. Ayrıntılar için "ignite scaffold list -help" bölümüne bakın.

```
ignite scaffold message [name] [field1] [field2] ... [flags]
```

#### Seçenekler

```
      --clear-cache        clear the build cache (advanced)
  -d, --desc string        description of the command
  -h, --help               help for message
      --module string      module to add the message into. Default: app's main module
      --no-simulation      disable CRUD simulation scaffolding
  -p, --path string        path of the app (default ".")
  -r, --response strings   response fields
      --signer string      label for the message signer (default: creator)
  -y, --yes                answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blok zinciri, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite scaffold modülü

Özel Cosmos SDK modülü

#### Özet

Yeni bir Cosmos SDK modülüne iskele kurun.

Cosmos SDK modüler bir çerçevedir ve her bağımsız işlevsellik parçası ayrı bir modülde uygulanır. Varsayılan olarak blok zinciriniz bir dizi standart Cosmos SDK modülünü içe aktarır. Blok zincirinizin özel işlevselliğini uygulamak için bir modül iskelesi oluşturun ve uygulamanızın mantığını uygulayın.

Bu komut aşağıdakileri yapar:

* "proto/" içinde modülün protokol tampon dosyalarını içeren bir dizin oluşturur
* "x/" içinde modülün şablon Go kodunu içeren bir dizin oluşturur
* "app/app.go" dosyasını değiştirerek yeni oluşturulan modülü içe aktarır
* "testutil/keeper/" içinde test için bir keeper oluşturma mantığını içeren bir dosya oluşturur amaçlar

Bu komut, "app/app.go" gerekli varsayılan yer tutuculara sahip olmasa bile modül iskelesi ile devam edecektir. Yer tutucular eksikse, modülü içe aktarmak için "app/app.go" dosyasını manuel olarak değiştirmeniz gerekecektir. Modülü içe aktaramazsa komutun başarısız olmasını istiyorsanız, "--require-registration" bayrağını kullanın.

IBC özellikli bir modülü iskelelemek için "--ibc" bayrağını kullanın. IBC özellikli bir modül, IBC'ye özgü mantık ve IBC paketlerini "ignite scaffold packet" ile iskelelemek için yer tutucular eklenmiş normal bir modül gibidir.

Bir modül bir veya daha fazla başka modüle bağımlı olabilir ve onların kaleci yöntemlerini içe aktarabilir. Bağımlılığı olan bir modülü iskelelemek için "--dep" bayrağını kullanın

Örneğin, yeni özel modülünüz "foo" hesaplar arasında belirteç göndermeyi gerektiren bir işleve sahip olabilir. Belirteçleri gönderme yöntemi "bank" modül bekçisinde tanımlanmıştır. Aşağıdaki komutla "bank" bağımlılığı olan bir "foo" modülünü iskeleleyebilirsiniz:

```
ignite scaffold module foo --dep bank
```

Daha sonra "expected\_keepers.go" dosyasında "bank" tutucusundan hangi yöntemleri içe aktarmak istediğinizi tanımlayabilirsiniz.

Ayrıca, hem standart hem de özel modülleri (mevcut olmaları koşuluyla) içerebilen bir bağımlılık listesi ile bir modülü iskeleleyebilirsiniz:

```
ignite scaffold module bar --dep foo,mint,account,FeeGrant
```

Not: "--dep" bayrağı uygulamanıza üçüncü taraf modülleri yüklemez, sadece yeni özel modülünüzün hangi mevcut modüllere bağlı olduğunu belirten ekstra kod oluşturur.

Bir Cosmos SDK modülü parametrelere (veya "params") sahip olabilir. Parametreler, blok zincirinin oluşumunda ayarlanabilen ve blok zinciri çalışırken değiştirilebilen değerlerdir. Parametrelerin bir örneği "mint" modülünün "Enflasyon oranı değişikliği "dir. Bir modül, param isimlerinin bir listesini kabul eden "--params" bayrağı kullanılarak paramlar ile iskele haline getirilebilir. Varsayılan olarak parametreler "string" tipindedir, ancak her parametre için bir tip belirtebilirsiniz. Örneğin:

```
ignite scaffold module foo --params baz:uint,bar:bool
```

Modüller, bağımlılıklar ve paramlar hakkında daha fazla bilgi edinmek için Cosmos SDK belgelerine bakın.

```
ignite scaffold module [name] [flags]
```

#### Seçenekler

```
      --clear-cache            clear the build cache (advanced)
      --dep strings            add a dependency on another module
  -h, --help                   help for module
      --ibc                    add IBC functionality
      --ordering string        channel ordering of the IBC module [none|ordered|unordered] (default "none")
      --params strings         add module parameters
  -p, --path string            path of the app (default ".")
      --require-registration   fail if module can't be registered
  -y, --yes                    answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

ignite scaffold - Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite scaffold paketi̇

Bir IBC paketi göndermek için mesaj

#### Özet

Bir IBC paketini belirli bir IBC özellikli Cosmos SDK modülünde iskeleleyin

```
ignite scaffold packet [packetName] [field1] [field2] ... --module [moduleName] [flags]
```

#### Seçenekler

```
      --ack strings     custom acknowledgment type (field1,field2,...)
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for packet
      --module string   IBC Module to add the packet into
      --no-message      disable send message scaffolding
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

ignite scaffold - Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite scaffold sorgusu

Bir blockchain'den veri almak için sorgu

```
ignite scaffold query [name] [request_field1] [request_field2] ... [flags]
```

#### Seçenekler

```
      --clear-cache        clear the build cache (advanced)
  -d, --desc string        description of the CLI to broadcast a tx with the message
  -h, --help               help for query
      --module string      module to add the query into. Default: app's main module
      --paginated          define if the request can be paginated
  -p, --path string        path of the app (default ".")
  -r, --response strings   response fields
  -y, --yes                answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite scaffold react

React web uygulaması şablonu

```
ignite scaffold react [flags]
```

#### Seçenekler

```
  -h, --help          help for react
  -p, --path string   path to scaffold content of the React app (default "./react")
  -y, --yes           answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blockchain, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite tekli scaffold

Tek bir konumda depolanan veriler için CRUD

```
ignite scaffold single NAME [field]... [flags]
```

#### Seçenekler

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for single
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blok zinciri, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite scaffold tipi

Tip tanımı

```
ignite scaffold type NAME [field]... [flags]
```

#### Seçenekler

```
      --clear-cache     clear the build cache (advanced)
  -h, --help            help for type
      --module string   specify which module to generate code in
      --no-message      skip generating message handling logic
      --no-simulation   skip simulation logic
  -p, --path string     path of the app (default ".")
      --signer string   label for the message signer (default: creator)
  -y, --yes             answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blok zinciri, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite iskele vue

Vue 3 web uygulaması şablonu

```
ignite scaffold vue [flags]
```

#### Seçenekler

```
  -h, --help          help for vue
  -p, --path string   path to scaffold content of the Vue.js app (default "./vue")
  -y, --yes           answers interactive yes/no questions with yes
```

#### Ayrıca Bakınız

[ignite scaffold](https://docs.ignite.com/references/cli#ignite-scaffold) - Yeni bir blok zinciri, modül, mesaj, sorgu ve daha fazlasını oluşturun

### ignite araçları

İleri düzey kullanıcılar için araçlar

#### Seçenekler

```
  -h, --help   help for tools
```

#### Ayrıca Bakınız

[ignite ](https://docs.ignite.com/references/cli#ignite)- Ignite CLI, blok zincirinizi Scaffold etmek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar

[ignite tools ibc-relayer](https://docs.ignite.com/references/cli#ignite-tools-ibc-relayer) - Bir IBC aktarıcısının TypeScript uygulaması

[ignite tools ibc-setup](https://docs.ignite.com/references/cli#ignite-tools-ibc-setup) - Bir aktarıcıyı hızlı bir şekilde kurmak için komutlar topluluğu

[ignite tools protoc](https://docs.ignite.com/references/cli#ignite-tools-protoc) - protoc komutunu çalıştırır

### ignite tools ibc-relayer

Bir IBC aktarıcısının TypeScript uygulaması

```
ignite tools ibc-relayer [--] [...] [flags]
```

**Örnekler**

```
ignite tools ibc-relayer -- -h
```

#### Seçenekler

```
  -h, --help   help for ibc-relayer
```

#### Ayrıca Bakınız

[ignite tools](https://docs.ignite.com/references/cli#ignite-tools) - İleri düzey kullanıcılar için araçlar

### ignite tools ibc-setup

Bir aktarıcıyı hızlı bir şekilde kurmak için komutlar topluluğu

```
ignite tools ibc-setup [--] [...] [flags]
```

**Örnekler**

```
ignite tools ibc-setup -- -h
ignite tools ibc-setup -- init --src relayer_test_1 --dest relayer_test_2
```

#### Seçenekler

```
  -h, --help   help for ibc-setup
```

#### Ayrıca Bakınız

[ignite tools](https://docs.ignite.com/references/cli#ignite-tools) - İleri düzey kullanıcılar için araçlar

### ignite tools protoc

protoc komutunu çalıştırın

#### Özet

Protoc komutu. Global protoc include klasörünü -I ile ayarlamanıza gerek yoktur, otomatik olarak işlenir

```
ignite tools protoc [--] [...] [flags]
```

**Örnekler**

```
ignite tools protoc -- --version
```

#### Seçenekler

```
  -h, --help   help for protoc
```

#### Ayrıca Bakınız

[ignite tools](https://docs.ignite.com/references/cli#ignite-tools) - İleri düzey kullanıcılar için araçlar

### ignite versiyonu

Geçerli yapı bilgilerini yazdırma

```
ignite version [flags]
```

#### Seçenekler

```
  -h, --help   help for version
```

#### Ayrıca Bakınız

* [ignite ](https://docs.ignite.com/references/cli#ignite)- Ignite CLI, blockchaininizi iskelelemek, test etmek, oluşturmak ve başlatmak için ihtiyacınız olan her şeyi sunar

****

### &#x20;<a href="#ignite-scaffold-chain" id="ignite-scaffold-chain"></a>

