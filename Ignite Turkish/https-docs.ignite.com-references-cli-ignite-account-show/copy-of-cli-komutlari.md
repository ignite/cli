# Copy of CLI Komutları

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
