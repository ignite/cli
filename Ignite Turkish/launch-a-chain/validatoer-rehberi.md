# Validatör Rehberi

Validatörler, Ignite Chain'de zincir launch'ları için genesis validatörleri olarak katılırlar.

Validatörler, Ignite'ta başlatılmak üzere launch edilmiş zincirleri listeleyebilir ve keşfedebilir.

**Çıktı**

```
Launch Id  Chain Id  Source                              Phase

3   example-1   https://github.com/ignite/example   coordinating
2   spn-10      https://github.com/tendermint/spn   launched
1   example-20  https://github.com/tendermint/spn   launching
```

* `Launch ID` , Ignite'ta zincirin benzersiz tanımlayıcısıdır. Bu, zincir launch'ı ile etkileşim kurmak için kullanılan kimliktir.
* `Chain ID` , başlatıldıktan sonra zincir ağının tanımlayıcısını temsil eder. Uygulamada benzersiz bir tanımlayıcı olmalıdır, ancak Ignite'ta benzersiz olması gerekmez.
* `Source`, projenin depo URL'sidir.
* `Phase`, zincir launch'ının mevcut aşamasıdır. Bir zincirin 3 farklı aşaması olabilir:
  * `coordinating`: zincirin validatörlerden istek almaya açık olduğu anlamına gelir
  * `launching`: zincirin artık istek almadığı ancak henüz başlatılmadığı anlamına gelir
  * `launched`: zincir ağının başlatıldığı anlamına gelir

***

Zincir koordinasyon aşamasındayken, validatörler zincir için bir genesis validatörü olmayı talep edebilir. Ignite CLI, validatör için bir node ayarlayabilen otomatik bir iş akışını ve node'ları için belirli bir kuruluma sahip ileri düzey kullanıcılar için bir iş akışını destekler.

#### Basit Akış <a href="#simple-flow" id="simple-flow"></a>

`ignite` validatör kurulumunu otomatik olarak halledebilir. Node'u başlatın ve varsayılan değerlerle bir gentx dosyası oluşturun:

**Çıktı**

```
✔ Source code fetched
✔ Blockchain set up
✔ Blockchain initialized
✔ Genesis initialized
? Staking amount 95000000stake
? Commission rate 0.10
? Commission max rate 0.20
? Commission max change rate 0.01
⋆ Gentx generated: /Users/lucas/spn/3/config/gentx/gentx.json
```

Şimdi, bir validatör olarak bir zincire katılmak için bir istek oluşturun ve yayınlayın:

```
ignite n chain join 3 --amount 100000000stake
```

join komutu, virgülle ayrılmış bir belirteç listesiyle birlikte `--amount` bayrağını kabul eder. Bayrak sağlanırsa, komut doğrulayıcının adresini belirli bir miktarla genesis'e bir hesap olarak eklemek için bir istek yayınlayacaktır.

**Çıktı**

```
? Peer's address 192.168.0.1:26656
✔ Source code fetched
✔ Blockchain set up
✔ Account added to the network by the coordinator!
✔ Validator added to the network by the coordinator!
```

***

#### Gelişmiş Akış <a href="#advanced-flow" id="advanced-flow"></a>

Daha gelişmiş bir kurulum (örn. özel `gentx`) kullanan validatörler, özel dosyaya işaret etmek için komutlarına ek bir bayrak sağlamalıdır:

```
ignite n chain join 3 --amount 100000000stake --gentx ~/chain/config/gentx/gentx.json
```

***

#### Basit Akış <a href="#simple-flow-1" id="simple-flow-1"></a>

Node'un son oluşumunu ve yapılandırmasını oluşturun:

**Çıktı**

```
✔ Source code fetched
✔ Blockchain set up
✔ Chain's binary built
✔ Genesis initialized
✔ Genesis built
✔ Chain is prepared for launch
```

Ardından, node'u başlatın:

```
exampled start --home ~/spn/3
```

***

#### Gelişmiş Akış <a href="#advanced-flow-1" id="advanced-flow-1"></a>

Zincir için son oluşumu getirin:

```
ignite n chain show genesis 3
```

**Çıktı**

```
✔ Source code fetched
✔ Blockchain set up
✔ Blockchain initialized
✔ Genesis initialized
✔ Genesis built
⋆ Genesis generated: ./genesis.json
```

Ardından, kalıcı eş listesini getirin:

```
ignite n chain show peers 3
```

**Çıktı**

```
⋆ Peer list generated: ./peers.txt
```

Getirilen genesis dosyası ve eş listesi manuel düğüm kurulumu için kullanılabilir.
